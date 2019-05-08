package sensorsapi

import (
	http "net/http"
	url "net/url"
	log "github.com/cihub/seelog"
	"encoding/json"
)

func GetFrontEndSensorsField(ck *http.Cookie,key string) string {
	sscrossdata,err := url.QueryUnescape(ck.Value)
	if err!=nil{
		log.Errorf("error in parsing url params,%s", err)
		return ""
	}
	log.Debugf("url values:%v",sscrossdata)
	values:=make(map[string]interface{})
	err=json.Unmarshal([]byte(sscrossdata), &values)
		if err!=nil{
		log.Errorf("error in unmarshalling url params,%s", err)
		return ""
	}
	v,ok:=values[key]
	if !ok{
		log.Errorf("no such field %s", key)
		return ""
	}
	log.Debugf("key %s value is:%s",key,v.(string))
	return v.(string)
}
