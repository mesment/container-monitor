package env

import (
	"container-monitor/internal/service"
	"github.com/joho/godotenv"
)

type sEnv struct{}

func newsEnv() *sEnv {
	return &sEnv{}
}

func init() {
	service.RegisterEnv(newsEnv())
}

func (s *sEnv) Load() error {
	return godotenv.Load()
}
