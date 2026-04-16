package svc

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"2cgrowth/goz/internal/nacosx"
	"2cgrowth/goz/rpc/user/user"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

func newNacosUserClient(c nacosx.Config) (user.UserClient, *grpc.ClientConn, naming_client.INamingClient, error) {
	serviceName := c.Service("user-rpc")
	groupName := c.Group()
	clusters := c.Clusters()

	namingClient, err := nacosx.NewNamingClient(c)
	if err != nil {
		return nil, nil, nil, err
	}

	discoverBuilder := manual.NewBuilderWithScheme("nacos-user-rpc")
	instances, err := namingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   groupName,
		Clusters:    clusters,
		HealthyOnly: true,
	})
	if err != nil {
		namingClient.CloseClient()
		return nil, nil, nil, fmt.Errorf("select nacos instances: %w", err)
	}

	initialAddrs := instancesToResolverAddrs(instances)
	if len(initialAddrs) == 0 {
		namingClient.CloseClient()
		return nil, nil, nil, fmt.Errorf("no healthy instance found for %s", serviceName)
	}
	discoverBuilder.InitialState(resolver.State{Addresses: initialAddrs})

	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("%s:///%s", discoverBuilder.Scheme(), serviceName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithResolvers(discoverBuilder),
	)
	if err != nil {
		namingClient.CloseClient()
		return nil, nil, nil, fmt.Errorf("dial user-rpc by nacos: %w", err)
	}

	if err := namingClient.Subscribe(&vo.SubscribeParam{
		ServiceName: serviceName,
		GroupName:   groupName,
		Clusters:    clusters,
		SubscribeCallback: func(services []model.Instance, subscribeErr error) {
			if subscribeErr != nil {
				logx.Errorf("nacos subscribe callback error for %s: %v", serviceName, subscribeErr)
				return
			}

			addrs := instancesToResolverAddrs(services)
			if len(addrs) == 0 {
				// Keep previous addresses to avoid forcing all requests to fail
				// when Nacos delivers a temporary empty snapshot.
				logx.Errorf("nacos returned empty healthy instances for %s", serviceName)
				return
			}

			discoverBuilder.UpdateState(resolver.State{Addresses: addrs})
		},
	}); err != nil {
		_ = conn.Close()
		namingClient.CloseClient()
		return nil, nil, nil, fmt.Errorf("subscribe nacos service updates: %w", err)
	}

	return user.NewUserClient(conn), conn, namingClient, nil
}

func instancesToResolverAddrs(instances []model.Instance) []resolver.Address {
	addrs := make([]resolver.Address, 0, len(instances))
	for _, instance := range instances {
		ip := strings.TrimSpace(instance.Ip)
		if ip == "" || instance.Port == 0 {
			continue
		}

		addrs = append(addrs, resolver.Address{
			Addr: net.JoinHostPort(ip, strconv.FormatUint(instance.Port, 10)),
		})
	}

	return addrs
}
