package authenticator

import (
	"github.com/chiyoi/oncorhynchus/internal/app/oncorhynchus/authenticator/verify"
	"github.com/chiyoi/oncorhynchus/pkg/sakana"
)

func Handler() sakana.Handler {
	c := sakana.NewCommandGroup("authenticator", "authenticator <command>", "interact with neko03authenticator")
	c.Welcome("command: oncorhynchus authenticator\n")

	c.Command("verify", verify.Handler(), "verify name")

	return c
}
