export interface PreservedComment {
	id: string;
	content: string;
}

// recast/jscodeshift drops or relocates some block comments whose contents
// happen to be valid JavaScript (e.g. tinygo's commented-out require() blocks).
// To survive the round-trip, we replace each block comment with a unique bare
// identifier so it appears in the AST as a normal ExpressionStatement, then
// swap the identifier back to the original comment after printing.
// This assumes block comments only occur at statement positions in the inputs
// we care about (Go and TinyGo wasm_exec.js both satisfy this).
export function preserveBlockComments(source: string): {
	source: string;
	preserved: PreservedComment[];
} {
	const preserved: PreservedComment[] = [];
	let counter = 0;
	const newSource = source.replace(/\/\*[\s\S]*?\*\//g, (match) => {
		const id = `__GENWASM_PRESERVE_${counter++}__`;
		preserved.push({ id, content: match });
		return id;
	});
	return { source: newSource, preserved };
}

export function restoreBlockComments(
	source: string,
	preserved: PreservedComment[],
): string {
	let result = source;
	for (const p of preserved) {
		// Match the identifier optionally followed by a recast-emitted `;`
		const re = new RegExp(`${p.id};?`);
		result = result.replace(re, p.content);
	}
	return result;
}

// recast emits a mix of tabs and 4-space indents around modified nodes when
// useTabs is on but the surrounding source uses tabs. Collapse runs of 4
// spaces that appear AFTER leading tabs into additional tabs. Leave any
// pre-tab whitespace alone so we don't mangle hand-written quirks like
// "<space><tabs>return …" that exist in upstream sources.
export function normalizeIndent(source: string): string {
	return source
		.split("\n")
		.map((line) => {
			const m = line.match(/^(\t*)( +)/);
			if (!m) return line;
			const tabs = m[1];
			const spaces = m[2];
			const extraTabs = Math.floor(spaces.length / 4);
			const remainingSpaces = spaces.length % 4;
			return (
				tabs +
				"\t".repeat(extraTabs) +
				" ".repeat(remainingSpaces) +
				line.slice(tabs.length + spaces.length)
			);
		})
		.join("\n");
}

export function matchTrailingNewline(
	originalSource: string,
	output: string,
): string {
	if (originalSource.endsWith("\n")) {
		return output.endsWith("\n") ? output : `${output}\n`;
	}
	return output.replace(/\n+$/, "");
}
