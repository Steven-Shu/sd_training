package manager
import(
	"sync"
	model "sensors_test/back_end/my_server/model"
	log "github.com/cihub/seelog"
	"os"
	"encoding/json"
	"io/ioutil"
	"time"
	sapi "sensors_test/back_end/sensorsapi"
	"crypto/md5"
	"strings"
	"errors"
	"encoding/hex"
)

var umMutex *sync.Mutex
var umInstance *UserManager

var dataFilePath="./data/users.d"

func init(){
	umMutex=new(sync.Mutex)
}


type UserManager struct{
	uMutex *sync.Mutex
	currentUsers map[string]*model.User
	aMutex *sync.Mutex
	allUsers map[string]*model.User
	lastWriteTime int64
	lastLength int
}

func GetUserManager()*UserManager{
	umMutex.Lock()
	defer umMutex.Unlock()
	if umInstance == nil{
		umInstance=new (UserManager)
		umInstance.uMutex=new(sync.Mutex)
		umInstance.aMutex=new(sync.Mutex)
		umInstance.currentUsers=make(map[string]*model.User)
		umInstance.initUsers()
	}
	return umInstance
}

func(um *UserManager)initUsers(){
	log.Info("initing users info..")
	umInstance.allUsers=make(map[string]*model.User)
	um.aMutex.Lock()
	defer um.aMutex.Unlock()
	file,err:=os.Open(dataFilePath)
	if err!=nil && os.IsNotExist(err){
		file,err=os.Create(dataFilePath)
		if err!=nil{
			log.Errorf("too many errors in creating a simple file, pelase crash your computer,error:%s", err)
			return 
		}
		file.Close()
		return
	}
	if err !=nil{
		log.Errorf("other errore occured when opening user file.error:%s",err)
		return
	}
	defer file.Close()
	data,err:=ioutil.ReadAll(file)
	if err!=nil{
		log.Errorf("error in reading file content,error :%s", err)
		return
	}
	if len(data)>0{
	err=json.Unmarshal(data,&(umInstance.allUsers))
		if err!=nil{
			log.Errorf("error in unmarshal uninstance,error:%s", err)
			return
		}
	}
	um.lastLength=len(umInstance.allUsers)
	um.lastWriteTime=time.Now().Unix()
	log.Infof("old users loaded.%d users exist",len(umInstance.allUsers))
}

func(um *UserManager) Login(mobile string,pwd string, firstId string) (string,error){
	log.Infof("try to login for user:%s with pwd:%s,with firstid:%s",mobile,pwd,firstId)
	um.uMutex.Lock()
	defer um.uMutex.Unlock()
	hash:=md5.New()
	hash.Write([]byte(mobile+pwd))
	uHash:=hex.EncodeToString(hash.Sum(nil))
	u,ok:=um.currentUsers[mobile]
	if ok{
		if !strings.EqualFold(uHash, u.UserProperties["pwdHash"].(string)){
			log.Error("wrong user name or password")
			// 记录登陆失败事件，自动采集时间
			sapi.GetSAInstance().Track(firstId, "WrongPasswordLogin", nil, true)
			err:=errors.New("Wrong user name or password")
			return "",err
		}else{
			log.Debug("user's already logined")
			return mobile,nil
		}
	}
	u,ok=um.allUsers[mobile]

	if ok{
		if !strings.EqualFold(uHash, u.UserProperties["pwdHash"].(string)){
			log.Error("wrong user name or password")
			// 记录登陆失败事件，自动采集时间
			sapi.GetSAInstance().Track(firstId, "WrongPasswordLogin", nil, true)
			err:=errors.New("Wrong user name or password")
			return "",err
		}
	}
	newUser:=new(model.User)
	newUser.Id=mobile
	newUser.UserProperties=make(map[string]interface{})
	newUser.UserProperties["lastLoginTime"]=time.Now()
	newUser.UserProperties["pwdHash"]=uHash
	sapi.GetSAInstance().ProfileSet(firstId, newUser.UserProperties, false)
	um.currentUsers[newUser.Id]=newUser

	//登陆进行id绑定
	//sapi.GetSAInstance().TrackSingUp(mobile, firstId)

	//记录登录成功事件
	sapi.GetSAInstance().Track(firstId, "SuccessLogin", nil, false)

	go um.addNewUsers(newUser)
	log.Infof("user login successfully, username:%s",mobile)
	return newUser.Id,nil
}

func(um *UserManager) addNewUsers(user *model.User){
	um.aMutex.Lock()
	defer um.aMutex.Unlock()
	_,ok:=um.allUsers[user.Id]
	if ok{
		return
	}else{
		um.allUsers[user.Id]=user
		um.writeUsers()
	}
}

func (um *UserManager) writeUsers() error{
	log.Info("try to write users file")
	um.aMutex.Lock()
	um.aMutex.Unlock()
	timeGap:=time.Now().Unix()-um.lastWriteTime
	if timeGap>60 || len(um.allUsers)-um.lastLength >= 10{
		log.Debug("time gap or length gap triggered writing..")
		data,err:=json.Marshal(um.allUsers)
		if err!=nil{
			log.Errorf("error in marshal users to json,error:%s",err)
			return err
		}
		err=ioutil.WriteFile(dataFilePath, data, os.ModePerm)
		if err!=nil{
			log.Errorf("error in opening user file,error:%s", err)
			return err
		}
		um.lastWriteTime=time.Now().Unix()
		um.lastLength=len(um.allUsers)
		log.Debug("user file written")
	}else{
		log.Debug("not ready to write user file")
	}
	return nil
}

func (um *UserManager) IsLogin(mobile string) (bool){
	um.uMutex.Lock()
	defer um.uMutex.Unlock()
	_,ok:=um.currentUsers[mobile]
	if ok{
		log.Debugf("user %s's login",mobile)
		return true
	}else{
		log.Debugf("user %s's not login",mobile)
		return false
	}
}






