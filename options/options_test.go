package options

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsNesting(t *testing.T) {
	assert := assert.New(t)
	type Options struct {
		Name    string `yoshi:"-n,--name"`
		Address struct {
			Street string `yoshi:"-s,--street"`
			City   string `yoshi:"-c,--city"`
		}
	}

	var opts Options
	err := options(reflect.ValueOf(&opts), []string{"-n", "mario", "-s", "123 main st", "-c", "new york"})
	assert.NoError(err)
	assert.Equal("mario", opts.Name)
	assert.Equal("123 main st", opts.Address.Street)
	assert.Equal("new york", opts.Address.City)
}

func TestOptionsNestingWithPositionals(t *testing.T) {
	assert := assert.New(t)
	type Options struct {
		Name    string `yoshi:"NAME"`
		Address struct {
			Street string `yoshi:"-s,--street"`
			City   string `yoshi:"-c,--city"`
		}
	}

	var opts Options
	err := options(reflect.ValueOf(&opts), []string{"mario", "-s", "123 main st", "-c", "new york"})
	assert.NoError(err)
	assert.Equal("mario", opts.Name)
	assert.Equal("123 main st", opts.Address.Street)
	assert.Equal("new york", opts.Address.City)
}

func TestOptionsNestingWithDefaults(t *testing.T) {
	assert := assert.New(t)
	type Options struct {
		Name    string `yoshi:"NAME"`
		Address struct {
			Street string `yoshi:"-s,--street"`
			City   string `yoshi:"-c,--city;;new york"`
		}
	}

	var opts Options
	assert.NoError(defaults(reflect.ValueOf(&opts)))
	assert.NoError(options(reflect.ValueOf(&opts), []string{"mario", "-s", "123 main st"}))
	assert.Equal("mario", opts.Name)
	assert.Equal("123 main st", opts.Address.Street)
	assert.Equal("new york", opts.Address.City)
}

func TestMultiplePositionals(t *testing.T) {
	assert := assert.New(t)
	type Options struct {
		Name    string `yoshi:"NAME"`
		Address struct {
			Street string `yoshi:"STREET"`
			City   string `yoshi:"CITY"`
		}
	}

	var opts Options
	assert.NoError(options(reflect.ValueOf(&opts), []string{"mario", "123 main st", "new york"}))
	assert.Equal("mario", opts.Name)
	assert.Equal("123 main st", opts.Address.Street)
	assert.Equal("new york", opts.Address.City)
}

func TestInvalidParamTypes(t *testing.T) {
	var opts string
	assert.Error(t, options(reflect.ValueOf(&opts), []string{}))

	type Options struct{}
	assert.Error(t, options(reflect.ValueOf(Options{}), []string{}))
}
