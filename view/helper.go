package view

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (res *Response) WithMessage(msg string) *Response {
	res.Message = msg

	return res
}

func (res *Response) WithError(err string, code ...int) *Response {
	res.Error = err
	res.Success = false
	res.HTTPCode = http.StatusBadRequest

	if len(code) > 0 {
		res.HTTPCode = code[0]
	}

	return res
}

func (res *Response) MarshalJSON() ([]byte, error) {
	mappingObj := make(map[string]any)

	if res.ExtraFields != nil {
		for k, v := range res.ExtraFields {
			mappingObj[k] = v
		}
	}

	mappingObj["data"] = res.Data

	if res.Error != "" {
		mappingObj["error"] = res.Error
	}

	if res.Message != "" {
		mappingObj["message"] = res.Message
	}

	if res.Success || !res.Success {
		mappingObj["success"] = res.Success
	}

	if res.HTTPCode != 0 {
		mappingObj["httpCode"] = res.HTTPCode
	}

	data, err := json.Marshal(mappingObj)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}
