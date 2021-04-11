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
