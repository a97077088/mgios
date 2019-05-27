package main

import (
	"fmt"
	"log"
	"time"
)

func user_log(_s string, _u *MGUser) {
	log.Printf("%s %s\n", _u.Value_for_key("username"), _s)
}

func user_task_exec(_u *MGUser) error {
	user_log("启动埋点服务", _u)
	crystal_tk, err := Crystal_aquireToken(_u, nil)
	if err != nil {
		return err
	}
	user_log(fmt.Sprintf("拿到crystal_token:%s", crystal_tk.Token), _u)
	crystal := New_Crystal_service_with_user_token(_u)
	err = crystal.Upload_with_datacollectionservice_cli(crystal.DataCollectionService_session_start(), nil)
	if err != nil {
		return err
	}
	user_log("session_start 成功", _u)

	err = crystal.Upload_with_dataeventservice(New_DataEventService_with_customevent_sdksessioninfo(
		Custom_events(
			New_customEvent_with_timestamp_du_eventname_eventparams(Crystal_Now(), "0", "firstLaunch", map[string]interface{}{
				"appVersion":    _u.Value_for_key("$APP-VERSION-CODE"),
				"time":          Crystal_Now(),
				"oldAppVersion": "",
			}),
		), crystal.SdkSessionInfo()), nil)
	user_log("firstLaunch 完成", _u)

	err = crystal.Upload_with_datacollectionservice_cli(New_DataCollectionService_with_sdksessioninfo_deviceinfo_exception_custominfo_sessionstart_sessionend(crystal.SdkSessionInfo(), nil, nil, []map[string]interface{}{
		map[string]interface{}{
			"index":          "",
			"timestamp":      time.Now().Format("2006-01-02 15:04:05"),
			"jid":            "",
			"contentType":    "",
			"channelID":      X_UP_CLIENT_CHANNEL_ID,
			"type":           "80000000",
			"title":          "",
			"userId":         "",
			"Result":         "0",
			"MG_MSG_TIME":    fmt.Sprintf("%d", time.Now().UnixNano()/1e6),
			"APP_PLAYER_KEY": "com.wondertek.hecmccmobile",
			"clientID":       _u.Value_for_key("$FCUUID"),
			"APP_NAME":       "咪咕视频",
			"location":       "",
			"account":        _u.Value_for_key("$FCUUID"),
			"chipSet":        _u.Value_for_key("$deviceModelName"),
			"rateType":       "",
		},
	}, nil, nil), nil)
	if err != nil {
		return err
	}
	user_log("<MGAuthenticationEvent:80000000>成功", _u)

	err = crystal.Upload_with_dataexposureservice(New_DataExposureService_with_customevents_sdksessioninfo(Custom_events(
		New_customEvent_with_timestamp_du_eventname_eventparams(Crystal_Now(), "1", "native_action", map[string]interface{}{
			"provinceCode": "",
			"promotionId":  "",
			"cityId":       "",
			"action":       `{\"params\":{\"extra\":{\"gktype\":0,\"isRedrain\":false,\"isCast\":false,\"xAxis\":0,\"isQuickPay\":false,\"isHorizon\":false,\"yAxis\":0,\"isH5Page\":false,\"originalActionType\":\"CLICK_TOP_TAB\",\"isRemote\":false,\"isPullDlna\":false},\"index\":\"1\",\"location\":\"965823946cac4bf3b78c7b99b76b728b#5a4c4b2b47c4498f80e4d7e57bdb67c1#c27e47dfcde243aba5c6f642717368b1#2c62468a9e2b4202834a883568951fce\",\"currentProgress\":0,\"pageID\":\"9aa6465a8b41451e85e4a292bf01ec87\",\"contentID\":\"9aa6465a8b41451e85e4a292bf01ec87\"},\"type\":\"EXPOSE_PROGRAM_DATA\"}`,
		}),
	), crystal.SdkSessionInfo()), nil)
	if err != nil {
		return err
	}
	user_log("<native_action>*15", _u)

	err = crystal.Upload_with_datacollectionservice_cli(New_DataCollectionService_with_sdksessioninfo_deviceinfo_exception_custominfo_sessionstart_sessionend(crystal.SdkSessionInfo(), nil, nil, []map[string]interface{}{
		map[string]interface{}{
			"result":        "0",
			"loadTime":      2303,
			"account":       _u.Value_for_key("$FCUUID"),
			"timestamp":     time.Now().Format("2006-01-02 15:04:05"),
			"type":          "10",
			"totalLoadTime": 21273,
		},
		map[string]interface{}{
			"loadTime": "1164",
			"number":   "30",
			"picIP":    "61.54.90.26",
			"type":     "19",
		},
	}, nil, nil), nil)
	if err != nil {
		return err
	}
	user_log("<firstPageDelay:10>,<picloadDelay:19>成功", _u)

	user_log("准备开始登录", _u)
	err = Login_with_u_p_user_cli(_u.Value_for_key("username"), _u.Value_for_key("userpass"), "焦作市", _u, nil)
	if err != nil {
		return err
	}
	user_log(fmt.Sprintf("登录成功,token:%s", _u.Value_for_key("userToken")), _u)

	return nil
}

