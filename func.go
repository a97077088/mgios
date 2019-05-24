package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"strconv"
	"strings"
	"time"
)

var sdkversion = "MGSDK3.2.5"
var appversion = "25040600"                                   //固定写死在app中
var SDKCEId = "7a9582bd-8660-4f1a-9d0f-8451f688c67b"          ///固定
var X_UP_CLIENT_CHANNEL_ID = "25040600-99000-200300020100001" //-[MGGlobalKey CHANNELID]   固定值正常取值 25040600-99000-200300020100001 还有一个可能是低版本25050506-99000-200300020100001

func getMiguSign(_cid string, _version string, _timestamp string) (string, string) {
	_s := fmt.Sprintf("%s%s%s", _timestamp, _cid, _version)
	md5val := hex.EncodeToString(Must_Md5_with_in([]byte(_s)))
	sale := fmt.Sprintf("%d", gofakeit.Number(10000000, 99999999))
	d4 := sale[0:4]
	d2 := sale[6:8]
	id2, err := strconv.Atoi(d2)
	if err != nil {
		return "", ""
	}
	sindex := GetMiguSign_arrindex[id2]
	ends := fmt.Sprintf("%s%s%s%s", md5val, sindex, "migu", d4)
	md5ends := hex.EncodeToString(Must_Md5_with_in([]byte(ends)))
	return md5ends, sale
}
func decodeauthfield_with_hexs_privatekey(_hexs string, _privatekey []byte) (string, error) {
	bys, err := hex.DecodeString(_hexs)
	if err != nil {
		return "", nil
	}
	byr, err := Rsa_decode_with_in_privatekey(bys, _privatekey)
	if err != nil {
		return "", nil
	}
	return string(byr), nil
}

//[IDMPThirdLoginMode checkKSIsValid:username:passWd:clientNonce:]
//sub_100248A64
func get_ks_with_auth_user(_mpparam map[string]string, _password string, _user IUser) ([]byte, error) {
	strtop := "MIGU_LOGIN_PW"
	upnonce := _mpparam["UP Nonce"]

	arg3 := fmt.Sprintf("%s%s", _mpparam["reqtime"], _user.Value_for_key("$MGSDK_UUID"))
	arg3 = hex.EncodeToString(Must_Md5_with_in([]byte(arg3)))
	//fmt.Printf("第一步md5:%s\n",arg3)

	var rs bytes.Buffer
	rs.WriteByte(0x01)
	rs.WriteString(strtop)
	n := int16(len(strtop))
	err := binary.Write(&rs, binary.BigEndian, n)
	if err != nil {
		return nil, err
	}
	rs.WriteString(upnonce)
	n = int16(len(upnonce))
	err = binary.Write(&rs, binary.BigEndian, n)
	if err != nil {
		return nil, err
	}
	rs.WriteString(arg3)
	n = int16(len(arg3))
	err = binary.Write(&rs, binary.BigEndian, n)
	if err != nil {
		return nil, err
	}
	fristkey := fmt.Sprintf("migupassport:%s:%s", _mpparam["username"], _password)
	fristkey = hex.EncodeToString(Must_Md5_with_in([]byte(fristkey)))
	//fmt.Printf("第一步hmackey:%s\n",fristkey)
	hmacrs := Hmac_sha256_with_in_key(rs.Bytes(), []byte(fristkey))
	return hmacrs[:16], nil
}

