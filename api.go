package main

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	. "github.com/a97077088/grequests"
	"github.com/brianvoe/gofakeit"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/*
新建接口快捷cli
*/
func cli_with_cli(_cli *Session) *Session {
	if _cli == nil {
		_cli = NewSession(nil)
	}
	return _cli
}

/*
登录接口
*/
func Login_with_u_p_user_cli(_u, _p string, _areaName string, _user IUser, _cli *Session) error {
	rauth, err := Client_authRequest(_u, _p, _user, nil)
	if err != nil {
		return err
	}
	//rauth:=map[string]string{}
	_, err = Login_migutokenforall(rauth, _p, _user, nil)
	if err != nil {
		return err
	}
	cr, err := Common_flow_service_province_code(_areaName, _user, nil)
	if err != nil {
		return err
	}
	_user.Set_key_value("cityId", cr.Body.CityId)
	_user.Set_key_value("ProvinceCode", cr.Body.ProvinceCode)
	return nil
}

/*
登录子接口
*/
func Client_authRequest(_username string, _password string, _user IUser, _cli *Session) (map[string]string, error) {
	strurl := "https://passport.migu.cn/client/authRequest"
	cli := cli_with_cli(_cli)
	//rcdata的实现位置 [UAUtil rcData]
	//net_oper应该是移动商代码，46002是移动gsm这有详细介绍http://blog.chinaunix.net/uid-20484604-id-1941290.html
	//net_country_iso 不用说了，cn
	//mnc具体不知道是什么，固定的不管他
	//loc_info 当前设备经纬度
	var rcdata = fmt.Sprintf(`{  "network_operator" : {    "net_oper" : "46002",    "net_country_iso" : "cn"  },  "os" : "iOS %s",  "mnc" : "02",  "wifi_ssid" : "%s",  "dev_model" : "%s",  "idfv" : "%s",  "open_udid" : "%s",  "app_name" : "com.wondertek.hecmccmobile",  "loc_info" : "%s"}`,
		_user.Value_for_key("$systemVersion"),
		base64.StdEncoding.EncodeToString([]byte(_user.Value_for_key("$wifi_ssid"))),
		strings.ReplaceAll(_user.Value_for_key("$deviceModelName"), " ", ""),
		_user.Value_for_key("$idfv"),
		_user.Value_for_key("$open_udid"),
		_user.Value_for_key("$loc_info"))
	encrcdata := base64.StdEncoding.EncodeToString(Must_Aes128_with_in_key([]byte(rcdata), []byte("15acb4a88285ed2f")))
	tm := time.Now().Unix()                                                                                                    //得到10位时间戳
	md5fmcnonce := hex.EncodeToString(Must_Md5_with_in([]byte((fmt.Sprintf("%d%s", tm, _user.Value_for_key("$MGSDK_UUID")))))) //10位时间戳+[IDMPDevice getDeviceID]
	btnonce, err := Rsa_with_in_publickey([]byte(md5fmcnonce), []byte(AuthRequestpubkey))
	if err != nil {
		return nil, err
	}
	nonce := hex.EncodeToString(btnonce)
	btpwd, err := Rsa_with_in_publickey([]byte(_password), []byte(AuthRequestpubkey))
	if err != nil {
		return nil, err
	}
	pwd := hex.EncodeToString(btpwd)
	msgid := hex.EncodeToString(Must_Md5_with_in([]byte(fmt.Sprintf("%s%s", time.Now().Format("20060102150405000"), _user.Value_for_key("$MGSDK_UUID")))))                                                                                                                                                                                            //MGSDKRequest msgId
	Authorization := fmt.Sprintf(`UP clientversion="1.0",sdkversion="MGUnionAuthSDK-iOSV4.6.0",appid="20300507",apptype="4",username="%s",isusernamecn="0",enccnonce="%s",encpasswd="%s",BTID="",Phone_ID="%s",msgid="%s",rcflag="1",sdkverifytype="2",secflag="1",tvflag="false"`, _username, nonce, pwd, _user.Value_for_key("$MGSDK_UUID"), msgid) //`UP clientversion="1.0",sdkversion="MGUnionAuthSDK-iOSV4.6.0",appid="20300507",apptype="4",username="15093717639",isusernamecn="0",enccnonce="3be146a9e42d0e54106485d9b274702ed0c5b25dd0289622ab162d9f439294ce28b58bc3076d6870e7d4058d6dc272070c41d79ca4dd06e414ec2567314b577a1f228d57eacd1d989a44665406de6e4f5ece3576f27a89cb3c52f03c466237fc2aefd33502da5b986fba7dc488fe15873827e21eab02e5600ba73118ddc4113b",encpasswd="8e9a1211155df0b57de28901d81a26e3cdb6ce0904a3154916218a92a51b7d3eeb2dfd45495f5bb0a1089528d0d2dc89da1f3a7182fbc79c0bfa3c4f4da6df8b2c1de22ce5fa4cb66c031520f325a3fff6f296d681a6d3f3e8d142f1fa6a2ae73b11f479c34a6fce1bf73f38cc7c394d2abc124210207763823e3658ba7d3aeb",BTID="",Phone_ID="D876D0297EB04678A6EAC0EE0D1F0FCD",msgid="282a520abec20255f3fe9033f1260b95",rcflag="1",sdkverifytype="2",secflag="1",tvflag="false"`,
	//IDMPRSA_Encrypt_Decrypt addPrivateKey:
	//SecKeyRawSign
	signature := hex.EncodeToString(Must_Rsa_sign_sha256_with_in_privatekey([]byte(Authorization), []byte(Signature_privatekey)))
	r, err := cli.Get(strurl, &RequestOptions{
		Headers: map[string]string{ //[MGSDKUPRequest head_parameter]
			"Connection":    "keep-alive",
			"x-gbatype":     "PW_GBA",
			"signature":     signature,
			"encRcData":     encrcdata, //+[UACrypto AESEncryptWithKey:andData:]  这一组没什么用，不挂也可以登录
			"Authorization": Authorization,
		},
		UserAgent: "%E5%92%AA%E5%92%95%E8%A7%86%E9%A2%91" + fmt.Sprintf("/%s.5 CFNetwork/808.0.2 Darwin/16.0.0", _user.Value_for_key("$APP-VERSION-CODE")),
	})
	if err != nil {
		return nil, New_neterr_with_error(err)
	}
	defer r.Close()
	if r.Header.Get("Resultcode") != "103000" {
		if r.Header.Get("Resultcode") == "103507" {
			return nil, errors.New("密码尝试次数过多")
		}
		br, _ := base64.StdEncoding.DecodeString(r.Header.Get("WWW-Authenticate"))
		return nil, errors.New(string(br))
	}

	rs := map[string]string{}
	spr := strings.Split(r.Header.Get("WWW-Authenticate"), ",")
	for _, v := range spr {
		spv := strings.Split(v, "=")
		if len(spv) == 2 {
			rs[spv[0]] = strings.ReplaceAll(spv[1], "\"", "")
		}
	}
	rs["mac"] = r.Header.Get("mac")
	rs["reqtime"] = strconv.FormatUint(uint64(tm), 10)
	return rs, nil
}

