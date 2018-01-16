package openshift

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"k8s.io/kubernetes/pkg/api/unversioned"
)

// StatusError is an error intended for consumption by a REST API server; it can also be
// reconstructed by clients from a REST response. Public to allow easy type switches.
type StatusError struct {
	ErrStatus unversioned.Status
}

// Error implements the Error interface.
func (e *StatusError) Error() string {
	return e.ErrStatus.Message
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
//
// The error type will be *RateLimitError for rate limit exceeded errors,
// and *TwoFactorAuthError for two-factor authentication errors.
func CheckApiStatus(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	// openshift returns 401 with a plain text but not ErrStatus json, so we hacked this response text.
	if r.StatusCode == http.StatusUnauthorized {
		errorResponse := &StatusError{}
		errorResponse.ErrStatus.Code = http.StatusUnauthorized
		errorResponse.ErrStatus.Message = http.StatusText(http.StatusUnauthorized)
		return errorResponse
	}

	errorResponse := &StatusError{}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		//clog.Errorf("%s", data)
		json.Unmarshal(data, &errorResponse.ErrStatus)
	}

	return errorResponse
}
