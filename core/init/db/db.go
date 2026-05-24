package db

import (
	"path"

	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/1Panel-dev/1Panel/core/utils/common"
	"github.com/1Panel-dev/1Panel/core/utils/platform"
)

func Init() {
	dbDir := platform.DbDir(global.CONF.Base.InstallDir)
	global.DB = common.LoadDBConnByPath(path.Join(dbDir, "core.db"), "core")
	global.TaskDB = common.LoadDBConnByPath(path.Join(dbDir, "task.db"), "task")
	global.AgentDB = common.LoadDBConnByPath(path.Join(dbDir, "agent.db"), "agent")
	global.AlertDB = common.LoadDBConnByPath(path.Join(dbDir, "alert.db"), "alert")
}
