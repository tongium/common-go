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

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToUpper(snake)
}

const prefix = "APP_"

// Support only exported field and type is string, bool, int, or float
func Load(T interface{}) error {
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
		name := prefix + ToSnakeCase(elem.Type().Field(i).Name)

		var required bool
		if value, ok := tag.Lookup("required"); ok && value == "true" {
			required = true
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
				field.SetBool(value == "true")

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
