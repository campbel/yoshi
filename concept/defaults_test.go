package concept

import (
	"reflect"
	"testing"
)

type T struct {
	A   string `yoshi-def:"foo"`
	B   string
	C   int               `yoshi-def:"42"`
	D   int64             `yoshi-def:"42"`
	E   bool              `yoshi-def:"true"`
	F   []string          `yoshi-def:"foo,bar,baz"`
	FF  []bool            `yoshi-def:"true,false,true"`
	FFF []int             `yoshi-def:"1,2,3"`
	G   map[string]string `yoshi-def:"foo=bar,bar=qux"`
	GG  map[string]int    `yoshi-def:"foo=1,bar=2"`
	GGG map[string]bool   `yoshi-def:"foo=true,bar=false"`
}

func TestSetDefaults(t *testing.T) {
	var tt T
	eval(&tt)
	if tt.A != "foo" {
		t.Errorf("expected tt.A to be foo, got %s", tt.A)
	}
	if tt.B != "" {
		t.Errorf("expected tt.B to be empty, got %s", tt.B)
	}
	if tt.C != 42 {
		t.Errorf("expected tt.C to be 42, got %d", tt.C)
	}
	if tt.D != 42 {
		t.Errorf("expected tt.D to be 42, got %d", tt.D)
	}
	if tt.E != true {
		t.Errorf("expected tt.E to be true, got %t", tt.E)
	}
	if !reflect.DeepEqual(tt.F, []string{"foo", "bar", "baz"}) {
		t.Errorf("expected tt.F to be [foo bar baz], got %s", tt.F)
	}
	if !reflect.DeepEqual(tt.FF, []bool{true, false, true}) {
		t.Errorf("expected tt.FF to be [true false true], got %v", tt.FF)
	}
	if !reflect.DeepEqual(tt.FFF, []int{1, 2, 3}) {
		t.Errorf("expected tt.FFF to be [1 2 3], got %v", tt.FFF)
	}
	if !reflect.DeepEqual(tt.G, map[string]string{"foo": "bar", "bar": "qux"}) {
		t.Errorf("expected tt.G to be {foo: bar, bar: qux}, got %s", tt.G)
	}
	if !reflect.DeepEqual(tt.GG, map[string]int{"foo": 1, "bar": 2}) {
		t.Errorf("expected tt.GG to be {foo: 1, bar: 2}, got %v", tt.GG)
	}
	if !reflect.DeepEqual(tt.GGG, map[string]bool{"foo": true, "bar": false}) {
		t.Errorf("expected tt.GGG to be {foo: true, bar: false}, got %v", tt.GGG)
	}
}
