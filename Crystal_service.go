package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/a97077088/grequests"
	"net/url"
	"strings"
	"time"
)

type Crystal_service struct {
	user                  IUser
	dataCollectionService DataCollectionService
	dataEventService      DataEventService
	dataExposureService   DataExposureService
}

func New_Crystal_service_with_user_token(_user IUser) *Crystal_service {
	c := &Crystal_service{user: _user}
	return c
}

func (this *Crystal_service) Get_crystal_token() string {
	return this.user.Value_for_key("crystal_token")
}
func (this *Crystal_service) Getsessionid() string {
	fcuuid := this.user.Value_for_key("$FCUUID")
	vts := this.user.Env_for_key("sessionts")
	var ts = int64(0)
	if vts != nil {
		ts = vts.(int64)
	}
	if time.Since(time.Unix(ts/1000, 0)).Seconds() >= 1 {
		//session过期是300秒
		ts = time.Now().UnixNano() / 1e6
		this.user.Set_env_value("sessionts", ts)
	}
	return fmt.Sprintf("%s%d", fcuuid, ts)
}
func (this *Crystal_service) DeviceInfoJson() map[string]interface{} {
	ts := time.Now().UnixNano() / 1e6
	rjson := map[string]interface{}{
		"imei":           this.user.Value_for_key("$FCUUID"),
		"udid":           this.user.Value_for_key("$FCUUID"),
		"installationID": this.user.Value_for_key("$FCUUID"),
		"phoneMode":      this.user.Value_for_key("$deviceModelName"),
		"phoneBrand":     "apple",
		"idfa":           this.user.Value_for_key("$idfa"),
		"idfv":           this.user.Value_for_key("$idfv"),
		"appVersion":     this.user.Value_for_key("$APP-VERSION-CODE"),
		"apppkg":         "com.wondertek.hecmccmobile",
		"os":             "iOS",
		"appchannel":     X_UP_CLIENT_CHANNEL_ID,
		"userId":         "",
		"osversion":      this.user.Value_for_key("$systemVersion"),
		"sdkversion":     sdkversion,
		"uploadTs":       ts,
	}
	return rjson
}
func (this *Crystal_service) SdkSessionInfo() map[string]interface{} {
	rjson := this.DeviceInfoJson()
	rjson["networkType"] = "WIFI"
	rjson["promotion"] = ""
	rjson["accountType"] = ""
	rjson["sessionId"] = this.Getsessionid()
	rjson["clientId"] = this.user.Value_for_key("$FCUUID")
	rjson["account"] = ""
	rjson["MG_SCORE_TIME"] = fmt.Sprintf("%d", rjson["uploadTs"].(int64))
	rjson["sdkpkg"] = ""
	return rjson
}
func (this *Crystal_service) Data_tojsonstring(_e interface{}) string {
	bt, err := json.Marshal(this)
	if err != nil {
		panic(err)
	}
	return string(bt)
}

func (this *Crystal_service) DataCollectionService_upload_with_datacollectionservice_cli(_d *DataCollectionService, _cli *Session) error {
	cli := cli_with_cli(_cli)
	s := this.Data_tojsonstring(_d)
	strsign := url.QueryEscape(base64.StdEncoding.EncodeToString(Must_Rsa_sign_md5_with_in_privatekey(Must_Md5_with_in([]byte(s)), []byte(MGStatic_rsaString_privatekey))))
	strurl := fmt.Sprintf("http://crystal.miguvideo.com/legacy/shm_video_interface/dataCollectionService?sign=%s", strsign)
	r, err := cli.Post(strurl, &RequestOptions{
		RequestBody: strings.NewReader(s),
		Headers: map[string]string{
			"Content-Type": "application/json",
			"token":        this.Get_crystal_token(),
		},
		UserAgent: fmt.Sprintf("MiguVideo/%s (iPhone; iOS %s; Scale/2.00)", this.user.Value_for_key("$APP-VERSION-CODE"), this.user.Value_for_key("$systemVersion")),
	})
	if err != nil {
		return New_neterr_with_error(err)
	}
	defer r.Close()
	var rs Api_R
	err = r.JSON(&rs)
	if err != nil {
		return New_neterr_with_error(err)
	}
	if rs.ResultCode != API_SUCCESS {
		return errors.New(rs.ResultDesc)
	}
	return nil
}
func (this *Crystal_service) DataCollectionService_session_start() *DataCollectionService {
	ts := time.Now().UnixNano() / 1e6
	rjson := map[string]interface{}{
		"os":                    "iOS",
		"imei":                  this.user.Value_for_key("$FCUUID"),
		"imsi":                  "",
		"appPackageName":        "com.wondertek.hecmccmobile",
		"idfa":                  this.user.Value_for_key("$idfa"),
		"idfv":                  this.user.Value_for_key("$idfv"),
		"currentAppVersionCode": this.user.Value_for_key("$APP-VERSION-CODE"),
		"currentOSVersion":      this.user.Value_for_key("$systemVersion"),
		"currentAppVersionName": "咪咕视频",
		"phoneNumber":           "(null)",
		"sessionId":             this.Getsessionid(),
		"udid":                  this.user.Value_for_key("$FCUUID"),
		"userId":                "",
		"account":               "",
		"startTs":               fmt.Sprintf("%d", ts),
		"clientId":              this.user.Value_for_key("$FCUUID"),
	}
	d := New_DataCollectionService_with_sdksessioninfo_deviceinfo_exception_custominfo_sessionstart_sessionend(
		this.SdkSessionInfo(),
		this.DeviceInfoJson(),
		nil,
		nil,
		rjson,
		nil)
	return d
}
