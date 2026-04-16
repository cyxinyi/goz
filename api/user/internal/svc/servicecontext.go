// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"2cgrowth/goz/api/user/internal/config"
	"2cgrowth/goz/rpc/user/user"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type ServiceContext struct {
	Config        config.Config
	UserRpc       user.UserClient
	userRpcConn   *grpc.ClientConn
	userRpcNaming naming_client.INamingClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	if c.UserRpcNacos.Enabled() {
		userRpc, conn, naming, err := newNacosUserClient(c.UserRpcNacos)
		if err != nil {
			panic(fmt.Errorf("init nacos user-rpc client: %w", err))
		}

		return &ServiceContext{
			Config:        c,
			UserRpc:       userRpc,
			userRpcConn:   conn,
			userRpcNaming: naming,
		}
	}

	client := zrpc.MustNewClient(c.UserRpc)

	return &ServiceContext{
		Config:  c,
		UserRpc: user.NewUserClient(client.Conn()),
	}
}

func (s *ServiceContext) Close() {
	if s.userRpcNaming != nil {
		s.userRpcNaming.CloseClient()
	}
	if s.userRpcConn != nil {
		_ = s.userRpcConn.Close()
	}
}
