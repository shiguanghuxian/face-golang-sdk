package sdk

import (
	"encoding/json"
	"errors"
	"net/http"
)

/**
 * 文档地址：https://console.faceplusplus.com.cn/documents/4887586
 * 将两个人脸进行比对，来判断是否为同一个人，返回比对结果置信度和不同误识率下的阈值。
 * 支持传入图片或 face_token 进行比对。使用图片时会自动选取图片中检测到人脸尺寸最大的一个人脸。
 */

const compareAPIURL = APIBaseURL + "/compare"

// CompareFaceResponse 人脸对比响应数据
type CompareFaceResponse struct {
	FaceResponse
	Confidence float32            `json:"confidence"`
	Thresholds map[string]float32 `json:"thresholds"`
	ImageId1   string             `json:"image_id1"`
	ImageId2   string             `json:"image_id2"`
	Faces1     []*Face            `json:"faces1"`
	Faces2     []*Face            `json:"faces2"`
}

// FaceCompare 人脸比对对象
type FaceCompare struct {
	FaceRequest
}

// Compare 构建一个人脸比对对象
func (sdk *FaceSDK) Compare(options ...map[string]interface{}) (*FaceCompare, error) {
	faceCompare := new(FaceCompare)
	if len(options) == 0 {
		faceCompare.options = make(map[string]interface{}, 0)
	} else {
		faceCompare.options = options[0]
	}
	// 添加api key信息
	faceCompare.options["api_key"] = sdk.APIKey
	faceCompare.options["api_secret"] = sdk.APISecret

	faceCompare.request = sdk.getHTTPRequest().
		Post(compareAPIURL).
		Type("multipart")

	return faceCompare, nil
}

// SetFace1 设置第一张照片
// dt可以是(face_token1|image_url1|image_file1|image_base64_1)
func (fc *FaceCompare) SetFace1(face1, dt string) *FaceCompare {
	if dt == "image_file1" {
		fc.request.SendFile(face1, "", "image_file1")
	} else {
		fc.options[dt] = face1
	}
	return fc
}

// SetFace2 设置第二张照片
// dt可以是(face_token2|image_url2|image_file2|image_base64_2)
func (fc *FaceCompare) SetFace2(face2, dt string) *FaceCompare {
	if dt == "image_file2" {
		fc.request.SendFile(face2, "", "image_file2")
	} else {
		fc.options[dt] = face2
	}
	return fc
}

// SetOption 设置请求参数
func (fc *FaceCompare) SetOption(key string, val interface{}) *FaceCompare {
	fc.options[key] = val
	return fc
}

// SetOptionMap 通过map设置请求参数
func (fc *FaceCompare) SetOptionMap(options map[string]interface{}) *FaceCompare {
	for key, val := range options {
		fc.options[key] = val
	}
	return fc
}

// End 发送请求获取结果
func (fc *FaceCompare) End() (*CompareFaceResponse, string, error) {
	resp, body, errs := fc.request.SendMap(fc.options).End()
	if len(errs) > 0 {
		return nil, "", errors.New("请求接口错误:" + errs[0].Error())
	}
	// 判断响应是否成功
	if resp.StatusCode != http.StatusOK {
		return nil, "", NewFaceError(resp.StatusCode, body)
	}
	// 解析body为对象
	compareFaceResponse := new(CompareFaceResponse)
	err := json.Unmarshal([]byte(body), compareFaceResponse)
	if err != nil {
		return nil, "", err
	}
	return compareFaceResponse, body, nil
}
