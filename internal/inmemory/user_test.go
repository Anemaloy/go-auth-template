package inmemory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUserStorage(t *testing.T) {
	storage := NewUserStorage()
	assert.Len(t, storage.users, 0)
	assert.Zero(t, storage.lastId)
}
