package platform

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	oasisInstallEnv = "OASIS_INSTALL_DIR"

	linuxInstallDir = "/opt"
	linuxPanelDir   = "1panel"
	linuxCtlFile    = "/usr/local/bin/1pctl"
	linuxSocketPath = "/etc/1panel/agent.sock"
)

func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

func DefaultInstallDir() string {
	if !IsDarwin() {
		return linuxInstallDir
	}
	if dir := strings.TrimSpace(os.Getenv(oasisInstallEnv)); dir != "" {
		return filepath.Clean(dir)
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, "Library", "Application Support", "Oasis")
	}
	return filepath.Join(os.TempDir(), "Oasis")
}

func ConfigDir() string {
	if IsDarwin() {
		return filepath.Join(DefaultInstallDir(), "conf")
	}
	return filepath.Join(linuxInstallDir, linuxPanelDir, "conf")
}

func AppConfigPath() string {
	return filepath.Join(ConfigDir(), "app.yaml")
}

func CtlFile() string {
	if IsDarwin() {
		return filepath.Join(DefaultInstallDir(), "bin", "oasisctl")
	}
	return linuxCtlFile
}

func DataDir(installDir string) string {
	installDir = normalizeInstallDir(installDir)
	if IsDarwin() {
		return installDir
	}
	return filepath.Join(installDir, linuxPanelDir)
}

func DbDir(installDir string) string {
	return filepath.Join(DataDir(installDir), "db")
}

func LogDir(installDir string) string {
	return filepath.Join(DataDir(installDir), "log")
}

func TaskLogDir(installDir string) string {
	return filepath.Join(LogDir(installDir), "task")
}

func TmpDir(installDir string) string {
	return filepath.Join(DataDir(installDir), "tmp")
}

func AppsDir(installDir string) string {
	return filepath.Join(DataDir(installDir), "apps")
}

func ResourceDir(installDir string) string {
	return filepath.Join(DataDir(installDir), "resource")
}

func AppResourceDir(installDir string) string {
	return filepath.Join(ResourceDir(installDir), "apps")
}

func RunDir(installDir string) string {
	if IsDarwin() {
		return filepath.Join(normalizeInstallDir(installDir), "run")
	}
	return filepath.Join(DataDir(installDir), "runtime")
}

func RecycleBinDir(installDir string) string {
	if IsDarwin() {
		return filepath.Join(TmpDir(installDir), "recycle-bin")
	}
	return "/.1panel_clash"
}

func AgentSocketPath(installDir string) string {
	if IsDarwin() {
		return filepath.Join(RunDir(installDir), "agent.sock")
	}
	return linuxSocketPath
}

func normalizeInstallDir(installDir string) string {
	if strings.TrimSpace(installDir) == "" {
		return DefaultInstallDir()
	}
	return filepath.Clean(installDir)
}
