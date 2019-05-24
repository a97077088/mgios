package main

type DataExposureService struct {
	CustomEvent    []map[string]interface{} `json:"customEvent"`
	SdkSessionInfo map[string]interface{}   `json:"sdkSessionInfo"`
}
