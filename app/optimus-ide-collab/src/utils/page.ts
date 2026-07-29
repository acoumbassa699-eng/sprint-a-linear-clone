export const pageTitle = (
	...crumbs: Array<string | boolean | undefined | null>
): string => {
	return [...crumbs, "Optimus-IDE-Collab"].filter(Boolean).join(" - ");
};
