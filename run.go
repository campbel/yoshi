package yoshi

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Config struct {
	HelpWriter io.Writer
}

var defaultConfig = Config{HelpWriter: os.Stdout}

type Yoshi struct {
	name   string
	config Config
}

func New(name string) *Yoshi {
	return &Yoshi{name: name, config: defaultConfig}
}

func (y *Yoshi) WithConfig(config Config) *Yoshi {
	y.config = config
	return y
}

func (y *Yoshi) Run(v any, args ...string) error {
	ctx := newContext(y.name)
	err := ctx.run(v, args...)
	switch err := err.(type) {
	case *userError:
		fmt.Fprint(y.config.HelpWriter, help(err.loc, err.err, ctx.chain...))
		if err.err == errHelp {
			return nil
		}
		return err.err
	default:
		return err
	}
}

var defaultContext = New(filepath.Base(os.Args[0]))

func Run(v any, args ...string) error {
	return defaultContext.Run(v, args...)
}

type runContext struct {
	chain []string
}

func newContext(name string) *runContext {
	return &runContext{chain: []string{name}}
}

func (ctx *runContext) run(v any, args ...string) error {
	if len(args) == 0 {
		args = os.Args[1:]
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Func:
		return ctx.handleFunc(val, args...)
	case reflect.Struct:
		return ctx.handleStruct(val, args...)
	default:
		return fmt.Errorf("expected a function or struct, got %s", val.Kind())
	}
}

func (ctx *runContext) handleFunc(val reflect.Value, args ...string) error {
	if val.Kind() != reflect.Func {
		return fmt.Errorf("expected a function, got %s", val.Kind())
	}
	fnType := val.Type()
	var parameters []reflect.Value
	for i := 0; i < fnType.NumIn(); i++ {
		paramType := fnType.In(i)
		if paramType.Kind() != reflect.Struct {
			return fmt.Errorf("expected a struct, got %s", paramType.Kind())
		}
		param := reflect.New(paramType)
		if err := defaults(param); err != nil {
			return err
		}
		if err := options(param, args...); err != nil {
			return userErr(paramType, err)
		}
		parameters = append(parameters, param.Elem())
	}
	ret := val.Call(parameters)
	if len(ret) > 0 {
		if err, ok := ret[0].Interface().(error); ok {
			return err
		}
	}
	return nil
}

func (ctx *runContext) handleStruct(val reflect.Value, args ...string) error {
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got %s", val.Kind())
	}
	fields := reflect.VisibleFields(val.Type())
	for _, field := range fields {
		if args[0] != strings.ToLower(field.Name) {
			continue
		}
		fieldVal := val.FieldByName(field.Name)
		switch fieldVal.Kind() {
		case reflect.Struct:
			ctx.chain = append(ctx.chain, strings.ToLower(field.Name))
			return ctx.handleStruct(fieldVal, args[1:]...)
		case reflect.Func:
			ctx.chain = append(ctx.chain, strings.ToLower(field.Name))
			return ctx.handleFunc(fieldVal, args[1:]...)
		default:
			return fmt.Errorf("expected a struct or function, got %s", fieldVal.Kind())
		}
	}
	return userErrf(val.Type(), "command not found: %s", args[0])
}
