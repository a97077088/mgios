package main

import (
	"fmt"
)

func main() {

	err := func() error {
		_user := "15771786564"
		_pwd := "pp786564"
		mguser := New_mguser_with_map(map[string]string{
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

		err := Login_with_u_p_user_cli(_user, _pwd, "焦作市", mguser, nil)
		if err != nil {
			return err
		}

		fmt.Printf("登录成功 用户:%s 地区:%s_%s token:%s\n", mguser.Value_for_key("mobile"), mguser.Value_for_key("ProvinceCode"), mguser.Value_for_key("cityId"), mguser.Value_for_key("userToken"))
		r, err := Playurl_v1_play_url("654961181", 4, "656373955", "27", mguser, nil)
		if err != nil {
			return err
		}
		fmt.Println(r)
		return nil
	}()
	if err != nil {
		panic(err)
	}

}
