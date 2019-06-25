package server

import (
	log "github.com/cihub/seelog"
	"github.com/drone/routes"
	"net/http"
	"sync"
)

func DownloadApp(w http.ResponseWriter, r *http.Request) {

}

func Start(configs map[string]string, wg *sync.WaitGroup) {
	log.Info("start dami server")
	mux := routes.New()
	service_root := configs["dami_service_root"]
	front_root := configs["dami_front_root"]
	service_port := configs["dami_server_port"]
	mux.Post(service_root+"dl", DownloadApp)
	http.Handle(service_root, mux)
	http.Handle(front_root, http.StripPrefix(front_root, http.FileServer(http.Dir("./dami_app/front/dami/"))))
	log.Debugf("ready to listen on %s", service_port)
	err := http.ListenAndServe(":"+service_port, nil)
	if err != nil {
		log.Error(err)
	}
	log.Debug("dami server ended")
	wg.Done()
}
