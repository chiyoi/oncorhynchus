package trinity

import (
	"os"
	"path/filepath"

	"github.com/chiyoi/go/pkg/logs"
	"github.com/chiyoi/go/pkg/sakana"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/commands/listen"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/common/data"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/config"
)

const (
	Name        = "trinity"
	Usage       = "trinity <command> ..."
	Description = "The trinity command-line interface."
)

func Main() {
	data.Load()
	defer data.Save()

	logs.SetOutput(LogFile())

	c := Command()
	c.ServeArgs(os.Args[1:])
}

func Command() *sakana.Command {
	c := sakana.NewCommand(Name, Usage, Description)
	c.Welcome("Nyan~")
	c.Command(listen.Command())
	return c
}

func LogFile() *os.File {
	f, err := os.Create(filepath.Join(config.DirData, "log.txt"))
	if err != nil {
		logs.Fatal("cannot create log file")
	}
	return f
}
