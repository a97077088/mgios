package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/a97077088/grequests"
	"net/url"
	"strings"
)

type DataCollectionService struct {
	user           IUser
	SdkSessionInfo map[string]interface{}   `json:"sdkSessionInfo"`
	DeviceInfo     map[string]interface{}   `json:"deviceInfo"`
	Exception      []map[string]interface{} `json:"exception"`
	CustomInfo     []map[string]interface{} `json:"customInfo"`
	SessionStart   map[string]interface{}   `json:"sessionStart"`
	SessionEnd     map[string]interface{}   `json:"sessionEnd"`
}

func (this *DataCollectionService) To_jsonstring() string {
	this.SdkSessionInfo = SdkSessionInfo_with_user(this.user)
	bt, err := json.Marshal(this)
	if err != nil {
		panic(err)
	}
	return string(bt)
}
func (this *DataCollectionService) Upload(_cli *Session) error {
	cli := cli_with_cli(_cli)
	s := this.To_jsonstring()
	strsign := url.QueryEscape(base64.StdEncoding.EncodeToString(Must_Rsa_sign_md5_with_in_privatekey(Must_Md5_with_in([]byte(s)), []byte(MGStatic_rsaString_privatekey))))
	strurl := fmt.Sprintf("http://crystal.miguvideo.com/legacy/shm_video_interface/dataCollectionService?sign=%s", strsign)
	r, err := cli.Post(strurl, &RequestOptions{
		RequestBody: strings.NewReader(s),
		Headers: map[string]string{
			"Content-Type": "application/json",
			"token":        Get_crystal_token_with_user(this.user),
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

func New_DataCollectionService_with_user_deviceinfo_exception_custominfo_sessionstart_sessionend(
	_user IUser,
	_deviceinfo map[string]interface{},
	_exception []map[string]interface{},
	_custominfo []map[string]interface{},
	_sessionstart map[string]interface{},
	_sessionend map[string]interface{}) *DataCollectionService {
	d := &DataCollectionService{}
	d.user = _user
	d.DeviceInfo = _deviceinfo
	if d.DeviceInfo == nil {
		d.DeviceInfo = map[string]interface{}{}
	}
	d.Exception = _exception
	if d.Exception == nil {
		d.Exception = []map[string]interface{}{}
	}
	d.CustomInfo = _custominfo
	if d.CustomInfo == nil {
		d.CustomInfo = []map[string]interface{}{}
	}
	d.SessionStart = _sessionstart
	if d.SessionStart == nil {
		d.SessionStart = map[string]interface{}{}
	}
	d.SessionEnd = _sessionend
	if d.SessionEnd == nil {
		d.SessionEnd = map[string]interface{}{}
	}
	return d
}
