package my_server
import (
	"net/http"
	"sync"
	log "github.com/cihub/seelog"
	"io/ioutil"
	routes "github.com/drone/routes"
	manager "sensors_test/back_end/my_server/manager"
	"encoding/json"
	sapi "sensors_test/back_end/sensorsapi"
)

var lsMutex=new(sync.Mutex)
var lsInstance *LoginService

func init(){
	log.Debug("init the login service...")
	lsMutex.Lock()
	defer lsMutex.Unlock()
	lsInstance = new(LoginService)
	GetMyHTTPServer().AddServiceInstance(lsInstance)
}


type LoginService struct{

}

func(ls LoginService) IsLogin(w http.ResponseWriter, r *http.Request){
	log.Debug("user login get from front end")
	err:=r.ParseForm()
	if err!=nil{
		log.Errorf("error in parsing form,%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		routes.ServeJson(w, map[string]string{"result":"failed"})
		return 
	}
	params:=r.Form
	log.Debugf("params from login:%s",params)
	isLogin:=manager.GetUserManager().IsLogin(params.Get("mobile"))
	routes.ServeJson(w, map[string]interface{}{"result":"success","isLogin":isLogin})
}


func(ls LoginService) UserLogin(w http.ResponseWriter, r *http.Request){
	log.Debug("user login from front end")
	bodyBytes,err:=ioutil.ReadAll(r.Body)
	if err!=nil{
		log.Errorf("error in ready request body, error:%s",err)
		w.WriteHeader(http.StatusInternalServerError)
		routes.ServeJson(w, map[string]string{"error":"invalid request param"})
		return
	}
	params:=make(map[string]interface{})
	results:=make(map[string]string)
	err=json.Unmarshal(bodyBytes, &params)
	if err!=nil{
		log.Errorf("error in parsing params from request body,err:%s",err)
		routes.ServeJson(w, map[string]string{"error":"invalid request param"})
		return
	}
	log.Debugf("all params:%v",params)
	username:=params["params"].(map[string]interface{})["username"].(string)
	log.Debugf("params from login:%v",username)
	password:=params["params"].(map[string]interface{})["password"].(string)
	log.Debugf("params from login:%v",password)
	sensorsCookie,err:=r.Cookie("sensorsdata2015jssdkcross")
	distinct_id:=sapi.GetFrontEndSensorsField(sensorsCookie, "distinct_id")
	// first_id:=sapi.GetFrontEndSensorsField(sensorsCookie,"first_id")
	// if !strings.EqualFold("", first_id){
	// 	distinct_id=first_id
	// }
	userid,err:=manager.GetUserManager().Login(username,password,distinct_id)
	if err!=nil{
		log.Errorf("error in login,%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		results["result"]="failed"
		routes.ServeJson(w,results)
		return
	}
	log.Debug("login successful")
	results["result"]="success"
	results["login_id"]=userid
	log.Debugf("result to front end:%v",results)
	routes.ServeJson(w, results)
}

func (ls LoginService)RegisterServices(m *routes.RouteMux){
	log.Debug("register login services...")
	m.Post("/edu/userLogin", ls.UserLogin)
	m.Get("/edu/userLogin",ls.IsLogin)
	log.Debug("login services registered.")
}
