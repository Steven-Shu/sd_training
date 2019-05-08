package sensorsapi

import (
	"sync"
	"errors"
	log "github.com/cihub/seelog"
	sdk "github.com/sensorsdata/sa-sdk-go"
	cs "github.com/sensorsdata/sa-sdk-go/consumers"
)

var initOnceOnly =new(sync.Once)

var DefaultSAServerURL = "http://10.19.197.78:8106/sa?project=fujie"
var ProjectName="fujie"
var TimeOut=10000
var ConsumerTypeNow=Logging
var fileName="./logAgent/my_log_agent_log.log"
var hourcut=true
var debugWrite=false

var msaaInstance *MySAAnalytics

type MySAAnalytics struct{
	sda sdk.SensorsAnalytics
	spMutex *sync.Mutex
	superProperties map[string]interface{}
}

func GetSAInstance() *MySAAnalytics{
	initOnceOnly.Do(func(){
		msaaInstance=new(MySAAnalytics)
		msaaInstance.spMutex=new(sync.Mutex)
		msaaInstance.superProperties=make(map[string]interface{})
		msaaInstance.initSDK()
	})
	return msaaInstance
}

/** InitSASDK :
start the sa sdk
*/
func(msda *MySAAnalytics) initSDK() {
	log.Info("starting sa sdk")
	// 根据配置类型获取不同的consumer
	var consumer cs.Consumer
	var err error
	switch ConsumerTypeNow{
	case Logging:
		consumer,err=sdk.InitLoggingConsumer(fileName, hourcut)
	case Default:
		consumer,err=sdk.InitConcurrentLoggingConsumer(fileName, hourcut)
	case Debug:
		consumer,err=sdk.InitDebugConsumer(DefaultSAServerURL, debugWrite, TimeOut)
	default:
		log.Errorf("unsupported consumer type, if you wanna use batch/default consumer, please DIY")
		err=errors.New("unsupported consumer type")
	}

	if err !=nil {
		log.Errorf("error in initiating sensors analystic...error:%s",err)
		return 
	}
	//...
	// 使用 Consumer 来构造 SensorsAnalytics 对象
	msda.sda = sdk.InitSensorsAnalytics(consumer, ProjectName, false)
	msda.RegisterSuperProperties(map[string]interface{}{"appName":"神澈教育"})
	log.Debugf("sa initilization finished with url:%s finished",DefaultSAServerURL)
}

func (msda *MySAAnalytics) Track(distinctId string, eventName string, properties map[string]interface{},isLoginId bool){
	log.Debug("track sign up to sensors sdk for id:%s",distinctId)
	fullProperties:=msda.getFullPropeties(properties)
	msda.sda.Track(distinctId, eventName, fullProperties, isLoginId)
}

func (msda *MySAAnalytics) TrackSingUp(distinctId string,originalId string){
	log.Debug("track sign up to sensors sdk for id:%s",distinctId)
	msda.sda.TrackSignup(distinctId, originalId)
}

func(msda *MySAAnalytics)ProfileSetOnce(distinctId string,properties map[string]interface{},isLoginId bool){
	msda.sda.ProfileSetOnce(distinctId, properties, isLoginId)
}

func(msda *MySAAnalytics)ProfileSet(distinctId string,properties map[string]interface{},isLoginId bool){
	msda.sda.ProfileSet(distinctId, properties, isLoginId)
}

func(msda *MySAAnalytics)ProfileIncrement(distinctId string,properties map[string]interface{},isLoginId bool){
	msda.sda.ProfileIncrement(distinctId, properties, isLoginId)
}

func (msda *MySAAnalytics) RegisterSuperProperties(properties map[string]interface{}){
	msda.spMutex.Lock()
	defer msda.spMutex.Unlock()
	for k,v :=range properties{
		msda.superProperties[k]=v
	}
	log.Info("super propeties updated")
}

func (msda *MySAAnalytics) ClearSuperProperties(){
	msda.spMutex.Lock()
	defer msda.spMutex.Unlock()
	msda.superProperties=map[string]interface{}{}
}

func (msda *MySAAnalytics) getFullPropeties(properties map[string]interface{}) map[string]interface{}{
	fullProperties:=map[string]interface{}(msda.superProperties)
	for k,v :=range properties {
		fullProperties[k]=v
	}

	return fullProperties
}

func (msda *MySAAnalytics) Close(){
	log.Info("sensors data analystic closing...")
	msda.sda.Close()
	log.Info("sensors data analystic closed.")
}
