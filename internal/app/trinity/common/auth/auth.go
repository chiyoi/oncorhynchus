package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"github.com/chiyoi/go/pkg/az/authentication"
	"github.com/chiyoi/go/pkg/logs"
	"github.com/chiyoi/go/pkg/neko"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/common/data"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/config"
)

func Login() (err error) {
	u, err := url.Parse(config.AzureADLoginRedirectURI)
	if err != nil {
		return
	}

	ch := make(chan authentication.Token)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", u.Port()),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code, err := authentication.GetCode(r)
			if err != nil {
				logs.Warning(err)
				neko.InternalServerError(w, "Login failed.")
				return
			}

			token, err := authentication.RedeemCode(code, authentication.Endpoint{
				Token: config.EndpointMicrosoftIdentityToken,
			}, authentication.Config{
				ClientID:    config.AzureADClientID,
				Scopes:      config.AzureADLoginScopes,
				RedirectURI: config.AzureADLoginRedirectURI,
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
		err = exec.Command("open", authentication.LoginURL(authentication.Endpoint{
			Authorization: config.EndpointMicrosoftIdentityAuthorize,
		}, authentication.Config{
			ClientID:    config.AzureADClientID,
			Scopes:      config.AzureADLoginScopes,
			RedirectURI: config.AzureADLoginRedirectURI,
		})).Start()
	default:
		err = errors.New("unsupported platform")
	}
	if err != nil {
		return
	}

	token := <-ch
	data.SetToken(token)
	return
}

func Refresh(token authentication.Token) authentication.Token {
	token, err := authentication.Refresh(token.RefreshToken, authentication.Endpoint{
		Token: config.EndpointMicrosoftIdentityToken,
	}, authentication.Config{
		ClientID: config.AzureADClientID,
	})
	if err != nil {
		fmt.Println("error while refreshing token:", err)
		os.Exit(1)
	}
	return token
}
