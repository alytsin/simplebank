package security

import (
	"github.com/alytsin/simplebank/internal/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordHash(t *testing.T) {
	p := new(Password)
	hash, err := p.Hash(util.RandomString(10))
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	p = new(Password)
	hash, err = p.Hash(util.RandomString(100))
	assert.Error(t, err)
	assert.Empty(t, hash)
}

func TestPasswordVerify(t *testing.T) {
	pwd := util.RandomString(10)
	p := new(Password)
	hash, err := p.Hash(pwd)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.True(t, p.Verify(hash, pwd))
	assert.False(t, p.Verify(util.RandomString(10), pwd))
}
