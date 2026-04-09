// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"2cgrowth/goz/api/user/internal/config"
	"2cgrowth/goz/rpc/user/user"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user.UserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	client := zrpc.MustNewClient(c.UserRpc)

	return &ServiceContext{
		Config:  c,
		UserRpc: user.NewUserClient(client.Conn()),
	}
}
