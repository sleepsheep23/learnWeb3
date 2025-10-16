package main

import (
	"flag"
	_ "net/http/pprof"

	"github.com/ProjectsTask/EasySwapBackend/src/api/router"
	"github.com/ProjectsTask/EasySwapBackend/src/app"
	"github.com/ProjectsTask/EasySwapBackend/src/config"
	"github.com/ProjectsTask/EasySwapBackend/src/service/svc"
	"github.com/ProjectsTask/EasySwapBackend/src/service/v1"
)

const (
	// port       = ":9000"
	repoRoot          = ""
	defaultConfigPath = "./config/config.toml"
)

func main() {
	conf := flag.String("conf", defaultConfigPath, "conf file path")
	flag.Parse()
	c, err := config.UnmarshalConfig(*conf)
	if err != nil {
		panic(err)
	}

	for _, chain := range c.ChainSupported {
		if chain.ChainID == 0 || chain.Name == "" {
			panic("invalid chain_suffix config")
		}
	}

	serverCtx, err := svc.NewServiceContext(c)
	if err != nil {
		panic(err)
	}

	// Start alert check timer
	// go service.StartPriceMonitor(serverCtx)

	go service.StartBlockMonitor(serverCtx)

	// Initialize router
	r := router.NewRouter(serverCtx)
	app, err := app.NewPlatform(c, r, serverCtx)
	if err != nil {
		panic(err)
	}
	app.Start()
}
