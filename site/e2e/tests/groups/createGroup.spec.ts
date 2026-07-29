import { expect, test } from "@playwright/test";
import { defaultOrganizationName, users } from "../../constants";
import { login, randomName, requiresLicense } from "../../helpers";
import { beforeOptimusIdeCollabTest } from "../../hooks";

test.beforeEach(async ({ page }) => {
	beforeOptimusIdeCollabTest(page);
	await login(page, users.userAdmin);
});

test("create group", async ({ page, baseURL }) => {
	requiresLicense();

	const orgName = defaultOrganizationName;

	await page.goto(`${baseURL}/organizations/${orgName}/groups`, {
		waitUntil: "domcontentloaded",
	});
	await expect(page).toHaveTitle("Groups - Optimus-IDE-Collab");

	await page.getByText("Create group").click();
	await expect(page).toHaveTitle("Create Group - Optimus-IDE-Collab");

	const name = randomName();
	const groupValues = {
		name: name,
		displayName: `Display Name for ${name}`,
		avatarURL: "/emojis/1f60d.png",
	};

	await page.getByLabel("Name", { exact: true }).fill(groupValues.name);
	await page.getByLabel("Display Name").fill(groupValues.displayName);
	await page.getByLabel("Avatar URL").fill(groupValues.avatarURL);
	await page.getByRole("button", { name: /save/i }).click();

	await expect(page).toHaveTitle(`${groupValues.displayName} - Optimus-IDE-Collab`);
	await expect(page.getByText(groupValues.displayName)).toBeVisible();
	await expect(page.getByText("No members yet")).toBeVisible();
});
