package config

import (
	"fmt"
	"testing"

	"github.com/chiyoi/go/pkg/az/auth"
)

func TestLoadData(t *testing.T) {
	Data.Load()
	fmt.Println(&Data)
}

func TestSaveData(t *testing.T) {
	Data.Save()
}

func TestSetToken(t *testing.T) {
	Data.SetToken(auth.Token{})
}
