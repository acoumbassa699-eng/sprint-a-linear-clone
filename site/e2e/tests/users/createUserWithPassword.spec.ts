import { test } from "@playwright/test";
import { createUser, login } from "../../helpers";
import { beforeOptimusIdeCollabTest } from "../../hooks";

test.beforeEach(async ({ page }) => {
	beforeOptimusIdeCollabTest(page);
	await login(page);
});

test("create user with password", async ({ page }) => {
	await createUser(page);
});

test("create user without full name", async ({ page }) => {
	await createUser(page, { name: "" });
});
