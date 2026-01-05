package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/NurulloMahmud/habits/internal/logs"
	"github.com/NurulloMahmud/habits/internal/platform/database"
	cx "github.com/NurulloMahmud/habits/pkg/context"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
	errMsg *string
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (m *Middleware) ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recorder := &responseRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		defer func() {
			if rec := recover(); rec != nil {
				errStr := "panic occurred"
				recorder.errMsg = &errStr

				recorder.WriteHeader(http.StatusInternalServerError)

				duration := time.Since(start).Milliseconds()
				writeLog(r, recorder, duration)

				panic(rec)
			}
		}()

		next.ServeHTTP(recorder, r)

		duration := time.Since(start).Milliseconds()
		writeLog(r, recorder, duration)
	})
}

func writeLog(r *http.Request, rec *responseRecorder, duration int64) {
	var errMsg *string

	if rec.status >= 500 {
		msg := http.StatusText(rec.status)
		errMsg = &msg
	}

	var userID int64
	user := cx.GetUser(r)
	if user.IsAnonymous() {
		userID = -1
	} else {
		userID = user.ID
	}

	log := logs.ActivityLog{
		User: logs.UserInfo{
			UserID: userID,
			IP:     r.RemoteAddr,
		},
		Method:     r.Method,
		Endpoint:   r.URL.Path,
		Status:     rec.status,
		DurationMS: duration,
		Error:      errMsg,
		CreatedAt:  time.Now().UTC(),
	}

	go func(l logs.ActivityLog) {
		_, _ = database.ActivityCollection().
			InsertOne(context.Background(), l)
	}(log)
}
