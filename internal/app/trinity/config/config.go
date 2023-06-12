package config

import (
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/chiyoi/go/pkg/logs"
	"github.com/chiyoi/go/pkg/sakana"
)

const (
	EndpointMicrosoftIdentityAuthorize = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
	EndpointMicrosoftIdentityToken     = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
)

var (
	AzureADClientID         = "e5a68652-2fed-4508-ad85-02e7a966307f"
	AzureADLoginRedirectURI = "http://localhost:10123"
	AzureADLoginScopes      = []string{"User.Read", "offline_access"}
)

var (
	Timeout = time.Second * 20

	EndpointNeko03 = "https://neko03.redriver-a01bf37e.japaneast.azurecontainerapps.io/"
	QueryKeyToken  = "token"

	DirData = func() string {
		u, err := user.Current()
		if err != nil {
			logs.Panic(err)
		}
		return filepath.Join(u.HomeDir, ".oncorhynchus", "trinity")
	}()
	PathData = filepath.Join(DirData, "data.json")
)

func init() {
	if err := os.MkdirAll(filepath.Dir(PathData), 0700); err != nil {
		logs.Panic(err)
	}
}

func EndpointAuthedNeko03(token string) (u *url.URL) {
	u, err := url.Parse(EndpointNeko03)
	if err != nil {
		logs.Error(err)
		sakana.InternalError(os.Stderr)
	}

	q := u.Query()
	q.Set(QueryKeyToken, token)
	u.RawQuery = q.Encode()
	return
}
