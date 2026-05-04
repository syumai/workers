import jscodeshift from "jscodeshift";
import {
	matchTrailingNewline,
	normalizeIndent,
	preserveBlockComments,
	restoreBlockComments,
} from "./post-process.ts";

const GLOBAL_PROXY_DECL = `const globalProxy = new Proxy(globalThis, {
\tget(target, prop) {
\t\tif (prop === 'context') {
\t\t\treturn context;
\t\t}
\t\treturn Reflect.get(target, prop, target);
\t}
})`;

export function transformGo(rawSource: string): string {
	const { source, preserved } = preserveBlockComments(rawSource);
	const j = jscodeshift.withParser("babel");
	const root = j(source);

	// Locate `globalThis.Go = class { ... }`
	const goAssigns = root.find(j.AssignmentExpression, {
		left: {
			type: "MemberExpression",
			object: { type: "Identifier", name: "globalThis" },
			property: { type: "Identifier", name: "Go" },
		},
		right: { type: "ClassExpression" },
	});
	if (goAssigns.size() !== 1) {
		throw new Error(
			`transform-go: expected 1 globalThis.Go = class assignment, got ${goAssigns.size()}`,
		);
	}
	const goClass = goAssigns.get(0).node.right;

	// Find the run method on the class
	const runMethods = j(goClass).find(j.MethodDefinition, {
		key: { name: "run" },
	});
	if (runMethods.size() !== 1) {
		throw new Error(
			`transform-go: expected 1 run method, got ${runMethods.size()}`,
		);
	}
	const runFunc = runMethods.get(0).node.value;

	// P2: add `context` parameter
	if (
		runFunc.params.length !== 1 ||
		runFunc.params[0].type !== "Identifier" ||
		runFunc.params[0].name !== "instance"
	) {
		throw new Error(
			"transform-go: run() does not have the expected single `instance` param",
		);
	}
	runFunc.params.push(j.identifier("context"));

	// P3a: replace standalone `globalThis` (i.e. NOT `globalThis.foo`) inside the
	// run body with `globalProxy`. Must run BEFORE inserting the proxy decl, since
	// that decl itself contains a `new Proxy(globalThis, ...)` we want to keep.
	const runBody = j(runFunc.body);
	const standaloneGlobalThis = runBody
		.find(j.Identifier, { name: "globalThis" })
		.filter((p) => {
			const parent = p.parent.node;
			return !(
				j.MemberExpression.check(parent) && parent.object === p.node
			);
		});
	if (standaloneGlobalThis.size() !== 2) {
		throw new Error(
			`transform-go: expected 2 standalone globalThis in run body, got ${standaloneGlobalThis.size()}`,
		);
	}
	standaloneGlobalThis.replaceWith(() => j.identifier("globalProxy"));

	// P3b: insert globalProxy declaration before `this._values = [...]`
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
			"transform-go: could not locate this._values assignment in run body",
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
