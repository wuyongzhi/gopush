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
	"crypto/md5"
	"encoding/hex"
)

const JPushServerUrl string = "http://api.jpush.cn:8800/v2/push"
const JPushServerUrlSsl string = "https://api.jpush.cn:443/v2/push"

type Request struct {
	url.Values
}

type Message struct {
	BuilderId int		`json:"n_builder_id"`
	Title string		`json:"n_title"`
	Content string		`json:"n_content"`
	Extras string		`json:"n_extras"`
}

type JPushMsgId interface {}

var InvalidMsgId JPushMsgId = nil

type Response struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgId   JPushMsgId `json:"msg_id"`
}

type Response2 struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

}


func NewRequest() *Request {
	m := Request{}
	m.Values = make(map[string][]string, 8)
	return &m
}

func (r *Response) IsOk() bool {
	return r.ErrCode == 0
}

func (r *Response) IsFailed() bool {
	return r.ErrCode != 0
}

func (m *Request) Set(key, value string){
	m.Values.Set(key, value)
}

func (m *Request) SetInt(key string, value int){
	m.Set(key, strconv.Itoa(value))
}

func (m *Request) SendNo(sendno int){
	m.SetInt("sendno", sendno)
}

func (m *Request) AppKey(app_key string) {
	m.Set("app_key", app_key)
}

const (
	ReceiverTypeTag            int = 2
	ReceiverTypeAlias              = 3
	ReceiverTypeBoardcast          = 4
	ReceiverTypeRegistrationID     = 5
)

const (
	MsgTypeNotify = 1
	MsgTypeCustom = 2
)

//	可以是以下值:
//		ReceiverTypeAlias
// 		ReceiverTypeTag
// 		ReceiverTypeBoardcast
// 		ReceiverTypeRegistrationID
func (m *Request) ReceiverType(receiver_type int)  {
	m.SetInt("receiver_type", receiver_type)
}

func (m *Request) ReceiverValue(receiver_values ...string) {
	m.Set("receiver_value", strings.Join(receiver_values, ","))
}

//允许传递认证码自行认证，也可以在调用Send 时，传递有效的 master_secret 参数来生成认证码
func (m *Request) VerificationCode(verification_code string)  {
	m.Set("verification_code", verification_code)
}

//可以是以下值：
//
// 	MsgTypeNotify
// 	MsgTypeCustom
func (m *Request) MsgType(msg_type int)  {
	m.SetInt("msg_type", msg_type)
}

func (m *Request) MsgContent(n_builder_id int, n_title, n_content, n_extras string) {
	msg := Message{n_builder_id, n_title, n_content, n_extras}
	bytes, _ := json.Marshal(msg)
	m.Set("msg_content", string(bytes))
}

func (m *Request) SendDescription(send_description string)  {
	m.Set("send_description", send_description)
}

//按可变参数，挨个传递“平台”，方法会用逗号将它们拼起来
func (m *Request) Platform(platforms ...string)  {
	m.Set("platform", strings.Join(platforms, ","))
}

// 仅IOS  适用 0 开发环境；1 生产环境
func (m *Request) APNSProduction(apns_production int)  {
	m.SetInt("apns_production", apns_production)
}

func (m *Request) TimeToLive(time_to_live int)  {
	m.SetInt("time_to_live", time_to_live)
}

func (m *Request) OverrideMsgId(override_msg_id string)  {
	m.Set("override_msg_id", override_msg_id)
}

func (m *Request) Sign(master_secret string)  {
	src := m.Values.Get("sendno") + m.Values.Get("receiver_type") + m.Values.Get("receiver_value") + master_secret
//	fmt.Println(src)
	sum := md5.Sum([]byte(src))
	verification_code := hex.EncodeToString(sum[:])
	m.Values.Set("verification_code", verification_code)

}

func (m *Request) send(url string) (*Response, error) {


	resp, err := defaultHttpClient.PostForm(url, m.Values)
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
		return nil, errors.New(err.Error() + " response: \n" + string(bytes))
	}

	// 如果失败，转换为 go 的 error
	if jpushResponse.IsFailed() {
		return &jpushResponse, errors.New(strconv.Itoa(jpushResponse.ErrCode) + ", " + jpushResponse.ErrMsg)
	}

	return &jpushResponse, nil

}


// 使用 http 协议
func (m *Request) Send() (*Response, error) {
	return m.send(JPushServerUrl)
}

// 使用 https 协议
func (m *Request) SendSecure() (*Response, error) {
	return m.send(JPushServerUrlSsl)
}

var defaultHttpClient *utils.HttpClient

func init() {
	timeout, _ := time.ParseDuration("10s")
	defaultHttpClient = utils.NewHttpClient(20, timeout, timeout, false)
	//defaultHttpClient.
}
