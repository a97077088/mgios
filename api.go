package main

import (
	. "github.com/a97077088/grequests"
)

func Login_with_u_p_user_cli(_u, _p string, _areaName string, _user IUser, _cli *Session) error {
	rauth, err := client_authRequest(_u, _p, _user, nil)
	if err != nil {
		return err
	}
	//rauth:=map[string]string{}
	_, err = login_migutokenforall(rauth, _p, _user, nil)
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
