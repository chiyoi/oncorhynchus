package oncorhynchus

import (
	"net/url"
	"os"
)

type defineConfig struct {
	AuthenticatorEndpoint *url.URL

	PublicKey string
}

var Config = defineConfig{
	AuthenticatorEndpoint: &url.URL{
		Scheme: "http",
		Host:   "localhost:80",
	},
}

func init() {
	Config.PublicKey = os.Getenv("PUBLIC_KEY")
}
