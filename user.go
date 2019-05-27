package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"sync"
)

type IUser interface {
	Value_for_key(_key string) string
	Set_key_value(_key string, _value string)
	Env_for_key(_key string) interface{}
	Set_env_value(_key string, _value interface{})
}
type BaseUser struct {
	data sync.Map
	env  sync.Map
}

func (this *BaseUser) Value_for_key(_key string) string {
	if v, isok := this.data.Load(_key); isok {
		return v.(string)
	} else {
		return ""
	}
}
func (this *BaseUser) Set_key_value(_key string, _value string) {
	this.data.Store(_key, _value)
}
func (this *BaseUser) Env_for_key(_key string) interface{} {
	v, b := this.env.Load(_key)
	if b == false {
		return nil
	}
	return v
}
func (this *BaseUser) Set_env_value(_key string, _value interface{}) {
	this.env.Store(_key, _value)
}

func (this *BaseUser) Data_with_map(_map map[string]string) {
	for k, v := range _map {
		this.data.Store(k, v)
	}
}
func (this *BaseUser) Data_with_s(_s string) error {
	mp := map[string]string{}
	err := json.Unmarshal([]byte(_s), &mp)
	if err != nil {
		return err
	}
	this.Data_with_map(mp)
	return nil
}

type MGUser struct {
	BaseUser
}

func New_mguser_with_map(_mp map[string]string) *MGUser {
	mguser := &MGUser{
		BaseUser{sync.Map{}, sync.Map{}},
	}
	mguser.Data_with_map(_mp)
	return mguser
}
func New_mguser() *MGUser {
	return &MGUser{BaseUser{sync.Map{}, sync.Map{}}}
}
func New_users_with_file(_fname string) ([]*MGUser, error) {
	rs := make([]*MGUser, 0)
	btall, err := ioutil.ReadFile(_fname)
	if err != nil {
		return nil, err
	}
	sreader := bufio.NewReader(bytes.NewReader(btall))
	for {
		line, _, err := sreader.ReadLine()
		if err == io.EOF {
			break
		}
		u := MGUser{}
		err = u.Data_with_s(string(line))
		if err != nil {
			return nil, err
		}
		rs = append(rs, &u)
	}
	return rs, nil
}
