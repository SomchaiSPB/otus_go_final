package unit

import (
	"otus_go_final/internal/services"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImageService(t *testing.T) {
	const (
		validURL   = "https://media.istockphoto.com/photos/asian-woman-holding-covid-rapid-test-and-waiting-for-results-picture-id1345296778" //nolint
		inValidURL = "http//google.com"
	)

	resizer := new(services.JpegResizer)

	t.Run("test validate no error", func(t *testing.T) {
		props := services.NewImageProperty(300, 300, validURL, nil)

		sut := services.NewProcessService(props, resizer)

		err := sut.Validate()

		require.Equal(t, 200, sut.ResponseCode)

		require.NoError(t, err)
	})

	t.Run("test validate returns error", func(t *testing.T) {
		props := services.NewImageProperty(300, 300, inValidURL, nil)

		sut := services.NewProcessService(props, resizer)

		err := sut.Validate()

		require.Error(t, err)
	})

	t.Run("test proxy call func", func(t *testing.T) {
		props := services.NewImageProperty(300, 300, validURL, nil)

		sut := services.NewProcessService(props, resizer)

		res, err := sut.ProxyRequest()

		require.NoError(t, err)
		require.NotEmpty(t, res)
	})

	t.Run("test service invoke", func(t *testing.T) {
		props := services.NewImageProperty(300, 300, validURL, nil)

		sut := services.NewProcessService(props, resizer)

		res, err := sut.ProxyRequest()

		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
