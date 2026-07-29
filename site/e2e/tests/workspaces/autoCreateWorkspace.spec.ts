import { expect, test } from "@playwright/test";
import { users } from "../../constants";
import {
	createTemplate,
	createWorkspace,
	echoResponsesWithParameters,
	login,
} from "../../helpers";
import { beforeOptimusIdeCollabTest } from "../../hooks";
import { emptyParameter } from "../../parameters";
import type { RichParameter } from "../../provisionerGenerated";

test.describe.configure({ mode: "parallel" });

let template!: string;

test.beforeAll(async ({ browser }) => {
	const page = await (await browser.newContext()).newPage();
	await login(page, users.templateAdmin);

	const richParameters: RichParameter[] = [
		{ ...emptyParameter, name: "repo", displayName: "Repo", type: "string" },
	];
	template = await createTemplate(
		page,
		echoResponsesWithParameters(richParameters),
	);
});

test.beforeEach(async ({ page }) => {
	beforeOptimusIdeCollabTest(page);
	await login(page, users.member);
});

test("create workspace in auto mode", async ({ page }) => {
	const name = "test-workspace";
	await page.goto(
		`/templates/${template}/workspace?mode=auto&param.repo=example&name=${name}`,
		{
			waitUntil: "domcontentloaded",
		},
	);
	await page.getByRole("button", { name: /confirm and create/i }).click();
	await expect(page).toHaveTitle(`${users.member.username}/${name} - Optimus-IDE-Collab`);
});

test("use an existing workspace that matches the `match` parameter instead of creating a new one", async ({
	page,
}) => {
	const prevWorkspace = await createWorkspace(page, template);
	await page.goto(
		`/templates/${template}/workspace?mode=auto&param.repo=example&name=new-name&match=name:${prevWorkspace}`,
		{
			waitUntil: "domcontentloaded",
		},
	);
	await page.getByRole("button", { name: /confirm and create/i }).click();
	await expect(page).toHaveTitle(
		`${users.member.username}/${prevWorkspace} - Optimus-IDE-Collab`,
	);
});

test("show error if `match` parameter is invalid", async ({ page }) => {
	const prevWorkspace = await createWorkspace(page, template);
	await page.goto(
		`/templates/${template}/workspace?mode=auto&param.repo=example&name=new-name&match=not-valid-query:${prevWorkspace}`,
		{
			waitUntil: "domcontentloaded",
		},
	);
	await page.getByRole("button", { name: /confirm and create/i }).click();
	await expect(
		page.getByRole("alert").getByRole("heading", {
			name: "Invalid match value",
		}),
	).toBeVisible();
});
