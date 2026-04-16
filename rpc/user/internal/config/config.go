package config

import (
	"2cgrowth/goz/internal/nacosx"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Nacos nacosx.Config `json:",optional"`
}

type UnifiedConfig struct {
	Rpc Config
}
