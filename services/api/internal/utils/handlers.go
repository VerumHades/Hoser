package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

// loggingResponseWriter captures outgoing status and body
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           bytes.NewBuffer(nil),
	}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

// ColorLogMiddleware logs both incoming and outgoing HTTP traffic with color
func ColorLogMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)

		// 🟢 Log incoming request
		fmt.Printf("%s-->%s %s%s%s %s\n", colorCyan, colorReset, colorBlue, r.Method, colorReset, r.URL.Path)
		for name, values := range r.Header {
			for _, v := range values {
				fmt.Printf("   %s%s%s: %s%s%s\n", colorCyan, name, colorReset, colorGray, v, colorReset)
			}
		}

		// Process request
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)

		// 🎨 Color by status
		statusColor := colorGreen
		switch {
		case lrw.statusCode >= 500:
			statusColor = colorRed
		case lrw.statusCode >= 400:
			statusColor = colorYellow
		case lrw.statusCode >= 300:
			statusColor = colorCyan
		}

		// 🟣 Log outgoing response
		fmt.Printf("%s<--%s %s%s%s %s%d%s (%v)\n",
			colorCyan, colorReset,
			colorBlue, r.Method, colorReset,
			statusColor, lrw.statusCode, colorReset,
			duration,
		)

		for name, values := range lrw.Header() {
			for _, v := range values {
				fmt.Printf("   %s%s%s: %s%s%s\n", colorCyan, name, colorReset, colorGray, v, colorReset)
			}
		}

		body := lrw.body.String()
		if len(body) > 0 {
			max := 500
			if len(body) > max {
				body = body[:max] + "... (truncated)"
			}
			fmt.Printf("   %sResponse:%s %s%s%s\n", colorCyan, colorReset, colorGray, body, colorReset)
		}
	})
}
