import { transformWasmExec } from "./transform-wasm-exec.ts";

// Top-level Node.js-only blocks shipped by upstream TinyGo. Each entry pins
// the exact opening substring; `closes` says how many "\n\t}" lines to skip
// past — TextEncoder is paired with TextDecoder under a single /* */, so it
// needs closes=2.
const NODEJS_BLOCKS: Array<{ opening: string; closes: number }> = [
	{
		opening: '\tif (!global.require && typeof require !== "undefined") {',
		closes: 1,
	},
	{ opening: "\tif (!global.fs && global.require) {", closes: 1 },
	{ opening: "\tif (!global.crypto) {", closes: 1 },
	{ opening: "\tif (!global.TextEncoder) {", closes: 2 },
	{
		opening: "\tif (\n\t\tglobal.require &&\n\t\tglobal.require.main === module &&",
		closes: 1,
	},
];

interface Range {
	start: number;
	end: number;
}

function findBlockRanges(source: string): Range[] {
	const ranges: Range[] = [];
	for (const block of NODEJS_BLOCKS) {
		const start = source.indexOf(block.opening);
		if (start === -1) {
			throw new Error(
				`transform-tinygo: could not find block opening: ${block.opening.slice(0, 50)}…`,
			);
		}
		let pos = start;
		for (let i = 0; i < block.closes; i++) {
			const closeNL = source.indexOf("\n\t}", pos);
			if (closeNL === -1) {
				throw new Error(
					`transform-tinygo: could not find matching close for: ${block.opening.slice(0, 50)}…`,
				);
			}
			pos = closeNL + 3;
		}
		ranges.push({ start, end: pos });
	}
	return ranges;
}

function commentOutNodejsBlocks(rawSource: string): string {
	const ranges = findBlockRanges(rawSource).sort((a, b) => b.start - a.start);
	let result = rawSource;
	for (const m of ranges) {
		const block = result.slice(m.start, m.end);
		result = `${result.slice(0, m.start)}\t/*\n${block}\n\t*/${result.slice(m.end)}`;
	}
	return result;
}

export function transformTinygo(rawSource: string): string {
	return transformWasmExec(commentOutNodejsBlocks(rawSource), {
		flavor: "tinygo",
		globalName: "global",
		expectedStandaloneGlobals: 1,
	});
}
