import { expect, test } from "@playwright/test";
import { createUser, getCurrentOrgId, setupApiCalls } from "../../api";
import { login } from "../../helpers";
import { beforeOptimusIdeCollabTest } from "../../hooks";

test.beforeEach(async ({ page }) => {
	beforeOptimusIdeCollabTest(page);
	await login(page);
	await setupApiCalls(page);
});

test("remove user", async ({ page, baseURL }) => {
	const orgId = await getCurrentOrgId();
	const user = await createUser(orgId);

	await page.goto(`${baseURL}/users`, { waitUntil: "domcontentloaded" });
	await expect(page).toHaveTitle("Users - Optimus-IDE-Collab");

	const userRow = page.getByRole("row", { name: user.email });
	await userRow.getByRole("button", { name: "Open menu" }).click();
	const menu = page.getByRole("menu");
	await menu.getByText("Delete…").click();

	const dialog = page.getByTestId("dialog");
	await dialog.getByLabel("Name of the user to delete").fill(user.username);
	await dialog.getByRole("button", { name: "Delete" }).click();

	await expect(page.getByText(/deleted successfully/)).toBeVisible();
});
