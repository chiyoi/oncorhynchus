package trinity

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/chiyoi/go/pkg/logs"
	"github.com/chiyoi/go/pkg/sakana"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/commands/listen"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/commands/post"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/config"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/data"
)

const (
	Name        = "trinity"
	Usage       = "trinity <command> ..."
	Description = "The trinity command-line interface."
)

func Main() {
	defer func() {
		if err := recover(); err != nil {
			os.Exit(1)
		}
	}()

	f, clean := LogFile()
	defer clean()

	logs.SetOutput(f)
	sakana.SetLogFile(f)
	logs.PrependPrefix(fmt.Sprintf("[%d] ", os.Getpid()))

	data.Load()
	defer data.Save()

	c := Command()
	c.ServeArgs(nil, os.Args[1:])
}

func Command() (c *sakana.Command) {
	c = sakana.NewCommand(Name)

	c.Welcome("Nyan~")
	c.Summary(Usage, Description)
	c.Command(listen.Command())
	c.Command(post.Command())

	c.Work(sakana.HandlerFunc(func(w io.Writer, args []string) {
		if len(args) == 0 {
			sakana.UsageError(w, "Subcommand is needed.\n"+c.Usage())
		}
	}))
	return c
}

func LogFile() (f *os.File, clean func()) {
	f, err := os.OpenFile(filepath.Join(config.DirData, "log.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		logs.Error("cannot create log file")
		sakana.InternalError(os.Stderr)
	}
	return f, func() { f.Close() }
}
