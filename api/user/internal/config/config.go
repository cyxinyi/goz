// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"2cgrowth/goz/internal/nacosx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	UserRpc      zrpc.RpcClientConf `json:",optional"`
	UserRpcNacos nacosx.Config      `json:",optional"`
}

type UnifiedConfig struct {
	Api Config
}
