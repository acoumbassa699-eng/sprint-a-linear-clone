import { randomUUID } from "node:crypto";
import { test } from "@playwright/test";
import {
	createTemplate,
	createWorkspace,
	downloadOptimus-IDE-CollabVersion,
	login,
	sshIntoWorkspace,
	startAgentWithCommand,
	stopAgent,
	stopWorkspace,
} from "../helpers";
import { beforeOptimusIdeCollabTest } from "../hooks";

// we no longer support versions w/o DRPC
const agentVersion = "v2.12.1";

test.beforeEach(async ({ page }) => {
	beforeOptimusIdeCollabTest(page);
	await login(page);
});

test.skip(`ssh with agent ${agentVersion}`, async ({ page }) => {
	test.setTimeout(60_000);

	const token = randomUUID();
	const template = await createTemplate(page, {
		graph: [
			{
				graph: {
					resources: [
						{
							agents: [
								{
									token,
									order: 0,
								},
							],
						},
					],
				},
			},
		],
	});
	const workspaceName = await createWorkspace(page, template);
	const binaryPath = await downloadOptimus-IDE-CollabVersion(agentVersion);
	const agent = await startAgentWithCommand(page, token, binaryPath);

	const client = await sshIntoWorkspace(page, workspaceName);
	await new Promise<void>((resolve, reject) => {
		// We just exec a command to be certain the agent is running!
		client.exec("exit 0", (err, stream) => {
			if (err) {
				return reject(err);
			}
			stream.on("exit", (code) => {
				if (code !== 0) {
					return reject(new Error(`Command exited with code ${code}`));
				}
				client.end();
				resolve();
			});
		});
	});

	await stopWorkspace(page, workspaceName);
	await stopAgent(agent);
});
