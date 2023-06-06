package listen

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
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

func Command() (name string, h sakana.Handler, description string) {
	return Name, Handler(), Description
}

func Handler() sakana.Handler {
	c := sakana.NewCommand(Name, Usage, Description)
	c.Welcome("command: trinity listen")
	c.Work(Work)
	return c
}

func Work(*flag.FlagSet) {
	if data.Data.Token.Expired() {
		data.SetToken(auth.Refresh(data.Data.Token))
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	token := data.Data.Token.AccessToken
	errorCount := 0

	logs.Info("start listener")
	fmt.Println("Listening.")
	for func() bool {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			select {
			case <-ctx.Done():
			case s := <-sig:
				logs.Info("stop:", s)
				cancel()
			}
		}()

		ms, err := PollUpdate(ctx, token)
		if err != nil {
			var re *kitsune.ResponseError
			if errors.As(err, &re) && re.StatusCode == http.StatusGatewayTimeout {
				logs.Info("polling timeout")
				return true
			}

			select {
			case <-ctx.Done():
				fmt.Print("Canceled.")
				return false
			default:
				logs.Warning(err)
				if errorCount++; errorCount > 3 {
					logs.Error("too many errors")
					sakana.InternalError()
				}
				return true
			}
		}

		errorCount = 0
		for _, m := range ms {
			PrintMessage(m)
		}
		return true
	}() {
	}
}

func PollUpdate(ctx context.Context, token string) (ms []namedMessage, err error) {
	logs.Info("polling update")
	var resp neko03ResponseFetch
	if err = (&kitsune.Client{Endpoint: endpointPollUpdate()}).RoundTrip(ctx, nil, &resp); err != nil {
		return
	}

	logs.Info("received update")
	return resp.Messages, nil
}

func PrintMessage(m namedMessage) {
	logs.Info("received message:", m.ID)
	fmt.Printf("%s (%s) - %s\n",
		m.SenderName.DisplayName,
		m.SenderName.UserPrincipalName,
		time.UnixMilli(m.Timestamp).UTC().Format(time.RFC822),
	)

	for _, p := range m.Content {
		switch p.Type {
		case trinity.ParagraphTypeText:
			fmt.Println(p.Data)
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
