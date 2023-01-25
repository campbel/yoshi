package concept

import (
	"reflect"
	"strconv"
	"strings"
)

func Run(app any, args ...string) {
	Eval(app)
	parseCommand(reflect.ValueOf(app), args)
}

func parseCommand(command reflect.Value, args []string) {
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
		parseOptions(command.Elem().FieldByName("Options").Addr(), args[:i])
	}
	if hasRun {
		runFunc := command.Elem().FieldByName("Run")
		if hasOpt {
			runFunc.Call([]reflect.Value{
				command.Elem().FieldByName("Options"),
			})
		} else {
			runFunc.Call([]reflect.Value{})
		}
	}
	if cmd != "" {
		parseCommand(command.Elem().FieldByName(cmd).Addr(), args[i+1:])
	}
}

func parseOptions(options reflect.Value, args []string) {
	flagMap := flagMap(args)
	fields := reflect.VisibleFields(options.Elem().Type())
	for _, field := range fields {
		flags := parseTag("flag", field.Tag.Get("yoshi"))
		if len(flags) == 0 {
			continue
		}
		var value string
		for _, flag := range flags {
			if value = flagMap[flag]; value != "" {
				break
			}
		}
		if value == "" {
			continue
		}
		prop := options.Elem().FieldByName(field.Name)
		switch prop.Kind() {
		case reflect.String:
			prop.SetString(value)
		case reflect.Int:
			i, _ := strconv.Atoi(value)
			prop.SetInt(int64(i))
		case reflect.Bool:
			b, _ := strconv.ParseBool(value)
			prop.SetBool(b)
		case reflect.Slice:
			prop.Set(reflect.MakeSlice(prop.Type(), 0, 0))
			for _, v := range strings.Split(value, ",") {
				switch prop.Type().Elem().Kind() {
				case reflect.String:
					prop.Set(reflect.Append(prop, reflect.ValueOf(v)))
				case reflect.Int:
					i, _ := strconv.Atoi(v)
					prop.Set(reflect.Append(prop, reflect.ValueOf(i)))
				case reflect.Bool:
					b, _ := strconv.ParseBool(v)
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
					i, _ := strconv.Atoi(p2[1])
					prop.SetMapIndex(reflect.ValueOf(p2[0]), reflect.ValueOf(i))
				case reflect.Bool:
					b, _ := strconv.ParseBool(p2[1])
					prop.SetMapIndex(reflect.ValueOf(p2[0]), reflect.ValueOf(b))
				}
			}
		}
	}
}

func flagMap(args []string) map[string]string {
	flags := map[string]string{}
	for i := 0; i < len(args); i += 2 {
		flags[args[i]] = args[i+1]
	}
	return flags
}

func parseTag(typ, tag string) []string {
	parts := strings.Split(tag, ";")
	for _, setting := range parts {
		p2 := strings.Split(setting, "=")
		key := p2[0]
		if key == typ {
			return strings.Split(p2[1], ",")
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
