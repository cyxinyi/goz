package svc

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"2cgrowth/goz/internal/nacosx"
	"2cgrowth/goz/rpc/user/internal/config"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/logx"
)

func RegisterToNacos(c config.Config) (func(), error) {
	if !c.Nacos.Enabled() {
		return func() {}, nil
	}

	namingClient, err := nacosx.NewNamingClient(c.Nacos)
	if err != nil {
		return nil, fmt.Errorf("new nacos naming client: %w", err)
	}

	listenHost, listenPort, err := parseListenOn(c.ListenOn)
	if err != nil {
		namingClient.CloseClient()
		return nil, err
	}

	registerIP, err := resolveRegisterIP(c.Nacos.RegisterIP, listenHost)
	if err != nil {
		namingClient.CloseClient()
		return nil, err
	}

	serviceName := c.Nacos.Service(c.Name)
	groupName := c.Nacos.Group()
	clusterName := c.Nacos.ClusterName
	ephemeral := c.Nacos.EphemeralEnabled()

	ok, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          registerIP,
		Port:        listenPort,
		Weight:      c.Nacos.RegisterWeight(),
		Enable:      true,
		Healthy:     true,
		Metadata:    c.Nacos.Metadata,
		ClusterName: clusterName,
		ServiceName: serviceName,
		GroupName:   groupName,
		Ephemeral:   ephemeral,
	})
	if err != nil {
		namingClient.CloseClient()
		return nil, fmt.Errorf("register rpc instance to nacos: %w", err)
	}
	if !ok {
		namingClient.CloseClient()
		return nil, fmt.Errorf("register rpc instance to nacos returned false")
	}

	logx.Infof("registered rpc service to nacos: service=%s group=%s ip=%s port=%d", serviceName, groupName, registerIP, listenPort)

	return func() {
		_, deregisterErr := namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          registerIP,
			Port:        listenPort,
			Cluster:     clusterName,
			ServiceName: serviceName,
			GroupName:   groupName,
			Ephemeral:   ephemeral,
		})
		if deregisterErr != nil {
			logx.Errorf("deregister rpc instance from nacos failed: %v", deregisterErr)
		}
		namingClient.CloseClient()
	}, nil
}

func parseListenOn(listenOn string) (string, uint64, error) {
	host, port, err := net.SplitHostPort(strings.TrimSpace(listenOn))
	if err != nil {
		return "", 0, fmt.Errorf("invalid ListenOn %q: %w", listenOn, err)
	}

	parsedPort, err := strconv.ParseUint(port, 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("invalid listen port %q: %w", port, err)
	}

	return host, parsedPort, nil
}

func resolveRegisterIP(registerIP string, listenHost string) (string, error) {
	registerIP = strings.TrimSpace(registerIP)
	if registerIP != "" {
		return registerIP, nil
	}

	listenHost = strings.TrimSpace(listenHost)
	if listenHost != "" && listenHost != "0.0.0.0" && listenHost != "::" {
		return listenHost, nil
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("list network interfaces: %w", err)
	}

	for _, iface := range interfaces {
		if (iface.Flags & net.FlagUp) == 0 {
			continue
		}
		if (iface.Flags & net.FlagLoopback) != 0 {
			continue
		}

		addrs, addrErr := iface.Addrs()
		if addrErr != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if ipNet.IP == nil || ipNet.IP.IsLoopback() {
				continue
			}
			if ip4 := ipNet.IP.To4(); ip4 != nil {
				return ip4.String(), nil
			}
		}
	}

	return "", fmt.Errorf("cannot determine register ip, please set Rpc.Nacos.RegisterIP")
}
