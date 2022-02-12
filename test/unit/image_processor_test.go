package unit

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"image"
	"io/ioutil"
	"log"
	"otus_go_final/internal/services"
	"testing"
)

func TestImageResizer(t *testing.T) {
	props := services.NewImageProperty(50, 50, "")

	t.Run("test image resize success", func(t *testing.T) {
		file, err := ioutil.ReadFile("../static/snowshoe.jpg")

		require.NoError(t, err)

		m, format, err := image.Decode(bytes.NewReader(file))

		require.NoError(t, err)

		sut := services.NewImageProcessor(format, m, props)

		res, err := sut.Resize()

		require.NoError(t, err)

		resizedInfo, _, err := image.DecodeConfig(bytes.NewReader(res))

		if err != nil {
			log.Println(err)
		}

		require.Equal(t, 50, resizedInfo.Width)
		require.Equal(t, 50, resizedInfo.Height)
	})
}
