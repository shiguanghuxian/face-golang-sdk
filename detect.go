package sdk

import (
	"encoding/json"
	"errors"
	"net/http"
)

/**
 * 文档地址：https://console.faceplusplus.com.cn/documents/4888373
 * 传入图片进行人脸检测和人脸分析
 * 可以检测图片内的所有人脸，对于每个检测出的人脸，会给出其唯一标识 face_token，可用于后续的人脸分析、人脸比对等操作。对于正式 API Key，支持指定图片的某一区域进行人脸检测。
 * 本 API 支持对检测到的人脸直接进行分析，获得人脸的关键点和各类属性信息。对于试用 API Key，最多只对人脸框面积最大的 5 个人脸进行分析，其他检测到的人脸可以使用 Face Analyze API 进行分析。对于正式 API Key，支持分析所有检测到的人脸。
 */

const detectAPIURL = APIBaseURL + "/detect"

// DetectFaceResponse 人脸检测响应数据
type DetectFaceResponse struct {
	FaceResponse
	ImageId string  `json:"image_id"` // 被检测的图片在系统中的标识
	Faces   []*Face `json:"faces"`    // 被检测出的人脸数组，具体包含内容见下文 注：如果没有检测出人脸则为空数组
}

// FaceDetect 人脸检测和人脸分析对象
type FaceDetect struct {
	FaceRequest
}

// Detect 构建一个人脸检测和人脸分析对象
func (sdk *FaceSDK) Detect(options ...map[string]interface{}) (*FaceDetect, error) {
	faceDetect := new(FaceDetect)

	if len(options) == 0 {
		faceDetect.options = make(map[string]interface{}, 0)
	} else {
		faceDetect.options = options[0]
	}
	// 添加api key信息
	faceDetect.options["api_key"] = sdk.APIKey
	faceDetect.options["api_secret"] = sdk.APISecret

	faceDetect.request = sdk.getHTTPRequest().
		Post(detectAPIURL).
		Type("multipart")

	return faceDetect, nil
}

// SetImage 设置图片信息
// dt可以是(image_url|image_file|image_base64)
func (fd *FaceDetect) SetImage(img, dt string) *FaceDetect {
	if dt == "image_file" {
		fd.request.SendFile(img, "", "image_file")
	} else {
		fd.options[dt] = img
	}
	return fd
}

// SetOption 设置请求参数
func (fd *FaceDetect) SetOption(key string, val interface{}) *FaceDetect {
	fd.options[key] = val
	return fd
}

// SetOptionMap 通过map设置请求参数
func (fd *FaceDetect) SetOptionMap(options map[string]interface{}) *FaceDetect {
	for key, val := range options {
		fd.options[key] = val
	}
	return fd
}

// End 发送请求获取结果
func (fd *FaceDetect) End() (*DetectFaceResponse, string, error) {
	resp, body, errs := fd.request.SendMap(fd.options).End()
	if len(errs) > 0 {
		return nil, "", errors.New("请求接口错误:" + errs[0].Error())
	}
	// 判断响应是否成功
	if resp.StatusCode != http.StatusOK {
		return nil, "", NewFaceError(resp.StatusCode, body)
	}
	// 解析body为对象
	detectFaceResponse := new(DetectFaceResponse)
	err := json.Unmarshal([]byte(body), detectFaceResponse)
	if err != nil {
		return nil, "", err
	}
	return detectFaceResponse, body, nil
}
