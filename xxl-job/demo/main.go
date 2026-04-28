package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	xxl "github.com/xxl-job/xxl-job-executor-go"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "xxl-job/demo/etc/demo.yaml", "xxl-job demo config file")

type Config struct {
	Name     string         `json:",optional"`
	Admin    AdminConfig    `json:",optional"`
	Executor ExecutorConfig `json:",optional"`
}

type AdminConfig struct {
	ServerAddr  string `json:",optional"`
	AccessToken string `json:",optional"`
}

type ExecutorConfig struct {
	RegistryKey string `json:",optional"`
	IP          string `json:",optional"`
	Port        string `json:",optional"`
}

func main() {
	flag.Parse()

	var c Config
	conf.MustLoad(*configFile, &c)
	applyDefaults(&c)

	executor := xxl.NewExecutor(
		xxl.ServerAddr(c.Admin.ServerAddr),
		xxl.AccessToken(c.Admin.AccessToken),
		xxl.RegistryKey(c.Executor.RegistryKey),
		xxl.ExecutorIp(c.Executor.IP),
		xxl.ExecutorPort(c.Executor.Port),
		xxl.SetLogger(&xxlLogAdapter{}),
	)
	executor.Init()
	executor.RegTask("demo.hello", helloTask)
	executor.RegTask("demo.sleep", sleepTask)

	errCh := make(chan error, 1)
	go func() {
		errCh <- executor.Run()
	}()

	logx.Infof(
		"xxl-job executor started name=%s admin=%s registryKey=%s listen=%s:%s handlers=[demo.hello,demo.sleep]",
		c.Name, c.Admin.ServerAddr, c.Executor.RegistryKey, c.Executor.IP, c.Executor.Port,
	)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case sig := <-sigCh:
		logx.Infof("received signal %s, stopping xxl-job executor", sig)
		executor.Stop()
		select {
		case err := <-errCh:
			if err != nil && !errors.Is(err, context.Canceled) {
				logx.Errorf("xxl-job executor stopped with error: %v", err)
			}
		case <-time.After(3 * time.Second):
			logx.Infof("xxl-job executor stop timeout, force exit")
		}
	case err := <-errCh:
		if err != nil {
			panic(fmt.Errorf("xxl-job executor run failed: %w", err))
		}
	}
}

func helloTask(_ context.Context, req *xxl.RunReq) string {
	logx.Infof(
		"[xxl-job] demo.hello executed jobId=%d logId=%d handler=%s params=%s shard=%d/%d",
		req.JobID, req.LogID, req.ExecutorHandler, req.ExecutorParams, req.BroadcastIndex, req.BroadcastTotal,
	)
	return fmt.Sprintf("hello from go executor, params=%s", req.ExecutorParams)
}

func sleepTask(ctx context.Context, req *xxl.RunReq) string {
	waitMs := parseSleepMs(req.ExecutorParams, 3000)
	timer := time.NewTimer(time.Duration(waitMs) * time.Millisecond)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		logx.Errorf("[xxl-job] demo.sleep canceled jobId=%d logId=%d err=%v", req.JobID, req.LogID, ctx.Err())
		return "canceled by timeout or kill"
	case <-timer.C:
		logx.Infof("[xxl-job] demo.sleep finished jobId=%d logId=%d waitMs=%d", req.JobID, req.LogID, waitMs)
		return fmt.Sprintf("slept %d ms", waitMs)
	}
}

func parseSleepMs(raw string, fallback int) int {
	v := strings.TrimSpace(raw)
	if v == "" {
		return fallback
	}
	ms, err := strconv.Atoi(v)
	if err != nil || ms <= 0 {
		return fallback
	}
	return ms
}

func applyDefaults(c *Config) {
	if c.Name == "" {
		c.Name = "xxl-job-demo"
	}
	//bdot注册了环境变量
	c.Admin.ServerAddr = resolveServerAddr(c.Admin.ServerAddr)
	if c.Admin.ServerAddr == "" {
		c.Admin.ServerAddr = "http://127.0.0.1:8087/xxl-job-admin"
	}
	//bdot注册了环境变量
	c.Admin.AccessToken = resolveAccessToken(c.Admin.AccessToken)
	if c.Admin.AccessToken == "" {
		c.Admin.AccessToken = "default_token"
	}
	if c.Executor.RegistryKey == "" {
		c.Executor.RegistryKey = "cursor-xxl-job-demo"
	}
	c.Executor.IP = resolveExecutorIP(c.Executor.IP)
	if c.Executor.Port == "" {
		c.Executor.Port = "9999"
	}
}

func resolveExecutorIP(configIP string) string {
	if podIP := strings.TrimSpace(os.Getenv("POD_IP")); podIP != "" {
		return podIP
	}
	if ip := strings.TrimSpace(configIP); ip != "" {
		return ip
	}
	return "127.0.0.1"
}

func resolveServerAddr(serverAddr string) string {
	if serverAddress := strings.TrimSpace(os.Getenv("XXL_JOB_ADMIN_ADDRESS")); serverAddress != "" {
		return serverAddress
	}
	return serverAddr
}

func resolveAccessToken(accessToken string) string {
	if token := strings.TrimSpace(os.Getenv("XXL_JOB_ACCESS_TOKEN")); token != "" {
		return token
	}
	return accessToken
}

type xxlLogAdapter struct{}

func (l *xxlLogAdapter) Info(format string, a ...interface{}) {
	logx.Infof("[xxl-lib] "+format, a...)
}

func (l *xxlLogAdapter) Error(format string, a ...interface{}) {
	logx.Errorf("[xxl-lib] "+format, a...)
}
