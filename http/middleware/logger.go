package middleware

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jabardigitalservice/golog/constant"

	gologlogger "github.com/jabardigitalservice/golog/logger"
)

func Logger(logger *gologlogger.Logger, data *gologlogger.LoggerData) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				ww        = NewWrapResponseWriter(w, r.ProtoMajor)
				ts        = time.Now().UTC()
				host      = r.Host
				uri       = r.RequestURI
				userAgent = r.UserAgent()
				ctx       = r.Context()
			)

			defer func() {
				var (
					remoteIP, _, err = net.SplitHostPort(r.RemoteAddr)
					scheme           = "http"
					method           = r.Method
					duration         = time.Since(ts)
					addr             = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
				)

				if err != nil {
					remoteIP = r.RemoteAddr
				}
				if r.TLS != nil {
					scheme = "https"
				}

				var (
					respStatus     = ww.Status()
					respStatusText = http.StatusText(respStatus)
				)

				if ctx.Value(constant.CtxUserIDKey) != nil {
					data.UserID = ctx.Value(constant.CtxUserIDKey).(string)
				}

				if ctx.Value(constant.CtxSessionIDKey) != nil {
					data.SessionID = ctx.Value(constant.CtxSessionIDKey).(string)
				}

				if ctx.Value(constant.CtxClientIDKey) != nil {
					data.ClientID = ctx.Value(constant.CtxClientIDKey).(string)
				}

				data.Category = gologlogger.LoggerRouter
				data.Duration = int64(duration)
				data.Method = fmt.Sprintf("[%s] %s", method, r.URL.Path)
				data.AdditionalInfo = map[string]interface{}{
					"http_host":         host,
					"http_uri":          uri,
					"http_proto":        r.Proto,
					"http_method":       method,
					"http_scheme":       scheme,
					"http_addr":         addr,
					"remote_addr":       remoteIP,
					"user_agent":        userAgent,
					"resp_elapsed_ms":   duration.String(),
					"resp_bytes_length": ww.BytesWritten(),
					"resp_status":       respStatus,
					"ts":                ts.Format(time.RFC3339),
					"resp_body":         ww.Body(),
				}

				if respStatus >= 200 && respStatus < 300 {
					logger.Info(data, respStatusText)
				} else {
					err := errors.New(respStatusText)
					logger.Error(data, err)
				}

			}()

			h.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
