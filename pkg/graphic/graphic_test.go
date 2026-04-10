package graphic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSimpleImageURL(t *testing.T) {
	img, err := GenerateSimpleImageURL(Red)
	assert.NoError(t, err)
	assert.NotEqual(t, img, "data:image/png;base64,")
}
