package common

import (
	"encoding/json"
	"io"
	"net/http"
)

// TO BE WRAPPED WITH PROPER ERROR CODE OUTSIDE WHEREVER USED
func ParseBody(r *http.Request, body any) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, body)
	if err != nil {
		return err
	}
	return nil
}
