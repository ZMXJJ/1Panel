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

func RunDir(installDir string) string {
	if IsDarwin() {
		return filepath.Join(normalizeInstallDir(installDir), "run")
	}
	return filepath.Join(DataDir(installDir), "runtime")
}

func TmpDir(installDir string) string {
	return filepath.Join(DataDir(installDir), "tmp")
}

func SecretDir(installDir string) string {
	return filepath.Join(DataDir(installDir), "secret")
}

func GeoDir(installDir string) string {
	return filepath.Join(DataDir(installDir), "geo")
}

func AgentSocketPath(installDir string) string {
	if IsDarwin() {
		return filepath.Join(RunDir(installDir), "agent.sock")
	}
	return linuxSocketPath
}

func EnsureCoreDirs(installDir string) error {
	dirs := []string{
		ConfigDir(),
		DbDir(installDir),
		LogDir(installDir),
		filepath.Join(LogDir(installDir), "task"),
		RunDir(installDir),
		TmpDir(installDir),
		SecretDir(installDir),
		GeoDir(installDir),
	}
	if IsDarwin() {
		dirs = append(dirs,
			filepath.Join(normalizeInstallDir(installDir), "apps"),
			filepath.Join(normalizeInstallDir(installDir), "resource", "apps", "mac"),
			filepath.Join(normalizeInstallDir(installDir), "backup"),
			filepath.Join(normalizeInstallDir(installDir), "bin"),
		)
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return err
		}
	}
	return nil
}

func normalizeInstallDir(installDir string) string {
	if strings.TrimSpace(installDir) == "" {
		return DefaultInstallDir()
	}
	return filepath.Clean(installDir)
}
