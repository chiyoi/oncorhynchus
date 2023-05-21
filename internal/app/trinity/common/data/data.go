package data

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/chiyoi/go/pkg/az/authentication"
	"github.com/chiyoi/go/pkg/logs"
	"github.com/chiyoi/oncorhynchus/internal/app/trinity/config"
)

var mu sync.RWMutex
var Data struct {
	Token authentication.Token `json:"token"`
}

func SetToken(token authentication.Token) {
	defer Save()

	mu.Lock()
	defer mu.Unlock()
	Data.Token = token
}

func Load() {
	mu.Lock()
	defer mu.Unlock()

	f, err := os.Open(config.PathData)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		logs.Panic(err)
	}
	json.NewDecoder(f).Decode(&Data)
}

func Save() {
	mu.RLock()
	defer mu.RUnlock()

	f, err := os.Create(config.PathData)
	if err != nil {
		logs.Panic(err)
	}
	json.NewEncoder(f).Encode(&Data)
}
