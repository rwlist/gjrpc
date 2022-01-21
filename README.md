# gjrpc

Go JSON RPC with code generation from go declarations. This project is in very early alpha and works only for simplest
cases.

You can look at [this example](./examples/1_simple_service/README.md), to see what is supported now.

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

## Goal

The goal of this project is to create tooling which will allow to create and implement API which is easy to use from
the browser (TypeScript) and server code in Go and other languages. To do that, JSON RPC is a nice choice, which
is easy to use and easy to implement.

This tooling will allow to create API from go declarations, and also generate code for servers and clients in other 
languages.