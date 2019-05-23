package main

import (
	"fmt"
	"testing"
)

func Test_getMiguSign(t *testing.T) {
	mguser := New_mguser()
	mguser.Data_with_map(map[string]string{
		"$loc_info":         "113.05,35.07",
		"$MGSDK_UUID":       "D876D0297EB04678A6EAC0EE0D1F0FCD",
		"$open_udid":        "9A73ADCA-A513-4AC1-B2FF-A8C62F188453",
		"$wifi_ssid":        "ASUS",
		"$deviceModelName":  "iPhone 6",
		"$systemVersion":    "10.0.2",
		"$ua":               "iPhone7,2",
		"$FCUUID":           "bbb",                                  //strings.ToLower(strings.ReplaceAll(gofakeit.UUID(),"-","")),
		"$idfa":             "39B0FAE6-76FB-43B5-8482-8DDBB255D826", //gofakeit.UUID(),
		"$idfv":             "BA61AD69-F048-41CD-A15C-996935A993D3", //gofakeit.UUID(),
		"$APP-VERSION-CODE": "5.6.3",
	})
	r, err := Crystal_aquireToken(mguser, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}
