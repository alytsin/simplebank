package token

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateVerifyToken(t *testing.T) {

	m, err := NewPasetoMaker("")
	assert.Nil(t, err)
	assert.NotEmpty(t, m)

	payload := NewPayload("alexey")
	token, err := m.CreateToken(payload, time.Second)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	verifiedPayload, err := m.VerifyToken(token)
	assert.Nil(t, err)
	assert.Equal(t, payload, verifiedPayload)

	time.Sleep(time.Millisecond * 1001)
	verifiedPayload, err = m.VerifyToken(token)
	assert.Error(t, err)
	assert.Empty(t, verifiedPayload)

	p, err := m.VerifyToken("abc")
	assert.Error(t, err)
	assert.Empty(t, p)
}
