package model

import "testing"

func TestEstablishConnection(t *testing.T) {
	if _, err := EstablishConnection(); err != nil {
		t.Errorf(err.Error())
	}
}
