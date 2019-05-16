package main

import (
	"encoding/json"
)

type IUser interface {
	Value_for_key(_key string)string
	Set_key_value(_key string,_value string)
}
type BaseUser struct {
	data map[string]string
}

func (this* BaseUser)Value_for_key(_key string)string{
	return this.data[_key]
}
func (this* BaseUser)Set_key_value(_key string,_value string){
	this.data[_key]=_value
}

func (this* BaseUser)Data_with_map(_map map[string]string)error{
	for k,v:=range _map{
		this.data[k]=v
	}
	return nil
}
func (this* BaseUser)Data_with_s(_s string)error{
	err:=json.Unmarshal([]byte(_s),&this.data)
	if err!=nil{
		return err
	}
	return nil
}
type MGUser struct {
	BaseUser
}

func New_mguser()*MGUser{
	return &MGUser{BaseUser{make(map[string]string)}}
}