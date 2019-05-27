package main

type DataEventService struct {
	CustomEvent    []map[string]interface{} `json:"customEvent"`
	SdkSessionInfo map[string]interface{}   `json:"sdkSessionInfo"`
}

func New_DataEventService_with_customevent_sdksessioninfo(
	_customevent []map[string]interface{},
	_sdksessioninfo map[string]interface{},
) *DataEventService {
	d := &DataEventService{
		[]map[string]interface{}{},
		map[string]interface{}{},
	}
	if _customevent != nil {
		d.CustomEvent = _customevent
	}
	if _sdksessioninfo != nil {
		d.SdkSessionInfo = _sdksessioninfo
	}
	return d
}
