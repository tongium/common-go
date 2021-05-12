package properties

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
var nameMatch = regexp.MustCompile("name:([A-Za-z0-9_]+)")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToUpper(snake)
}

type Config struct {
	Prefix string
}

// Support only exported field and type is string, bool, int, or float
func Load(T interface{}) error {
	cfg := Config{
		Prefix: "APP_",
	}

	return load(cfg, T)
}

func LoadWithConfig(T interface{}, cfg Config) error {
	return load(cfg, T)
}

func load(cfg Config, T interface{}) error {
	v := reflect.ValueOf(T)
	if v.Kind() != reflect.Ptr {
		return errors.New("given argument is not pointer")
	}

	elem := v.Elem()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if !field.CanSet() {
			continue
		}

		tag := elem.Type().Field(i).Tag
		name := cfg.Prefix + ToSnakeCase(elem.Type().Field(i).Name)

		var required bool
		if value, ok := tag.Lookup("prop"); ok {
			if strings.Contains(value, "require") {
				required = true
			}

			if n := nameMatch.FindString(value); n != "" {
				ss := strings.Split(n, ":")
				name = ss[len(ss)-1]
			}
		}

		value, ok := os.LookupEnv(name)
		if required && !ok {
			return fmt.Errorf("env '%s' not set", name)
		}

		if ok {
			switch field.Type().Name() {
			case "string":
				field.SetString(value)

			case "int", "int8", "int16", "int32", "int64":
				if i, err := strconv.Atoi(value); err == nil {
					field.SetInt(int64(i))
				} else {
					return fmt.Errorf("env '%s' expect integer but got '%s'", name, value)
				}

			case "bool":
				field.SetBool(value == "true" || value == "enabled" || value == "enable")

			case "float32", "float64":
				if f, err := strconv.ParseFloat(value, 64); err == nil {
					field.SetFloat(f)
				} else {
					return fmt.Errorf("env '%s' expect float but got '%s'", name, value)
				}

			default:
				// do nothing
			}
		}
	}

	return nil
}
