package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKana(t *testing.T) {
	assert.True(t, IsSmall('ォ'))
	assert.True(t, IsSmall('ぁ'))
	assert.False(t, IsSmall('ア'))
	assert.False(t, IsSmall('え'))
	assert.Equal(t, MapSmallToBig('ォ'), 'オ')
	assert.Equal(t, MapSmallToBig('ゃ'), 'や')
	assert.Equal(t, MapSmallToBig('ぃ'), 'い')
	assert.Equal(t, GetFirstKana("へんたい"), 'へ')
	assert.Equal(t, GetFirstKana("キス"), 'キ')
	assert.Equal(t, GetFirstKana("ラ"), 'ラ')
	assert.Equal(t, GetFirstKana("ー"), 'ー')
	assert.Equal(t, GetLastKana("へんたい"), 'い')
	assert.Equal(t, GetLastKana("キス"), 'ス')
	assert.Equal(t, GetLastKana("ラ"), 'ラ')
	assert.Equal(t, GetLastKana("スキー"), 'ー')
	assert.NotEqual(t, GetLastKana("しゅしょ"), 'ょ')
	assert.Equal(t, GetLastKana("しゅしょ"), 'よ')
	assert.Equal(t, GetFirstKana("ラ"), GetLastKana("ラ"))
}