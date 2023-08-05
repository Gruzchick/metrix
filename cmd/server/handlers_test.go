package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricsHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Тест допустимого значения счётчика counter",
			url:  "http://localhost:8080/update/counter/test-counter/1.2",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Тест допустимого значения типа счётчика",
			url:  "http://localhost:8080/update/cou-nter/test-counter/1",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.url, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			updateMetricsByJSONHandler(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, res.StatusCode, test.want.code)
			// получаем и проверяем тело запроса
			err := res.Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		})
	}
}
