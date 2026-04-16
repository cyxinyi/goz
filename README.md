# goz

## Shared config

API and RPC can now share one yaml config file:

- `api/user/etc/user-api.yaml`
- `rpc/user/etc/user-all.yaml`

Run examples:

- API: `cd api/user && go run . -f etc/user-api.yaml`
- RPC: `cd rpc/user && go run . -f etc/user-all.yaml`

## Nacos service discovery

- `Rpc.Nacos` is used by `user-rpc` to register itself in Nacos.
- `Api.UserRpcNacos` is used by `user-api` to discover `user-rpc`.
- `Api.UserRpc.Endpoints` remains as a fallback when Nacos is not configured.

Default sample config points to `127.0.0.1:8848`, service name `user-rpc`.
