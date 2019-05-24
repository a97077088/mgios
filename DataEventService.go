package main

type DataEventService struct {
	CustomEvent    []map[string]interface{} `json:"customEvent"`
	SdkSessionInfo map[string]interface{}   `json:"sdkSessionInfo"`
}
