package unit

import (
	"bytes"
	"image"
	"io/ioutil"
	"log"
	"otus_go_final/internal/services"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImageResizer(t *testing.T) {
	props := services.NewImageProperty(50, 50, "", nil)

	t.Run("test image resize success", func(t *testing.T) {
		file, err := ioutil.ReadFile("../static/snowshoe.jpg")

		require.NoError(t, err)

		m, _, err := image.Decode(bytes.NewReader(file))

		require.NoError(t, err)

		sut := services.NewImageResizer(props)

		res, err := sut.Resize(m)

		require.NoError(t, err)

		resizedInfo, _, err := image.DecodeConfig(bytes.NewReader(res))
		if err != nil {
			log.Println(err)
		}

		require.Equal(t, 50, resizedInfo.Width)
		require.Equal(t, 50, resizedInfo.Height)
	})
}
