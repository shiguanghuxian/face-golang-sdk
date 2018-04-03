package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/parnurzeal/gorequest"
)

const (
	APIBaseURL = "https://api-cn.faceplusplus.com/facepp/v3"
)

// FaceSDK Face++ sdk 对象
type FaceSDK struct {
	APIKey    string
	APISecret string
	Debug     bool // 是否调试
}

// FaceRequest 请求操作对象
type FaceRequest struct {
	options map[string]interface{}
	request *gorequest.SuperAgent
}

/**
 * NewFaceSDK 创建一个sdk对象，用于操作sdk
 * @param apiKey API Key
 * @param apiSecret API Secret
 * @param formal 是否是正式版 API Key
 */
func NewFaceSDK(apiKey, apiSecret string, debug ...bool) (*FaceSDK, error) {
	if apiKey == "" || apiSecret == "" {
		return nil, errors.New("API Key 和 API Secret 不能为空")
	}
	faceSDK := &FaceSDK{
		APIKey:    apiKey,
		APISecret: apiSecret,
	}
	if len(debug) > 0 {
		faceSDK.Debug = debug[0]
	}
	return faceSDK, nil
}

// 获取http请求对象
func (sdk *FaceSDK) getHTTPRequest() *gorequest.SuperAgent {
	superAgent := gorequest.New()
	superAgent.Debug = sdk.Debug
	return superAgent
}

// FaceResponse 接口返回数据结构体
type FaceResponse struct {
	RequestId    string `json:"request_id"`    // 用于区分每一次请求的唯一的字符串。此字符串可以用于后续数据反查。
	TimeUsed     int    `json:"time_used"`     // 整个请求所花费的时间，单位为毫秒。
	ErrorMessage string `json:"error_message"` // 当请求失败时才会返回此字符串，具体返回内容见后续错误信息章节。否则此字段不存在。
}

// FaceError 用于返回错误信息给调用方
type FaceError struct {
	Code         int    `json:"code"`          // 错误状态码
	ErrorMessage string `json:"error_message"` // 接口返回错误
	Message      string `json:"message"`       // 错误描述
}

// Error 输出错误信息为字符串
func (fe *FaceError) Error() string {
	return fmt.Sprintf("code:%d,error:%s,message:%s", fe.Code, fe.ErrorMessage, fe.Message)
}

