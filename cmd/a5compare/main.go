package main

import (
	"a5go"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type compareRequest struct {
	LonLatCases       []lonLatCase       `json:"lonLatCases"`
	CellCases         []cellCase         `json:"cellCases"`
	SphericalCapCases []sphericalCapCase `json:"sphericalCapCases"`
}

type lonLatCase struct {
	Name       string     `json:"name"`
	LonLat     [2]float64 `json:"lonLat"`
	Resolution int        `json:"resolution"`
}

type cellCase struct {
	CellHex string `json:"cellHex"`
}

type sphericalCapCase struct {
	CellHex string  `json:"cellHex"`
	Radius  float64 `json:"radius"`
}

type tsResponse struct {
	LonLatCases       []tsLonLatResult `json:"lonLatCases"`
	CellCases         []tsCellResult   `json:"cellCases"`
	SphericalCapCases []tsCapResult    `json:"sphericalCapCases"`
}

type tsLonLatResult struct {
	Name               string       `json:"name"`
	Resolution         int          `json:"resolution"`
	LonLat             [2]float64   `json:"lonLat"`
	LonLatToCellHex    string       `json:"lonLatToCellHex"`
	CellToLonLat       [2]float64   `json:"cellToLonLat"`
	CellToSpherical    [2]float64   `json:"cellToSpherical"`
	CellToBoundary1    [][2]float64 `json:"cellToBoundary1"`
	CellToBoundaryAuto [][2]float64 `json:"cellToBoundaryAuto"`
}

type tsCellResult struct {
	CellHex            string       `json:"cellHex"`
	CellToLonLat       [2]float64   `json:"cellToLonLat"`
	CellToSpherical    [2]float64   `json:"cellToSpherical"`
	CellToBoundary1    [][2]float64 `json:"cellToBoundary1"`
	CellToBoundaryAuto [][2]float64 `json:"cellToBoundaryAuto"`
}

type tsCapResult struct {
	CellHex      string   `json:"cellHex"`
	Radius       float64  `json:"radius"`
	CompactedHex []string `json:"compactedHex"`
	FlatHex      []string `json:"flatHex"`
}

type placesFixture struct {
	Features []struct {
		Properties struct {
			Name string `json:"name"`
		} `json:"properties"`
		Geometry struct {
			Coordinates [2]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

type capFixture struct {
	SphericalCap []struct {
		CellID string  `json:"cellId"`
		Radius float64 `json:"radius"`
	} `json:"sphericalCap"`
}

func main() {
	var (
		tsRepo      = flag.String("ts-repo", "", "path to a local checkout of the official TypeScript a5 repo")
		pointsLimit = flag.Int("points", 50, "number of populated-place points to compare")
		maxRes      = flag.Int("max-res", 8, "highest resolution to compare for lon/lat indexing")
	)
	flag.Parse()

	repoPath := strings.TrimSpace(*tsRepo)
	if repoPath == "" {
		repoPath = strings.TrimSpace(os.Getenv("A5_TS_REPO"))
	}
	if repoPath == "" {
		fatalf("missing TypeScript repo path; pass --ts-repo /path/to/a5 or set A5_TS_REPO")
	}

	request, err := buildRequest(*pointsLimit, *maxRes)
	if err != nil {
		fatalf("%v", err)
	}

	response, err := runTS(repoPath, request)
	if err != nil {
		fatalf("%v", err)
	}

	failures := compareResults(request, response)
	if len(failures) > 0 {
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Printf("comparison passed: %d lon/lat cases, %d cell cases, %d spherical-cap cases\n",
		len(request.LonLatCases), len(request.CellCases), len(request.SphericalCapCases))
}

func buildRequest(pointsLimit, maxRes int) (compareRequest, error) {
	root, err := os.Getwd()
	if err != nil {
		return compareRequest{}, err
	}

	var places placesFixture
	loadJSONFile(filepath.Join(root, "testdata", "data", "ne_50m_populated_places_nameonly.json"), &places)
	testIDsPath := filepath.Join(root, "testdata", "test-ids.json")
	capPath := filepath.Join(root, "testdata", "fixtures", "traversal", "cap.json")

	var testIDs []string
	loadJSONFile(testIDsPath, &testIDs)

	var caps capFixture
	loadJSONFile(capPath, &caps)

	request := compareRequest{}
	limit := min(pointsLimit, len(places.Features))
	for i := 0; i < limit; i++ {
		feature := places.Features[i]
		for resolution := 1; resolution <= maxRes; resolution++ {
			request.LonLatCases = append(request.LonLatCases, lonLatCase{
				Name:       feature.Properties.Name,
				LonLat:     feature.Geometry.Coordinates,
				Resolution: resolution,
			})
		}
	}

	for _, cellHex := range testIDs {
		request.CellCases = append(request.CellCases, cellCase{CellHex: cellHex})
	}

	for _, c := range caps.SphericalCap {
		request.SphericalCapCases = append(request.SphericalCapCases, sphericalCapCase{
			CellHex: c.CellID,
			Radius:  c.Radius,
		})
	}

	return request, nil
}

func runTS(tsRepo string, request compareRequest) (tsResponse, error) {
	tsRepoAbs, err := filepath.Abs(tsRepo)
	if err != nil {
		return tsResponse{}, err
	}

	scriptPath, err := filepath.Abs("scripts/ts_compare_runner.mjs")
	if err != nil {
		return tsResponse{}, err
	}

	distPath := filepath.Join(tsRepoAbs, "dist", "a5.cjs")
	if _, err := os.Stat(distPath); err != nil {
		return tsResponse{}, fmt.Errorf("TypeScript build not found at %s; run `npm install` and `npm run build` in %s first", distPath, tsRepoAbs)
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return tsResponse{}, err
	}

	cmd := exec.Command("node", scriptPath, tsRepoAbs)
	cmd.Stdin = bytes.NewReader(payload)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return tsResponse{}, fmt.Errorf("ts runner failed: %s", strings.TrimSpace(stderr.String()))
		}
		return tsResponse{}, err
	}

	var response tsResponse
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		return tsResponse{}, fmt.Errorf("decode ts response: %w", err)
	}
	return response, nil
}

func compareResults(request compareRequest, response tsResponse) []string {
	failures := []string{}
	if len(request.LonLatCases) != len(response.LonLatCases) {
		failures = append(failures, fmt.Sprintf("lon/lat case count mismatch: go=%d ts=%d", len(request.LonLatCases), len(response.LonLatCases)))
		return failures
	}
	if len(request.CellCases) != len(response.CellCases) {
		failures = append(failures, fmt.Sprintf("cell case count mismatch: go=%d ts=%d", len(request.CellCases), len(response.CellCases)))
		return failures
	}
	if len(request.SphericalCapCases) != len(response.SphericalCapCases) {
		failures = append(failures, fmt.Sprintf("spherical-cap case count mismatch: go=%d ts=%d", len(request.SphericalCapCases), len(response.SphericalCapCases)))
		return failures
	}

	for i, tc := range request.LonLatCases {
		ts := response.LonLatCases[i]
		cellID := a5go.LonLatToCell(a5go.LonLat(tc.LonLat), tc.Resolution)
		cellHex := a5go.U64ToHex(cellID)
		if cellHex != ts.LonLatToCellHex {
			failures = append(failures, fmt.Sprintf("lonLatToCell mismatch for %s r=%d: go=%s ts=%s", tc.Name, tc.Resolution, cellHex, ts.LonLatToCellHex))
			continue
		}
		failures = append(failures, compareCellOutputs("lon/lat "+tc.Name, cellID, ts.CellToLonLat, ts.CellToSpherical, ts.CellToBoundary1, ts.CellToBoundaryAuto)...)
	}

	for i, tc := range request.CellCases {
		ts := response.CellCases[i]
		cellID, err := a5go.HexToU64(tc.CellHex)
		if err != nil {
			failures = append(failures, fmt.Sprintf("invalid cell hex %s: %v", tc.CellHex, err))
			continue
		}
		failures = append(failures, compareCellOutputs("cell "+tc.CellHex, cellID, ts.CellToLonLat, ts.CellToSpherical, ts.CellToBoundary1, ts.CellToBoundaryAuto)...)
	}

	for i, tc := range request.SphericalCapCases {
		ts := response.SphericalCapCases[i]
		cellID, err := a5go.HexToU64(tc.CellHex)
		if err != nil {
			failures = append(failures, fmt.Sprintf("invalid cap cell hex %s: %v", tc.CellHex, err))
			continue
		}
		compacted := a5go.SphericalCap(cellID, tc.Radius)
		gotCompacted := make([]string, len(compacted))
		for j, cell := range compacted {
			gotCompacted[j] = a5go.U64ToHex(cell)
		}
		if diff := compareStringSlices(gotCompacted, ts.CompactedHex); diff != "" {
			failures = append(failures, fmt.Sprintf("sphericalCap compacted mismatch for %s radius=%.2f: %s", tc.CellHex, tc.Radius, diff))
			continue
		}

		targetRes := a5go.GetResolution(cellID)
		flat := a5go.Uncompact(compacted, targetRes)
		gotFlat := make([]string, len(flat))
		for j, cell := range flat {
			gotFlat[j] = a5go.U64ToHex(cell)
		}
		if diff := compareStringSlices(gotFlat, ts.FlatHex); diff != "" {
			failures = append(failures, fmt.Sprintf("sphericalCap flat mismatch for %s radius=%.2f: %s", tc.CellHex, tc.Radius, diff))
		}
	}

	return failures
}

func compareCellOutputs(label string, cellID uint64, wantLonLat, wantSpherical [2]float64, wantBoundary1, wantBoundaryAuto [][2]float64) []string {
	failures := []string{}
	gotLonLat := a5go.CellToLonLat(cellID)
	if !closePair([2]float64(gotLonLat), wantLonLat, 1e-9) {
		failures = append(failures, fmt.Sprintf("%s cellToLonLat mismatch: go=%v ts=%v", label, gotLonLat, wantLonLat))
	}

	gotSpherical := a5go.CellToSpherical(cellID)
	if !closePair([2]float64(gotSpherical), wantSpherical, 1e-9) {
		failures = append(failures, fmt.Sprintf("%s cellToSpherical mismatch: go=%v ts=%v", label, gotSpherical, wantSpherical))
	}

	gotBoundary1 := a5go.CellToBoundary(cellID, a5go.CellBoundaryOptions{ClosedRing: true, Segments: 1})
	if diff := compareBoundary(gotBoundary1, wantBoundary1, 1e-9); diff != "" {
		failures = append(failures, fmt.Sprintf("%s cellToBoundary(1) mismatch: %s", label, diff))
	}

	gotBoundaryAuto := a5go.CellToBoundary(cellID, a5go.CellBoundaryOptions{ClosedRing: true, AutoSegments: true})
	if diff := compareBoundary(gotBoundaryAuto, wantBoundaryAuto, 1e-9); diff != "" {
		failures = append(failures, fmt.Sprintf("%s cellToBoundary(auto) mismatch: %s", label, diff))
	}

	return failures
}

func compareBoundary(got []a5go.LonLat, want [][2]float64, tol float64) string {
	if len(got) != len(want) {
		return fmt.Sprintf("length go=%d ts=%d", len(got), len(want))
	}
	for i := range got {
		if !closePair([2]float64(got[i]), want[i], tol) {
			return fmt.Sprintf("index %d go=%v ts=%v", i, got[i], want[i])
		}
	}
	return ""
}

func compareStringSlices(got, want []string) string {
	if len(got) != len(want) {
		return fmt.Sprintf("length go=%d ts=%d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			return fmt.Sprintf("index %d go=%s ts=%s", i, got[i], want[i])
		}
	}
	return ""
}

func closePair(got, want [2]float64, tol float64) bool {
	return math.Abs(got[0]-want[0]) <= tol && math.Abs(got[1]-want[1]) <= tol
}

func loadJSONFile(path string, target any) {
	data, err := os.ReadFile(path)
	if err != nil {
		fatalf("read %s: %v", path, err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		fatalf("unmarshal %s: %v", path, err)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
