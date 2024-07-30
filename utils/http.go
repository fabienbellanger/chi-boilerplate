package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// Error status codes
const (
	StatusBadRequest                   = 400
	StatusUnauthorized                 = 401
	StatusPaymentRequired              = 402
	StatusForbidden                    = 403
	StatusNotFound                     = 404
	StatusMethodNotAllowed             = 405
	StatusNotAcceptable                = 406
	StatusProxyAuthRequired            = 407
	StatusRequestTimeout               = 408
	StatusConflict                     = 409
	StatusGone                         = 410
	StatusLengthRequired               = 411
	StatusPreconditionFailed           = 412
	StatusRequestEntityTooLarge        = 413
	StatusRequestURITooLong            = 414
	StatusUnsupportedMediaType         = 415
	StatusRequestedRangeNotSatisfiable = 416
	StatusExpectationFailed            = 417
	StatusTeapot                       = 418
	StatusMisdirectedRequest           = 421
	StatusUnprocessableEntity          = 422
	StatusLocked                       = 423
	StatusFailedDependency             = 424
	StatusTooEarly                     = 425
	StatusUpgradeRequired              = 426
	StatusPreconditionRequired         = 428
	StatusTooManyRequests              = 429
	StatusRequestHeaderFieldsTooLarge  = 431
	StatusUnavailableForLegalReasons   = 451

	StatusInternalServerError           = 500
	StatusNotImplemented                = 501
	StatusBadGateway                    = 502
	StatusServiceUnavailable            = 503
	StatusGatewayTimeout                = 504
	StatusHTTPVersionNotSupported       = 505
	StatusVariantAlsoNegotiates         = 506
	StatusInsufficientStorage           = 507
	StatusLoopDetected                  = 508
	StatusNotExtended                   = 510
	StatusNetworkAuthenticationRequired = 511
)

// HTTPError represents an HTTP error.
type HTTPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Err     error       `json:"-"`
}

// NewHTTPError returns a new HTTPError.
func NewHTTPError(code int, message string, details interface{}, err error) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
		Details: details,
		Err:     err,
	}
}

// Error returns the error message.
func (e *HTTPError) Error() string {
	return e.Message
}

// SendError sends the error to the client.
func (e *HTTPError) SendError(w http.ResponseWriter) {
	// TODO: Manage logging
	log.Printf("Error: %v\n", e.Err)

	res, err := json.Marshal(e)
	if err != nil {
		Err500(w, err, "error when encoding the response", nil)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	w.Write(res)
}

func Err(w http.ResponseWriter, status int, err error, msg string, details interface{}) {
	e := NewHTTPError(status, msg, details, err)
	e.SendError(w)
}

func Err400(w http.ResponseWriter, err error, msg string, details interface{}) {
	Err(w, StatusBadRequest, err, msg, details)
}

func Err500(w http.ResponseWriter, err error, msg string, details interface{}) {
	Err(w, StatusInternalServerError, err, msg, details)
}

func JSON(w http.ResponseWriter, data interface{}) {
	res, err := json.Marshal(data)
	if err != nil {
		Err500(w, err, "error when encoding the response", nil)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}
