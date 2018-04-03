package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

/**
 * 文档地址：https://console.faceplusplus.com.cn/documents/4888391
 * 创建一个人脸的集合 FaceSet，用于存储人脸标识 face_token。一个 FaceSet 能够存储 1,000 个 face_token
 */

const facesetAPIURL = APIBaseURL + "/faceset"

// FaceSetBaseFaceResponse 基础响应结构
type FaceSetBaseFaceResponse struct {
	FaceResponse
	FacesetToken  string           `json:"faceset_token"`  // FaceSet 的标识
	OuterId       string           `json:"outer_id"`       // 用户自定义的 FaceSet 标识，如果未定义则返回值为空
	FaceCount     int              `json:"face_count"`     // 操作结束后 FaceSet 中的 face_token 总数量
	FailureDetail []*FailureDetail `json:"failure_detail"` // 无法被加入 FaceSet 的 face_token 以及原因
}

// FailureDetail 不能被添加的原因
type FailureDetail struct {
	Reason    string `json:"reason"`     // 不能被添加的原因，包括 INVALID_FACE_TOKEN 人脸表示不存在 ，QUOTA_EXCEEDED 已达到 FaceSet 存储上限
	FaceToken string `json:"face_token"` // 人脸标识
}

// FaceSetCreateFaceResponse 创建人脸集合FaceSet响应数据
type FaceSetCreateFaceResponse struct {
	FaceSetBaseFaceResponse
	FaceAdded int `json:"face_added"` // 本次操作成功加入 FaceSet的face_token 数量
}

// FaceSetAddFaceFaceResponse  添加人脸标识face_token响应数据，和创建FaceSet相同
type FaceSetAddFaceFaceResponse struct {
	FaceSetCreateFaceResponse
}

// FaceSetRemoveFaceFaceResponse 移除一个FaceSet中的某些或者全部face_token响应
type FaceSetRemoveFaceFaceResponse struct {
	FaceSetBaseFaceResponse
	FaceRemoved int `json:"face_removed"` // 成功从FaceSet中移除的face_token数量
}

// FaceSetUpdateFaceFaceResponse 更新一个人脸集合的属性响应
type FaceSetUpdateFaceFaceResponse struct {
	FaceResponse
	FacesetToken string `json:"faceset_token"` // FaceSet 的标识
	OuterId      string `json:"outer_id"`      // 用户自定义的 FaceSet 标识，如果未定义则返回值为空
}

// FaceSetGetDetailFaceFaceResponse 获取一个 FaceSet 的所有信息响应
type FaceSetGetDetailFaceFaceResponse struct {
	FaceResponse
	FacesetToken string   `json:"faceset_token"` // FaceSet 的标识
	OuterId      string   `json:"outer_id"`      // 用户自定义的 FaceSet 标识，如果未定义则返回值为空
	DisplayName  string   `json:"display_name"`  // 人脸集合的名字
	UserData     string   `json:"user_data"`     // 自定义用户信息
	Tags         string   `json:"tags"`          // 自定义标签
	FaceCount    int      `json:"face_count"`    // 操作结束后 FaceSet 中的 face_token 总数量
	FaceTokens   []string `json:"face_tokens"`   // face_token的数组
	Next         string   `json:"next"`          // 用于进行下一次请求。返回值表示排在此次返回的所有 face_token 之后的下一个 face_token 的序号
}

// FaceSetDeleteFaceFaceResponse 删除一个人脸集合响应
type FaceSetDeleteFaceFaceResponse struct {
	FaceResponse
	FacesetToken string `json:"faceset_token"` // FaceSet 的标识
	OuterId      string `json:"outer_id"`      // 用户自定义的 FaceSet 标识，如果未定义则返回值为空
}

// FaceSetGetFaceSetsFaceFaceResponse 获取某一 API Key 下的 FaceSet 列表响应
type FaceSetGetFaceSetsFaceFaceResponse struct {
	FaceResponse
	Facesets []*Faceset `json:"facesets"` // 该 API Key 下的 FaceSet 信息
	Next     string     `json:"next"`     // 用于进行下一次请求。返回值表示排在此次返回的所有 face_token 之后的下一个 face_token 的序号
}

// FaceSetRequest FaceSet管理对象
type FaceSetRequest struct {
	FaceRequest
	response interface{}
}

