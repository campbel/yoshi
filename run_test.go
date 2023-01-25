package yoshi

import (
	"strings"
	"testing"
)

type NuApp struct {
	Options struct {
		Verbose bool `yoshi-flag:"-v,--verbose"`
	}
	Server struct {
		Options ServerOptions
		Run     func(ServerOptions)
	}
	Client struct {
		Options ClientOptions
		Run     func(ClientOptions)
	}
}

type ClientOptions struct {
	Name     string            `yoshi-flag:"-n,--name"      yoshi-desc:"name of client"`
	Insecure bool              `yoshi-flag:"-k,--insecure"  yoshi-desc:"insecure connection"`
	Count    int               `yoshi-flag:"-c,--count"     yoshi-desc:"number of times to run" yoshi-def:"1"`
	ListStr  []string          `yoshi-flag:"-l,--list"      yoshi-desc:"list of strings"`
	ListInt  []int             `yoshi-flag:"-i,--int-list"  yoshi-desc:"list of ints"           yoshi-def:"1,2,3"`
	ListBool []bool            `yoshi-flag:"-b,--bool-list" yoshi-desc:"list of bools"`
	DictStr  map[string]string `yoshi-flag:"-d,--dict"      yoshi-desc:"dict of strings"`
	DictInt  map[string]int    `yoshi-flag:"-e,--int-dict"  yoshi-desc:"dict of ints"`
	DictBool map[string]bool   `yoshi-flag:"-f,--bool-dict" yoshi-desc:"dict of bools"          yoshi-def:"a=true,b=false,c=true"`
}

type ServerOptions struct {
	Port string `yoshi-flag:"-p,--port"`
}

func TestConcept(t *testing.T) {
	count := 0
	app := new(NuApp)
	app.Client.Run = func(options ClientOptions) {
		if options.Name != "bob" {
			t.Error("name should be bob")
		}
		count++
	}
	app.Server.Run = func(options ServerOptions) {
		t.Fatal("server should not run")
	}

	args := "--verbose true "
	args += "client --name bob --insecure --count 1 "
	args += "--list 1,2,3 --int-list 1,2,3 --bool-list true,false,true "
	args += "--dict a=d,b=e,c=f --int-dict a=1,b=2,c=3 --bool-dict a=true,b=false,c=true"
	run("tt", app, strings.Split(args, " ")...)
	if count != 1 {
		t.Error("count should be 2")
	}
}

func TestHelp(t *testing.T) {
	app := new(NuApp)
	run("tt", app, "client", "--help")
}
