# face-golang-sdk

face++的golang版 sdk

示例清参考 `examples` 文件夹

也可以参考小程序人脸登录服务端[https://github.com/shiguanghuxian/face-login](https://github.com/shiguanghuxian/face-login)

示例:

```
	// 创建一个sdk对象
	faceSDK, err := sdk.NewFaceSDK(APIKey, APISecret)
	log.Println(err)
	// 创建人脸检测对象
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
	log.Println("年龄：", dr.Faces[0].Attributes.Age.Value)
```
