# common-go

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
![Go: Version](https://img.shields.io/github/go-mod/go-version/tongium/common-go)
![Tag](https://img.shields.io/github/v/tag/tongium/common-go)
![Test](https://github.com/tongium/common-go/actions/workflows/test.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Common code in Golang HTTP Server

## Properties loader

Set struct value from environment

### Usage

example:

```sh
export APP_NAME='Apple'
export APP_NUMBER=1
export APP_DIGIT=0.99
export APP_SOME_WORD='Yes, it is from environment'
```

main.go:

```golang
package main

import (
	"fmt"

	"github.com/tongium/common-go/pkg/properties"
)

type Configuration struct {
	Name     string  ``
	Number   int     `required:"true"`
	Digit    float64 ``
	SomeWord string  ``
}

func main() {
	cfg := Configuration{}
	err := properties.Load(&cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg)
}
```

result:

```sh
{Apple 1 0.99 Yes, it is from environment}
```

Add `required:"true"` to return error if ENV not found


see more: [example](example/properties/main.go)

## Opentracing middleware

Set tags http.status_code, http.request_id and http.user_id

see more: [example](example/jeager/main.go) with [Echo](https://echo.labstack.com/)

run Jeager all in one

```sh
docker run -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest
```

then make a request

```
curl --location --request GET 'http://localhost:1323/' \
--header 'X-User-ID: 4493'
```

then visit http://localhost:16686