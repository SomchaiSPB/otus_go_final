package integration

import (
	"context"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"otus_go_final/internal/services"
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
	h := controllers.NewBaseHandler(&cfg)

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
		URL := "localhost:4000/fill/200/400/localhost:8000/snowshoe.jpg"

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, URL, nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("width", "200")
		rctx.URLParams.Add("height", "400")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.Index(w, req)

		res := w.Result()

		defer res.Body.Close()

		require.Equal(t, 200, res.StatusCode)
		require.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		require.NotEmpty(t, data)
	})

	t.Run("test target file is not an image but pdf", func(t *testing.T) {
		URL := "localhost:4000/fill/200/400/www.africau.edu/images/default/sample.pdf"

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, URL, nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("width", "200")
		rctx.URLParams.Add("height", "400")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.Index(w, req)

		res := w.Result()

		defer res.Body.Close()

		require.Equal(t, 400, res.StatusCode)
		require.Equal(t, "text/html; charset=UTF-8", w.Header().Get("Content-Type"))

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		require.Equal(t, services.ErrNotImageType.Error(), string(data))
	})

	t.Run("test server does not exist", func(t *testing.T) {
		URL := "localhost:4000/fill/200/400/www.notexists.su"

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, URL, nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("width", "200")
		rctx.URLParams.Add("height", "400")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.Index(w, req)

		res := w.Result()

		defer res.Body.Close()

		require.Equal(t, 500, res.StatusCode)
		require.Equal(t, "text/html; charset=UTF-8", w.Header().Get("Content-Type"))

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		require.Equal(t, services.ErrServerNotExists.Error(), string(data))
	})

	t.Run("test image not found", func(t *testing.T) {
		URL := "localhost:4000/fill/200/400/localhost:8000/no.jpg"

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, URL, nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("width", "200")
		rctx.URLParams.Add("height", "400")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.Index(w, req)

		res := w.Result()

		defer res.Body.Close()

		require.Equal(t, 404, res.StatusCode)
		require.Equal(t, "text/html; charset=UTF-8", w.Header().Get("Content-Type"))

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		require.Equal(t, services.ErrNotFound.Error(), string(data))
	})

	t.Run("test image is smaller than requested size", func(t *testing.T) {
		URL := "localhost:4000/fill/600/600/localhost:8000/coffee.jpg"

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, URL, nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("width", "600")
		rctx.URLParams.Add("height", "600")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.Index(w, req)

		res := w.Result()

		defer res.Body.Close()

		require.Equal(t, 400, res.StatusCode)
		require.Equal(t, "text/html; charset=UTF-8", w.Header().Get("Content-Type"))

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		require.Equal(t, services.ErrImageSizeViolation.Error(), string(data))
	})
}
