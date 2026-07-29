import type { Meta, StoryObj } from "@storybook/react-vite";
import { InlineMarkdown } from "./InlineMarkdown";

const meta: Meta<typeof InlineMarkdown> = {
	title: "components/Markdown/InlineMarkdown",
	component: InlineMarkdown,
};

export default meta;
type Story = StoryObj<typeof InlineMarkdown>;

export const WithFormatting: Story = {
	args: {
		children: "This supports **bold** and *italic* text.",
	},
};

export const WithLink: Story = {
	args: {
		children: "Read the [documentation](https://optimus-ide-collab.com/docs).",
	},
};

export const WithCode: Story = {
	args: {
		children: "Run `optimus-ide-collab templates push` to publish your template.",
	},
};
