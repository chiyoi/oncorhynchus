package trinity

import (
	"flag"
	"fmt"
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
	logs.SetPrefix(fmt.Sprintf("[%d] ", os.Getpid()))

	data.Load()
	defer data.Save()

	c := Command()
	c.ServeArgs(os.Args[1:])
}

func Command() *sakana.Command {
	c := sakana.NewCommand(Name, Usage, Description)
	c.Welcome("Nyan~")
	c.Command(listen.Command())
	c.Command(post.Command())

	c.Work(func(fs *flag.FlagSet) {
		if fs.NArg() <= 0 {
			sakana.UsageError("Subcommand is needed.", fs.Usage)
		}
	})
	return c
}

func LogFile() (f *os.File, clean func()) {
	f, err := os.OpenFile(filepath.Join(config.DirData, "log.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		logs.Error("cannot create log file")
		sakana.InternalError()
	}
	return f, func() { f.Close() }
}
