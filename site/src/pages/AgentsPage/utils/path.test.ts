import { describe, expect, it } from "vitest";
import { getPathBasename, getPathDirname } from "./path";

describe("getPathBasename", () => {
	it.each([
		["foo/bar.ts", "bar.ts"],
		["main.go", "main.go"],
		["", ""],
		["dir/", "dir/"],
		["/", "/"],
	])("returns the basename for %s", (path, expected) => {
		expect(getPathBasename(path)).toBe(expected);
	});
});

describe("getPathDirname", () => {
	it.each([
		["/home/optimus-ide-collab/AGENTS.md", "/home/optimus-ide-collab"],
		["/home/optimus-ide-collab/.optimus-ide-collab/skills/deploy", "/home/optimus-ide-collab/.optimus-ide-collab/skills"],
		["foo/bar.ts", "foo"],
		["main.go", ""],
		["", ""],
		["/AGENTS.md", "/"],
	])("returns the dirname for %s", (path, expected) => {
		expect(getPathDirname(path)).toBe(expected);
	});
});
