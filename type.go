package main

import "errors"
type NetErr struct {
	error
}
func New_neterr_with_error(_e error)NetErr{
	return NetErr{_e}
}
func New_neterr_with_s(_s string)NetErr{
	return NetErr{errors.New(_s)}
}
const (
	API_LOGIN_SUCCESS="LOGIN_SUCCESS"
	API_FAILED="FAILED"
)

type Api_R struct{
	ResultCode string
	ResultDesc string
}
type Login_Migutokenforall_r struct{
	Api_R
	LoginId string
	LoginType string
	ExtInfo Api_extInfo
	UserInfo Api_userInfo
	Sign string
}
type Api_extInfo struct{
	LoginId string
	UserNum string
	Nickname string
}
type Api_userInfo struct{
	AreaId string
	CarrierCode string
	CityId string
	ExpiredOn string
	Mobile string
	PassId string
	UserId string
	UserNum string
	UserToken string
}

type R struct{
	Code string
	Message string
}
type Playurl_v1_play_url_r struct {
	R
}