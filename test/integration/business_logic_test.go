package integration

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"otus_go_final/config"
	"otus_go_final/internal/controllers"
	"otus_go_final/internal/services"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
)

func TestBusinessLogic(t *testing.T) {
	type errorTestCases struct {
		description string
		input       struct {
			width  string
			height string
			URL    string
		}
		expectedError string
		expectedCode  int
	}

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

	for _, scenario := range []errorTestCases{
		{
			description: "test image is smaller than requested size",
			input: struct {
				width  string
				height string
				URL    string
			}{width: "600", height: "600", URL: "localhost:4000/fill/600/600/localhost:8000/coffee.jpg"},
			expectedError: services.ErrImageSizeViolation.Error(),
			expectedCode:  400,
		},
		{
			description: "test image not found",
			input: struct {
				width  string
				height string
				URL    string
			}{width: "200", height: "400", URL: "localhost:4000/fill/200/400/localhost:8000/no.jpg"},
			expectedError: services.ErrNotFound.Error(),
			expectedCode:  404,
		},
		{
			description: "test server does not exist",
			input: struct {
				width  string
				height string
				URL    string
			}{width: "200", height: "400", URL: "localhost:4000/fill/200/400/www.notexists.su"},
			expectedError: services.ErrServerNotExists.Error(),
			expectedCode:  500,
		},
		{
			description: "test target file is not an image but pdf",
			input: struct {
				width  string
				height string
				URL    string
			}{width: "200", height: "400", URL: "localhost:4000/fill/200/400/localhost:8000/sample.pdf"},
			expectedError: services.ErrNotImageType.Error(),
			expectedCode:  400,
		},
		{
			description: "test target file is not an image but pdf",
			input: struct {
				width  string
				height string
				URL    string
			}{width: "200", height: "400", URL: "localhost:4000/fill/200/400/localhost:8000/sample.pdf"},
			expectedError: services.ErrNotImageType.Error(),
			expectedCode:  400,
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, scenario.input.URL, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("width", scenario.input.width)
			rctx.URLParams.Add("height", scenario.input.height)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			h.Index(w, req)

			res := w.Result()

			defer res.Body.Close()

			require.Equal(t, scenario.expectedCode, res.StatusCode)
			require.Equal(t, "text/html; charset=UTF-8", w.Header().Get("Content-Type"))

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, scenario.expectedError, string(data))
		})
	}

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
}
