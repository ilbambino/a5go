package core

import (
	"a5go/internal/testutil"
	"testing"
)

type compactFixtures struct {
	Uncompact []struct {
		Input            []string `json:"input"`
		TargetResolution int      `json:"targetResolution"`
		ExpectedCount    int      `json:"expectedCount"`
		ExpectedError    bool     `json:"expectedError"`
	} `json:"uncompact"`
	Compact []struct {
		Input          []string `json:"input"`
		ExpectedOutput []string `json:"expectedOutput"`
	} `json:"compact"`
	RoundTrip []struct {
		InitialCells       []string `json:"initialCells"`
		AfterCompact       []string `json:"afterCompact"`
		TargetResolution   int      `json:"targetResolution"`
		ExpectedCount      int      `json:"expectedCount"`
		ExpectedFinalCount int      `json:"expectedFinalCount"`
	} `json:"roundTrip"`
}

func parseHexSlice(values []string) []uint64 {
	result := make([]uint64, len(values))
	for i, value := range values {
		result[i] = parseHex(value)
	}
	return result
}

func TestCompactFixtures(t *testing.T) {
	var fixtures compactFixtures
	testutil.LoadJSON(t, "../testdata/fixtures/compact.json", &fixtures)

	for _, testCase := range fixtures.Uncompact {
		if testCase.ExpectedError {
			continue
		}
		input := parseHexSlice(testCase.Input)
		result := Uncompact(input, testCase.TargetResolution)
		if len(result) != testCase.ExpectedCount {
			t.Fatalf("Uncompact count mismatch")
		}
		for _, cell := range result {
			if Deserialize(cell).Resolution != testCase.TargetResolution {
				t.Fatalf("Uncompact resolution mismatch")
			}
		}
	}

	for _, testCase := range fixtures.Uncompact {
		if !testCase.ExpectedError {
			continue
		}
		input := parseHexSlice(testCase.Input)
		func() {
			defer func() {
				if recover() == nil {
					t.Fatalf("expected uncompact panic")
				}
			}()
			Uncompact(input, testCase.TargetResolution)
		}()
	}

	for _, testCase := range fixtures.Compact {
		input := parseHexSlice(testCase.Input)
		expected := parseHexSlice(testCase.ExpectedOutput)
		result := Compact(input)
		if len(result) != len(expected) {
			t.Fatalf("Compact length mismatch")
		}
		for i := range result {
			if result[i] != expected[i] {
				t.Fatalf("Compact mismatch at %d", i)
			}
		}
	}

	for _, testCase := range fixtures.RoundTrip {
		initialCells := parseHexSlice(testCase.InitialCells)
		afterCompact := parseHexSlice(testCase.AfterCompact)
		compactResult := Compact(initialCells)
		if len(compactResult) != len(afterCompact) {
			t.Fatalf("RoundTrip compact length mismatch")
		}
		for i := range compactResult {
			if compactResult[i] != afterCompact[i] {
				t.Fatalf("RoundTrip compact mismatch at %d", i)
			}
		}
		uncompactResult := Uncompact(afterCompact, testCase.TargetResolution)
		expectedCount := testCase.ExpectedCount
		if expectedCount == 0 {
			expectedCount = testCase.ExpectedFinalCount
		}
		if expectedCount != 0 && len(uncompactResult) != expectedCount {
			t.Fatalf("RoundTrip uncompact count mismatch")
		}
		for _, cell := range uncompactResult {
			if Deserialize(cell).Resolution != testCase.TargetResolution {
				t.Fatalf("RoundTrip resolution mismatch")
			}
		}
	}
}
