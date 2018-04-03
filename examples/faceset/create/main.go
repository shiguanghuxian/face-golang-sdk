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
	faceSDK, err := sdk.NewFaceSDK(APIKey, APISecret, true)
	log.Println(err)
	faceSet, err := faceSDK.FaceSet(map[string]interface{}{
		"display_name": "user1",
		"outer_id":     "user1",
	})
	log.Println(err)

	cr, body, err := faceSet.Create().End()

	log.Println(err)
	log.Println(body)
	js, _ := json.Marshal(cr)
	log.Println(string(js))
}
