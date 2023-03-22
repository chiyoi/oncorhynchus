package oncorhynchus

import (
	"github.com/chiyoi/oncorhynchus/internal/oncorhynchus/authenticator"
	"github.com/chiyoi/oncorhynchus/pkg/sakana"
)

const welcome = `

Welcome to oncorhynchus!

`

func Handler() sakana.Handler {
	c := sakana.NewCommandGroup("oncorhynchus", "oncorhynchus <command> ...", "cli app")
	c.Welcome(welcome)

	c.Command("authenticator", authenticator.Handler(), "interact with neko03authenticator")

	return c
}
