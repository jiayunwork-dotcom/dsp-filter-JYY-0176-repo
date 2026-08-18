package design

import "encoding/json"

func ParseSpecJSON(data []byte) (*DesignSpec, error) {
	var spec DesignSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, &Error{Code: ErrBadKind, Message: "invalid spec json: " + err.Error()}
	}
	return &spec, nil
}

func (s *DesignSpec) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}
