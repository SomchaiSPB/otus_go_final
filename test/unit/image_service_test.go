package unit

import (
	"net/http"
	"otus_go_final/internal/services"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImageService(t *testing.T) {
	const (
		validURL   = "https://media.istockphoto.com/photos/asian-woman-holding-covid-rapid-test-and-waiting-for-results-picture-id1345296778"
		inValidURL = "http//google.com"
	)

	t.Run("test validate no error", func(t *testing.T) {
		props := services.NewImageProperty(300, 300, validURL)

		headers := http.Header{}

		sut := services.NewProcessService(props, headers)

		code, err := sut.Validate()

		require.Equal(t, 200, code)

		require.NoError(t, err)
	})

	t.Run("test validate returns error", func(t *testing.T) {
		props := services.NewImageProperty(300, 300, inValidURL)

		headers := http.Header{}

		sut := services.NewProcessService(props, headers)

		code, err := sut.Validate()

		require.Nil(t, 200, code)
		require.Error(t, err)
	})

	t.Run("test proxy call func", func(t *testing.T) {
		props := services.NewImageProperty(300, 300, validURL)

		headers := http.Header{}

		sut := services.NewProcessService(props, headers)

		res, err := sut.ProxyRequest()

		require.NoError(t, err)
		require.NotEmpty(t, res)
	})

	t.Run("test service invoke", func(t *testing.T) {
		props := services.NewImageProperty(300, 300, validURL)

		headers := http.Header{}

		sut := services.NewProcessService(props, headers)

		res, err := sut.ProxyRequest()

		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
