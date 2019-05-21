package main

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	. "github.com/a97077088/grequests"
	"github.com/brianvoe/gofakeit"
	"math/big"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var appversion = "25040600"                                   //固定写死在app中
var SDKCEId = "7a9582bd-8660-4f1a-9d0f-8451f688c67b"          ///固定
var X_UP_CLIENT_CHANNEL_ID = "25040600-99000-200300020100001" //-[MGGlobalKey CHANNELID]   固定值正常取值 25040600-99000-200300020100001 还有一个可能是低版本25050506-99000-200300020100001
var pubkey = ""
var signature_privatekey = ""
var rsasign_privatekey = ""
var uadecode_privatekey = ""

//去补码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:length-unpadding]
}

//使用PKCS7进行填充，IOS也是7
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func gorsa_decode(_in []byte, _privatekey []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(_privatekey)
	if block == nil {
		return nil, errors.New("解析pem失败")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	b, err := rsa.DecryptPKCS1v15(rand.Reader, priv, _in)
	if err != nil {
		return nil, err
	}
	return b, nil
}
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func gorsasha256_with_private(_in []byte, _privatekey []byte) string {
	//解密pem格式的公钥
	block, _ := pem.Decode([]byte(_privatekey))
	if block == nil {
		return ""
	}
	// 解析公钥
	privatekey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return ""
	}
	rng := rand.Reader
	hashed := sha256.Sum256(_in)
	signature, err := rsa.SignPKCS1v15(rng, privatekey, crypto.SHA256, hashed[:])
	if err != nil {
		return ""
	}
	signatureHex := hex.EncodeToString(signature)
	return signatureHex

}
func goaes128_with_byte(_byte []byte, _key []byte) string {
	block, _ := aes.NewCipher(_key)
	blockSize := block.BlockSize()
	origData := PKCS7Padding(_byte, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted)
}
func goaes128decode_with_base64(_s string, _key []byte) []byte {
	_crypted, err := base64.StdEncoding.DecodeString(_s)
	if err != nil {
		return []byte("")
	}
	block, _ := aes.NewCipher(_key)
	blockMode := cipher.NewCBCDecrypter(block, []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
	origData := make([]byte, len(_crypted))
	blockMode.CryptBlocks(origData, _crypted)
	origData = PKCS7UnPadding(origData)
	return origData
}

//md5实现
func gomd5(_s string) string {
	h := md5.New()
	h.Write([]byte(_s))
	return hex.EncodeToString(h.Sum(nil))
}

// rsa加密
func gorsaenc_with_publickey(origData []byte, _pubkey string) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode([]byte(_pubkey))
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

var arrindex []string

func gorsaenc_with_private(data []byte, _pritekey []byte) string {
	e := data
	block, _ := pem.Decode(_pritekey)
	privkey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	encData, _ := gorsaenc_with_privatekey_data_imp(privkey, e)
	return base64.StdEncoding.EncodeToString(encData)
}
func gorsaenc_with_privatekey_data_imp(priv *rsa.PrivateKey, data []byte) (enc []byte, err error) {
	var (
		ErrInputSize  = errors.New("input size too large")
		ErrEncryption = errors.New("encryption error")
	)
	k := (priv.N.BitLen() + 7) / 8
	tLen := len(data)
	// rfc2313, section 8:
	// The length of the data D shall not be more than k-11 octets
	if tLen > k-11 {
		err = ErrInputSize
		return
	}
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-tLen-1; i++ {
		em[i] = 0xff
	}
	copy(em[k-tLen:k], data)
	c := new(big.Int).SetBytes(em)
	if c.Cmp(priv.N) > 0 {
		err = ErrEncryption
		return
	}
	var m *big.Int
	var ir *big.Int
	if priv.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, priv.D, priv.N)
	} else {
		// We have the precalculated values needed for the CRT.
		m = new(big.Int).Exp(c, priv.Precomputed.Dp, priv.Primes[0])
		m2 := new(big.Int).Exp(c, priv.Precomputed.Dq, priv.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, priv.Primes[0])
		}
		m.Mul(m, priv.Precomputed.Qinv)
		m.Mod(m, priv.Primes[0])
		m.Mul(m, priv.Primes[1])
		m.Add(m, m2)

		for i, values := range priv.Precomputed.CRTValues {
			prime := priv.Primes[2+i]
			m2.Exp(c, values.Exp, prime)
			m2.Sub(m2, m)
			m2.Mul(m2, values.Coeff)
			m2.Mod(m2, prime)
			if m2.Sign() < 0 {
				m2.Add(m2, prime)
			}
			m2.Mul(m2, values.R)
			m.Add(m, m2)
		}
	}

	if ir != nil {
		// Unblind.
		m.Mul(m, ir)
		m.Mod(m, priv.N)
	}
	enc = m.Bytes()
	return
}
func getMiguSign(_cid string, _version string, _timestamp string) (string, string) {
	_s := fmt.Sprintf("%s%s%s", _timestamp, _cid, _version)
	md5val := gomd5(_s)
	sale := fmt.Sprintf("%d", gofakeit.Number(10000000, 99999999))
	d4 := sale[0:4]
	d2 := sale[6:8]
	id2, err := strconv.Atoi(d2)
	if err != nil {
		return "", ""
	}
	sindex := arrindex[id2]
	ends := fmt.Sprintf("%s%s%s%s", md5val, sindex, "migu", d4)
	md5ends := gomd5(ends)
	return md5ends, sale
}
func decodeauthfield_with_hexs_privatekey(_hexs string, _privatekey []byte) (string, error) {
	bys, err := hex.DecodeString(_hexs)
	if err != nil {
		return "", nil
	}
	byr, err := gorsa_decode(bys, _privatekey)
	if err != nil {
		return "", nil
	}
	return string(byr), nil
}
func gohmac_sha256(_key []byte, _in []byte) []byte {
	h := hmac.New(sha256.New, _key)
	h.Write(_in)
	return h.Sum(nil)
}
func gosha256(_in []byte) []byte {
	h := sha256.New()
	h.Write(_in)
	return h.Sum(nil)
}

