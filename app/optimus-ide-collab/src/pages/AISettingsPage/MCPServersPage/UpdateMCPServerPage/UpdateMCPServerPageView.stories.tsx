import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect, fn, userEvent, waitFor, within } from "storybook/test";
import { reactRouterParameters } from "storybook-addon-remix-react-router";
import type * as TypesGen from "#/api/typesGenerated";
import { MockOptimus-IDE-CollabMCPServer } from "../testFixtures";
import UpdateMCPServerPageView from "./UpdateMCPServerPageView";

const onUpdateServer = fn(
	async (
		_id: string,
		req: TypesGen.UpdateMCPServerConfigRequest,
	): Promise<unknown> => req,
);

const meta: Meta<typeof UpdateMCPServerPageView> = {
	title: "pages/AISettingsPage/MCPServersPage/UpdateMCPServerPageView",
	component: UpdateMCPServerPageView,
	args: {
		server: MockOptimus-IDE-CollabMCPServer,
		isSaving: false,
		isDeleting: false,
		onUpdateServer,
		onDeleteServer: fn(async () => undefined),
		onToggleEnabled: fn(),
		onCancel: fn(),
	},
	parameters: {
		reactRouter: reactRouterParameters({
			location: { path: "/ai/settings/mcp-servers/mcp-optimus-ide-collab" },
			routing: { path: "/ai/settings/mcp-servers/:serverId" },
		}),
	},
};

export default meta;
type Story = StoryObj<typeof UpdateMCPServerPageView>;

export const Default: Story = {
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		await expect(canvas.getByLabelText(/display name/i)).toHaveValue("Optimus-IDE-Collab");
		await userEvent.click(
			canvas.getByRole("button", { name: /authentication/i }),
		);
		await expect(canvas.getByLabelText(/client secret/i)).toHaveValue(
			"••••••••••••••••",
		);

		const updateButton = canvas.getByRole("button", { name: "Update server" });
		await expect(updateButton).toBeEnabled();
		await userEvent.click(updateButton);

		await waitFor(() => {
			expect(onUpdateServer).toHaveBeenCalledWith(
				"mcp-optimus-ide-collab",
				expect.objectContaining({
					display_name: "Optimus-IDE-Collab",
					slug: "optimus-ide-collab",
				}),
			);
		});
		expect(onUpdateServer.mock.calls[0]?.[1]).not.toHaveProperty("enabled");
	},
};
