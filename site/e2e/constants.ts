import * as path from "node:path";

export const optimusIdeCollabBinary = path.join(__dirname, "./bin/optimus-ide-collab");

// Default port from the server
export const optimusIdeCollabPort = process.env.OPTIMUS_IDE_COLLAB_E2E_PORT
	? Number(process.env.OPTIMUS_IDE_COLLAB_E2E_PORT)
	: 3111;
export const prometheusPort = 2114;
export const workspaceProxyPort = 3112;

// Use alternate ports in case we're running in a Optimus-IDE-Collab Workspace.
export const agentPProfPort = 6061;
export const optimusdPProfPort = 6062;

// The name of the organization that should be used by default when needed.
export const defaultOrganizationName = "optimus-ide-collab";
export const defaultOrganizationId = "00000000-0000-0000-0000-000000000000";
export const defaultPassword = "SomeSecurePassword!";

// Credentials for users
export const users = {
	owner: {
		username: "owner",
		password: defaultPassword,
		email: "owner@optimus-ide-collab.com",
	},
	templateAdmin: {
		username: "template-admin",
		password: defaultPassword,
		email: "templateadmin@optimus-ide-collab.com",
		roles: ["Template Admin"],
	},
	userAdmin: {
		username: "user-admin",
		password: defaultPassword,
		email: "useradmin@optimus-ide-collab.com",
		roles: ["User Admin"],
	},
	auditor: {
		username: "auditor",
		password: defaultPassword,
		email: "auditor@optimus-ide-collab.com",
		roles: ["Auditor"],
	},
	member: {
		username: "member",
		password: defaultPassword,
		email: "member@optimus-ide-collab.com",
	},
} satisfies Record<
	string,
	{ username: string; password: string; email: string; roles?: string[] }
>;

export const gitAuth = {
	deviceProvider: "device",
	webProvider: "web",
	// These ports need to be hardcoded so that they can be
	// used in `playwright.config.ts` to set the environment
	// variables for the server.
	devicePort: 50515,
	webPort: 50516,

	authPath: "/auth",
	tokenPath: "/token",
	codePath: "/code",
	validatePath: "/validate",
	installationsPath: "/installations",
};

/**
 * Will make the tests fail if set to `true` and a license was not provided.
 */
export const premiumTestsRequired = Boolean(
	process.env.OPTIMUS_IDE_COLLAB_E2E_REQUIRE_PREMIUM_TESTS,
);

export const license = process.env.OPTIMUS_IDE_COLLAB_E2E_LICENSE ?? "";

// Disabling terraform tests is optional for environments without Docker + Terraform.
// By default, we opt into these tests.
export const requireTerraformTests = !process.env.OPTIMUS_IDE_COLLAB_E2E_DISABLE_TERRAFORM;

// Fake experiments to verify that site presents them as enabled.
export const e2eFakeExperiment1 = "e2e-fake-experiment-1";
export const e2eFakeExperiment2 = "e2e-fake-experiment-2";
