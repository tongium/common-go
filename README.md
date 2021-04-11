# common-go

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Test](https://github.com/tongium/common-go/actions/workflows/test.yml/badge.svg)

Common code in Golang HTTP Server

## Opentracing middleware

Set tags http.status_code, http.request_id and http.user_id

[example](example/jeager/main.go) with [Echo](https://echo.labstack.com/)

run Jeager all in one

```
docker run -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest
```

then visit http://localhost:16686/

## Properties loader

Set struct value from environment

### Usage

Example:

```
package main

import (
	"fmt"

	"github.com/tongium/common-go/pkg/properties"
)

type Configuration struct {
	Name        string
	Number      int     `required:"true"`
	Digit       float64
    SomeWord    string
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

value of cfg.Name = os.Getenv("APP_NAME")
value of cfg.SomeWord = os.Getenv("APP_SOME_WORD")
Add require tag to return error if ENV not found

```golang
RequiredValue string `required:"true"`
```

[example](example/properties/main.go)