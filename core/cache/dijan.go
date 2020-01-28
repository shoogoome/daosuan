package cache

import (
	"daosuan/utils"
	"fmt"
	"github.com/shoogoome/godijan"
)

var Dijan godijan.GoDijan

func InitDijan() {
	Dijan = godijan.NewGoDijanConnection(
		fmt.Sprintf("%s:%s", utils.GlobalConfig.Dijan.Host, utils.GlobalConfig.Dijan.Port),
		nil)
}