func user_task(_u *MGUser) {
	err := user_task_exec(_u)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	us, err := New_users_with_file("用户.txt")
	if err != nil {
		panic(err)
	}
	for _, u := range us {
		go func(_u *MGUser) {
			user_task(_u)
		}(u)
	}

	for {

		time.Sleep(time.Second * 1)
	}

	//
	//err := func() error {
	//	_user := "15771786564"
	//	_pwd := "pp786564"
	//
	//	err := Login_with_u_p_user_cli(_user, _pwd, "焦作市", mguser, nil)
	//	if err != nil {
	//		return err
	//	}
	//
	//	fmt.Printf("登录成功 用户:%s 地区:%s_%s token:%s\n", mguser.Value_for_key("mobile"), mguser.Value_for_key("ProvinceCode"), mguser.Value_for_key("cityId"), mguser.Value_for_key("userToken"))
	//
	//	//获取token
	//	_, err = Crystal_aquireToken(mguser, nil)
	//	if err != nil {
	//		return err
	//	}
	//	c := New_Crystal_service_with_user_token(mguser)
	//	e := c.DataExposureService_Events([]map[string]interface{}{
	//		map[string]interface{}{
	//			"timestamp": "1558436308531",
	//			"du":        "1",
	//			"eventName": "native_action",
	//			"eventParams": map[string]interface{}{
	//				"provinceCode": "",
	//				"promotionId":  "",
	//				"cityId":       "",
	//				"action":       "{\"params\":{\"extra\":{\"gktype\":0,\"isRedrain\":false,\"isCast\":false,\"xAxis\":0,\"isQuickPay\":false,\"isHorizon\":false,\"yAxis\":0,\"isH5Page\":false,\"originalActionType\":\"CLICK_TOP_TAB\",\"isRemote\":false,\"isPullDlna\":false},\"index\":\"1\",\"location\":\"965823946cac4bf3b78c7b99b76b728b#5a4c4b2b47c4498f80e4d7e57bdb67c1#c27e47dfcde243aba5c6f642717368b1#2c62468a9e2b4202834a883568951fce\",\"currentProgress\":0,\"pageID\":\"9aa6465a8b41451e85e4a292bf01ec87\",\"contentID\":\"9aa6465a8b41451e85e4a292bf01ec87\"},\"type\":\"EXPOSE_PROGRAM_DATA\"}",
	//			},
	//		},
	//	})
	//	err = c.Upload_with_dataexposureservice(e, nil)
	//	fmt.Println(err)
	//
	//	return nil
	//
	//	//r, err := Playurl_v1_play_url("654961181", 4, "656373955", "27", mguser, nil)
	//	//if err != nil {
	//	//	return err
	//	//}
	//	//fmt.Println(r)
	//	return nil
	//}()
	//if err != nil {
	//	panic(err)
	//}

}
