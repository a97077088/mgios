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

		//获取token
		_, err := Crystal_aquireToken(mguser, nil)
		if err != nil {
			return err
		}
		c := New_Crystal_service_with_user_token(mguser)
		e := c.DataExposureService_Events([]map[string]interface{}{
			map[string]interface{}{
				"timestamp": "1558436308531",
				"du":        "1",
				"eventName": "native_action",
				"eventParams": map[string]interface{}{
					"provinceCode": "",
					"promotionId":  "",
					"cityId":       "",
					"action":       "{\"params\":{\"extra\":{\"gktype\":0,\"isRedrain\":false,\"isCast\":false,\"xAxis\":0,\"isQuickPay\":false,\"isHorizon\":false,\"yAxis\":0,\"isH5Page\":false,\"originalActionType\":\"CLICK_TOP_TAB\",\"isRemote\":false,\"isPullDlna\":false},\"index\":\"1\",\"location\":\"965823946cac4bf3b78c7b99b76b728b#5a4c4b2b47c4498f80e4d7e57bdb67c1#c27e47dfcde243aba5c6f642717368b1#2c62468a9e2b4202834a883568951fce\",\"currentProgress\":0,\"pageID\":\"9aa6465a8b41451e85e4a292bf01ec87\",\"contentID\":\"9aa6465a8b41451e85e4a292bf01ec87\"},\"type\":\"EXPOSE_PROGRAM_DATA\"}",
				},
			},
		})
		err = c.Upload_with_dataexposureservice(e, nil)
		fmt.Println(err)

		return nil
		err = Login_with_u_p_user_cli(_user, _pwd, "焦作市", mguser, nil)
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
