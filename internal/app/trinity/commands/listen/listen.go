package listen

import (
	"context"
	"errors"
	"fmt"
	"io"
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

func Command() (name, description string, h sakana.Handler) {
	return Name, Description, Handler()
}

func Handler() sakana.Handler {
	c := sakana.NewCommand(Name)
	c.Welcome("command: trinity listen")
	c.Summary(Usage, Description)

	c.Work(sakana.HandlerFunc(func(w io.Writer, _ []string) {
		if data.Data.Token.Expired() {
			data.SetToken(auth.Refresh(data.Data.Token))
		}

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		token := data.Data.Token.AccessToken
		errorCount := 0

		logs.Info("start listener")
		fmt.Fprintln(w, "Listening.")
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
					fmt.Fprint(w, "Canceled.")
					return false
				default:
					logs.Warning(err)
					if errorCount++; errorCount > 3 {
						logs.Error("too many errors")
						sakana.InternalError(w)
					}
					return true
				}
			}

			errorCount = 0
			for _, m := range ms {
				PrintMessage(w, m)
			}
			return true
		}() {
		}
	}))
	return c
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

func PrintMessage(w io.Writer, m namedMessage) {
	logs.Info("received message:", m.ID)
	if _, err := fmt.Fprintf(w, "%s (%s) - %s\n",
		m.SenderName.DisplayName,
		m.SenderName.UserPrincipalName,
		time.UnixMilli(m.Timestamp).UTC().Format(time.RFC822),
	); err != nil {
		logs.Error(err)
		sakana.InternalError(w)
	}

	for _, p := range m.Content {
		switch p.Type {
		case trinity.ParagraphTypeText:
			if _, err := fmt.Fprintln(w, p.Data); err != nil {
				logs.Error(err)
				sakana.InternalError(w)
			}
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
