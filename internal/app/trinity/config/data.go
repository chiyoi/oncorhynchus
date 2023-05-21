package config

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/chiyoi/go/pkg/az/auth"
	"github.com/chiyoi/go/pkg/logs"
)

var Data data

type data struct {
	mu    sync.RWMutex
	Token auth.Token `json:"token"`
}

func (d *data) SetToken(token auth.Token) {
	defer d.Save()

	d.mu.Lock()
	defer d.mu.Unlock()
	d.Token = token
}

func (d *data) Load() {
	d.mu.Lock()
	defer d.mu.Unlock()

	f, err := os.Open(PathData)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		logs.Panic(err)
	}
	json.NewDecoder(f).Decode(d)
}

func (d *data) Save() {
	d.mu.RLock()
	defer d.mu.RUnlock()

	f, err := os.Create(PathData)
	if err != nil {
		logs.Panic(err)
	}
	json.NewEncoder(f).Encode(&d)
}
