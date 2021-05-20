package model

import (
	"testing"

	"github.com/traPtitech/Jomon/ent"
)

func setup(t *testing.T) (*ent.Client, error) {
	client, err := SetupTestEntClient(t)
	if err != nil {
		return nil, err
	}
	return client, nil
}
