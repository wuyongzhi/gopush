package jpush

import (
	"encoding/json"
	"errors"
	"github.com/wuyongzhi/gopush/utils"
	"io/ioutil"
	"net/http"
	_ "net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const JPushServerUrl string = "http://api.jpush.cn:8800/v2/push"

type Message struct {
	url.Values
}

type Response struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgId   string `json:"msg_id"`
}

func (r *Response) IsOk() bool {
	return r.ErrCode == 0
}

func (r *Response) IsFailed() bool {
	return r.ErrCode != 0
}


func (m *Message) Set(key, value string) *Message {
	m.Set(key, value)
	return m
}

func (m *Message) SetInt(key string, value int) *Message {
	return m.Set(key, strconv.Itoa(value))
}

func (m *Message) SendNo(sendno int) *Message {
	return m.SetInt("sendno", sendno)
}

func (m *Message) AppKey(app_key string) *Message {
	return m.Set("app_key", app_key)
}

const (
	ReceiverTypeTag            int = 2
	ReceiverTypeAlias              = 3
	ReceiverTypeBoardcast          = 4
	ReceiverTypeRegistrationID     = 5
)

const (
	MessageTypeNotify = 1
	MessageTypeCustom = 2
)

//	可以是以下值:
//		ReceiverTypeAlias
// 		ReceiverTypeTag
// 		ReceiverTypeBoardcast
// 		ReceiverTypeRegistrationID
func (m *Message) ReceiverType(receiver_type int) *Message {
	return m.SetInt("receiver_type", receiver_type)
}

func (m *Message) ReceiverValue(receiver_values ... string) *Message {
	return m.Set("receiver_value", strings.Join(receiver_values, ","))
}

//允许传递认证码自行认证，也可以在调用Send 时，传递有效的 master_secret 参数来生成认证码
func (m *Message) VerificationCode(verification_code string) *Message {
	return m.Set("verification_code", verification_code)
}

//可以是以下值：
//
// 	MessageTypeNotify
// 	MessageTypeCustom
func (m *Message) MsgType(msg_type int) *Message {
	return m.SetInt("msg_type", msg_type)
}

func (m *Message) MsgContent(msg_content string) *Message {
	return m.Set("msg_content", msg_content)
}

func (m *Message) SendDescription(send_description string) *Message {
	return m.Set("send_description", send_description)
}

//按可变参数，挨个传递“平台”，方法会用逗号将它们拼起来
func (m *Message) Platform(platforms ...string) *Message {
	return m.Set("platform", strings.Join(platforms, ","))
}

// 仅IOS  适用 0 开发环境；1 生产环境
func (m *Message) APNSProduction(apns_production int) *Message {
	return m.SetInt("apns_production", apns_production)
}

func (m *Message) TimeToLive(time_to_live int) *Message {
	return m.SetInt("time_to_live", time_to_live)
}

func (m *Message) OverrideMsgId(override_msg_id string) *Message {
	return m.Set("override_msg_id", override_msg_id)
}

func (m *Message) Sign(master_secret *string) bool {

	return true
}

func (m *Message) Send(master_secret *string) (*Response, error) {

	//进行认证
	if master_secret != nil {
		m.Sign(master_secret)
	}

	resp, err := defaultHttpClient.PostForm(JPushServerUrl, m.Values)
	if err != nil {
		return nil, err
	}


	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(resp.StatusCode) + resp.Status)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jpushResponse Response

	//	responseContent := string(bytes)
	err = json.Unmarshal(bytes, &jpushResponse)
	if err != nil {
		return nil, errors.New(err.Error() + " response: \n"  + string(bytes))
	}

	if jpushResponse.IsFailed() {
		return &jpushResponse, errors.New(strconv.Itoa(jpushResponse.ErrCode) + ", " + jpushResponse.ErrMsg)
	}

	return &jpushResponse, nil
}



var defaultHttpClient *utils.HttpClient

func init() {
	timeout, _ := time.ParseDuration("10s")
	defaultHttpClient = utils.NewHttpClient(20, timeout, timeout, false)
	//defaultHttpClient.
}
