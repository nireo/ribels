package league_api

import "net/http"

type Error struct {
	Message    string
	StatusCode int
}

func (e *Error) Error() string {
	return e.Message
}

var (
	ErrBadRequest = Error{
		Message:    "Bad request",
		StatusCode: http.StatusBadRequest,
	}
	ErrUnauthorized = Error{
		Message:    "Unauthorized request",
		StatusCode: http.StatusUnauthorized,
	}
	ErrForbidden = Error{
		Message:    "Forbidden request",
		StatusCode: http.StatusForbidden,
	}
	ErrNotFound = Error{
		Message:    "Not found",
		StatusCode: http.StatusNotFound,
	}
	ErrRateLimitExceeded = Error{
		Message:    "rate limit exceeded",
		StatusCode: http.StatusTooManyRequests,
	}
	ErrInternalServerError = Error{
		Message:    "internal server error",
		StatusCode: http.StatusInternalServerError,
	}
	ErrBadGateway = Error{
		Message:    "bad gateway",
		StatusCode: http.StatusBadGateway,
	}
	ErrServiceUnavailable = Error{
		Message:    "service unavailable",
		StatusCode: http.StatusServiceUnavailable,
	}
	ErrGatewayTimeout = Error{
		Message:    "gateway timeout",
		StatusCode: http.StatusGatewayTimeout,
	}
	StatusToError = map[int]Error{
		http.StatusBadRequest:          ErrBadRequest,
		http.StatusUnauthorized:        ErrUnauthorized,
		http.StatusForbidden:           ErrForbidden,
		http.StatusNotFound:            ErrNotFound,
		http.StatusTooManyRequests:     ErrRateLimitExceeded,
		http.StatusInternalServerError: ErrInternalServerError,
		http.StatusBadGateway:          ErrBadGateway,
		http.StatusServiceUnavailable:  ErrServiceUnavailable,
		http.StatusGatewayTimeout:      ErrGatewayTimeout,
	}
)
