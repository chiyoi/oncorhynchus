package data

import (
	"encoding/json"
	"fmt"
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

	logs.Info("loading data")
	f, err := os.Open(config.PathData)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		logs.Error(err)
		fmt.Println("Failed to load data.")
		logs.Fatal("exit")
	}
	json.NewDecoder(f).Decode(&Data)
	logs.Info("data loaded")
}

func Save() {
	mu.RLock()
	defer mu.RUnlock()

	logs.Info("saving data")
	f, err := os.Create(config.PathData)
	if err != nil {
		logs.Error(err)
		fmt.Println("Failed to save data.")
		logs.Fatal("exit")
	}
	json.NewEncoder(f).Encode(&Data)
	logs.Info("data saved")
}
