package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/chiyoi/go/pkg/az/authentication"
	"github.com/chiyoi/go/pkg/logs"
	"github.com/chiyoi/go/pkg/neko"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/config"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/data"
)

func Login() (err error) {
	logs.Info("login")
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
	logs.Info("login succeeded")
	return
}

func Refresh(token authentication.Token) authentication.Token {
	logs.Info("refreshing token")
	token, err := authentication.Refresh(token.RefreshToken, authentication.Endpoint{
		Token: config.EndpointMicrosoftIdentityToken,
	}, authentication.Config{
		ClientID: config.AzureADClientID,
	})
	if err != nil {
		logs.Error(err)
		fmt.Println("Failed to authenticate.")
		logs.Fatal("exit")
	}
	logs.Info("token refreshed")
	return token
}
