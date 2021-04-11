# common-go

![Test](https://github.com/tongium/common-go/actions/workflows/test.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Common code in Golang HTTP Server

## Jaeger middleware

Set tags http.status_code, http.request_id and http.user_id

[example](example/jeager/main.go)

## Properties loader

Set struct value from environment

### Usage

Add require tag to return error if ENV not found

```golang
RequiredValue string `required:"true"`
```

[example](example/properties/main.go)