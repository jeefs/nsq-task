package nsqQueue

import (
	"cgm_manager/internal/glucose/model"
	glucoseServer "cgm_manager/internal/glucose/service"
	"cgm_manager/utils/qiniu"
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/qiniu/go-sdk/v7/storage"
	"log"
	"os"
)

type MyHandler struct {
	Title string
}

// HandleMessage 是需要实现的处理消息的方法
func (m *MyHandler) HandleMessage(msg *nsq.Message) (err error) {
	fmt.Printf("%s recv from %v, msg:%v\n", m.Title, msg.NSQDAddress, string(msg.Body))
	return
}

// 计算设备平均mard值
type CalculateMardHandler struct {
	DeviceType int8   `json:"deviceType"`
	Mac        string `json:"mac"`
	UserId     string `json:"userId"`
}

func (c *CalculateMardHandler) HandleMessage(msg *nsq.Message) error {
	var res CalculateMardHandler
	if len(msg.Body) == 0 {
		return nil
	}
	err := json.Unmarshal(msg.Body, &res)
	if err != nil {
		return nil
	}
	logger := GetTaskLogger(CalculateMardTopic)
	perId, err := glucoseServer.CalculateKetoneDeviceMard(res.Mac, res.UserId)
	if err != nil {
		if logger != nil {
			logger.Printf("计算设备平均mard值失败:mac:%v,userId:%v,res:nil,err:%v", res.Mac, res.UserId, err.Error())
		}
		return err
	} else {
		if logger != nil {
			log.Printf("计算设备平均mard值成功:mac:%v,userId:%v,res:%v,err:nil", res.Mac, res.UserId, perId)
		}
		return nil
	}
}

// 上传图片到七牛云
type UploadFileToQiniuHandler struct {
	RealFilePath string `json:"realFilePath"`
	PublicUrl    string `json:"publicUrl"`
	CheckId      string `json:"checkId"`
	UpToken      string `json:"upToken"`
}

func (u *UploadFileToQiniuHandler) HandleMessage(msg *nsq.Message) error {
	var res UploadFileToQiniuHandler
	if len(msg.Body) == 0 {
		return nil
	}
	err := json.Unmarshal(msg.Body, &res)
	if err != nil {
		return nil
	}
	file, err := os.OpenFile(LogPath+"/uploadFileToQiniu.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("打开任务队列日志文件失败,请求后端排除错误")
		panic(err.Error())
	}
	defer func() {
		file.Close()
	}()
	log.SetOutput(file)
	//指血图片上传任务
	cfg := storage.Config{}
	//空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	//是否使用https域名
	cfg.UseHTTPS = false
	//传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	//上传服务器文件到七牛云的任务
	fileInfo, err := os.Stat(res.RealFilePath)
	if err != nil {
		return err
	}
	key := fileInfo.Name() //上传后生成的文件名
	_, err = qiniu.FormUpload(res.RealFilePath, cfg, &key, res.UpToken)
	if err != nil {
		return err
	}
	checkModel := model.PersonGlucoseCheck{ //todo 暂时不清理服务器本地图片，测试稳定后再删除，防止用户图片丢失
		PerID:           res.CheckId,
		GlucoseValuePic: res.PublicUrl,
		StorageType:     2,
	}
	err = glucoseServer.UpdateGlucoseCheckValByPerId(checkModel) //更新七牛云链接到数据库
	if err != nil {
		log.Printf("上传七牛云图片失败:userId:%v,publicUrl:%v,err:%v", res.CheckId, res.PublicUrl, err.Error())
		return err
	} else {
		log.Printf("上传七牛云图片成功:userId:%v,publicUrl:%v,err:%v", res.CheckId, res.PublicUrl, nil)
		return nil
	}
}
