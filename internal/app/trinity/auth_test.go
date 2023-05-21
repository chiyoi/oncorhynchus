package trinity

import "testing"

func TestLogin(t *testing.T) {
	if err := Login(); err != nil {
		t.Fatal(err)
	}
}