//[IDMPThirdLoginMode checkKSIsValid:username:passWd:clientNonce:]
//sub_100248A64
func get_ks_with_auth_user(_mpparam map[string]string, _password string, _user IUser) ([]byte, error) {
	strtop := "MIGU_LOGIN_PW"
	upnonce := _mpparam["UP Nonce"]

	arg3 := gomd5(fmt.Sprintf("%s%s", _mpparam["reqtime"], _user.Value_for_key("$MGSDK_UUID")))
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
	fristkey := gomd5(fmt.Sprintf("migupassport:%s:%s", _mpparam["username"], _password))
	//fmt.Printf("第一步hmackey:%s\n",fristkey)
	hmacrs := gohmac_sha256([]byte(fristkey), rs.Bytes())
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
	lastkey1 := []byte(hex.EncodeToString(gosha256(bf.Bytes())[:16]))
	//fmt.Printf("拿到sha256key:%s\n",string(lastkey1))
	bf.Write(lastkey1)

	key1 := _anduser["KS"]
	key2 := gohmac_sha256([]byte(key1), bf.Bytes())[:16]
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

	hmacsha256sign := gohmac_sha256(key2, bf2.Bytes())
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

func cli_with_cli(_cli *Session) *Session {
	if _cli == nil {
		_cli = NewSession(nil)
	}
	return _cli
}
func playurl_v1_play_url(_contId string, _rateType int, _serialNo string, _clientProvinceCode string, _user IUser, _cli *Session) (*Playurl_v1_play_url_r, error) {
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

//+[IDMPUPMode getUPKSByUserName:andPassWd:sessionid:validcode:rcflag:tvflag:validtype:successBlock:failBlock:]
//sub_1002F5CA4
//[MGSecondAuthAPI requestUrl]
//总结一下migutoken的算法,简单几句话,
// 第一步拿到ks作为下一步计算密钥
// 第二步计算一下sha256作为下一步key，然后hmacsha256作为下一步key，反正一直往下推
func login_migutokenforall(_rauth map[string]string, _password string, _user IUser, _cli *Session) (*Login_Migutokenforall_o, error) {
	cli := cli_with_cli(_cli)

	passid, err := decodeauthfield_with_hexs_privatekey(_rauth["passid"], []byte(uadecode_privatekey))
	if err != nil {
		return nil, err
	}
	sqn, err := decodeauthfield_with_hexs_privatekey(_rauth["sqn"], []byte(uadecode_privatekey))
	if err != nil {
		return nil, err
	}
	msisdn, err := decodeauthfield_with_hexs_privatekey(_rauth["msisdn"], []byte(uadecode_privatekey))
	if err != nil {
		return nil, err
	}
	btid, err := decodeauthfield_with_hexs_privatekey(_rauth["BTID"], []byte(uadecode_privatekey))
	if err != nil {
		return nil, err
	}
	username, err := decodeauthfield_with_hexs_privatekey(_rauth["username"], []byte(uadecode_privatekey))
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
	md5jsonbody := gomd5(jsonbody)
	//这一步实现了java中的rsa_privateenc 真几把不好实现这个耗时一下午
	sign := gorsaenc_with_private([]byte(md5jsonbody), []byte(rsasign_privatekey))
	strurl := fmt.Sprintf("https://api.miguvideo.com/login/migutokenforall?clientId=%s&signType=RSA&sign=%s", clientid, sign)

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

func client_authRequest(_username string, _password string, _user IUser, _cli *Session) (map[string]string, error) {
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
	encrcdata := goaes128_with_byte([]byte(rcdata), []byte("15acb4a88285ed2f"))
	tm := time.Now().Unix()                                                           //得到10位时间戳
	md5fmcnonce := gomd5(fmt.Sprintf("%d%s", tm, _user.Value_for_key("$MGSDK_UUID"))) //10位时间戳+[IDMPDevice getDeviceID]
	btnonce, err := gorsaenc_with_publickey([]byte(md5fmcnonce), pubkey)
	if err != nil {
		return nil, err
	}
	nonce := hex.EncodeToString(btnonce)
	btpwd, err := gorsaenc_with_publickey([]byte(_password), pubkey)
	if err != nil {
		return nil, err
	}
	pwd := hex.EncodeToString(btpwd)
	msgid := gomd5(fmt.Sprintf("%s%s", time.Now().Format("20060102150405000"), _user.Value_for_key("$MGSDK_UUID")))                                                                                                                                                                                                                                   //MGSDKRequest msgId
	Authorization := fmt.Sprintf(`UP clientversion="1.0",sdkversion="MGUnionAuthSDK-iOSV4.6.0",appid="20300507",apptype="4",username="%s",isusernamecn="0",enccnonce="%s",encpasswd="%s",BTID="",Phone_ID="%s",msgid="%s",rcflag="1",sdkverifytype="2",secflag="1",tvflag="false"`, _username, nonce, pwd, _user.Value_for_key("$MGSDK_UUID"), msgid) //`UP clientversion="1.0",sdkversion="MGUnionAuthSDK-iOSV4.6.0",appid="20300507",apptype="4",username="15093717639",isusernamecn="0",enccnonce="3be146a9e42d0e54106485d9b274702ed0c5b25dd0289622ab162d9f439294ce28b58bc3076d6870e7d4058d6dc272070c41d79ca4dd06e414ec2567314b577a1f228d57eacd1d989a44665406de6e4f5ece3576f27a89cb3c52f03c466237fc2aefd33502da5b986fba7dc488fe15873827e21eab02e5600ba73118ddc4113b",encpasswd="8e9a1211155df0b57de28901d81a26e3cdb6ce0904a3154916218a92a51b7d3eeb2dfd45495f5bb0a1089528d0d2dc89da1f3a7182fbc79c0bfa3c4f4da6df8b2c1de22ce5fa4cb66c031520f325a3fff6f296d681a6d3f3e8d142f1fa6a2ae73b11f479c34a6fce1bf73f38cc7c394d2abc124210207763823e3658ba7d3aeb",BTID="",Phone_ID="D876D0297EB04678A6EAC0EE0D1F0FCD",msgid="282a520abec20255f3fe9033f1260b95",rcflag="1",sdkverifytype="2",secflag="1",tvflag="false"`,
	//IDMPRSA_Encrypt_Decrypt addPrivateKey:
	//SecKeyRawSign
	signature := gorsasha256_with_private([]byte(Authorization), []byte(signature_privatekey))
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

//初始化密钥,定位特征
//ApiHandler requestPlayURLWithParams:success:failure:
//MGTEncryptWrapper getEncryptWithTimestamp:
//进两层off_指向
func init() {
	arrindex = []string{
		"57e7f108bfe849d1b7cdfe8ee09d253c",
		"e20c2f9c8007461b823b568f7e2f1df8",
		"de6d4d4e764f4bd9b9d3c61482dc0e19",
		"aeba4e6604884a3aad4c9c1f8e19927d",
		"6dc9ecf9c4414f1a8ba43f340105bc5b",
		"1d88a1f67977431f8a0bd3e37ed458a0",
		"0b5b7a1418d040faad7cbfa029c91ce5",
		"42984878daf74ca7b827ededee590953",
		"117845d9550446bfa4f2c158cdd1fd73",
		"e3eb9ae03571433a99ea5099bce0fc80",
		"7fdb3be24a8a472982c871297f06cd84",
		"0dc13af15fcb4f59bf6fa73979c3922b",
		"11e72b06381d4e11b2d09138c5f438a7",
		"d63cbbafdadc468eaf08ed41eeb9a005",
		"882578c312224526b92ccf4ce8adbbeb",
		"1d0f73d53f4444488884bffb33632076",
		"17261d6d179940b48a0bcaed5d866dfe",
		"ccfb0cac55ef4a66bf1570bdc993a112",
		"49ca01d63676425a8a4da851237c3746",
		"e925ec0094534595995c7cc0ee5d497e",
		"8e1aa9aade3f44df985ed674df81636f",
		"7f163c264789422fbb4e94601e57b5d2",
		"7a0708a4a1514858a70d2925f9e652f6",
		"fd570de139d94a82af8534fd718e77f5",
		"f18a099403394af5b6704edebdefc808",
		"531010e9b1f84023890a5b117abdf5cd",
		"b85de9ea18a94270ab000a6bf470f3d2",
		"c482d9f60b8240909caa1ff93eaa614a",
		"b96a5c3c506c41079fb44974dab80692",
		"af1c58d2c4ef49d19c2dbc8999bc689a",
		"da22fc905b9c4253abfb896725b45563",
		"1c88a1a93b1d42ed87a67f6aa515d4c9",
		"7c3fa7851f21495d844d4d3ed10b7b3e",
		"cfc38b72bd0244839ccdc144f7884251",
		"8a5431431581404f906c8fbb61bb509f",
		"0114c865b5b143c69e665297d09cde21",
		"da0ff36022c94d82ae32b861c6203d0a",
		"deb557c0976c48c59fd2a8d597d92e78",
		"d5965cd090184a78a60ef4805055fbd9",
		"27a1a6c25862493ebec83e1fcb8783bd",
		"5e73077f1ce14d7cac6099d97f360c1d",
		"d52695923cc4467abc52e354560a7cb0",
		"a71e25f15165489ba16c34c2e45c435e",
		"41b228fc1acd4aa68f7669bfefaeceb8",
		"c29c6e07bc324670a6404362c399f7d4",
		"8da92c6775564d7f87ee1df9e367656b",
		"1c35312f70ed4291a686e1ffb11b0426",
		"07a423e0b15d43bdaadb0d538bdd0bda",
		"403c8c4730d042ca965bebd81d079153",
		"1d068f637c1a479886922ce6763a0c4f",
		"ba569bf27cd24a6faa3fc3478c8e87e4",
		"d02e5299c8d243b3a261abe72b979ab3",
		"4103d0198b5d4fef9985a5bd34d8122a",
		"8705d4615a824c9790d4d25792bfdc82",
		"73f42fe46c204b05a0d5222bc3e433fb",
		"ba78f5e849c240e48b60f9e49f8dd1be",
		"a6acc8e4942248c89cdd8ff4cd1b8aa0",
		"af69737f912a43c984b7c4f5ddbe7c99",
		"92eaec15b90e4c05a76bdce7f3bd2637",
		"06fb6827b8244cf495620699e1aa6eb1",
		"05c4e01d1e904cc18609c1ee3fe43a53",
		"9070d9eb650841c6b14f9868593ddd16",
		"4c8065d58f424ff6ab279c53557f1030",
		"1b3b1143d5ef472596d5988523229ee3",
		"d97918056fee43d9a1eaa02fd0d3b1cf",
		"82eafa984900413d91bd45472853c972",
		"7f01eab066754c689fdd05e46e6d825c",
		"34790025cb6f46bcae3fe7f9bd5a8489",
		"2f5011fb2a2648e5883774920c373a34",
		"72312fbeb2e444f8941fd650678ab4b6",
		"55f7860c0af04f29bc6f8a27762dd201",
		"4c963cb27bec45d3a81a47ae0bbe1e8e",
		"6d83f57cee6c408f819b8a7c89985b26",
		"56457f86ec2f43b0a6ff8de6e3d1bb05",
		"ffd8bd316a6948eca113bd232961c1c3",
		"315a49e821ee4128b992716cea3f65c3",
		"8a64cac647c5443ebc3c081e64695587",
		"841c9d63a887435fb6264a388fc2ca28",
		"0cd69540164c4c689299166b7774952b",
		"7ef85021241745a288641914587befa5",
		"e494dd0e4f4f42b49f027d3cc87b9d8f",
		"b7a2484c210540fa951c0e74e8f4abda",
		"eeaf6717c7474349afb76f8051ac7466",
		"342634215c904533827ca9558776622c",
		"0c3680ca49ad4d91a837f413a9ae6147",
		"8524acd7e2b64137a9d218178b63528b",
		"25be331bcc034a53bb698061b293844a",
		"b41d563d48ce483aae4f84491f30f910",
		"d5fe1432988e4dd5a0a0aad27f47e786",
		"2b5edde5b2cc4ea88be381647f5d38a4",
		"dee6a59c81ec4130af1652db14efdb77",
		"efea4dc7a8e4460a995fa55f19f73227",
		"b24e6adb29b346d8a3022d596ec9392f",
		"6a91a2649bd142cb89fd74bc45c04684",
		"db7014ddd71f4de58f1c308e6e79a312",
		"d120e89a5fc74228b70bb10e192ae325",
		"712ede6028b640d6b92f33d5694e089a",
		"77177921d86e4d4394bac284e7baef3b",
		"0001bce083b348d6a50aa1a1149d908f",
		"0e41c2b6b01a46b6952affb9f78ac2b8",
	}
	pubkey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDMQK5vjD6+FRf8/RQDgeyg8T/G
vLo3CD7ngcZ0l2fXuHh83gxgM1C7BNN1RgQFRXWFzMDLxbynm/ZA3GsntY+OiN8S
Bza2SEnir+p0jrsZ+WNhaRWch4FpAxnfcnfxrJL26ZkfYaxS85pKby/IQuSwIswV
om0F06RIZreK3rQDrwIDAQAB
-----END PUBLIC KEY-----`
	signature_privatekey = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQDN2CQ+ZcsbHVrLSJ+/v0BA6kgoa+riowfZErlGdlMjUhf4UTos
1ZzKfavQjctIQjlEStxBWAZ+gtmtYakUmjWFBupks6RWcR/ho+VEePpzxhSis8M6
zHDhf9Dl8u6uvSAdk+CM3NgGle0bC8WnMxy2k9k3MSrX93AKZV36DCMtUQIDAQAB
AoGACBp+MsanHEYnkOEnCNFqoiOW+6Bj+tAYOv91s8RsuXM95lSsSZ+PMJmJ7gfm
/M0+m+Wmjhv9BXX5Q84Ybes0OBpS9qk2Rq6mQooXqo+6BaPlJb/UD160ULiQJIA7
P2x2XX/Z8xW7goq6r13i4VOZj4GHxRnlfvtCyKKso2U4qMECQQD2bBrY+1O1Du7D
B6jjJC0rnZMQllIlex9zftvIgr69qvQcWoCKukBUjmx8gw7Pall33sRgJbhMj0yf
qLBn6GplAkEA1dhJcJhZ9euFsjWnjXG5g5CyRO5HPJ9Kfoqw04aQX7oHHRXUyVNU
H4L9FsbXQBYS+rXrtJUcas1Ns7inhamyfQJAXwvndxXJfaaa1ULZE3NasN4AYX95
g9dvlB60Kyyy4XlU2rLVrayVL4gXtBbg2YPNqnyUBfnGklEbXuGz2QA+OQJAUOiQ
gMVj5CPEZfTe/Ck3I4wvptzwnwM10ELxPFcBcPaVkm+cHsAkZ/fLgj6hWmH/tFP4
Zk60fcRHzePjnjLikQJAcaztmnB32RZNP3T0EWDT/ByFKjUl53epXX8YIfxipa1M
sWH66CeJ3S1o8G3CORXOPSJbhDCRBQnWLH+Dp3LovQ==
-----END RSA PRIVATE KEY-----`
	rsasign_privatekey = `-----BEGIN PRIVATE KEY-----
MIICXAIBAAKBgQDBFOnukK/bomnQB6rSupd7Cgh71Pwbylgbm2kuWRmH25BadjqggMsoU/i4uycsjj
oG6CwqXugeq1JnUV7hF6jAiWz7FsYrtESPAHnOq40C7vpSSrPt1q7oHMnpLReYVHrbcDc+uEhE5Bbrad
C38c32S5jxjpodn87v/5zLeCFqqwIDAQABAoGAHiRcrLCLs5b0O9Sml0Un1r5nOqWyQchh3tVxIxonwS
zGqUihuGLC1GXfgz88S1lct61RD8BHLlqCf7yVOkCOSFMBSKSg9VBFrZb+7ogWxmGPpjDL4yRehwuLRL
U/kdLOjLN0tub/bnNtZ/z9zP/SihVQmLL/dPbPCEMJmFmZOBECQQDfLWkoDYylpRbR7RG8tKysDbNyLW
oQI+2zSGSY40+jv/tAl4dk2JUA3yjx0ynPF83gX7Ogabz/g5ZyN0lv6CQDAkEA3Xpn0ps/N5WKB9pLAe
EO0LK5zvCo2gsUlVYsadwuAerBTwqGqlRNQSEBFrw5GGOWcbZZ/X+ff/Dkv3gW7uMiOQJBAJzHkuK/Of
m3kuNgLiCpr8+iRFhGTQcPplFW3syEixLWsBN9H3EZsPsOyf2vwOQprcgcktxyy4GYGB0ed6l0o68CP1
TIhvYYeBrPJAfdiNkvTnyV7otJlVni4/5G/rmHkBUryNR5MxQBMZG9EK2jYDf156GgLNLeCDDyBp9FWY
IZcQJBANRxXgq6IUXKpgvTFCHeTCHoNztOvN4XS8MYWkZK34a/1/mGLs9U+K1lnBnPho+ZtEpYyW5VZk
awSHjDplzlM8U=
-----END PRIVATE KEY-----`
	uadecode_privatekey = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQDN2CQ+ZcsbHVrLSJ+/v0BA6kgoa+riowfZErlGdlMjUhf4UTos
1ZzKfavQjctIQjlEStxBWAZ+gtmtYakUmjWFBupks6RWcR/ho+VEePpzxhSis8M6
zHDhf9Dl8u6uvSAdk+CM3NgGle0bC8WnMxy2k9k3MSrX93AKZV36DCMtUQIDAQAB
AoGACBp+MsanHEYnkOEnCNFqoiOW+6Bj+tAYOv91s8RsuXM95lSsSZ+PMJmJ7gfm
/M0+m+Wmjhv9BXX5Q84Ybes0OBpS9qk2Rq6mQooXqo+6BaPlJb/UD160ULiQJIA7
P2x2XX/Z8xW7goq6r13i4VOZj4GHxRnlfvtCyKKso2U4qMECQQD2bBrY+1O1Du7D
B6jjJC0rnZMQllIlex9zftvIgr69qvQcWoCKukBUjmx8gw7Pall33sRgJbhMj0yf
qLBn6GplAkEA1dhJcJhZ9euFsjWnjXG5g5CyRO5HPJ9Kfoqw04aQX7oHHRXUyVNU
H4L9FsbXQBYS+rXrtJUcas1Ns7inhamyfQJAXwvndxXJfaaa1ULZE3NasN4AYX95
g9dvlB60Kyyy4XlU2rLVrayVL4gXtBbg2YPNqnyUBfnGklEbXuGz2QA+OQJAUOiQ
gMVj5CPEZfTe/Ck3I4wvptzwnwM10ELxPFcBcPaVkm+cHsAkZ/fLgj6hWmH/tFP4
Zk60fcRHzePjnjLikQJAcaztmnB32RZNP3T0EWDT/ByFKjUl53epXX8YIfxipa1M
sWH66CeJ3S1o8G3CORXOPSJbhDCRBQnWLH+Dp3LovQ==
-----END RSA PRIVATE KEY-----`
}
