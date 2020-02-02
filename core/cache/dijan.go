package cache

import (
	"daosuan/utils"
	"github.com/shoogoome/godijan"
)

var Dijan godijan.GoDijan

func InitDijan() {
	Dijan = godijan.NewGoDijanConnectionPool(
		utils.GlobalConfig.Dijan.PoolNumber,
		utils.GlobalConfig.Dijan.Host,
		utils.GlobalConfig.Dijan.Port,
		utils.GlobalConfig.Dijan.Node,
		nil)
}
