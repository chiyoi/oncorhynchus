package verify

import (
	"fmt"
	"os"

	"github.com/chiyoi/oncorhynchus"
	"github.com/chiyoi/oncorhynchus/internal/api"
	"github.com/chiyoi/oncorhynchus/pkg/sakana"
)

var config = &oncorhynchus.Config

type Request struct {
	Name string `json:"name"`
}

type Response struct {
	Passed bool `json:"passed"`
}

func Handler() sakana.Handler {
	c := sakana.NewCommand("verify", "[<options>]", "verify if name exist")
	c.Welcome("command: oncorhynchus authenticator verify\n")

	var name, password string

	c.StringVar(&name, "n", "", "")
	c.StringVar(&name, "name", "", "")
	c.OptionUsage([]string{"-n", "--name"}, true, "name")

	c.StringVar(&password, "p", "", "")
	c.StringVar(&password, "password", "", "")
	c.OptionUsage([]string{"-p", "--password"}, true, "password")

	c.Work(func() {
		if name == "" || password == "" {
			fmt.Fprintln(c.Output(), "name or password is missing")
			c.Usage()
			os.Exit(5)
		}

		var resp Response
		if err := api.RoundTrip(config.AuthenticatorEndpoint.JoinPath("verify").String(), Request{
			Name: name,
		}, &resp); err != nil {
			fmt.Fprintln(c.Output(), err)
			os.Exit(1)
		}

		fmt.Println("passed:", resp.Passed)
	})

	c.Example(
		"oncorhynchus authenticator verify -n nacho -p ***",
		"verify a neko with name `nacho` and password `***`",
	)

	return c
}
