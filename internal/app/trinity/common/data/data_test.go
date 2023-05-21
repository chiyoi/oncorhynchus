package data

import (
	"fmt"
	"testing"

	"github.com/chiyoi/go/pkg/az/authentication"
)

func TestLoad(t *testing.T) {
	Load()
	fmt.Println(&Data)
}

func TestSave(t *testing.T) {
	Save()
}

func TestSetToken(t *testing.T) {
	SetToken(authentication.Token{})
}
