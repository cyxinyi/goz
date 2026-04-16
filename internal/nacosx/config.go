package nacosx

import (
	"fmt"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

const defaultGroupName = "DEFAULT_GROUP"

type Config struct {
	Host        string            `json:",optional"`
	Port        uint64            `json:",optional"`
	NamespaceId string            `json:",optional"`
	GroupName   string            `json:",optional"`
	ClusterName string            `json:",optional"`
	ServiceName string            `json:",optional"`
	Username    string            `json:",optional"`
	Password    string            `json:",optional"`
	RegisterIP  string            `json:",optional"`
	Weight      float64           `json:",optional"`
	Ephemeral   *bool             `json:",optional"`
	Metadata    map[string]string `json:",optional"`
}

func (c Config) Enabled() bool {
	return strings.TrimSpace(c.Host) != "" && c.Port > 0
}

func (c Config) Group() string {
	if strings.TrimSpace(c.GroupName) == "" {
		return defaultGroupName
	}

	return c.GroupName
}

func (c Config) Service(defaultService string) string {
	if strings.TrimSpace(c.ServiceName) == "" {
		return defaultService
	}

	return c.ServiceName
}

func (c Config) Clusters() []string {
	if strings.TrimSpace(c.ClusterName) == "" {
		return nil
	}

	return []string{c.ClusterName}
}

func (c Config) RegisterWeight() float64 {
	if c.Weight <= 0 {
		return 1
	}

	return c.Weight
}

func (c Config) EphemeralEnabled() bool {
	if c.Ephemeral == nil {
		return true
	}

	return *c.Ephemeral
}

func NewNamingClient(c Config) (naming_client.INamingClient, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("nacos is not enabled: host=%q port=%d", c.Host, c.Port)
	}

	return clients.NewNamingClient(
		// Minimal setup for service discovery and registration.
		// Optional paths/logging/auth can be extended here if needed.
		vo.NacosClientParam{
			ClientConfig: &constant.ClientConfig{
				NamespaceId: c.NamespaceId,
				Username:    c.Username,
				Password:    c.Password,
				LogLevel:    "warn",
			},
			ServerConfigs: []constant.ServerConfig{
				{
					IpAddr: c.Host,
					Port:   c.Port,
				},
			},
		},
	)
}
