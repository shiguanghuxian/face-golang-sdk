package sdk

/**
 * 文档地址：https://console.faceplusplus.com.cn/documents/4888383
 * 传入在 Detect API 检测出的人脸标识 face_token，分析得出人脸关键点，人脸属性信息。一次调用最多支持分析 5 个人脸。
 */

const faceAPIURL = APIBaseURL + "/face"

// FaceAPIRequest FaceSet管理对象
type FaceAPIRequest struct {
	FaceRequest
	response interface{}
}

// Face 人脸分析对象
func (sdk *FaceSDK) Face(options ...map[string]interface{}) (*FaceAPIRequest, error) {
	faceAPIRequest := new(FaceAPIRequest)
	if len(options) == 0 {
		faceAPIRequest.options = make(map[string]interface{}, 0)
	} else {
		faceAPIRequest.options = options[0]
	}
	// 添加api key信息
	faceAPIRequest.options["api_key"] = sdk.APIKey
	faceAPIRequest.options["api_secret"] = sdk.APISecret

	faceAPIRequest.request = sdk.getHTTPRequest().
		Type("multipart")

	return faceAPIRequest, nil
}

// SetOption 设置请求参数
func (far *FaceAPIRequest) SetOption(key string, val interface{}) *FaceAPIRequest {
	far.options[key] = val
	return far
}

// SetOptionMap 通过map设置请求参数
func (far *FaceAPIRequest) SetOptionMap(options map[string]interface{}) *FaceAPIRequest {
	for key, val := range options {
		far.options[key] = val
	}
	return far
}
