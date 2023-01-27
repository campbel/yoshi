package yoshi

import "strings"

type args []*arg

type arg struct {
	command string
	flag    string
	value   string
}

func parseArgs(arguments []string) args {
	var parsed args
	for i := 0; i < len(arguments); {
		if !strings.HasPrefix(arguments[i], "-") {
			parsed = append(parsed, &arg{command: arguments[i]})
			i++
			continue
		}
		flag := arguments[i]
		if i+1 >= len(arguments) {
			parsed = append(parsed, &arg{flag: flag})
			break
		}
		if strings.HasPrefix(arguments[i+1], "-") {
			parsed = append(parsed, &arg{flag: flag})
			i++
			continue
		}
		parsed = append(parsed, &arg{flag: flag, value: arguments[i+1]})
		i += 2
	}
	return parsed
}
