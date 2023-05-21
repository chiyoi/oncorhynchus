package trinity

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/chiyoi/go/pkg/az/auth"
	"github.com/chiyoi/go/pkg/logs"
	"github.com/chiyoi/go/pkg/neko"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/config"
)

var (
	Scopes = []string{"User.Read", "offline_access"}
)

func Login() (err error) {
	u, err := url.Parse(config.AzureADRedirectURI)
	if err != nil {
		return
	}

	ch := make(chan auth.Token)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", u.Port()),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code, err := auth.GetCode(r)
			if err != nil {
				logs.Warning(err)
				neko.InternalServerError(w, "Login failed.")
				return
			}

			token, err := auth.RedeemCode(code, auth.Endpoint{
				Token: config.EndpointMicrosoftIdentityToken,
			}, auth.Config{
				ClientID:    config.AzureADClientID,
				Scopes:      Scopes,
				RedirectURI: config.AzureADRedirectURI,
			})
			if err != nil {
				logs.Warning(err)
				neko.InternalServerError(w, "Get token failed.")
				return
			}

			fmt.Fprintln(w, "Login success.")
			ch <- token
		}),
	}

	go neko.StartServer(srv, false)
	defer neko.StopServer(srv)

	switch runtime.GOOS {
	case "darwin":
		err = exec.Command("open", auth.LoginURL(auth.Endpoint{
			Authorization: config.EndpointMicrosoftIdentityAuthorize,
		}, auth.Config{
			ClientID:    config.AzureADClientID,
			Scopes:      Scopes,
			RedirectURI: config.AzureADRedirectURI,
		})).Start()
	default:
		err = errors.New("unsupported platform")
	}
	if err != nil {
		return
	}

	token := <-ch
	config.Data.SetToken(token)
	return
}
