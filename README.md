# goz

## Shared config

API and RPC can now share one yaml config file:

- `etc/user-all.yaml`

Run examples:

- API: `cd api/user && go run . -f ../../etc/user-all.yaml`
- RPC: `cd rpc/user && go run . -f ../../etc/user-all.yaml`
