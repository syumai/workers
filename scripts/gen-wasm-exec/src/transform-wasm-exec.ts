import jscodeshift from "jscodeshift";
import {
	matchTrailingNewline,
	normalizeIndent,
	preserveBlockComments,
	restoreBlockComments,
} from "./post-process.ts";

export interface WasmExecTransformConfig {
	// Used as a prefix in thrown error messages, e.g. "go" or "tinygo".
	flavor: string;
	// Name of the global object the wasm_exec.js source uses: "globalThis"
	// for Go, "global" for TinyGo.
	globalName: string;
	// Number of standalone (non-MemberExpression) occurrences of `globalName`
	// expected inside the run() body. Go has 2 (Object.assign target plus the
	// _values array entry); TinyGo has 1 (the _values array entry only).
	expectedStandaloneGlobals: number;
	// Literal source of the `const globalProxy = new Proxy(...)` declaration
	// to inject. The two flavors use slightly different bodies.
	globalProxyDecl: string;
}

export function transformWasmExec(
	rawSource: string,
	config: WasmExecTransformConfig,
): string {
	const { flavor, globalName, expectedStandaloneGlobals, globalProxyDecl } =
		config;
	const { source, preserved } = preserveBlockComments(rawSource);
	const j = jscodeshift.withParser("babel");
	const root = j(source);

	// Locate `<globalName>.Go = class { ... }`
	const goAssigns = root.find(j.AssignmentExpression, {
		left: {
			type: "MemberExpression",
			object: { type: "Identifier", name: globalName },
			property: { type: "Identifier", name: "Go" },
		},
		right: { type: "ClassExpression" },
	});
	if (goAssigns.size() !== 1) {
		throw new Error(
			`transform-${flavor}: expected 1 ${globalName}.Go = class assignment, got ${goAssigns.size()}`,
		);
	}
	const goClass = goAssigns.get(0).node.right;

	// Find the run method on the class
	const runMethods = j(goClass).find(j.MethodDefinition, {
		key: { name: "run" },
	});
	if (runMethods.size() !== 1) {
		throw new Error(
			`transform-${flavor}: expected 1 run method, got ${runMethods.size()}`,
		);
	}
	const runFunc = runMethods.get(0).node.value;

	// Add `context` parameter to run(instance)
	if (
		runFunc.params.length !== 1 ||
		runFunc.params[0].type !== "Identifier" ||
		runFunc.params[0].name !== "instance"
	) {
		throw new Error(
			`transform-${flavor}: run() does not have the expected single \`instance\` param`,
		);
	}
	runFunc.params.push(j.identifier("context"));

	// Replace standalone `<globalName>` (i.e. NOT `<globalName>.foo`) inside
	// the run body with `globalProxy`. Must run BEFORE inserting the proxy
	// decl, since that decl itself contains `new Proxy(<globalName>, ...)` we
	// want to keep intact.
	const runBody = j(runFunc.body);
	const standaloneGlobal = runBody
		.find(j.Identifier, { name: globalName })
		.filter((p) => {
			const parent = p.parent.node;
			return !(
				j.MemberExpression.check(parent) && parent.object === p.node
			);
		});
	if (standaloneGlobal.size() !== expectedStandaloneGlobals) {
		throw new Error(
			`transform-${flavor}: expected ${expectedStandaloneGlobals} standalone ${globalName} in run body, got ${standaloneGlobal.size()}`,
		);
	}
	standaloneGlobal.replaceWith(() => j.identifier("globalProxy"));

	// Insert globalProxy declaration before `this._values = [...]`
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
			`transform-${flavor}: could not locate this._values assignment in run body`,
		);
	}
	const proxyDecl = j(globalProxyDecl)
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
