package gwent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadCardsCSV(t *testing.T) {
	cards := loadcards()
	assert.True(t, len(cards) > 0)
}
