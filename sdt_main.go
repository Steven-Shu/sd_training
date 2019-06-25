package main

import (
	json "encoding/json"
	log "github.com/cihub/seelog"
	ioutil "io/ioutil"
	"os"
	ad "sd_training/ad_site/server"
	dami "sd_training/dami_app/server"
	scedu "sd_training/shenche_site/back_end"
	"sync"
)

func main() {
	replaceLogger()
	configs := loadConfig("./etc/sd_config.json")
	if nil == configs {
		log.Errorf("error in initing my little poor server, check the last error information, or call")
		log.Error("exit with error")
		return
	}
	log.Info(configs)
	wg := sync.WaitGroup{}
	wg.Add(3)
	go ad.Start(configs, &wg)
	go dami.Start(configs, &wg)
	go scedu.Start(configs, &wg)

	log.Info("all servers started")
	wg.Wait()
	log.Info("main ended.")
}

func replaceLogger() {
	logger, err := log.LoggerFromConfigAsFile("./seelog.xml")

	if err != nil {
		log.Critical("err parsing config log file", err)
		return
	}
	log.ReplaceLogger(logger)

}

func loadConfig(config_file string) map[string]string {
	log.Debugf("loading global config of %s", config_file)
	file, err := os.Open(config_file)
	if err != nil {
		log.Errorf("error in loading config file,%s", err)
		return nil
	}
	defer file.Close()
	configBytes, _ := ioutil.ReadAll(file)
	configs := make(map[string]string)
	err = json.Unmarshal(configBytes, &configs)
	if err != nil {
		log.Errorf("error in parsing config json,%s", err)
		return nil
	}
	log.Info("config loading successfully")
	return configs
}
