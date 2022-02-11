package integration

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"otus_go_final/config"
	"otus_go_final/internal/controllers"
)

func TestBusinessLogic(t *testing.T) {
	cfg := config.Config{
		Port:     "4000",
		Capacity: 10,
	}

	t.Run("test health check ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/check", nil)
		w := httptest.NewRecorder()
		controllers.Check(w, req)
		res := w.Result()

		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		require.Equal(t, "ok", string(data))
	})

	t.Run("test image processing success", func(t *testing.T) {
		h := controllers.NewBaseHandler(&cfg)

		URL := "localhost:4000/fill/200/400/nginx/snowshoe.jpg"

		req := httptest.NewRequest(http.MethodGet, URL, nil)

		w := httptest.NewRecorder()

		h.Index(w, req)

		res := w.Result()

		defer res.Body.Close()

		require.Equal(t, 200, res.StatusCode)
		require.Equal(t, "application/octet-stream", res.Header.Get("Content-Type"))

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		require.NotEmpty(t, data)
	})

	t.Run("test target file is not an image but pdf", func(t *testing.T) {
		h := controllers.NewBaseHandler(&cfg)

		URL := "localhost:4000/fill/200/400/www.africau.edu/images/default/sample.pdf"

		req := httptest.NewRequest(http.MethodGet, URL, nil)

		w := httptest.NewRecorder()

		h.Index(w, req)

		res := w.Result()

		defer res.Body.Close()

		require.Equal(t, 200, res.StatusCode)
		require.Equal(t, "text/html; charset=UTF-8", res.Header.Get("Content-Type"))

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		require.Equal(t, "target file is not an image", string(data))
	})
}
