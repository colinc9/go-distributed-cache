package main

import (
	cache "github.com/colinc9/go-distributed-cache/pkg"
	"github.com/colinc9/go-distributed-cache/pkg/network"
)

func main() {
	cache.SetUpAndRun()
	network.PeriodicSdkDiscovery()
}