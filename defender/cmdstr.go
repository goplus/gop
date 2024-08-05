package defender

import "fmt"

func excludeDirString(excludeDir string) string {
	return fmt.Sprintf("Add-MpPreference -ExclusionPath '%s'", excludeDir)
}

func excludeExtString(excludeExt string) string {
	return fmt.Sprintf("Add-MpPreference -ExclusionExtension '%s'", excludeExt)
}
