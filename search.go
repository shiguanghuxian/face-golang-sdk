package sdk

import (
	"encoding/json"
	"errors"
	"net/http"
)

/**
 * 文档地址：https://console.faceplusplus.com.cn/documents/4888381
 * 在一个已有的 FaceSet 中找出与目标人脸最相似的一张或多张人脸，返回置信度和不同误识率下的阈值。
 * 支持传入图片或 face_token 进行人脸搜索。使用图片进行搜索时会选取图片中检测到人脸尺寸最大的一个人脸。
 */

const searchAPIURL = APIBaseURL + "/search"

// SearchFaceResponse 搜索接口返响应数据
type SearchFaceResponse struct {
	FaceResponse
	Results    []*SearchResults   `json:"results"`    // 搜索结果对象数组
	Thresholds map[string]float32 `json:"thresholds"` // 一组用于参考的置信度阈值
	ImageId    string             `json:"image_id"`   // 传入的图片在系统中的标识
	Faces      []*Face            `json:"faces"`      // 传入的图片中检测出的人脸数组
}

// SearchResults 搜索结果对象
type SearchResults struct {
	FaceToken  string  `json:"face_token"` // 从 FaceSet 中搜索出的一个人脸标识 face_token
	Confidence float32 `json:"confidence"` // 比对结果置信度，范围 [0,100]，小数点后3位有效数字，数字越大表示两个人脸越可能是同一个人
	UserId     string  `json:"user_id"`    // 用户提供的人脸标识，如果未提供则为空
}

// SearchRequest 人脸搜索对象
type SearchRequest struct {
	FaceRequest
}

// Search 构建一个人脸比对对象
func (sdk *FaceSDK) Search(options ...map[string]interface{}) (*SearchRequest, error) {
	searchRequest := new(SearchRequest)
	if len(options) == 0 {
		searchRequest.options = make(map[string]interface{}, 0)
	} else {
		searchRequest.options = options[0]
	}
	// 添加api key信息
	searchRequest.options["api_key"] = sdk.APIKey
	searchRequest.options["api_secret"] = sdk.APISecret

	searchRequest.request = sdk.getHTTPRequest().
		Post(searchAPIURL).
		Type("multipart")

	return searchRequest, nil
}

// SetFace 设置要搜索的图片
// dt可以是(face_token|image_url|image_file|image_base64)
func (sc *SearchRequest) SetFace(face, dt string) *SearchRequest {
	if dt == "image_file" {
		sc.request.SendFile(face, "", "image_file")
	} else {
		sc.options[dt] = face
	}
	return sc
}

// SetFaceSet 设置要查找的faceset
// dt可以是(faceset_token|outer_id)
func (sc *SearchRequest) SetFaceSet(set, dt string) *SearchRequest {
	sc.options[dt] = set
	return sc
}

// SetOption 设置请求参数
func (sc *SearchRequest) SetOption(key string, val interface{}) *SearchRequest {
	sc.options[key] = val
	return sc
}

// SetOptionMap 通过map设置请求参数
func (sc *SearchRequest) SetOptionMap(options map[string]interface{}) *SearchRequest {
	for key, val := range options {
		sc.options[key] = val
	}
	return sc
}

// End 发送请求获取结果
func (sc *SearchRequest) End() (*SearchFaceResponse, string, error) {
	resp, body, errs := sc.request.SendMap(sc.options).End()
	if len(errs) > 0 {
		return nil, "", errors.New("请求接口错误:" + errs[0].Error())
	}
	// 判断响应是否成功
	if resp.StatusCode != http.StatusOK {
		return nil, "", NewFaceError(resp.StatusCode, body)
	}
	// 解析body为对象
	searchFaceResponse := new(SearchFaceResponse)
	err := json.Unmarshal([]byte(body), searchFaceResponse)
	if err != nil {
		return nil, "", err
	}
	return searchFaceResponse, body, nil
}
