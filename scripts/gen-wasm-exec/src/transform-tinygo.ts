import jscodeshift from "jscodeshift";
import {
	matchTrailingNewline,
	normalizeIndent,
	preserveBlockComments,
	restoreBlockComments,
} from "./post-process.ts";

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

const GLOBAL_PROXY_DECL = `const globalProxy = new Proxy(global, {
\tget(target, prop) {
\t\tif (prop === 'context') {
\t\t\treturn context;
\t\t}
\t\treturn Reflect.get(...arguments);
\t}
})`;

export function transformTinygo(rawSource: string): string {
	const wrapped = commentOutNodejsBlocks(rawSource);
	const { source, preserved } = preserveBlockComments(wrapped);
	const j = jscodeshift.withParser("babel");
	const root = j(source);

	// Locate `global.Go = class { ... }`
	const goAssigns = root.find(j.AssignmentExpression, {
		left: {
			type: "MemberExpression",
			object: { type: "Identifier", name: "global" },
			property: { type: "Identifier", name: "Go" },
		},
		right: { type: "ClassExpression" },
	});
	if (goAssigns.size() !== 1) {
		throw new Error(
			`transform-tinygo: expected 1 global.Go = class assignment, got ${goAssigns.size()}`,
		);
	}
	const goClass = goAssigns.get(0).node.right;

	// Find the run method
	const runMethods = j(goClass).find(j.MethodDefinition, {
		key: { name: "run" },
	});
	if (runMethods.size() !== 1) {
		throw new Error(
			`transform-tinygo: expected 1 run method, got ${runMethods.size()}`,
		);
	}
	const runFunc = runMethods.get(0).node.value;

	// P1: add `context` parameter
	if (
		runFunc.params.length !== 1 ||
		runFunc.params[0].type !== "Identifier" ||
		runFunc.params[0].name !== "instance"
	) {
		throw new Error(
			"transform-tinygo: run() does not have the expected single `instance` param",
		);
	}
	runFunc.params.push(j.identifier("context"));

	// P2a: replace standalone `global` inside the run body with `globalProxy`.
	// Must run BEFORE inserting the proxy decl (which references `global`).
	// In TinyGo's run body, the only standalone `global` is in the _values array.
	const runBody = j(runFunc.body);
	const standaloneGlobal = runBody
		.find(j.Identifier, { name: "global" })
		.filter((p) => {
			const parent = p.parent.node;
			return !(
				j.MemberExpression.check(parent) && parent.object === p.node
			);
		});
	if (standaloneGlobal.size() !== 1) {
		throw new Error(
			`transform-tinygo: expected 1 standalone global in run body, got ${standaloneGlobal.size()}`,
		);
	}
	standaloneGlobal.replaceWith(() => j.identifier("globalProxy"));

	// P2b: insert globalProxy declaration before `this._values = [...]`
	const valuesStmt = runBody
		.find(j.AssignmentExpression, {
			left: {
				type: "MemberExpression",
				object: { type: "ThisExpression" },
				property: { name: "_values" },
			},
		})
		.closest(j.ExpressionStatement);
	if (valuesStmt.size() !== 1) {
		throw new Error(
			"transform-tinygo: could not locate this._values assignment in run body",
		);
	}
	const proxyDecl = j(GLOBAL_PROXY_DECL)
		.find(j.VariableDeclaration)
		.get(0).node;
	valuesStmt.insertBefore(proxyDecl);

	let output = root.toSource({
		lineTerminator: "\n",
		useTabs: true,
		reuseWhitespace: true,
	});
	output = restoreBlockComments(output, preserved);
	output = normalizeIndent(output);
	output = matchTrailingNewline(rawSource, output);
	return output;
}
