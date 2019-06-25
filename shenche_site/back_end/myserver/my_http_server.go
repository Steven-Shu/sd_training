package myserver

import (
	http "net/http"
	"os"
	"path/filepath"
	sa "sd_training/shenche_site/back_end/sensorsapi"
	"sync"

	seelog "github.com/cihub/seelog"
	routes "github.com/drone/routes"
)

var s_mutex *sync.Mutex
var hsInstance *MyHTTPServer

func init() {
	if s_mutex == nil {
		seelog.Info("first time to init http server...")
		s_mutex = new(sync.Mutex)
		hsInstance = GetMyHTTPServer()
	}
}

type MyHTTPServer struct {
	mhsMutex    *sync.Mutex
	my_services []MyService
}

func GetMyHTTPServer() *MyHTTPServer {
	s_mutex.Lock()
	defer s_mutex.Unlock()
	if hsInstance == nil {
		seelog.Info("get my http server instance...")
		hsInstance = new(MyHTTPServer)
		hsInstance.mhsMutex = new(sync.Mutex)
		hsInstance.my_services = make([]MyService, 0)
	}
	return hsInstance
}

func (mhs *MyHTTPServer) AddServiceInstance(ser MyService) {
	seelog.Debug("add my service to register queue")
	s_mutex.Lock()
	defer s_mutex.Unlock()
	mhs.my_services = append(mhs.my_services, ser)
	seelog.Debugf("current services number:%d", len(mhs.my_services))
}

func (mhs *MyHTTPServer) RegisterServices(m *routes.RouteMux) {
	seelog.Debug("register all the services...")
	s_mutex.Lock()
	defer s_mutex.Unlock()
	for _, s := range mhs.my_services {
		s.RegisterServices(m)
	}
}

func (mhs *MyHTTPServer) StartServer(configs map[string]string) {
	seelog.Info("start http server...")
	sa.GetSAInstance(configs)

	mux := routes.New()
	mhs.RegisterServices(mux)

	service_root := configs["scedu_service_root"]
	front_root := configs["scedu_front_root"]
	service_port := configs["scedu_server_port"]
	filePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	seelog.Info(filePath)
	http.Handle(service_root, mux)
	http.Handle(front_root, http.StripPrefix(front_root, http.FileServer(http.Dir("./shenche_site/front_end/scedu/"))))
	seelog.Debugf("ready to listen on %s", service_port)
	err := http.ListenAndServe(":"+service_port, nil)
	if err != nil {
		seelog.Error(err)
	}
}
