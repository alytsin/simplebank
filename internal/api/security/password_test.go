package security

import (
	"github.com/alytsin/simplebank/internal/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordHash(t *testing.T) {
	p := Password(util.RandomString(10))
	hash, err := p.Hash()
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	p = Password(util.RandomString(100))
	hash, err = p.Hash()
	assert.Error(t, err)
	assert.Empty(t, hash)
}

func TestPasswordVerify(t *testing.T) {
	p := Password(util.RandomString(10))
	hash, err := p.Hash()
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.True(t, p.Verify(hash))
	assert.False(t, p.Verify(util.RandomString(10)))
}
