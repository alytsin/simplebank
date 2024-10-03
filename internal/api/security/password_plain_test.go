package security

import (
	"github.com/alytsin/simplebank/internal/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlainPasswordHash(t *testing.T) {
	pwd := util.RandomString(10)
	p := new(PasswordPlain)
	hash, err := p.Hash(pwd)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Equal(t, hash, pwd)
}

func TestPlainPasswordVerify(t *testing.T) {
	pwd := util.RandomString(10)
	p := new(PasswordPlain)
	hash, err := p.Hash(pwd)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Equal(t, hash, pwd)
	assert.True(t, p.Verify(hash, pwd))
	assert.False(t, p.Verify(util.RandomString(10), pwd))
}
