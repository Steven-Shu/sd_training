package main
import (
	ms "./my_server"
	log "github.com/cihub/seelog"
)

func main(){

	replaceLogger()
	log.Info("scedu main's running")
	myServer:=ms.GetMyHTTPServer()
	myServer.StartServer()
	
}

func replaceLogger() {
	logger, err := log.LoggerFromConfigAsFile("./seelog.xml")

	if err != nil {
		log.Critical("err parsing config log file", err)
		return
	}
	log.ReplaceLogger(logger)
}