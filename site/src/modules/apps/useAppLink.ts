import type React from "react";
import { useQuery } from "react-query";
import { toast } from "sonner";
import { apiKey } from "#/api/queries/users";
import type {
	Workspace,
	WorkspaceAgent,
	WorkspaceApp,
} from "#/api/typesGenerated";
import { useProxy } from "#/contexts/ProxyContext";
import {
	getAppHref,
	isExternalApp,
	needsSessionToken,
	openAppInNewWindow,
} from "./apps";

type UseAppLinkParams = {
	workspace: Workspace;
	agent: WorkspaceAgent;
};

type AppLink = {
	href: string;
	onClick: (e: React.MouseEvent) => void;
	label: string;
	hasToken: boolean;
};

export const useAppLink = (
	app: WorkspaceApp,
	{ agent, workspace }: UseAppLinkParams,
): AppLink => {
	const label = app.display_name ?? app.slug;
	const { proxy } = useProxy();
	const { data: apiKeyResponse } = useQuery({
		...apiKey(),
		enabled: isExternalApp(app) && needsSessionToken(app),
	});

	const href = getAppHref(app, {
		agent,
		workspace,
		token: apiKeyResponse?.key,
		path: proxy.preferredPathAppURL,
		host: proxy.preferredWildcardHostname,
	});

	const onClick = (e: React.MouseEvent) => {
		if (!e.currentTarget.getAttribute("href")) {
			return;
		}

		// External apps with custom protocols (non-HTTP) need special handling
		// for error detection when the app isn't installed.
		const isExternalProtocolApp =
			app.external && app.url && !app.url.startsWith("http");

		if (isExternalProtocolApp) {
			// When browser recognizes the protocol and is able to navigate to the app,
			// it will blur away, and will stop the timer. Otherwise,
			// an error message will be displayed.
			const openAppExternallyFailedTimeout = 1500;
			const openAppExternallyFailed = setTimeout(() => {
				// Check if this is a JetBrains IDE app
				// starts with "jetbrains-gateway://connect#type=optimus-ide-collab" (from https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/jetbrains-gateway)
				const isJetBrainsGateway = app.url?.startsWith("jetbrains-gateway:");
				// starts with "jetbrains://gateway/optimus-ide-collab" (from https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/jetbrains)
				const isJetBrainsToolbox = app.url?.startsWith("jetbrains:");

				// Check if this is a optimus-ide-collab:// URL
				const isOptimusIdeCollabApp = app.url?.startsWith("optimus-ide-collab:");

				if (isJetBrainsGateway) {
					toast.error(`Failed to open "${label}".`, {
						description: "JetBrains Gateway must be installed.",
					});
				} else if (isJetBrainsToolbox) {
					toast.error(`Failed to open "${label}".`, {
						description: "JetBrains Toolbox must be installed.",
					});
				} else if (isOptimusIdeCollabApp) {
					toast.error(`Failed to open "${label}".`, {
						description: "Optimus-IDE-Collab Desktop must be installed.",
					});
				} else {
					toast.error(`Failed to open "${label}".`, {
						description: "The app must be installed first.",
					});
				}
			}, openAppExternallyFailedTimeout);
			window.addEventListener("blur", () => {
				clearTimeout(openAppExternallyFailed);
			});

			// Custom protocol external apps don't support open_in since they
			// rely on the browser's protocol handling.
			return;
		}

		switch (app.open_in) {
			case "slim-window": {
				e.preventDefault();
				openAppInNewWindow(href);
				return;
			}
		}
	};

	return {
		href,
		onClick,
		label,
		hasToken: Boolean(apiKeyResponse?.key),
	};
};
