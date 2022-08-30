package pkg

import (
	"github.com/colinc9/go-distributed-cache/pkg/config"
	"github.com/colinc9/go-distributed-cache/pkg/controller"
)

func SetUpAndRun() {
	config.GetDefaultInsCfg()
	controller.Run()
}