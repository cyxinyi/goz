// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"

	"2cgrowth/goz/api/user/internal/config"
	"2cgrowth/goz/api/user/internal/handler"
	"2cgrowth/goz/api/user/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	var uc config.UnifiedConfig
	if err := conf.Load(*configFile, &uc); err == nil && uc.Api.Name != "" {
		c = uc.Api
	} else {
		conf.MustLoad(*configFile, &c)
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