// Face 数组中单个元素的结构
type Face struct {
	FaceToken     string `json:"face_token"` // 人脸的标识
	FaceRectangle struct {
		Top    int `json:"top"`
		Left   int `json:"left"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"face_rectangle"` // 人脸矩形框的位置，包括以下属性
	Landmark   map[string]*Landmark `json:"landmark"`   // 人脸的关键点坐标数组
	Attributes *Attributes          `json:"attributes"` // 人脸属性特征，具体包含的信息见下表
}

// Landmark 关键点坐标
type Landmark struct {
	X interface{}
	Y interface{}
}

// Attributes 人脸属性特征
type Attributes struct {
	Gender struct {
		Value string `json:"value"` // 性别值
	} `json:"gender"` // 性别 Male男|Female女
	Age struct {
		Value int `json:"value"` // 年龄值
	} `json:"age"` // 年龄分析结果。返回值为一个非负整数
	Smile struct {
		Value     float32 `json:"value"`     // 值为一个 [0,100] 的浮点数，小数点后3位有效数字。数值越大表示笑程度高
		Threshold float32 `json:"threshold"` // 代表笑容的阈值，超过该阈值认为有笑容
	} `json:"smile"` // 笑容分析结果
	Headpose struct {
		PitchAngle float64 `json:"pitch_angle"` // 抬头
		RollAngle  float64 `json:"roll_angle"`  // 旋转（平面旋转）
		YawAngle   float64 `json:"yaw_angle"`   // 摇头
	} `json:"headpose"` // 人脸姿势分析结果
	Eyestatus struct {
		LeftEyeStatus  map[string]float32 `json:"left_eye_status"`  // 左眼的状态
		RightEyeStatus map[string]float32 `json:"right_eye_status"` // 右眼的状态
	} `json:"eyestatus"` // 眼睛状态信息
	Emotion     map[string]float32 `json:"emotion"` // 情绪识别结果, map下标对应值 anger|愤怒 disgust|厌恶 fear|恐惧 happiness|高兴 neutral|平静 sadness|伤心 surprise|惊讶
	Facequality struct {
		Value     float32 `json:"value"`     // 值为人脸的质量判断的分数，是一个浮点数，范围 [0,100]，小数点后 3 位有效数字
		Threshold float32 `json:"threshold"` // 表示人脸质量基本合格的一个阈值，超过该阈值的人脸适合用于人脸比对
	} `json:"facequality"` // 人脸质量判断结果
	Ethnicity struct {
		Value string `json:"value"` // 人种值
	} `json:"ethnicity"` // 人种分析结果 Asian|亚洲人 White|白人 Black|黑人
	Beauty struct {
		MaleScore   float32 `json:"male_score"`   // 男性认为的此人脸颜值分数。值越大，颜值越高
		FemaleScore float32 `json:"female_score"` // 女性认为的此人脸颜值分数。值越大，颜值越高
	} `json:"beauty"` // 颜值识别结果
	Mouthstatus map[string]float32 `json:"mouthstatus"` // 嘴部状态信息
	Eyegaze     struct {
		LeftEyeGaze  map[string]float32 `json:"left_eye_gaze"`  // 左眼的位置与视线状态
		RightEyeGaze map[string]float32 `json:"right_eye_gaze"` // 右眼的位置与视线状态
	} `json:"eyegaze"` // 眼球位置与视线方向信息
	Skinstatus map[string]float32 `json:"skinstatus"` // 面部特征识别结果 health:健康 stain:色斑 acne:青春痘 dark_circle:黑眼圈
}

/**
 * NewFaceError 创建一个错误
 * @param code http 错误码
 * @param body 接口返回的body
 */
func NewFaceError(code int, body string) (err *FaceError) {
	var errorMessage = ""
	err = &FaceError{
		Code: code,
	}
	if code != http.StatusRequestEntityTooLarge {
		raceResponse := new(FaceResponse)
		er := json.Unmarshal([]byte(body), raceResponse)
		if er != nil {
			errorMessage = fmt.Sprintf("body解析错误:%s", string(body))
			err.ErrorMessage = errorMessage
			return err
		}
		errorMessage = raceResponse.ErrorMessage
	}

	err.ErrorMessage = errorMessage

	switch code {
	case http.StatusUnauthorized: // 401
		err.Message = "api_key和api_secret不匹配"
		break
	case http.StatusForbidden: // 403
		if errorMessage == "CONCURRENCY_LIMIT_EXCEEDED" {
			err.Message = "并发数超过限制"
		} else {
			msgs := strings.Split(errorMessage, ":")
			if len(msgs) != 2 {
				err.Message = "api_key没有调用本API的权限"
			} else {
				if msgs[1] == "Denied by Client" {
					err.Message = "用户自己禁止该api_key调用"
				} else if msgs[1] == "Denied by Admin" {
					err.Message = "管理员禁止该api_key调用"
				} else {
					err.Message = "由于账户余额不足禁止调用"
				}
			}
		}
		break
	case http.StatusBadRequest: // 400
		if errorMessage == "COEXISTENCE_ARGUMENTS" {
			err.Message = "同时传入了要求是二选一或多选一的参数，如有特殊说明则不返回此错误"
		} else {
			msgs := strings.Split(errorMessage, ":")
			if len(msgs) == 2 {
				if msgs[0] == "MISSING_ARGUMENTS" {
					err.Message = fmt.Sprintf("%s:%s", "缺少某个必选参数", msgs[1])
				} else {
					err.Message = fmt.Sprintf("%s:%s", "某个参数解析出错", msgs[1])
				}
			}
		}
		break
	case http.StatusRequestEntityTooLarge: // 413
		err.Message = "客户发送的请求大小超过了2MB限制，该错误的返回格式为纯文本，不是json格式"
		break
	case http.StatusNotFound: // 404
		err.Message = "所调用的API不存在"
		break
	case http.StatusInternalServerError: // 500
		err.Message = "服务器内部错误，当此类错误发生时请再次请求，如果持续出现此类错误，请及时联系技术支持团队。"
		break
	default:
		err.Message = "未知错误类型"
	}
	return err
}
