package main

import (
	"fmt"

	_ "github.com/joho/godotenv/autoload"
	"github.com/tongium/common-go/pkg/properties"
)

type Configuration struct {
	Name    string  ``
	Enabled bool    ``
	Number  int     `required:"true"`
	Digit   float64 ``
}

func main() {
	cfg := Configuration{}
	err := properties.Load(&cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg)
}
