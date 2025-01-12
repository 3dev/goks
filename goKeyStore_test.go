package goKeyStore

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestOpen(t *testing.T) {

	ks, err := Open("nonexistenfile.goKey", "fakepasskey")
	assert.Nil(t, ks)
	assert.Error(t, err)
}

func TestCreate(t *testing.T) {

	ks, err := New("test1.goKey", "test1passkey")
	assert.NotNil(t, ks)
	assert.NoError(t, err)

	ks, err = Open("test1.goKey", "test1passkey")
	assert.NotNil(t, ks)
	assert.NoError(t, err)
}
