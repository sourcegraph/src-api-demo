# Sourcegraph API Automation Demo

> This is not an offical Sourcegraph project and is not supported by Sourcegraph.

This repository contains a demo of how to build automation with the Sourcegraph API.

## Usage

Before getting started:

- You must be a Site Administrator to run this demo.
- Know the GraphQL endpoint of your Sourcegraph instance, e.g., `https://src.acme.com/.api/graphql`, `https://acme.sourcegraphcloud.com/.api/graphql`
- Create an access token, https://docs.sourcegraph.com/cli/how-tos/creating_an_access_token

Then, configure environment variables:

```sh
export SRC_ACCESS_TOKEN="sgp_<REDACTED>"
export SRC_ENDPOINT="https://src.acme.com/.api/graphql"
```

### List code host connections

```sh
go run ./cmd/src-api-demo/ extsvc list
```

### Update a code host connection token

> Only GitHub is supported. You can implement your own.

```sh
go run ./cmd/src-api-demo/ extsvc update-token -id 'id' 'ghp_<REDACTED>'
```

## Development

A lot of code in this repository involves configuring Sourcegraph through its [GraphQL API](https://docs.sourcegraph.com/api/graphql).

To do so we use [`genqlient`](https://github.com/Khan/genqlient), which uses GraphQL queries and [the upstream Sourcegraph GraphQL API schema](#upgrading-the-sourcegraph-graphql-schema) to generate type-safe Go code used to interact with the instance GraphQL endpoint. The generated package is located in `srcgql`, and all operations are defined in [`srcgql/operations.graphql`](srcgql/operations.graphql).

### Interacting with the Sourcegraph GraphQL API

All GraphQL queries should be first added to [`srcgql/operations.graphql`](./srcgql/operations.graphql). To accept arguments, you must declare GraphQL arguments in your operation, for example:

```gql
# Docstrings about operation also gets added to the generated functions
mutation SendTestEmail($to: String!) {
  # ... fields
}
```

### Upgrading the Sourcegraph GraphQL schema

To upgrade the GraphQL schema used for code generation to a newer release, bump the version in 

- [`srcgql/gen.go`](./srcgql/gen.go)
- [`srcconf/gen.go](./srcconf/gen.go)

and run the following command:

```sh
go generate ./...
```
