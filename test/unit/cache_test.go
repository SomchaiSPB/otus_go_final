package unit

import (
	"github.com/stretchr/testify/require"
	imagecache "otus_go_final/internal/cache"
	"testing"
)

func TestCache(t *testing.T) {
	t.Run("test add/get/clear cache", func(t *testing.T) {
		var ok bool
		sut := imagecache.NewCache(3)

		ok = sut.Set("1", "first value")
		require.False(t, ok)
		ok = sut.Set("2", "second value")
		require.False(t, ok)
		ok = sut.Set("3", "third value")
		require.False(t, ok)
		ok = sut.Set("4", "fourth value")
		require.False(t, ok)

		res, ok := sut.Get("1")
		require.False(t, ok)
		require.Nil(t, res)

		res, ok = sut.Get("2")
		require.True(t, ok)
		require.Equal(t, "second value", res)

		res, ok = sut.Get("3")
		require.True(t, ok)
		require.Equal(t, "third value", res)

		res, ok = sut.Get("4")
		require.True(t, ok)
		require.Equal(t, "fourth value", res)

		sut.Clear()

		res, ok = sut.Get("4")
		require.False(t, ok)
		require.Nil(t, res)
	})
}
