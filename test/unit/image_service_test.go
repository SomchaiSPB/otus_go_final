package unit

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"otus_go_final/config"
	"otus_go_final/internal/controllers"
	"otus_go_final/internal/services"
	"testing"
)

func TestImageService(t *testing.T) {
	cfg := &config.Config{
		Port:     "4000",
		Capacity: 10,
	}
	props := services.NewImageProperty(300, 300, "data/snowshoe.jpg")

	headers := http.Header{}

	sut := services.NewProcessService(props, headers)

	_ = sut

	// TODO Write UNIT tests for service logic
	t.Run("test service success", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "http://localhost:4000/fill/600/600/localhost:8000/snowshoe.jpg", nil)
		w := httptest.NewRecorder()

		handler := controllers.NewBaseHandler(cfg)

		handler.Index(w, req)
		res := w.Result()

		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)

		require.NoError(t, err)

		fmt.Println(data, res.Body)
	})
}
