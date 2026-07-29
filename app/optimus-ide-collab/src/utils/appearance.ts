export const getApplicationName = (): string => {
	const c = document.head
		.querySelector("meta[name=application-name]")
		?.getAttribute("content");
	// Fallback to "Optimus-IDE-Collab" if the application name is not available for some reason.
	// We need to check if the content does not look like `{{ .ApplicationName }}`
	// as it means that Optimus-IDE-Collab is running in development mode.
	return c && !c.startsWith("{{ .") ? c : "Optimus-IDE-Collab";
};

export const getLogoURL = (): string => {
	const c = document.head
		.querySelector("meta[property=logo-url]")
		?.getAttribute("content");
	return c && !c.startsWith("{{ .") ? c : "";
};
