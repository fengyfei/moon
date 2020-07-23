package request

import (
	"io"
	"net/http"
)

// Get - Get Request
func Get(url string, body io.Reader) *http.Request {
	req, err := http.NewRequest("GET", url, body)

	if err != nil {
		return nil
	}

	return req
}

// SetAjaxHeader -
func SetAjaxHeader(req *http.Request, refer string) {
	req.Header.Set("X-Request-Type", "ajax")
	req.Header.Set("X-Requested-With", "XMLHTTPRequest")
	req.Header.Set("Referer", refer)
}
