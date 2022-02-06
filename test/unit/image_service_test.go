package unit

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"otus_go_final/internal/controllers"
	"otus_go_final/internal/services"
	"testing"
)

func TestImageService(t *testing.T) {
	props := services.NewImageProperty(300, 300, "data/snowshoe.jpg")

	headers := http.Header{}

	sut := services.NewProcessService(props, headers)

	// TODO Write UNIT tests for service logic
	t.Run("test service success", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "localhost:4000/fill/600/600/assets.imgix.net/hp/snowshoe.jpg?auto=compress&w=900&h=600&fit=crop", nil)
		w := httptest.NewRecorder()

		controllers.Index(w, req)
		res := w.Result()

		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)

		require.NoError(t, err)

		fmt.Println(data, res.Body)
	})
}
