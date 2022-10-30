package pkg

import (
	"github.com/colinc9/go-distributed-cache/pkg/config"
	"github.com/colinc9/go-distributed-cache/pkg/controller"
	"github.com/colinc9/go-distributed-cache/pkg/network"
)

func SetUpAndRun() {
	config.GetDefaultInsCfg()
	network.PeriodicSdkDiscovery()
	controller.Run()
}