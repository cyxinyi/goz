// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"awesomeProject3/go-zero/api/user/internal/config"
	"awesomeProject3/go-zero/rpc/user/user"
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
