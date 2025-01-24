// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

type (
	IEnv interface {
		Load() error
	}
)

var (
	localEnv IEnv
)

func Env() IEnv {
	if localEnv == nil {
		panic("implement not found for interface IEnv, forgot register?")
	}
	return localEnv
}

func RegisterEnv(i IEnv) {
	localEnv = i
}
