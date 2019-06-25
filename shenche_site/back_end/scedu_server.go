package back_end

import (
	log "github.com/cihub/seelog"
	ms "sd_training/shenche_site/back_end/myserver"
	"sync"
)

func Start(configs map[string]string, wg *sync.WaitGroup) {
	log.Info("scedu main's running")
	myServer := ms.GetMyHTTPServer()
	myServer.StartServer(configs)
	log.Debug("scedu server ended.")
	wg.Done()
}
