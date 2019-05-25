package main

type DataCollectionService struct {
	SdkSessionInfo map[string]interface{}   `json:"sdkSessionInfo"`
	DeviceInfo     map[string]interface{}   `json:"deviceInfo"`
	Exception      []map[string]interface{} `json:"exception"`
	CustomInfo     []map[string]interface{} `json:"customInfo"`
	SessionStart   map[string]interface{}   `json:"sessionStart"`
	SessionEnd     map[string]interface{}   `json:"sessionEnd"`
}

func New_DataCollectionService_with_sdksessioninfo_deviceinfo_exception_custominfo_sessionstart_sessionend(
	_sdksessioninfo map[string]interface{},
	_deviceinfo map[string]interface{},
	_exception []map[string]interface{},
	_custominfo []map[string]interface{},
	_sessionstart map[string]interface{},
	_sessionend map[string]interface{}) *DataCollectionService {
	d := &DataCollectionService{
		map[string]interface{}{},
		map[string]interface{}{},
		[]map[string]interface{}{},
		[]map[string]interface{}{},
		map[string]interface{}{},
		map[string]interface{}{},
	}
	if _sdksessioninfo != nil {
		d.SdkSessionInfo = _sdksessioninfo
	}
	if _deviceinfo != nil {
		d.DeviceInfo = _deviceinfo
	}
	if _exception != nil {
		d.Exception = _exception
	}
	if _custominfo != nil {
		d.CustomInfo = _custominfo
	}
	if _sessionstart != nil {
		d.SessionStart = _sessionstart
	}
	if _sessionend != nil {
		d.SessionEnd = _sessionend
	}
	return d
}