// Faceset FaceSet数组中单个元素的结构
type Faceset struct {
	FacesetToken string `json:"faceset_token"` // FaceSet 的标识
	OuterId      string `json:"outer_id"`      // 用户自定义的 FaceSet 标识，如果未定义则返回值为空
	DisplayName  string `json:"display_name"`  // FaceSet的名字，如果未提供为空
	Tags         string `json:"tags"`          // FaceSet的标签，如果未提供为空
}

// FaceSet 创建一个FaceSet操作对象
func (sdk *FaceSDK) FaceSet(options ...map[string]interface{}) (*FaceSetRequest, error) {
	faceSetRequest := new(FaceSetRequest)
	if len(options) == 0 {
		faceSetRequest.options = make(map[string]interface{}, 0)
	} else {
		faceSetRequest.options = options[0]
	}
	// 添加api key信息
	faceSetRequest.options["api_key"] = sdk.APIKey
	faceSetRequest.options["api_secret"] = sdk.APISecret

	faceSetRequest.request = sdk.getHTTPRequest()
	faceSetRequest.request.Debug = sdk.Debug

	return faceSetRequest, nil
}

// SetOption 设置请求参数
func (fsr *FaceSetRequest) SetOption(key string, val interface{}) *FaceSetRequest {
	fsr.options[key] = val
	return fsr
}

// SetOptionMap 通过map设置请求参数
func (fsr *FaceSetRequest) SetOptionMap(options map[string]interface{}) *FaceSetRequest {
	for key, val := range options {
		fsr.options[key] = val
	}
	return fsr
}

// Create 创建一个人脸集合
func (fsr *FaceSetRequest) Create() *FaceSetRequest {
	urlStr := fmt.Sprintf("%s/create", facesetAPIURL)
	log.Println(urlStr)
	fsr.request.Post(urlStr)
	fsr.response = new(FaceSetCreateFaceResponse)
	return fsr
}

// AddFace 添加人脸标识 face_token到FaceSet
func (fsr *FaceSetRequest) AddFace() *FaceSetRequest {
	urlStr := fmt.Sprintf("%s/addface", facesetAPIURL)
	fsr.request.Post(urlStr)
	fsr.response = new(FaceSetAddFaceFaceResponse)
	return fsr
}

// RemoveFace 移除一个FaceSet中的某些或者全部face_token
func (fsr *FaceSetRequest) RemoveFace() *FaceSetRequest {
	urlStr := fmt.Sprintf("%s/removeface", facesetAPIURL)
	fsr.request.Post(urlStr)
	fsr.response = new(FaceSetRemoveFaceFaceResponse)
	return fsr
}

// Update 更新一个人脸集合的属性
func (fsr *FaceSetRequest) Update() *FaceSetRequest {
	urlStr := fmt.Sprintf("%s/update", facesetAPIURL)
	fsr.request.Post(urlStr)
	fsr.response = new(FaceSetUpdateFaceFaceResponse)
	return fsr
}

// GetDetail 更新一个人脸集合的属性
func (fsr *FaceSetRequest) GetDetail() *FaceSetRequest {
	urlStr := fmt.Sprintf("%s/getdetail", facesetAPIURL)
	fsr.request.Post(urlStr)
	fsr.response = new(FaceSetGetDetailFaceFaceResponse)
	return fsr
}

// Delete 删除一个人脸集合
func (fsr *FaceSetRequest) Delete() *FaceSetRequest {
	urlStr := fmt.Sprintf("%s/delete", facesetAPIURL)
	fsr.request.Post(urlStr)
	fsr.response = new(FaceSetDeleteFaceFaceResponse)
	return fsr
}

// GetFaceSets 获取某一 API Key 下的 FaceSet 列表
func (fsr *FaceSetRequest) GetFaceSets() *FaceSetRequest {
	urlStr := fmt.Sprintf("%s/getfacesets", facesetAPIURL)
	fsr.request.Post(urlStr)
	fsr.response = new(FaceSetGetFaceSetsFaceFaceResponse)
	return fsr
}

// End 发送请求获取结果
func (fsr *FaceSetRequest) End() (interface{}, string, error) {
	log.Println(fsr.options)
	log.Println(fsr.request.Url)

	resp, body, errs := fsr.request.Type("multipart").SendMap(fsr.options).End()
	if len(errs) > 0 {
		return nil, "", errors.New("请求接口错误:" + errs[0].Error())
	}
	// 判断响应是否成功
	if resp.StatusCode != http.StatusOK {
		return nil, "", NewFaceError(resp.StatusCode, body)
	}
	// 解析body为对象
	err := json.Unmarshal([]byte(body), fsr.response)
	if err != nil {
		return nil, "", err
	}
	return fsr.response, body, nil
}
