// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"goz/api/user/internal/config"
	"goz/rpc/user/user"
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
