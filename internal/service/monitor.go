// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
)

type (
	IContainerMonitor interface {
		Start()
		Monitor(ctx context.Context, t *gtime.Time) error
	}
)

var (
	localContainerMonitor IContainerMonitor
)

func ContainerMonitor() IContainerMonitor {
	if localContainerMonitor == nil {
		panic("implement not found for interface IContainerMonitor, forgot register?")
	}
	return localContainerMonitor
}

func RegisterContainerMonitor(i IContainerMonitor) {
	localContainerMonitor = i
}
