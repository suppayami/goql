# goql

goql generate GraphQL Schema from a relational database (only MySQL supported at the moment) and serve as GraphQL API, supports CRUD.

## Installation

Install dependencies with [dep](https://github.com/golang/dep)

```
dep ensure
```

## Usage
`go run main.go -e > schema.graphql` - Export GraphQL schema to file

`go run main.go -s` - Serve resolver
