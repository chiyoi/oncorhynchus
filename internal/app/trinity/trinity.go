package trinity

import (
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
	f, clean := LogFile()
	defer clean()

	logs.SetOutput(f)
	logs.SetPrefix(fmt.Sprintf("[%d] ", os.Getpid()))

	// sig := make(chan os.Signal, 1)
	// signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	// go func() {
	// 	s := <-sig
	// 	logs.Info("stop:", s)
	// 	fmt.Println("Stop.")
	// 	os.Exit(0)
	// }()

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
