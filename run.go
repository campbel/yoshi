package yoshi

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func Run(name string, app any, args ...string) {
	eval(app)
	ctx := &context{}
	if err := ctx.parseCommand(reflect.ValueOf(app), args); err != nil {
		if err.Error() == "help" {
			fmt.Print(help(ctx.currentCommand, nil, append([]string{name}, ctx.commandList...)...))
			os.Exit(0)
		}
		fmt.Print(help(ctx.currentCommand, err, append([]string{name}, ctx.commandList...)...))
		os.Exit(1)
	}
}

type context struct {
	commandList    []string
	currentCommand reflect.Type
}

func (ctx *context) parseCommand(command reflect.Value, args []string) error {
	ctx.currentCommand = command.Elem().Type()
	var (
		hasRun = false
		hasOpt = false
		cmds   []string
	)
	fields := reflect.VisibleFields(command.Elem().Type())
	for _, field := range fields {
		switch field.Name {
		case "Options":
			hasOpt = true
		case "Run":
			hasRun = true
		default:
			cmds = append(cmds, field.Name)
		}
	}
	cmd, i := firstSubCommand(cmds, args)
	if hasOpt {
		if err := parseOptions(command.Elem().FieldByName("Options").Addr(), args[:i]); err != nil {
			return err
		}
	}
	if hasRun {
		runFunc := command.Elem().FieldByName("Run")
		if !runFunc.IsNil() {
			var retVal []reflect.Value
			if hasOpt {
				retVal = runFunc.Call([]reflect.Value{
					command.Elem().FieldByName("Options"),
				})
			} else {
				retVal = runFunc.Call([]reflect.Value{})
			}
			if len(retVal) > 0 && retVal[0].Interface() != nil {
				return retVal[0].Interface().(error)
			}
		}
	}
	if cmd != "" {
		ctx.commandList = append(ctx.commandList, cmd)
		return ctx.parseCommand(command.Elem().FieldByName(cmd).Addr(), args[i+1:])
	}
	return nil
}

func parseOptions(options reflect.Value, args []string) error {
	// If help, return with help error
	for _, arg := range args {
		if arg == "--help" {
			return errors.New("help")
		}
	}
	flagToFieldMap := make(map[string]reflect.StructField)
	fields := reflect.VisibleFields(options.Elem().Type())
	for _, field := range fields {
		flags := strings.Split(field.Tag.Get("yoshi-flag"), ",")
		if len(flags) == 0 {
			continue
		}
		for _, flag := range flags {
			flagToFieldMap[flag] = field
		}
	}

	for i := 0; i < len(args); i += 2 {
		arg := args[i]
		field, ok := flagToFieldMap[arg]
		if !ok {
			return fmt.Errorf("unknown flag: %s", arg)
		}
		var value string
		if len(args) > i+1 {
			value = args[i+1]
		}
		prop := options.Elem().FieldByName(field.Name)
		switch prop.Kind() {
		case reflect.String:
			prop.SetString(value)
		case reflect.Int:
			i, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			prop.SetInt(int64(i))
		case reflect.Bool:
			b, err := strconv.ParseBool(value)
			if err != nil {
				prop.SetBool(true)
				i--
			} else {
				prop.SetBool(b)
			}
		case reflect.Slice:
			prop.Set(reflect.MakeSlice(prop.Type(), 0, 0))
			for _, v := range strings.Split(value, ",") {
				switch prop.Type().Elem().Kind() {
				case reflect.String:
					prop.Set(reflect.Append(prop, reflect.ValueOf(v)))
				case reflect.Int:
					i, err := strconv.Atoi(v)
					if err != nil {
						return err
					}
					prop.Set(reflect.Append(prop, reflect.ValueOf(i)))
				case reflect.Bool:
					b, err := strconv.ParseBool(v)
					if err != nil {
						return err
					}
					prop.Set(reflect.Append(prop, reflect.ValueOf(b)))
				}
			}
		case reflect.Map:
			prop.Set(reflect.MakeMap(prop.Type()))
			for _, v := range strings.Split(value, ",") {
				p2 := strings.Split(v, "=")
				switch prop.Type().Elem().Kind() {
				case reflect.String:
					prop.SetMapIndex(reflect.ValueOf(p2[0]), reflect.ValueOf(p2[1]))
				case reflect.Int:
					i, err := strconv.Atoi(p2[1])
					if err != nil {
						return err
					}
					prop.SetMapIndex(reflect.ValueOf(p2[0]), reflect.ValueOf(i))
				case reflect.Bool:
					b, err := strconv.ParseBool(p2[1])
					if err != nil {
						return err
					}
					prop.SetMapIndex(reflect.ValueOf(p2[0]), reflect.ValueOf(b))
				}
			}
		}
	}
	return nil
}

func firstSubCommand(cmds, args []string) (string, int) {
	for i, arg := range args {
		for _, sub := range cmds {
			if arg == strings.ToLower(sub) {
				return sub, i
			}
		}
	}
	return "", len(args)
}
