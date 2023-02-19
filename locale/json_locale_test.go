package locale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFromDir(t *testing.T) {
	l, err := NewFromDir("testdata")
	assert.NoError(t, err)
	assert.NotNil(t, l)
	assert.Equal(t, "Selam", l.Message(Turkish, "sample_message"))
	assert.Equal(t, "Selam Ali", l.Message(Turkish, "formatted_message", "Ali"))
}
