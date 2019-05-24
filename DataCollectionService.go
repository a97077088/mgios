package main

type DataCollectionService struct {
	SdkSessionInfo map[string]interface{}   `json:"sdkSessionInfo"`
	DeviceInfo     map[string]interface{}   `json:"deviceInfo"`
	Exception      []map[string]interface{} `json:"exception"`
	CustomInfo     []map[string]interface{} `json:"customInfo"`
	SessionStart   map[string]interface{}   `json:"sessionStart"`
	SessionEnd     map[string]interface{}   `json:"sessionEnd"`
}
