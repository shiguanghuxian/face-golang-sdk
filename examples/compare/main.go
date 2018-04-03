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
	faceSDK, err := sdk.NewFaceSDK(APIKey, APISecret)
	log.Println(err)
	compare, err := faceSDK.Compare()
	log.Println(err)
	cr, body, err := compare.
		SetFace1("./demo-pic39.jpg", "image_file1").
		SetFace2("./demo-pic33.jpg", "image_file2").
		End()

	log.Println(err)
	log.Println(body)
	js, _ := json.Marshal(cr)
	log.Println(string(js))
}
