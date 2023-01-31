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
				{key: "-n", value: "my name"},
			},
		},
		{
			name: "multiple",
			args: []string{"-n", "my name", "-a", "my address"},
			parsed: []arg{
				{key: "-n", value: "my name"},
				{key: "-a", value: "my address"},
			},
		},
		{
			name: "multiple with command",
			args: []string{"-n", "my name", "-a", "my address", "call", "-n", "123", "message", "-t", "hello, world"},
			parsed: []arg{
				{key: "-n", value: "my name"},
				{key: "-a", value: "my address"},
				{value: "call"},
				{key: "-n", value: "123"},
				{value: "message"},
				{key: "-t", value: "hello, world"},
			},
		},
		{
			name: "trick bool parser",
			args: []string{"message", "-t", "hello, world", "-b", "foobar"},
			parsed: []arg{
				{value: "message"},
				{key: "-t", value: "hello, world"},
				{key: "-b", value: "true"},
				{value: "foobar"},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parsed := parseArgs(tc.args, "-b")
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
