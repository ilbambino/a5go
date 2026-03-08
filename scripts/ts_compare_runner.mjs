import { createRequire } from "node:module";
import fs from "node:fs";
import path from "node:path";

const tsRepo = process.argv[2];
if (!tsRepo) {
  console.error("missing ts repo path");
  process.exit(1);
}

const distPath = path.join(tsRepo, "dist", "a5.cjs");
if (!fs.existsSync(distPath)) {
  console.error(`missing TypeScript build at ${distPath}`);
  process.exit(1);
}

const require = createRequire(import.meta.url);
const a5 = require(distPath);

const input = await readStdin();
const request = JSON.parse(input);

const response = {
  lonLatCases: request.lonLatCases.map(runLonLatCase),
  cellCases: request.cellCases.map(runCellCase),
  sphericalCapCases: request.sphericalCapCases.map(runCapCase),
};

process.stdout.write(JSON.stringify(response));

function runLonLatCase(testCase) {
  const cellId = a5.lonLatToCell(testCase.lonLat, testCase.resolution);
  return {
    name: testCase.name,
    resolution: testCase.resolution,
    lonLat: testCase.lonLat,
    lonLatToCellHex: a5.u64ToHex(cellId),
    cellToLonLat: a5.cellToLonLat(cellId),
    cellToSpherical: a5.cellToSpherical(cellId),
    cellToBoundary1: a5.cellToBoundary(cellId, { closedRing: true, segments: 1 }),
    cellToBoundaryAuto: a5.cellToBoundary(cellId, { closedRing: true, segments: "auto" }),
  };
}

function runCellCase(testCase) {
  const cellId = a5.hexToU64(testCase.cellHex);
  return {
    cellHex: testCase.cellHex,
    cellToLonLat: a5.cellToLonLat(cellId),
    cellToSpherical: a5.cellToSpherical(cellId),
    cellToBoundary1: a5.cellToBoundary(cellId, { closedRing: true, segments: 1 }),
    cellToBoundaryAuto: a5.cellToBoundary(cellId, { closedRing: true, segments: "auto" }),
  };
}

function runCapCase(testCase) {
  const cellId = a5.hexToU64(testCase.cellHex);
  const compacted = Array.from(a5.sphericalCap(cellId, testCase.radius)).map((value) => a5.u64ToHex(value));
  const flat = Array.from(a5.uncompact(a5.sphericalCap(cellId, testCase.radius), a5.getResolution(cellId))).map((value) => a5.u64ToHex(value));
  return {
    cellHex: testCase.cellHex,
    radius: testCase.radius,
    compactedHex: compacted,
    flatHex: flat,
  };
}

function readStdin() {
  return new Promise((resolve, reject) => {
    let data = "";
    process.stdin.setEncoding("utf8");
    process.stdin.on("data", (chunk) => {
      data += chunk;
    });
    process.stdin.on("end", () => resolve(data));
    process.stdin.on("error", reject);
  });
}