/*
+[IDMPUPMode getUPKSByUserName:andPassWd:sessionid:validcode:rcflag:tvflag:validtype:successBlock:failBlock:]
sub_1002F5CA4
[MGSecondAuthAPI requestUrl]
总结一下migutoken的算法,简单几句话,
第一步拿到ks作为下一步计算密钥
第二步计算一下sha256作为下一步key，然后hmacsha256作为下一步key，反正一直往下推
登录子接口1
*/
func Login_migutokenforall(_rauth map[string]string, _password string, _user IUser, _cli *Session) (*Login_Migutokenforall_o, error) {
	cli := cli_with_cli(_cli)

	passid, err := decodeauthfield_with_hexs_privatekey(_rauth["passid"], []byte(Login_migutokenforall_uadecode_privatekey))
	if err != nil {
		return nil, err
	}
	sqn, err := decodeauthfield_with_hexs_privatekey(_rauth["sqn"], []byte(Login_migutokenforall_uadecode_privatekey))
	if err != nil {
		return nil, err
	}
	msisdn, err := decodeauthfield_with_hexs_privatekey(_rauth["msisdn"], []byte(Login_migutokenforall_uadecode_privatekey))
	if err != nil {
		return nil, err
	}
	btid, err := decodeauthfield_with_hexs_privatekey(_rauth["BTID"], []byte(Login_migutokenforall_uadecode_privatekey))
	if err != nil {
		return nil, err
	}
	username, err := decodeauthfield_with_hexs_privatekey(_rauth["username"], []byte(Login_migutokenforall_uadecode_privatekey))
	if err != nil {
		return nil, err
	}

	mpparam := map[string]string{
		"sqn":            sqn,
		"mac":            _rauth["mac"],
		"implicit":       _rauth["implicit"],
		"isLocalNum":     "0",
		"msisdn":         msisdn,
		"passid":         passid,
		"authtype":       _rauth["authtype"],
		"expiretime":     _rauth["expiretime"],
		"mg_global_flag": "1",
		"username":       username,
		"getKSWay":       "UP",
		"email":          _rauth["email"],
		"BTID":           btid,
		"UP Nonce":       _rauth["UP Nonce"],
		"reqtime":        _rauth["reqtime"],
	}

	//MIGU_LOGIN_PW
	//migu-token-vers   这一组是测试数据
	//mpparam["UP Nonce"]="XP6PStlT2g2iW0McapRFv0UhwiTa2Rqz"
	//mpparam["sqn"]="86843200"
	//mpparam["reqtime"]="1557918216"
	//mpparam["BTID"]="MzQ5NzBEQUNGNUQyMzU3NjlC@http://passport.migu.cn/n0001/@6810d2a850d08056bcedbd1864159aee"
	//mpparam["passid"]="73378307247815"
	//mpparam["username"]="15093717639"
	//_password="a1231231"

	//itm,_:=strconv.Atoi(mpparam["reqtime"])
	ks, err := get_ks_with_auth_user(mpparam, _password, _user)
	if err != nil {
		return nil, err
	}
	mpparam["KS"] = string(ks)
	bymigutoken, err := getTokenWithUserName_appId_uuid_andUser_(passid, "203005", _user.Value_for_key("$MGSDK_UUID"), mpparam)
	if err != nil {
		return nil, err
	}
	migutoken := strings.ToUpper(hex.EncodeToString(bymigutoken))
	timespamp := time.Now().Format("20060102150405")
	//fmt.Println(timespamp)
	//timespamp:=""
	//timespamp="20190513014826"
	//[MGGlobalKey setCurrentEnvironment:]设置，是固定的
	clientid := "7a9582bd-8660-4f1a-9d0f-8451f688c67b"
	jsonbody := fmt.Sprintf(`{"deviceId":"%s","timestamp":%s,"extInfo":{"IDFA":"%s"},"miguToken":"%s"}`,
		_user.Value_for_key("$idfv"),
		timespamp,
		_user.Value_for_key("$idfa"),
		migutoken)
	md5jsonbody := hex.EncodeToString(Must_Md5_with_in([]byte(jsonbody)))
	//这一步实现了java中的rsa_privateenc 真几把不好实现这个耗时一下午
	strsign := base64.StdEncoding.EncodeToString(Must_Rsa_with_in_privatekey([]byte(md5jsonbody), []byte(Login_migutokenforall_rsasign_privatekey)))
	strurl := fmt.Sprintf("https://api.miguvideo.com/login/migutokenforall?clientId=%s&signType=RSA&sign=%s", clientid, strsign)

	r, err := cli.Post(strurl, &RequestOptions{
		Headers: map[string]string{
			"Connection":   "keep-alive",
			"appType":      "4",
			"sourceId":     "203005",
			"Content-Type": "application/json",
		},
		RequestBody: strings.NewReader(jsonbody),
		UserAgent:   fmt.Sprintf("MiguVideo/%s (iPhone; iOS %s; Scale/2.00)", _user.Value_for_key("$APP-VERSION-CODE"), _user.Value_for_key("$systemVersion")),
	})
	if err != nil {
		return nil, New_neterr_with_error(err)
	}
	defer r.Close()
	var rs Login_Migutokenforall_r
	err = r.JSON(&rs)
	if err != nil {
		return nil, err
	}
	if rs.ResultCode != API_LOGIN_SUCCESS {
		return nil, errors.New(rs.ResultDesc)
	}

	_user.Set_key_value("sign", rs.Sign)
	_user.Set_key_value("cityId", rs.UserInfo.CityId)
	_user.Set_key_value("areaId", rs.UserInfo.AreaId)
	_user.Set_key_value("expiredOn", rs.UserInfo.ExpiredOn)
	_user.Set_key_value("mobile", rs.UserInfo.Mobile)
	_user.Set_key_value("passId", rs.UserInfo.PassId)
	_user.Set_key_value("userId", rs.UserInfo.UserId)
	_user.Set_key_value("carrierCode", rs.UserInfo.CarrierCode)
	_user.Set_key_value("userNum", rs.UserInfo.UserNum)
	_user.Set_key_value("userToken", rs.UserInfo.UserToken)
	_user.Set_key_value("userInfo", fmt.Sprintf(`{"areaId":"%s","cityId":"%s","expiredOn":"%s","mobile":"%s","passId":"%s","userId":"%s","carrierCode":"%s","userNum":"%s","userToken":"%s"}`,
		_user.Value_for_key("areaId"),
		_user.Value_for_key("cityId"),
		_user.Value_for_key("expiredOn"),
		_user.Value_for_key("mobile"),
		_user.Value_for_key("passId"),
		_user.Value_for_key("userId"),
		_user.Value_for_key("carrierCode"),
		_user.Value_for_key("userNum"),
		_user.Value_for_key("userToken"),
	))
	return &rs.Login_Migutokenforall_o, nil
}

