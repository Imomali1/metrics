package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

type Handler struct {
	Storage *MemStorage
}

func ValidateURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		contentType := r.Header.Get("Content-Type")
		if contentType != "text/plain" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		path := r.URL.RequestURI()
		path = strings.Trim(path, "/")
		params := strings.Split(path, "/")
		if len(params) != 4 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(strings.TrimLeft(r.URL.RequestURI(), "/"), "/")
	// params => [update, <ТИП_МЕТРИКИ>, <ИМЯ_МЕТРИКИ>, <ЗНАЧЕНИЕ_МЕТРИКИ>]
	metricType, metricName, metricStrValue := params[1], params[2], params[3]
	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch metricType {
	case "gauge":
		gauge, err := strconv.ParseFloat(metricStrValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.Storage.Gauge[metricName] = gauge
	case "counter":
		counter, err := strconv.ParseInt(metricStrValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.Storage.Counter[metricName] += counter
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	m := NewMemStorage()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//if request.Method != http.MethodGet || request.URL.Path != "/" {
		//	http.NotFound(writer, request)
		//	return
		//}
		writer.Write([]byte(`Hello World! This is homepage!`))
	})
	h := Handler{
		Storage: m,
	}
	mux.Handle("/update/", ValidateURL(http.HandlerFunc(h.Update)))
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server started...")
}
