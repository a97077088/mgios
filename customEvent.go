package main

import ()

type CustomEvent struct {
	Timestamp   string
	Du          string
	EventName   string
	EventParams map[string]interface{}
}

func Custom_events(_events ...map[string]interface{}) []map[string]interface{} {
	rmp := []map[string]interface{}{}
	for _, it := range _events {
		rmp = append(rmp, it)
	}
	return rmp
}
func New_customEvent_with_timestamp_du_eventname_eventparams(_timestamp string, _du string, _eventname string, _eventparams map[string]interface{}) map[string]interface{} {
	rmp := map[string]interface{}{
		"timestamp":   _timestamp,
		"du":          _du,
		"eventName":   _eventname,
		"eventParams": _eventparams,
	}
	return rmp
}