//详细算法sub_100369414
func getTokenWithUserName_appId_uuid_andUser_(_username string, _appid string, _uuid string, _anduser map[string]string) ([]byte, error) {
	//fmt.Printf("拿到ks:%s\n",hex.EncodeToString([]byte(_anduser["KS"])))
	strtop := "migu-token-version3.0"
	//分割后获取BTID的@前部分
	strbtidtop := strings.Split(_anduser["BTID"], "@")[0]
	strpassid := _anduser["passid"] //"73378307247815"
	var bf bytes.Buffer
	bf.WriteByte(0x01)
	bf.WriteString(strtop)
	err := binary.Write(&bf, binary.BigEndian, int16(len(strtop)))
	if err != nil {
		return nil, err
	}
	bf.WriteString(strbtidtop)
	err = binary.Write(&bf, binary.BigEndian, int16(len(strbtidtop)))
	if err != nil {
		return nil, err
	}
	bf.WriteString(strpassid)
	err = binary.Write(&bf, binary.BigEndian, int16(len(strpassid)))
	if err != nil {
		return nil, err
	}
	bf.WriteString(_appid)
	err = binary.Write(&bf, binary.BigEndian, int16(len(_appid)))
	if err != nil {
		return nil, err
	}
	//此处拿到第一次的lastkey
	lastkey1 := []byte(hex.EncodeToString(Sha256_with_in(bf.Bytes())[:16]))
	//fmt.Printf("拿到sha256key:%s\n",string(lastkey1))
	bf.Write(lastkey1)

	key1 := _anduser["KS"]
	key2 := Hmac_sha256_with_in_key(bf.Bytes(), []byte(key1))[:16]
	//fmt.Printf("拿到token_hmackey1:%s\n",hex.EncodeToString(key2))

	isqn, err := strconv.Atoi(_anduser["sqn"])
	if err != nil {
		return nil, err
	}
	isqn++ //sqn要+1
	//详细算法sub_10036998C
	var bf2 bytes.Buffer
	//1
	bf2.WriteByte(0x01)
	bf2.WriteByte(0x00)
	bf2.WriteByte(0x01)
	bf2.WriteByte(0x33)
	//2
	bf2.WriteByte(0x2)
	bf2.WriteByte(0x00)
	bf2.WriteByte(byte(len(_anduser["BTID"])))
	bf2.WriteString(_anduser["BTID"])
	//3
	bf2.WriteByte(0x3)
	bf2.WriteByte(0x00)
	bf2.WriteByte(0x04)
	err = binary.Write(&bf2, binary.BigEndian, int32(isqn))
	if err != nil {
		return nil, err
	}
	//fmt.Println(hex.EncodeToString(bf2.Bytes()))
	//4
	bf2.WriteByte(0x4)
	bf2.WriteByte(0x00)
	bf2.WriteByte(byte(len(_appid)))
	bf2.WriteString(_appid)
	//5
	bf2.WriteByte(0x5)
	bf2.WriteByte(0x00)
	bf2.WriteByte(0x20) //写入32位
	bf2.WriteString(_uuid)

	hmacsha256sign := Hmac_sha256_with_in_key(bf2.Bytes(), key2)
	//fmt.Printf("拿到hmac r:%s\n",hex.EncodeToString(hmacsha256sign))

	//FF
	bf2.WriteByte(0xFF)
	bf2.WriteByte(0x00)
	bf2.WriteByte(0x20)
	bf2.Write(hmacsha256sign)

	//写入头部
	var rs bytes.Buffer
	rs.WriteByte(0x84)
	rs.WriteByte(0x84)
	rs.Write(bf2.Bytes())

	//fmt.Printf("migutoken:%s\n",hex.EncodeToString(rs.Bytes()))
	return rs.Bytes(), nil
}

func SessionStart_with_user(_user IUser) map[string]interface{} {
	ts := time.Now().UnixNano() / 1e6
	rjson := map[string]interface{}{
		"os":                    "iOS",
		"imei":                  _user.Value_for_key("$FCUUID"),
		"imsi":                  "",
		"appPackageName":        "com.wondertek.hecmccmobile",
		"idfa":                  _user.Value_for_key("$idfa"),
		"idfv":                  _user.Value_for_key("$idfv"),
		"currentAppVersionCode": _user.Value_for_key("$APP-VERSION-CODE"),
		"currentOSVersion":      _user.Value_for_key("$systemVersion"),
		"currentAppVersionName": "咪咕视频",
		"phoneNumber":           "(null)",
		"sessionId":             Getsessionid_with_user(_user),
		"udid":                  _user.Value_for_key("$FCUUID"),
		"userId":                "",
		"account":               "",
		"startTs":               fmt.Sprintf("%d", ts),
		"clientId":              _user.Value_for_key("$FCUUID"),
	}
	return rjson
}
func DeviceInfoJson_with_user(_user IUser) map[string]interface{} {
	ts := time.Now().UnixNano() / 1e6
	rjson := map[string]interface{}{
		"imei":           _user.Value_for_key("$FCUUID"),
		"udid":           _user.Value_for_key("$FCUUID"),
		"installationID": _user.Value_for_key("$FCUUID"),
		"phoneMode":      _user.Value_for_key("$deviceModelName"),
		"phoneBrand":     "apple",
		"idfa":           _user.Value_for_key("$idfa"),
		"idfv":           _user.Value_for_key("$idfv"),
		"appVersion":     _user.Value_for_key("$APP-VERSION-CODE"),
		"apppkg":         "com.wondertek.hecmccmobile",
		"os":             "iOS",
		"appchannel":     X_UP_CLIENT_CHANNEL_ID,
		"userId":         "",
		"osversion":      _user.Value_for_key("$systemVersion"),
		"sdkversion":     sdkversion,
		"uploadTs":       ts,
	}
	return rjson
}
func SdkSessionInfo_with_user(_user IUser) map[string]interface{} {
	rjson := DeviceInfoJson_with_user(_user)
	rjson["networkType"] = "WIFI"
	rjson["promotion"] = ""
	rjson["accountType"] = ""
	rjson["sessionId"] = Getsessionid_with_user(_user)
	rjson["clientId"] = _user.Value_for_key("$FCUUID")
	rjson["account"] = ""
	rjson["MG_SCORE_TIME"] = fmt.Sprintf("%d", rjson["uploadTs"].(int64))
	rjson["sdkpkg"] = ""
	return rjson
}
func Getsessionid_with_user(_user IUser) string {
	fcuuid := _user.Value_for_key("$FCUUID")
	vts := _user.Env_for_key("sessionts")
	var ts = int64(0)
	if vts != nil {
		ts = vts.(int64)
	}
	if time.Since(time.Unix(ts/1000, 0)).Seconds() >= 1 {
		//session过期是300秒
		ts = time.Now().UnixNano() / 1e6
		_user.Set_env_value("sessionts", ts)
	}
	return fmt.Sprintf("%s%d", fcuuid, ts)
}
func Get_crystal_token_with_user(_user IUser) string {
	return _user.Value_for_key("crystal_token")
}
