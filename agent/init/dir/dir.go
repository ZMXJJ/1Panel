package dir

import (
	"path"

	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/1Panel-dev/1Panel/agent/utils/files"
	"github.com/1Panel-dev/1Panel/agent/utils/platform"
)

func Init() {
	fileOp := files.NewFileOp()
	baseDir := global.CONF.Base.InstallDir
	dataDir := platform.DataDir(baseDir)
	logDir := platform.LogDir(baseDir)
	appsDir := platform.AppsDir(baseDir)
	resourceDir := platform.ResourceDir(baseDir)
	appResourceDir := platform.AppResourceDir(baseDir)

	_, _ = fileOp.CreateDirWithPath(true, path.Join(dataDir, "docker/compose/"))

	global.Dir.BaseDir, _ = fileOp.CreateDirWithPath(true, baseDir)
	global.Dir.DataDir, _ = fileOp.CreateDirWithPath(true, dataDir)
	global.Dir.DbDir, _ = fileOp.CreateDirWithPath(true, platform.DbDir(baseDir))
	global.Dir.LogDir, _ = fileOp.CreateDirWithPath(true, logDir)
	global.Dir.TaskDir, _ = fileOp.CreateDirWithPath(true, platform.TaskLogDir(baseDir))
	global.Dir.TmpDir, _ = fileOp.CreateDirWithPath(true, platform.TmpDir(baseDir))

	global.Dir.AppDir, _ = fileOp.CreateDirWithPath(true, appsDir)
	global.Dir.ResourceDir, _ = fileOp.CreateDirWithPath(true, resourceDir)
	global.Dir.IconCacheDir, _ = fileOp.CreateDirWithPath(true, path.Join(resourceDir, "icon"))
	global.Dir.AppResourceDir, _ = fileOp.CreateDirWithPath(true, appResourceDir)
	global.Dir.AppInstallDir, _ = fileOp.CreateDirWithPath(true, appsDir)
	global.Dir.LocalAppResourceDir, _ = fileOp.CreateDirWithPath(true, path.Join(appResourceDir, "local"))
	global.Dir.LocalAppInstallDir, _ = fileOp.CreateDirWithPath(true, path.Join(appsDir, "local"))
	global.Dir.RemoteAppResourceDir, _ = fileOp.CreateDirWithPath(true, path.Join(appResourceDir, "remote"))
	global.Dir.CustomAppResourceDir, _ = fileOp.CreateDirWithPath(true, path.Join(appResourceDir, "custom"))
	global.Dir.OfflineAppResourceDir, _ = fileOp.CreateDirWithPath(true, path.Join(resourceDir, "offline"))
	if platform.IsDarwin() {
		_, _ = fileOp.CreateDirWithPath(true, path.Join(appResourceDir, "mac"))
	}
	global.Dir.RuntimeDir, _ = fileOp.CreateDirWithPath(true, platform.RunDir(baseDir))
	global.Dir.RecycleBinDir, _ = fileOp.CreateDirWithPath(true, platform.RecycleBinDir(baseDir))
	global.Dir.SSLLogDir, _ = fileOp.CreateDirWithPath(true, path.Join(logDir, "ssl"))
	global.Dir.McpDir, _ = fileOp.CreateDirWithPath(true, path.Join(dataDir, "mcp"))
	global.Dir.ConvertLogDir, _ = fileOp.CreateDirWithPath(true, path.Join(logDir, "convert"))
	global.Dir.TensorRTLLMDir, _ = fileOp.CreateDirWithPath(true, path.Join(dataDir, "ai/tensorrt_llm"))
	global.Dir.FirewallDir, _ = fileOp.CreateDirWithPath(true, path.Join(dataDir, "firewall"))
}
