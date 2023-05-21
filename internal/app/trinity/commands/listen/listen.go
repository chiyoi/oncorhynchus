package listen

import (
	"context"
	"flag"
	"net/http"
	"net/url"

	"github.com/chiyoi/go/pkg/kitsune"
	"github.com/chiyoi/go/pkg/logs"
	"github.com/chiyoi/go/pkg/sakana"
	"github.com/chiyoi/go/pkg/trinity"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/config"
)

const (
	Name        = "listen"
	Usage       = "listen"
	Description = "Listen to new messages."
)

var (
	EndpointPollUpdate = func() string {
		u, err := url.Parse(config.EndpointNeko03)
		if err != nil {
			logs.Panic(err)
		}

		return u.JoinPath("/trinity/fetch/update").String()
	}()
)

func Command() (name string, h sakana.Handler, usage string) {
	return Name, Handler(), Description
}

func Handler() sakana.Handler {
	c := sakana.NewCommand(Name, Usage, Description)
	c.Welcome("command: trinity listen")
	c.Work(Work)
	return c
}

func Work(*flag.FlagSet) {
	ch := make(chan trinity.Message)
	go pollUpdate()
}

func pollUpdate(token string, ch chan<- trinity.Message) {
	header := http.Header{}
	header.Set("Authorization", "Bearer "+token)
	c := kitsune.Client{
		Endpoint: EndpointPollUpdate,
		Header:   header,
	}

	for {
		logs.Info("polling update")
		var m trinity.Message
		if err := c.RoundTrip(context.Background(), nil, &m); err != nil {
			logs.Warning(err)
			continue
		}

		logs.Info("received update")
		ch <- m
	}
}
