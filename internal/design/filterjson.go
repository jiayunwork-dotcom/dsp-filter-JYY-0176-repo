package design

import "encoding/json"

func ParseFilterJSON(data []byte) (*Filter, error) {
	var f Filter
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, &Error{Code: ErrBadKind, Message: "invalid filter json: " + err.Error()}
	}
	return &f, nil
}

func (f *Filter) ToJSON() ([]byte, error) {
	return json.MarshalIndent(f, "", "  ")
}
