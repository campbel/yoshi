package yoshi

import (
	"testing"
)

func TestApp(t *testing.T) {

	tt := []struct {
		name      string
		args      []string
		rootArgs  []string
		someArgs  []string
		otherArgs []string
		help      bool
	}{
		{
			name:      "no args",
			args:      []string{},
			rootArgs:  []string{},
			someArgs:  []string{},
			otherArgs: []string{},
		},
		{
			name:      "root args",
			args:      []string{"-v"},
			rootArgs:  []string{"-v"},
			someArgs:  []string{},
			otherArgs: []string{},
		},
		{
			name:      "some args",
			args:      []string{"some", "-n", "foo"},
			rootArgs:  []string{},
			someArgs:  []string{"-n", "foo"},
			otherArgs: []string{},
		},
		{
			name:      "other args",
			args:      []string{"other", "-n", "foo"},
			rootArgs:  []string{},
			someArgs:  []string{},
			otherArgs: []string{"-n", "foo"},
		},
		{
			name:      "root and some args",
			args:      []string{"-v", "some", "-n", "foo"},
			rootArgs:  []string{"-v"},
			someArgs:  []string{"-n", "foo"},
			otherArgs: []string{},
		},
		{
			name:      "root and other args",
			args:      []string{"-v", "other", "-n", "foo"},
			rootArgs:  []string{"-v"},
			someArgs:  []string{},
			otherArgs: []string{"-n", "foo"},
		},
		{
			name:      "some and other args",
			args:      []string{"some", "-n", "foo", "other", "-n", "bar"},
			rootArgs:  []string{},
			someArgs:  []string{"-n", "foo", "other", "-n", "bar"},
			otherArgs: []string{},
		},
		{
			name:      "help",
			args:      []string{"--help"},
			rootArgs:  []string{"--help"},
			someArgs:  []string{},
			otherArgs: []string{},
			help:      true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var rootArgs, someArgs, otherArgs []string
			var help bool
			app := App(Run(func(args Args) {
				rootArgs = args
				help = help || args.Help()
			}))
			app.Sub("some", Run(func(args Args) {
				someArgs = args
				help = help || args.Help()
			}))
			app.Sub("other", Run(func(args Args) {
				otherArgs = args
				help = help || args.Help()
			}))
			app.Parse(tc.args)
			equal(t, rootArgs, tc.rootArgs)
			equal(t, someArgs, tc.someArgs)
			equal(t, otherArgs, tc.otherArgs)
			if tc.help != help {
				t.Errorf("got %v, want %v", help, tc.help)
			}
		})
	}
}

func equal(t *testing.T, a1, a2 []string) {
	if len(a1) != len(a2) {
		t.Errorf("got %v, want %v", a1, a2)
	}
	for i := range a1 {
		if a1[i] != a2[i] {
			t.Errorf("got %v, want %v", a1, a2)
		}
	}
}
