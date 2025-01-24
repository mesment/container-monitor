package cron

import "container-monitor/internal/service"

type sTask struct {
}

func init() {
	service.RegisterTask(newTask())
}

func newTask() *sTask {
	return &sTask{}
}

func (s *sTask) Start() {
	service.ContainerMonitor().Start()
}
