# face-golang-sdk

face++????sdk?????????????

????? `examples` ??

???????????? [https://github.com/shiguanghuxian/face-login](https://github.com/shiguanghuxian/face-login)

??:

```
	// ????sdk??
	faceSDK, err := sdk.NewFaceSDK(APIKey, APISecret)
	log.Println(err)
	// ????????api????
	detect, err := faceSDK.Detect()
	log.Println(err)
	// ????
	dr, body, err := detect.SetImage("./demo.jpg", "image_file").
		SetOption("return_attributes", "gender,age,smiling,headpose,facequality,blur,eyestatus,emotion,ethnicity,beauty,mouthstatus,eyegaze,skinstatus").
		SetOption("return_landmark", 1).
		End()

	log.Println(err)
	log.Println(body)
	js, _ := json.Marshal(dr)
	log.Println(string(js))
	log.Println("?????", dr.Faces[0].Attributes.Age.Value)
```
