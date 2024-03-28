package defender

import "fmt"

func excludeCmdString(excludeDir string) string {
	//excludeCmd := fmt.Sprintf("Add-MpPreference -ExclusionPath '%s' -ControlledFolderAccessAllowedApplications '%s'", v, filepath.Join(goproot, "/bin/gop.exe"))
	return fmt.Sprintf("Add-MpPreference -ExclusionPath '%s'", excludeDir)
}
