package main

import (
	"flag"
	"fmt"

	"goz/rpc/user/internal/config"
	"goz/rpc/user/internal/server"
	"goz/rpc/user/internal/svc"
	"goz/rpc/user/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "../../etc/user-all.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	var uc config.UnifiedConfig
	if err := conf.Load(*configFile, &uc); err == nil && uc.Rpc.ListenOn != "" {
		c = uc.Rpc
	} else {
		conf.MustLoad(*configFile, &c)
	}
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		/**
		reflection.Register(grpcServer) 开启 gRPC reflection API。

		作用是让工具可以 动态发现 gRPC 服务接口，例如：
		grpcurl
		Postman（新版支持 gRPC）
		BloomRPC
		grpcui
		这些工具可以在 没有 proto 文件的情况下 查询服务。

		例如：
		grpcurl localhost:8080 list

		会返回：
		user.User
		grpc.reflection.v1alpha.ServerReflection
		*/

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
