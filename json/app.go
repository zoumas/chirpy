package main

import (
	"fmt"
	"net/http"
)

type App struct {
	env            *Env
	fileserverHits int
}

type StatusCaptureResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *StatusCaptureResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (app *App) IncrementMetrics(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCaptureWriter := &StatusCaptureResponseWriter{ResponseWriter: w}
		handler.ServeHTTP(statusCaptureWriter, r)

		if statusCaptureWriter.statusCode != http.StatusMovedPermanently {
			app.fileserverHits++
		}
	})
}

func (app *App) ReportMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf("Hits: %d", app.fileserverHits)
	w.Write([]byte(body))
}

func (app *App) ResetMetrics(w http.ResponseWriter, _ *http.Request) {
	app.fileserverHits = 0

	w.WriteHeader(http.StatusOK)
}