/*
MGStatic testContent:keyValue:appSign:
获取埋点的token接口
*/
func Crystal_aquireToken(_user IUser, _cli *Session) (*Crystal_aqutokenireToken_o, error) {
	cli := cli_with_cli(_cli)
	//appKey //应该是固定的
	//info->md5("com.wondertek.hecmccmobile,iOS")=6190f274c5f5ff8338107aec4fba6143
	strinfo := hex.EncodeToString(Must_Md5_with_in([]byte(fmt.Sprintf("com.wondertek.hecmccmobile,iOS"))))
	strtimespamp := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	strpd := fmt.Sprintf(`{"appKey":"534B9120BEDF3C2A7713EF40B9D99530","info":"%s","timestamp":"%s"}`, strinfo, strtimespamp)
	strsign := base64.StdEncoding.EncodeToString(Must_Rsa_sign_md5_with_in_privatekey([]byte(hex.EncodeToString(Must_Md5_with_in([]byte(strpd)))), []byte(MGStatic_rsaString_privatekey)))
	strurl := fmt.Sprintf("http://crystal.miguvideo.com/aquireToken?sign=%s", url.QueryEscape(strsign))
	r, err := cli.Post(strurl, &RequestOptions{
		RequestBody: strings.NewReader(strpd),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		UserAgent: "%E5%92%AA%E5%92%95%E8%A7%86%E9%A2%91" + fmt.Sprintf("/%s.5 CFNetwork/808.0.2 Darwin/16.0.0", _user.Value_for_key("$APP-VERSION-CODE")),
	})
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var rs Crystal_aquireToken_r
	err = r.JSON(&rs)
	if err != nil {
		return nil, err
	}
	if rs.ResultCode != "SUCCESS" {
		return nil, errors.New(rs.ResultDesc)
	}
	_user.Set_key_value("crystal_token", rs.Token)
	return &rs.Crystal_aqutokenireToken_o, nil
}

