package listen

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chiyoi/go/pkg/kitsune"
	"github.com/chiyoi/go/pkg/logs"
	"github.com/chiyoi/go/pkg/sakana"
	"github.com/chiyoi/go/pkg/trinity"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/auth"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/config"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/data"
)

const (
	Name        = "listen"
	Usage       = "listen"
	Description = "Listen to new messages."
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
	ch := make(chan namedMessage)
	token := data.Data.Token
	if token.Expired() {
		token = auth.Refresh(token)
		data.SetToken(token)
	}

	logs.Info("start listener")
	fmt.Println("Start listener.")
	go PollUpdate(token.AccessToken, ch)

	for m := range ch {
		logs.Info("received message:", m.ID)
		fmt.Printf("%s (%s) - %s",
			m.SenderName.DisplayName,
			m.SenderName.UserPrincipalName,
			time.Unix(m.Timestamp, 0).UTC().Format(time.RFC822),
		)
		for _, p := range m.Content {
			switch p.Type {
			case trinity.ParagraphTypeText:
				fmt.Println(p.Data)
			}
		}
	}
}

func PollUpdate(token string, ch chan<- namedMessage) {
	fc := 0

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		logs.Info("stop:", s)
		cancel()
	}()

loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stop.")
			close(ch)
			break loop
		default:
		}
		time.Sleep(time.Second)
		logs.Info("polling update")
		var resp neko03ResponseFetch
		if err := (&kitsune.Client{Endpoint: endpointPollUpdate()}).RoundTrip(ctx, nil, &resp); err != nil {
			if fc++; fc > 3 {
				logs.Error("too many polling failures")
				sakana.InternalError()
			}
			logs.Warning(err)
			continue
		}
		fc = 0
		logs.Info("received update")
		for _, m := range resp.Messages {
			ch <- m
		}
	}
}

func endpointPollUpdate() string {
	return config.EndpointAuthedNeko03(
		data.Data.Token.AccessToken,
	).JoinPath(
		"/trinity/fetch/update",
	).String()
}

type namedMessage struct {
	trinity.Message
	SenderName trinity.Name `json:"sender_name"`
}

type neko03ResponseFetch struct {
	Messages []namedMessage `json:"messages"`
}
