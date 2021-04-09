package properties_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/tongium/common-go/pkg/properties"
)

type Configuration struct {
	Float         float64
	String        string
	Integer       int
	Boolean       bool
	CamalCase     string
	RequiredValue string `required:"true"`
}

func TestFloatShouldBeCorrect(t *testing.T) {
	expected := 0.01
	os.Setenv("APP_FLOAT", fmt.Sprint(expected))
	os.Setenv("APP_REQUIRED_VALUE", "test")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err != nil {
		t.Error("got error", err)
	}

	result := cfg.Float
	if result != expected {
		t.Errorf("expected %v but got %v", expected, result)
	}

	os.Unsetenv("APP_FLOAT")
	os.Unsetenv("APP_REQUIRED_VALUE")
}

func TestStringShouldBeCorrect(t *testing.T) {
	expected := "hello"
	os.Setenv("APP_STRING", expected)
	os.Setenv("APP_REQUIRED_VALUE", "test")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err != nil {
		t.Error("got error", err)
	}

	result := cfg.String
	if result != expected {
		t.Errorf("expected %v but got %v", expected, result)
	}

	os.Unsetenv("APP_STRING")
	os.Unsetenv("APP_REQUIRED_VALUE")
}

func TestBoolShouldBeCorrect(t *testing.T) {
	expected := "true"
	os.Setenv("APP_BOOLEAN", fmt.Sprint(expected))
	os.Setenv("APP_REQUIRED_VALUE", "test")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err != nil {
		t.Error("got error", err)
	}

	result := cfg.Boolean
	if !result {
		t.Errorf("expected %v but got %v", expected, result)
	}

	os.Unsetenv("APP_BOOLEAN")
	os.Unsetenv("APP_REQUIRED_VALUE")
}

func TestIntShouldBeCorrect(t *testing.T) {
	expected := 1
	os.Setenv("APP_INTEGER", fmt.Sprint(expected))
	os.Setenv("APP_REQUIRED_VALUE", "test")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err != nil {
		t.Error("got error", err)
	}

	result := cfg.Integer
	if result != expected {
		t.Errorf("expected %v but got %v", expected, result)
	}

	os.Unsetenv("APP_INTEGER")
	os.Unsetenv("APP_REQUIRED_VALUE")
}

func TestIntShouldBeError(t *testing.T) {
	os.Setenv("APP_INTEGER", "word")
	os.Setenv("APP_REQUIRED_VALUE", "test")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err == nil {
		t.Error("expected error but got", err)
	}

	os.Unsetenv("APP_INTEGER")
	os.Unsetenv("APP_REQUIRED_VALUE")
}

func TestFloatShouldBeError(t *testing.T) {
	os.Setenv("APP_FLOAT", "word")
	os.Setenv("APP_REQUIRED_VALUE", "test")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err == nil {
		t.Error("expected error but got", err)
	}

	os.Unsetenv("APP_FLOAT")
	os.Unsetenv("APP_REQUIRED_VALUE")
}

func TestRequiredValueShouldBeError(t *testing.T) {
	os.Unsetenv("APP_REQUIRED_VALUE")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err == nil {
		t.Error("expected error but got", err)
	}
}

func TestCamalCaseShouldBeCorrect(t *testing.T) {
	expected := "test"
	os.Setenv("APP_CAMAL_CASE", fmt.Sprint(expected))
	os.Setenv("APP_REQUIRED_VALUE", "test")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err != nil {
		t.Error("got error", err)
	}

	result := cfg.CamalCase
	if result != expected {
		t.Errorf("expected %v but got %v", expected, result)
	}

	os.Unsetenv("APP_CAMAL_CASE")
	os.Unsetenv("APP_REQUIRED_VALUE")
}
