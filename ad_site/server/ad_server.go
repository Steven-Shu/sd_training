package server

import (
	log "github.com/cihub/seelog"
	"github.com/drone/routes"
	"net/http"
	"sync"
)

func AdCall(w http.ResponseWriter, r *http.Request) {

}

func Start(configs map[string]string, wg *sync.WaitGroup) {
	log.Debug("start ad server")
	mux := routes.New()
	service_root := configs["ad_service_root"]
	front_root := configs["ad_front_root"]
	service_port := configs["ad_server_port"]
	mux.Post(service_root+"ad", AdCall)
	http.Handle(service_root, mux)
	http.Handle(front_root, http.StripPrefix(front_root, http.FileServer(http.Dir("./ad_site/front/toutiao/"))))
	log.Debugf("ready to listen on %s", service_port)
	err := http.ListenAndServe(":"+service_port, nil)
	if err != nil {
		log.Error(err)
	}
	wg.Done()
}
