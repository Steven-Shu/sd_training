package my_server

import (
	http "net/http"
 	"sync"
	sa "sensors_test/back_end/sensorsapi"

	seelog "github.com/cihub/seelog"
	routes "github.com/drone/routes"
)

var s_mutex *sync.Mutex
var hsInstance *MyHTTPServer

func init(){
	 if s_mutex == nil {
	 	seelog.Info("first time to init http server...")
	 	s_mutex = new(sync.Mutex)
	 	hsInstance=GetMyHTTPServer()
	}
}

type MyHTTPServer struct{
	mhsMutex *sync.Mutex
	my_services []MyService
}

func GetMyHTTPServer() *MyHTTPServer{
	s_mutex.Lock()
	defer s_mutex.Unlock()
	if hsInstance == nil{
		seelog.Info("get my http server instance...")
		hsInstance =new(MyHTTPServer)
		hsInstance.mhsMutex=new(sync.Mutex)
		hsInstance.my_services=make([]MyService,0)
	}
	return hsInstance
}

func (mhs *MyHTTPServer)AddServiceInstance(ser MyService){
	seelog.Debug("add my service to register queue")
	s_mutex.Lock()
	defer s_mutex.Unlock()
	mhs.my_services=append(mhs.my_services,ser)
	seelog.Debugf("current services number:%d",len(mhs.my_services))
}


func (mhs *MyHTTPServer)RegisterServices(m *routes.RouteMux){
	seelog.Debug("register all the services...")
	s_mutex.Lock()
	defer s_mutex.Unlock()
	for _,s:=range mhs.my_services{
		s.RegisterServices(m)
	}
}

func (mhs *MyHTTPServer)StartServer() {
	seelog.Info("start http server...")
	sa.GetSAInstance()

	mux := routes.New()
	mhs.RegisterServices(mux)

	http.Handle("/", mux)
	http.Handle("/scedu/", http.StripPrefix("/scedu/", http.FileServer(http.Dir("../front_end/scedu/"))))
	seelog.Debug("ready to listen on 8080")
	err:=http.ListenAndServe(":8080", nil)
	if err!=nil{
		seelog.Error(err)
	}
}


