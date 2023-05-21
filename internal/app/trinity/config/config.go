package config

import (
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/chiyoi/go/pkg/logs"
)

const (
	EndpointMicrosoftIdentityAuthorize = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
	EndpointMicrosoftIdentityToken     = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
)

var (
	AzureADClientID    = "e5a68652-2fed-4508-ad85-02e7a966307f"
	AzureADRedirectURI = "http://localhost:10123"
)

var (
	Timeout = time.Second * 20

	EndpointNeko03 = "https://api.neko03.moe/"

	PathData = func() string {
		u, err := user.Current()
		if err != nil {
			logs.Panic(err)
		}
		return filepath.Join(u.HomeDir, ".oncorhynchus", "trinity.json")
	}()
)

func init() {
	if err := os.MkdirAll(filepath.Dir(PathData), 0700); err != nil {
		logs.Panic(err)
	}
}
