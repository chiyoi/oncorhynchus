package post

import (
	"flag"
	"fmt"

	"github.com/chiyoi/go/pkg/kitsune"
	"github.com/chiyoi/go/pkg/logs"
	"github.com/chiyoi/go/pkg/sakana"
	"github.com/chiyoi/go/pkg/trinity"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/config"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/data"
)

const (
	Name        = "post"
	Usage       = "post -t <text-message>"
	Description = "Post a text message."
)

func Command() (name string, h sakana.Handler, usage string) {
	return Name, Handler(), Description
}

func Handler() sakana.Handler {
	c := sakana.NewCommand(Name, Usage, Description)
	c.Welcome("command: trinity post")

	text := c.FlagSet.String("t", "", "")
	c.FlagSet.StringVar(text, "text", "", "")
	c.OptionUsage([]string{"t", "text"}, true, "Text to post.")

	c.Work(func(fs *flag.FlagSet) {
		if *text == "" {
			logs.Warning("empty text")
			sakana.UsageError("Empty text.", fs.Usage)
		}

		req := neko03RequestPost{
			Content: []trinity.Paragraph{
				trinity.Text(*text),
			},
		}
		if err := kitsune.RoundTrip(endpointPost(), req, nil); err != nil {
			logs.Error("post failed:", err)
			sakana.InternalError()
			return
		}
		fmt.Println("Posted.")
	})
	return c
}

func endpointPost() string {
	return config.EndpointAuthedNeko03(
		data.Data.Token.AccessToken,
	).JoinPath(
		"/trinity/post",
	).String()
}

type neko03RequestPost struct {
	Content []trinity.Paragraph `json:"content"`
}
