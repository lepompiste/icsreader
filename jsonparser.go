package icsreader

import "encoding/json"

// ParseJSON : parse Events into JSON readable format
func (evs *Events) ParseJSON() string {
	data, err := json.Marshal(evs)
	if err == nil {
		return string(data)
	}
	return ""
}
