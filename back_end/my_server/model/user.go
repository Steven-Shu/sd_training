package model

type User struct{
	Id string
	NickName string
	UserProperties map[string]interface{}
}