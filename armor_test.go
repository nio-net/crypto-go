package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleArmor(t *testing.T) {
	blockType := "MINT TEST"
	data := []byte("love")
	armorStr := EncodeArmor(blockType, nil, data)

	blockType2, _, data2, err := DecodeArmor(armorStr)
	require.Nil(t, err, "%+v", err)
	assert.Equal(t, blockType, blockType2)
	assert.Equal(t, data, data2)
}
