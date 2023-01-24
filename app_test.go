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
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var rootArgs, someArgs, otherArgs []string
			app := App(func(args []string) {
				rootArgs = args
			})
			app.Sub("some", func(args []string) {
				someArgs = args
			})
			app.Sub("other", func(args []string) {
				otherArgs = args
			})
			app.Parse(tc.args)
			equal(t, rootArgs, tc.rootArgs)
			equal(t, someArgs, tc.someArgs)
			equal(t, otherArgs, tc.otherArgs)
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