/*
获取地区接口
*/
func Common_flow_service_province_code(_areaName string, _user IUser, _cli *Session) (*Common_flow_service_province_code_o, error) {
	cli := cli_with_cli(_cli)
	strurl := fmt.Sprintf("http://v.miguvideo.com/common/flow-service/province/code?areaName=%s", url.QueryEscape(_areaName))
	r, err := cli.Get(strurl, &RequestOptions{
		Headers: map[string]string{
			"appVersion":             appversion,
			"clientId":               _user.Value_for_key("$FCUUID"),
			"BUSS_ID":                strings.ReplaceAll(strings.ToUpper(gofakeit.UUID()), "-", ""),
			"SDKCEId":                SDKCEId,
			"X-UP-CLIENT-CHANNEL-ID": X_UP_CLIENT_CHANNEL_ID,
			"csessionId":             fmt.Sprintf("%s%s%s", _user.Value_for_key("$FCUUID"), _user.Value_for_key("$idfv"), strings.ToUpper(gofakeit.UUID())),
			"imei":                   _user.Value_for_key("$idfa"),
			"APP-VERSION-CODE":       _user.Value_for_key("$APP-VERSION-CODE"),
			"Phone-Info":             fmt.Sprintf("%s%s%s", strings.ReplaceAll(_user.Value_for_key("$deviceModelName"), " ", "%20"), url.QueryEscape("|"), _user.Value_for_key("$systemVersion")),
		},
		UserAgent: fmt.Sprintf("MiguVideo/%s (iPhone; iOS %s; Scale/2.00)", _user.Value_for_key("$APP-VERSION-CODE"), _user.Value_for_key("$systemVersion")),
	})
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var rs Common_flow_service_province_code_r
	err = r.JSON(&rs)
	if err != nil {
		return nil, err
	}
	if rs.Code != 200 {
		return nil, errors.New(rs.Message)
	}
	return &rs.Common_flow_service_province_code_o, err
}

