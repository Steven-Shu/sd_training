package model
import "time"


type Session struct{
	SessionId string
	CookieId string
	Expired time.Time
	Properties map[string]interface{}
}