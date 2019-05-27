package main

type DataExposureService struct {
	CustomEvent    []map[string]interface{} `json:"customEvent"`
	SdkSessionInfo map[string]interface{}   `json:"sdkSessionInfo"`
}

func New_DataExposureService_with_customevents_sdksessioninfo(
	_custominfo []map[string]interface{},
	_sdksessioninfo map[string]interface{},
) *DataExposureService {
	d := &DataExposureService{
		[]map[string]interface{}{},
		map[string]interface{}{},
	}
	if _custominfo != nil {
		d.CustomEvent = _custominfo
	}
	if _sdksessioninfo != nil {
		d.SdkSessionInfo = _sdksessioninfo
	}
	return d
}