/*
播放接口
*/
func Playurl_v1_play_url(_contId string, _rateType int, _serialNo string, _clientProvinceCode string, _user IUser, _cli *Session) (*Playurl_v1_play_url_r, error) {
	strtime := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	m5sign, salt := getMiguSign(_contId, appversion, strtime)
	cli := cli_with_cli(_cli)
	strurl := fmt.Sprintf("https://play.miguvideo.com/playurl/v1/play/playurl?audio=false&contId=%s&dolby=false&drm=true&mgdbId=&nt=4&os=%s&ott=false&rateType=%d&salt=%s&serialId=&serialNo=%s&sign=%s&timestamp=%s&ua=%s&vr=true", _contId, _user.Value_for_key("$systemVersion"), _rateType, salt, _serialNo, m5sign, strtime, url.QueryEscape(_user.Value_for_key("$ua")))
	r, err := cli.Get(strurl, &RequestOptions{
		Headers: map[string]string{ //MGBaseRequest mg_header
			"clientProvinceCode":     _user.Value_for_key("ProvinceCode"),                                                                                                                         //这里是一个城市代码，17,服务器返回
			"userToken":              _user.Value_for_key("userToken"),                                                                                                                            //登录服务器返回
			"userId":                 _user.Value_for_key("userId"),                                                                                                                               //登录服务器返回
			"Phone-Info":             fmt.Sprintf("%s%s%s", strings.ReplaceAll(_user.Value_for_key("$deviceModelName"), " ", "%20"), url.QueryEscape("|"), _user.Value_for_key("$systemVersion")), //fmt.sprintf("%s|%s",[NSObject deviceModelName],[[UIDevice currentDevice] systemVersion])  然后urlencode
			"clientId":               _user.Value_for_key("$FCUUID"),                                                                                                                              //[[MGUUID sharedInstance] uuid]->FCUUID uuidForDevice->UICKeyChainStore.keyChainStoreWithService_accessGroup_(NULL,NULL)->keychain.setString_forKey_("bbb","fc_uuidForDevice")
			"BUSS_ID":                strings.ReplaceAll(strings.ToUpper(gofakeit.UUID()), "-", ""),                                                                                               //+[MGBaseRequest getUUIDString]  随机一个uuid然后替换掉-字符,每个请求都不同
			"appVersion":             appversion,
			"X-UP-CLIENT-CHANNEL-ID": X_UP_CLIENT_CHANNEL_ID,          //-[MGGlobalKey CHANNELID]   固定值正常取值 25040600-99000-200300020100001 还有一个可能是低版本25050506-99000-200300020100001
			"clientCityId":           _user.Value_for_key("cityId"),   //由/login/migutokenforall  resultDesc->userInfo返回,服务器生成
			"imei":                   _user.Value_for_key("$idfa"),    //ASIdentifierManager.alloc().init().advertisingIdentifier().UUIDString()  理解为idfa
			"userInfo":               _user.Value_for_key("userInfo"), //由/login/migutokenforall resultDesc->userInfo返回,服务器生成
			"sign":                   _user.Value_for_key("sign"),     //由/login/migutokenforall  resultDesc->sign返回,服务器生成
			"Connection":             "keep-alive",
			"csessionId":             fmt.Sprintf("%s%s%s", _user.Value_for_key("$FCUUID"), _user.Value_for_key("$idfv"), strings.ToUpper(gofakeit.UUID())), //[MGAnaticsID currentUserSessionId]->MGRuntimeInfoHelper currentUserSessionId->[UIDevice currentUserSessionId]->fmt.sprintf("%s%s%s",FCUUID,idfv,newfcuuid)  理解为clientId idfv 然后一个随机uuid，每次启动会变
			"mobile":                 _user.Value_for_key("mobile"),
			"Accept-Language":        "zh-Hans-CN;q=1",
			"Accept":                 "*/*",
			"Accept-Encoding":        "deflate",
			"APP-VERSION-CODE":       _user.Value_for_key("$APP-VERSION-CODE"),
			"SDKCEId":                SDKCEId, //固定写死在app中
		},
		UserAgent: fmt.Sprintf("MiguVideo/%s (iPhone; iOS %s; Scale/2.00)", _user.Value_for_key("$APP-VERSION-CODE"), _user.Value_for_key("$systemVersion")),
	})
	if err != nil {
		return nil, New_neterr_with_error(err)
	}
	defer r.Close()
	var rs Playurl_v1_play_url_r
	err = r.JSON(&rs)
	if err != nil {
		return nil, New_neterr_with_error(err)
	}
	if rs.Code != "200" {
		return nil, errors.New(rs.Message)
	}
	return &rs, nil
}
