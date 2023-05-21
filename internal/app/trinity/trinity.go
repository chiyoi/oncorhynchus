package trinity

import (
	"os"

	"github.com/chiyoi/go/pkg/sakana"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/commands/listen"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/config"
)

const (
	Name        = "trinity"
	Usage       = "trinity <command> ..."
	Description = "The trinity command-line interface."
)

func Main() {
	config.Data.Load()
	defer config.Data.Save()

	c := Command()
	c.ServeArgs(os.Args[1:])
}

func Command() *sakana.Command {
	c := sakana.NewCommand(Name, Usage, Description)
	c.Welcome("Nyan~")
	c.Command(listen.Command())
	return c
}
