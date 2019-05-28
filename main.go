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
			"MG_MSG_TIME":    Crystal_Now(),
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
			"loadTime":      2303, //2198-2303
			"account":       _u.Value_for_key("$FCUUID"),
			"timestamp":     time.Now().Format("2006-01-02 15:04:05"),
			"type":          "10",
			"totalLoadTime": 21273, //21273-653098
		},
		map[string]interface{}{
			"loadTime": "1164",        //781-1164
			"number":   "30",          //18-30
			"picIP":    "61.54.90.26", //61.54.90.32,61.54.90.26
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

	//下面是每次播放动作
	err = crystal.Upload_with_dataexposureservice(New_DataExposureService_with_customevents_sdksessioninfo(Custom_events(
		New_customEvent_with_timestamp_du_eventname_eventparams(Crystal_Now(), "1", "native_action", map[string]interface{}{
			"provinceCode": _u.Value_for_key("ProvinceCode"),
			"promotionId":  "",
			"cityId":       _u.Value_for_key("cityId"),
			"action":       "{\"params\":{\"extra\":{\"gktype\":0,\"isRedrain\":false,\"isCast\":false,\"xAxis\":0,\"isQuickPay\":false,\"isHorizon\":false,\"yAxis\":0,\"isH5Page\":false,\"originalActionType\":\"JUMP_H5_BY_WEB_VIEW\",\"isRemote\":false,\"isPullDlna\":false},\"url\":\"https:\\/\\/m.miguvideo.com\\/mgs\\/share\\/member_order_tem\\/prd\\/index.html?pageId=011863ec6c814dfab269f5bb1bc8f685&platform=PLATFORM\",\"location\":\"DETAIL_PAGE_MOVIE#4172bb901749437da7d7a29949db7d10#028190ce176246a8a489c42b2e86545e#574b2ff6570e43d1b5e076c0cef5df7f\",\"currentProgress\":0,\"frameID\":\"main-frame\",\"contentID\":\"657047455\"},\"type\":\"EXPOSE_PROGRAM_DATA\",\"name\":\"JUMP_H5_BY_WEB_VIEW\"}",
		}),
	), crystal.SdkSessionInfo()), nil)
	user_log("<native_action:JUMP_H5_BY_WEB_VIEW>成功", _u)

	err = crystal.Upload_with_datacollectionservice_cli(New_DataCollectionService_with_sdksessioninfo_deviceinfo_exception_custominfo_sessionstart_sessionend(crystal.SdkSessionInfo(), nil, nil, []map[string]interface{}{
		map[string]interface{}{
			"result":    "0",
			"loadTime":  1288,
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			"source_id": "657047455", //这个应该是作品id
			"type":      "11",
			"account":   _u.Value_for_key("$FCUUID"),
		},
		map[string]interface{}{
			"Result":        "1",
			"account":       _u.Value_for_key("$FCUUID"),
			"ContentID":     "657047455", //这个应该是作品id
			"type":          "22",
			"loadTime":      434,
			"timestamp":     time.Now().Format("2006-01-02 15:04:05"),
			"playSessionID": "39ac36f4fdf739db5385c92b710a0edf", //这里是jid
		},
		map[string]interface{}{
			"clientID":                           _u.Value_for_key("$FCUUID"),
			"CurNetType":                         "Wifi",
			"MG_MSG_TIME":                        Crystal_Now(),
			"MG_MSG_MEDIA_INFO_AUDIO_CODEC":      "aac",
			"MG_MSG_MEDIA_INFO_TYPE":             "AV",
			"SubsessionServiceIP":                "182.118.23.11",
			"MG_MSG_MEDIA_INFO_AUDIO_CHANNELS":   "2",
			"location":                           "29efb2c6badd494c906000403f0d2728#b3e99ee9d9f7496dae3d02e2a57acb84#88452153e73641c4bd693a8291012117",
			"account":                            _u.Value_for_key("$FCUUID"),
			"type":                               "56000004",
			"MG_MSG_START_TIME":                  "1559036991483", //开始时间
			"contentType":                        "3",
			"MG_MSG_MEDIA_INFO_AUDIO_SAMPLERATE": "44100",
			"MG_TOTAL_LOAD_TIME":                 "1532.000000",
			"index":                              "1",
			"rateType":                           "3",
			"jid":                                "39ac36f4fdf739db5385c92b710a0edf", //jid
			"Status":                             "1",
			"MG_MSG_MEDIA_INFO_VIDEO_BITRATE":    "1449",
			"MG_MSG_PROGRAM_URL":                 "http://gslbmgspvod.miguvideo.com/depository_sjq/asset/zhengshi/5102/219/103/5102219103/media/5102219103_5005016096_91.mp4.m3u8?msisdn=3a98f8560a7c419c9b711a2e545575e7&mdspid=&spid=600058&netType=4&sid=5500743415&pid=2028600738&timestamp=20190528174948&Channel_ID=0116_25040600-99000-200300020100001&ProgramID=657047455&ParentNodeID=-99&preview=1&playseek=000000-000600&assertID=5500743415&client_ip=42.234.8.51&SecurityKey=20190528174948&imei=633D6707-7DC2-4055-98FC-7BC87486647E&promotionId=&mvid=5102219103&mcid=1000&mpid=130000119820&playurlVersion=SJ-A1-0.0.2&encrypt=a939c15fb3fc6e9244d7924892b47320&jid=39ac36f4fdf739db5385c92b710a0edf",
			"MG_MSG_MEDIA_INFO_VIDEO_FRAMERATE":  "25",
			"MG_MSG_GETURL_TIME":                 "1559036990423", //获取url时间
			"channelID":                          X_UP_CLIENT_CHANNEL_ID,
			"LastSession":                        " ",
			"PlaybackRate":                       "1.00",
			"Session":                            "39ac36f4fdf739db5385c92b710a0edf", //jid
			"MG_MSG_MEDIA_INFO_VIDEO_CODEC":      "h264",
			"MG_MSG_FFRAME_TIME":                 "1098",
			"contentId":                          "657047455", //作品
			"MG_MSG_MEDIA_INFO_VIDEO_RESOLUTION": "1280x720",
			"SourceID":                           "18a98758fe54d7809c64b529c4859ca4_1559036989941", //暂时没搞明白
			"PlayerID":                           "139569bb0_1559036989937",                        //PlayerID
			"title":                              "《我们都要好好的》第01集",
			"MG_URL_LOAD_TIME":                   "434",
			"chipSet":                            "iPhone 6",
			"userId":                             "",
			"timestamp":                          time.Now().Format("2006-01-02 15:04:05"),
		},
		map[string]interface{}{
			"loadTime": "120",
			"number":   "5",
			"picIP":    "61.54.90.32",
			"type":     "19",
		},
		map[string]interface{}{
			"Session":  "39ac36f4fdf739db5385c92b710a0edf",
			"PlayerID": "139569bb0_1559036989937",
			"userId":   "", "SubsessionServiceURL": "http://mgsp.vod.miguvideo.com:8088/depository_sjq/asset/zhengshi/5102/219/103/5102219103/media/5102219103_5005016096_91.mp4_0-19.ts?msisdn=3a98f8560a7c419c9b711a2e545575e7&mdspid=&spid=600058&netType=4&sid=5500743415&pid=2028600738&timestamp=20190528174948&Channel_ID=0116_25040600-99000-200300020100001&ProgramID=657047455&ParentNodeID=-99&preview=1&playseek=000000-000600&assertID=5500743415&client_ip=42.234.8.51&SecurityKey=20190528174948&imei=633D6707-7DC2-4055-98FC-7BC87486647E&promotionId=&mvid=5102219103&mcid=1000&mpid=130000119820&playurlVersion=SJ-A1-0.0.2&FreePlay=1&encrypt=0b72c39d0fe29dd8a358a2dc5604b193&jid=39ac36f4fdf739db5385c92b710a0edf&sjid=subsession_1559036990425&hls_type=2&HlsSubType=2&HlsProfileId=0&mtv_session=0b72c39d0fe29dd8a358a2dc5604b193", "title": "《我们都要好好的》第01集", "account": "3a98f8560a7c419c9b711a2e545575e7", "location": "29efb2c6badd494c906000403f0d2728#b3e99ee9d9f7496dae3d02e2a57acb84#88452153e73641c4bd693a8291012117", "clientID": "3a98f8560a7c419c9b711a2e545575e7", "index": "1", "MG_MSG_PROGRAM_URL": "http://gslbmgspvod.miguvideo.com/depository_sjq/asset/zhengshi/5102/219/103/5102219103/media/5102219103_5005016096_91.mp4.m3u8?msisdn=3a98f8560a7c419c9b711a2e545575e7&mdspid=&spid=600058&netType=4&sid=5500743415&pid=2028600738&timestamp=20190528174948&Channel_ID=0116_25040600-99000-200300020100001&ProgramID=657047455&ParentNodeID=-99&preview=1&playseek=000000-000600&assertID=5500743415&client_ip=42.234.8.51&SecurityKey=20190528174948&imei=633D6707-7DC2-4055-98FC-7BC87486647E&promotionId=&mvid=5102219103&mcid=1000&mpid=130000119820&playurlVersion=SJ-A1-0.0.2&encrypt=a939c15fb3fc6e9244d7924892b47320&jid=39ac36f4fdf739db5385c92b710a0edf",
			"MG_EVENT_PERIOD":       "15",
			"PlayDuration":          "14215",
			"SourceID":              "18a98758fe54d7809c64b529c4859ca4_1559036989941",
			"jid":                   "39ac36f4fdf739db5385c92b710a0edf",
			"type":                  "60000000",
			"channelID":             X_UP_CLIENT_CHANNEL_ID,
			"MG_MSG_TIME":           Crystal_Now(),
			"rateType":              "3",
			"PlaybackRate":          "1.00",
			"MG_MSG_PLAYER_VERSION": "9.2.0.25",
			"contentType":           "3",
			"chipSet":               "iPhone 6",
			"timestamp":             time.Now().Format("2006-01-02 15:04:05"),
			"contentId":             "657047455",
			"NetType":               "Wifi",
			"Subsession":            "subsession_1559036990425",
			"SubsessionServiceIP":   "182.118.23.11",
		},
	}, nil, nil), nil)
	user_log("<videoPageLoadDelay:11>,<getPlayUrlDuration:22>,<MGFirstVideoRenderEvent:56000004>,<MGPlayerDurationEvent:60000000>,<picloadDelay:19>", _u)
	//接下来一直发60000000保持连接的包

	//下边是断开连接
	err = crystal.Upload_with_datacollectionservice_cli(New_DataCollectionService_with_sdksessioninfo_deviceinfo_exception_custominfo_sessionstart_sessionend(crystal.SdkSessionInfo(), nil, nil, []map[string]interface{}{
		map[string]interface{}{
			"account":               _u.Value_for_key("$FCUUID"),
			"PlayerID":              "139569bb0_1559036989937",
			"Session":               "39ac36f4fdf739db5385c92b710a0edf",
			"userId":                "",
			"title":                 "《我们都要好好的》第01集",
			"location":              "29efb2c6badd494c906000403f0d2728#b3e99ee9d9f7496dae3d02e2a57acb84#88452153e73641c4bd693a8291012117",
			"DataUsage":             "95710773",
			"clientID":              _u.Value_for_key("$FCUUID"),
			"index":                 "1",
			"MG_MSG_PROGRAM_URL":    "http://gslbmgspvod.miguvideo.com/depository_sjq/asset/zhengshi/5102/219/103/5102219103/media/5102219103_5005016096_91.mp4.m3u8?msisdn=3a98f8560a7c419c9b711a2e545575e7&mdspid=&spid=600058&netType=4&sid=5500743415&pid=2028600738&timestamp=20190528174948&Channel_ID=0116_25040600-99000-200300020100001&ProgramID=657047455&ParentNodeID=-99&preview=1&playseek=000000-000600&assertID=5500743415&client_ip=42.234.8.51&SecurityKey=20190528174948&imei=633D6707-7DC2-4055-98FC-7BC87486647E&promotionId=&mvid=5102219103&mcid=1000&mpid=130000119820&playurlVersion=SJ-A1-0.0.2&encrypt=a939c15fb3fc6e9244d7924892b47320&jid=39ac36f4fdf739db5385c92b710a0edf",
			"PlayDuration":          "360784", //播放总时长应该是
			"SourceID":              "18a98758fe54d7809c64b529c4859ca4_1559036989941",
			"video_length":          "5838961", //视频长度
			"type":                  "70000000",
			"quit_point":            "360612",
			"jid":                   "39ac36f4fdf739db5385c92b710a0edf",
			"channelID":             X_UP_CLIENT_CHANNEL_ID,
			"MG_MSG_TIME":           Crystal_Now(),
			"rateType":              "3",
			"PlaybackRate":          "1.00",
			"MG_MSG_PLAYER_VERSION": "9.2.0.25",
			"contentType":           "3",
			"chipSet":               "iPhone 6",
			"timestamp":             time.Now().Format("2006-01-02 15:04:05"),
			"contentId":             "657047455",
			"Subsession":            "subsession_1559036990425",
			"SubsessionServiceIP":   "182.118.23.11",
		},
		map[string]interface{}{
			"Session":               "39ac36f4fdf739db5385c92b710a0edf",
			"PlayerID":              "139569bb0_1559036989937",
			"userId":                "",
			"SubsessionServiceURL":  "http://mgsp.vod.miguvideo.com:8088/depository_sjq/asset/zhengshi/5102/219/103/5102219103/media/5102219103_5005016096_91.mp4_0-37.ts?msisdn=3a98f8560a7c419c9b711a2e545575e7&mdspid=&spid=600058&netType=4&sid=5500743415&pid=2028600738&timestamp=20190528174948&Channel_ID=0116_25040600-99000-200300020100001&ProgramID=657047455&ParentNodeID=-99&preview=1&playseek=000000-000600&assertID=5500743415&client_ip=42.234.8.51&SecurityKey=20190528174948&imei=633D6707-7DC2-4055-98FC-7BC87486647E&promotionId=&mvid=5102219103&mcid=1000&mpid=130000119820&playurlVersion=SJ-A1-0.0.2&FreePlay=0&encrypt=a939c15fb3fc6e9244d7924892b47320&jid=39ac36f4fdf739db5385c92b710a0edf&sjid=subsession_1559036990425&hls_type=2&HlsSubType=2&HlsProfileId=0&mtv_session=a939c15fb3fc6e9244d7924892b47320",
			"Action":                "0",
			"account":               _u.Value_for_key("$FCUUID"),
			"title":                 "《我们都要好好的》第01集",
			"DataUsage":             "95710773",
			"location":              "29efb2c6badd494c906000403f0d2728#b3e99ee9d9f7496dae3d02e2a57acb84#88452153e73641c4bd693a8291012117",
			"clientID":              "3a98f8560a7c419c9b711a2e545575e7",
			"HostIP":                "42.234.8.51", //这个是我的ip
			"index":                 "1",
			"PlayDuration":          "360797",
			"SourceID":              "18a98758fe54d7809c64b529c4859ca4_1559036989941",
			"jid":                   "39ac36f4fdf739db5385c92b710a0edf",
			"type":                  "57000000",
			"channelID":             X_UP_CLIENT_CHANNEL_ID,
			"MG_MSG_TIME":           Crystal_Now(),
			"rateType":              "3",
			"PlaybackRate":          "1.00",
			"MG_MSG_PLAYER_VERSION": "9.2.0.25",
			"CurNetType":            "Wifi",
			"BeginTime":             "1559036990428",
			"contentType":           "3",
			"chipSet":               "iPhone 6",
			"contentId":             "657047455",
			"timestamp":             time.Now().Format("2006-01-02 15:04:05"),
			"EndTime":               "1559037351983",
			"NetType":               "Wifi",
			"Subsession":            "subsession_1559036990425",
			"SubsessionServiceIP":   "182.118.23.11",
		},
	}, nil, nil), nil)
	user_log("<MGPlayerShutdownEvent:70000000>,<MGTrafficStatusticsEvent:57000000>成功", _u)

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
