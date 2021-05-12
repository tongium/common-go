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
	RequiredValue string `prop:"require,name:CUSTOM_NAME"`
	BaseURL       string `prop:"require"`
}

func TestFloatShouldBeCorrect(t *testing.T) {
	os.Setenv("APP_BASE_URL", "world")
	defer os.Unsetenv("APP_BASE_URL")

	expected := 0.01
	os.Setenv("APP_FLOAT", fmt.Sprint(expected))
	os.Setenv("CUSTOM_NAME", "test")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err != nil {
		t.Error("got error", err)
	}

	result := cfg.Float
	if result != expected {
		t.Errorf("expected %v but got %v", expected, result)
	}

	os.Unsetenv("APP_FLOAT")
	os.Unsetenv("CUSTOM_NAME")
}

func TestStringShouldBeCorrect(t *testing.T) {
	os.Setenv("APP_BASE_URL", "world")
	defer os.Unsetenv("APP_BASE_URL")

	os.Setenv("APP_STRING", "hello")
	os.Setenv("CUSTOM_NAME", "my")
	os.Setenv("APP_BASE_URL", "world")
	defer os.Unsetenv("APP_BASE_URL")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err != nil {
		t.Error("got error", err)
	}

	if cfg.String != "hello" {
		t.Errorf("expected %v but got %v", "hello", cfg.String)
	}

	if cfg.BaseURL != "world" {
		t.Errorf("expected %v but got %v", "world", cfg.BaseURL)
	}

	defer os.Unsetenv("APP_STRING")
	defer os.Unsetenv("CUSTOM_NAME")

}

func TestBoolShouldBeCorrect(t *testing.T) {
	os.Setenv("APP_BASE_URL", "world")
	defer os.Unsetenv("APP_BASE_URL")

	expected := "true"
	os.Setenv("APP_BOOLEAN", fmt.Sprint(expected))
	os.Setenv("CUSTOM_NAME", "test")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err != nil {
		t.Error("got error", err)
	}

	result := cfg.Boolean
	if !result {
		t.Errorf("expected %v but got %v", expected, result)
	}

	os.Unsetenv("APP_BOOLEAN")
	os.Unsetenv("CUSTOM_NAME")
}

func TestIntShouldBeCorrect(t *testing.T) {
	os.Setenv("APP_BASE_URL", "world")
	defer os.Unsetenv("APP_BASE_URL")

	expected := 1
	os.Setenv("APP_INTEGER", fmt.Sprint(expected))
	os.Setenv("CUSTOM_NAME", "test")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err != nil {
		t.Error("got error", err)
	}

	result := cfg.Integer
	if result != expected {
		t.Errorf("expected %v but got %v", expected, result)
	}

	os.Unsetenv("APP_INTEGER")
	os.Unsetenv("CUSTOM_NAME")
}

func TestIntShouldBeError(t *testing.T) {
	os.Setenv("APP_BASE_URL", "world")
	defer os.Unsetenv("APP_BASE_URL")

	os.Setenv("APP_INTEGER", "word")
	os.Setenv("CUSTOM_NAME", "test")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err == nil {
		t.Error("expected error but got", err)
	}

	os.Unsetenv("APP_INTEGER")
	os.Unsetenv("CUSTOM_NAME")
}

func TestFloatShouldBeError(t *testing.T) {
	os.Setenv("APP_BASE_URL", "world")
	defer os.Unsetenv("APP_BASE_URL")

	os.Setenv("APP_FLOAT", "word")
	os.Setenv("CUSTOM_NAME", "test")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err == nil {
		t.Error("expected error but got", err)
	}

	os.Unsetenv("APP_FLOAT")
	os.Unsetenv("CUSTOM_NAME")
}

func TestRequiredValueShouldBeError(t *testing.T) {
	os.Setenv("APP_BASE_URL", "world")
	defer os.Unsetenv("APP_BASE_URL")

	os.Unsetenv("CUSTOM_NAME")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err == nil {
		t.Error("expected error but got", err)
	}
}

func TestCamalCaseShouldBeCorrect(t *testing.T) {
	os.Setenv("APP_BASE_URL", "world")
	defer os.Unsetenv("APP_BASE_URL")

	expected := "test"
	os.Setenv("APP_CAMAL_CASE", fmt.Sprint(expected))
	os.Setenv("CUSTOM_NAME", "test")

	cfg := Configuration{}
	if err := properties.Load(&cfg); err != nil {
		t.Error("got error", err)
	}

	result := cfg.CamalCase
	if result != expected {
		t.Errorf("expected %v but got %v", expected, result)
	}

	os.Unsetenv("APP_CAMAL_CASE")
	os.Unsetenv("CUSTOM_NAME")
}
