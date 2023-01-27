package yoshi

import "testing"

func TestParseArgs(t *testing.T) {
	tt := []struct {
		name   string
		args   []string
		parsed []arg
	}{
		{
			name:   "empty",
			args:   []string{},
			parsed: []arg{},
		},
		{
			name: "single",
			args: []string{"-n", "my name"},
			parsed: []arg{
				{flag: "-n", value: "my name"},
			},
		},
		{
			name: "multiple",
			args: []string{"-n", "my name", "-a", "my address"},
			parsed: []arg{
				{flag: "-n", value: "my name"},
				{flag: "-a", value: "my address"},
			},
		},
		{
			name: "multiple with command",
			args: []string{"-n", "my name", "-a", "my address", "call", "-n", "123", "message", "-t", "hello, world"},
			parsed: []arg{
				{flag: "-n", value: "my name"},
				{flag: "-a", value: "my address"},
				{command: "call"},
				{flag: "-n", value: "123"},
				{command: "message"},
				{flag: "-t", value: "hello, world"},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parsed := parseArgs(tc.args)
			if len(parsed) != len(tc.parsed) {
				t.Errorf("Expected %d args, got %d", len(tc.parsed), len(parsed))
			}
			for i, arg := range parsed {
				if *arg != tc.parsed[i] {
					t.Errorf("Expected %v, got %v", tc.parsed[i], arg)
				}
			}
		})
	}
}
