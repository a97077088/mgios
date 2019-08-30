package mgios

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestToken_refresh_migu(t *testing.T) {
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
		"userId":            "1225451333",
		"userToken":         "D8663E1329A50ECDDFA8",
		"miguToken":         "8484010001330200584E446B334E6A5A45515455784D544A4451305134526A517740687474703A2F2F70617373706F72742E6D6967752E636E2F6E303030312F4034333231333731333435356136363364326131326539643834353735616333300300040417A3140400063230333030350500203036433731453631453342433431383542414546353130323336353738434231FF00208EE15E5AF85C9190A6B0A525662A2A6F09CBB1D361960D94C8135BB62C7D7249",
	})
	o,err:=Token_refresh_migu(mguser,nil)
	if err!=nil{
		panic(err)
	}
	ntm,err:=strconv.ParseInt(o.UserInfo.ExpiredOn,10,64)
	if err!=nil{
		panic(err)
	}
	fmt.Printf("下次到期时间:%s\n",time.Unix(ntm/1000,0).Format("2006-01-02-15:04:05"))
}
func Test_login(t *testing.T) {
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
	err := Login_with_u_p_user_cli("15738558265", "a1231231", "焦作市", mguser, nil)
	if err != nil {
		panic(err)
	}
	user_log(
		fmt.Sprintf("登录成功,token:%s \nmigutoken:%s\n",
			mguser.Value_for_key("userToken"),
			mguser.Value_for_key("migutoken"),
		), mguser)
}
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
