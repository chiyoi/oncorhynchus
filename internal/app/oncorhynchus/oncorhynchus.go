package oncorhynchus

import (
	"github.com/chiyoi/oncorhynchus/internal/app/oncorhynchus/authenticator"
	"github.com/chiyoi/oncorhynchus/pkg/sakana"
)

const welcome = `

Welcome to oncorhynchus!

`

func Handler() sakana.Handler {
	c := sakana.NewCommandGroup("oncorhynchus", "oncorhynchus <command> ...", "command-line app")
	c.Welcome(welcome)

	c.Command("authenticator", authenticator.Handler(), "interact with neko03authenticator")

	return c
}
