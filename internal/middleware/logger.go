package middleware

import (
	"bytes"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var sensitiveFields = map[string]bool{
	"password":      true,
	"token":         true,
	"authorization": true,
	"api_key":       true,
	"secret":        true,
	"access_token":  true,
	"refresh_token": true,
}

type responseWriter struct {
	gin.ResponseWriter
	body      *bytes.Buffer
	limit     int
	truncated bool
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.capture(data)
	return w.ResponseWriter.Write(data)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.capture([]byte(s))
	return w.ResponseWriter.WriteString(s)
}

func (w *responseWriter) LoggedBody() string {
	if !w.truncated {
		return w.body.String()
	}

	return w.body.String() + "... [TRUNCATED]"
}

func (w *responseWriter) capture(data []byte) {
	if w.body.Len() >= w.limit {
		w.truncated = true
		return
	}

	remaining := w.limit - w.body.Len()
	if len(data) > remaining {
		w.body.Write(data[:remaining])
		w.truncated = true
		return
	}

	w.body.Write(data)
}

const (
	maxRequestBodySize     = 1 << 20
	maxResponseBodyLogSize = 64 << 10
)

type fileInfo struct {
	Field       string `json:"field"`
	FileName    string `json:"file_name"`
	Size        string `json:"size"`
	ContentType string `json:"content_type"`
}

func parseRequestBody(c *gin.Context) (body any, files []fileInfo) {
	contentType := c.GetHeader("Content-Type")

	switch {
	case strings.HasPrefix(contentType, "multipart/form-data"):
		formBody := make(map[string]any)
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil || c.Request.MultipartForm == nil {
			return
		}
		for k, v := range c.Request.MultipartForm.Value {
			if len(v) == 1 {
				formBody[k] = v[0]
			} else {
				formBody[k] = v
			}
		}
		for field, fs := range c.Request.MultipartForm.File {
			for _, f := range fs {
				files = append(files, fileInfo{
					Field:       field,
					FileName:    f.Filename,
					Size:        formatSize(f.Size),
					ContentType: f.Header.Get("Content-Type"),
				})
			}
		}
		body = formBody

	case strings.HasPrefix(contentType, "application/json"):
		raw, err := readBody(c)
		if err != nil {
			return
		}
		body = parseLoggedPayload(raw)

	case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
		formBody := make(map[string]any)
		raw, err := readBody(c)
		if err != nil {
			return
		}
		values, _ := url.ParseQuery(string(raw))
		for k, v := range values {
			if len(v) == 1 {
				formBody[k] = v[0]
			} else {
				formBody[k] = v
			}
		}
		body = formBody

	default:
		raw, err := readBody(c)
		if err != nil {
			return
		}
		if len(raw) > 0 {
			body = parseLoggedPayload(raw)
		}
	}

	return
}
func readBody(c *gin.Context) ([]byte, error) {
	raw, err := io.ReadAll(io.LimitReader(c.Request.Body, int64(maxRequestBodySize)))
	if err != nil {
		return nil, err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))
	return raw, nil
}

func parseResponseBody(contentType, raw string) any {
	switch {
	case strings.HasPrefix(contentType, "image/"),
		strings.HasPrefix(contentType, "video/"),
		strings.HasPrefix(contentType, "audio/"),
		strings.HasPrefix(contentType, "application/octet-stream"):
		return "[BINARY DATA]"

	case strings.HasPrefix(contentType, "application/json") ||
		strings.HasPrefix(strings.TrimSpace(raw), "{") ||
		strings.HasPrefix(strings.TrimSpace(raw), "["):
		return parseLoggedPayload([]byte(raw))

	default:
		return compactText(raw)
	}
}

func Logger(httpLogger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		reqBody, files := parseRequestBody(c)

		rw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
			limit:          maxResponseBodyLogSize,
		}
		c.Writer = rw

		c.Next()

		status := c.Writer.Status()

		respRaw := rw.LoggedBody()
		respContentType := c.Writer.Header().Get("Content-Type")
		respBody := parseResponseBody(respContentType, respRaw)

		var logEvent *zerolog.Event
		switch {
		case status >= 500:
			logEvent = httpLogger.Error()
		case status >= 400:
			logEvent = httpLogger.Warn()
		default:
			logEvent = httpLogger.Info()
		}

		e := logEvent.
			Str("trace_id", GetTraceID(c.Request.Context())).
			Str("request_id", GetRequestID(c)).
			// Request
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("query", c.Request.URL.RawQuery).
			Str("protocol", c.Request.Proto).
			Str("host", c.Request.Host).
			Str("client_ip", c.ClientIP()).
			Str("remote_addr", c.Request.RemoteAddr).
			Str("user_agent", c.Request.UserAgent()).
			Str("referer", c.Request.Referer()).
			Str("content_type", c.GetHeader("Content-Type")).
			Int64("content_length", c.Request.ContentLength).
			Interface("request_body", sanitizeValue(reqBody, sensitiveFields)).
			// Response
			Int("status", status).
			Int("bytes_out", c.Writer.Size()).
			Str("response_content_type", respContentType).
			Interface("response_body", respBody).
			Int64("latency_ms", time.Since(start).Milliseconds())

		if len(files) > 0 {
			e = e.Interface("files", files)
		}

		e.Msg("HTTP")
	}
}
