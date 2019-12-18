package goalisa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	a := assert.New(t)
	var dr Driver
	_, err := dr.Open("alisa://endpoint?a=1")
	a.NoError(err)
}
