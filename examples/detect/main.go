package main

import (
	"encoding/json"
	"log"

	sdk "github.com/shiguanghuxian/face-golang-sdk"
)

var (
	APIKey    = ""
	APISecret = ""
)

func main() {
	// 创建一个sdk对象
	faceSDK, err := sdk.NewFaceSDK(APIKey, APISecret)
	log.Println(err)
	// 创建一个人脸检测api调用对象
	detect, err := faceSDK.Detect()
	log.Println(err)
	// 设置参数
	dr, body, err := detect.SetImage("./demo.jpg", "image_file").
		SetOption("return_attributes", "gender,age,smiling,headpose,facequality,blur,eyestatus,emotion,ethnicity,beauty,mouthstatus,eyegaze,skinstatus").
		SetOption("return_landmark", 1).
		End()

	log.Println(err)
	log.Println(body)
	js, _ := json.Marshal(dr)
	log.Println(string(js))
	log.Println("预测年龄：", dr.Faces[0].Attributes.Age.Value)
}
