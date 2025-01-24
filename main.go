package main

import (
	"container-monitor/internal/cmd"
	_ "container-monitor/internal/logic"
	_ "container-monitor/internal/packed"
	"container-monitor/internal/service"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	service.Env().Load()
	service.Task().Start()
	cmd.Main.Run(gctx.GetInitCtx())
}
