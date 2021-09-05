# gjrpc
Go JSON RPC with code generation from go declarations. This project is in very early alpha and works only for simplest cases.

You can look at [123](./examples/1_simple_service/README.md)

## Code annotation glossary

- `gjrpc:service <rpc_path>`
  * Place: `<> type Service interface { [methods] }`
  * Used on service declaration, to register service in rpc protocol
- `gjrpc:method <rpc_path>`
  * Place: `type Service interface { <> func method1() ... }`
  * Used on method declaration, to register method in the service
- `gjrpc:handle-route <service_go_type>`
  * Place: `type Handlers struct { <> ServiceName Type ... }`
  * Used on field with service implementation, to generate appropriate router for this handler