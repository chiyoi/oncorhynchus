package auth

import (
	"fmt"

	"github.com/chiyoi/go/pkg/az/authentication"
	"github.com/chiyoi/go/pkg/logs"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/config"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/data"
)

func Login() (err error) {
	logs.Info("login")
	token, err := authentication.Login(authentication.Endpoint{
		Token:     config.EndpointMicrosoftIdentityToken,
		Authorize: config.EndpointMicrosoftIdentityAuthorize,
	}, authentication.Config{
		ClientID:    config.AzureADClientID,
		Scopes:      config.AzureADLoginScopes,
		RedirectURI: config.AzureADLoginRedirectURI,
	})
	if err != nil {
		return
	}

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
