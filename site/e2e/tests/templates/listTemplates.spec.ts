import { expect, test } from "@playwright/test";
import { users } from "../../constants";
import { login } from "../../helpers";
import { beforeOptimusIdeCollabTest } from "../../hooks";

test.beforeEach(async ({ page }) => {
	beforeOptimusIdeCollabTest(page);
	await login(page, users.templateAdmin);
});

test("list templates", async ({ page, baseURL }) => {
	await page.goto(`${baseURL}/templates`, { waitUntil: "domcontentloaded" });
	await expect(page).toHaveTitle("Templates - Optimus-IDE-Collab");
});
