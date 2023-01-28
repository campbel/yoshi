package yoshi

import "strings"

type args []*arg

type arg struct {
	command string
	flag    string
	value   string
}

// parseArgs parses the given arguments into a slice of arg structs.
// The boolFlags argument is a list of flags that should be treated as boolean
// flags. If a flag is in the boolFlags list, then the value of the flag is
// optional. If the value is not provided, then the value is assumed to be true.
// If the value is provided, then it must be either "true" or "false".
// It is necessary to know which flags are boolean flags because the parser
// needs to know whether to treat the next argument as a value or as a command.
func parseArgs(arguments []string, boolFlags ...string) args {
	var parsed args
LOOP:
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
		for _, boolFlag := range boolFlags {
			if flag == boolFlag {
				if arguments[i+1] == "true" || arguments[i+1] == "false" {
					parsed = append(parsed, &arg{flag: flag, value: arguments[i+1]})
					i += 2
					continue LOOP
				}
				parsed = append(parsed, &arg{flag: flag, value: "true"})
				i++
				continue LOOP
			}
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
