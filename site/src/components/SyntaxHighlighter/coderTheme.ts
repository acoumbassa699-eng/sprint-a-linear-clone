import { useTheme } from "@emotion/react";
import { useMonaco } from "@monaco-editor/react";
import { useEffect, useState } from "react";
export const useOptimus-IDE-CollabTheme = (): { isLoading: boolean; name: string } => {
	const [isLoading, setIsLoading] = useState(true);
	const monaco = useMonaco();
	const theme = useTheme();
	const name = "optimus-ide-collab";

	useEffect(() => {
		if (monaco) {
			monaco.editor.defineTheme(name, theme.monaco);
			setIsLoading(false);
		}
	}, [monaco, theme]);

	return {
		isLoading,
		name,
	};
};
