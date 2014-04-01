package xinge

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
	"log"
)

type ResultCommon struct {
	Code   int    `json:"ret_code"`
	ErrMsg string `json:"err_msg"`
	//	Result string	`json:"result"`
}

type PushTagsResult struct {
	ResultCommon
	Result struct {
		PushId string `json:"push_id"`
	}
}

type Engine struct {
	AccessId  int64
	SecretKey string
}

type Message struct {
	Title         string                 `json:"title"`
	Content       string                 `json:"content"`
	CustomContent map[string]interface{} `json:"custom_content"`
}

type Request struct {
	url.Values
	m      *Message
	secret string
}

func NewRequest() *Request {
	r := Request{}
	r.Values = make(map[string][]string)
	return &r
}

func (me *Request) SetMessage(m *Message) {
	me.m = m
}
func (me *Request) GetMessage() *Message {
	if me.m == nil {
		me.m = new(Message)
	}
	return me.m
}
func (me *Request) SetSecret(secret string) {
	me.secret = secret
}
func (me *Request) SetTitle(title string) {
	me.GetMessage().Title = title
}
func (me *Request) SetContent(content string) {
	me.GetMessage().Content = content
}
func (me *Request) SetCustomContent(key string, value interface{}) {
	me.GetMessage().CustomContent[key] = value
}

func (me *Request) SetAccessId(accessId int64) {
	me.Set("access_id", strconv.FormatInt(accessId, 10))
}
func (me *Request) SetTimestamp(timestamp int64) {
	me.Set("timestamp", strconv.FormatInt(timestamp, 10))
}
func (me *Request) SetValidTime(valid_time int) {
	me.Set("valid_time", strconv.Itoa(valid_time))
}

func (me *Request) SetMessageType(messageType int) {
	me.Set("message_type", strconv.Itoa(messageType))
}
func (me *Request) SetExpireTime(expireTime int) {
	me.Set("expire_time", strconv.Itoa(expireTime))
}
func (me *Request) SetMultiPkg(multiPkg int) {
	me.Set("multi_pkg", strconv.Itoa(multiPkg))
}
func (me *Request) SetSendTime(sendTime time.Time) {
	me.Set("send_time", sendTime.Format("2006-01-02 15:04:05"))
}
func (me *Request) SetEnvironment(environment int) {
	me.Set("environment", strconv.Itoa(environment))
}

var Host = "openapi.xg.qq.com"
var URI_SEND_TAGS = "/v2/push/tags_device"
var URI_SEND_DEVICE = "/v2/push/single_device"
var Method = "GET"

func (me *Request) prepare() {

	// 消息内容
	message, _ := json.Marshal(me.m)
	me.Set("message", string(message))

	// 如果没有timestamp参数，添加一个
	_, ok := me.Values["timestamp"]
	if !ok {
		now := time.Now()
		me.SetTimestamp(now.Unix())
	}

}

func (me *Request) PushTagsAnd(tags ...string) (*PushTagsResult, error) {

	//
	// 准备必须参数
	//
	me.prepare()

	// tags 操作类型
	me.Set("tags_op", "AND")

	// tags 列表
	tags_list, _ := json.Marshal(tags)
	me.Set("tags_list", string(tags_list))

	me.computeSign(URI_SEND_TAGS)

	var resp *http.Response
	var err error
	if Method == "POST" {
		resp, err = http.PostForm("http://"+Host+URI_SEND_TAGS, me.Values)
	} else {
		queryString := me.Values.Encode()
		url := "http://"+Host+URI_SEND_TAGS + "?" + queryString
		log.Println(url)
		resp, err = http.Get(url)
	}

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	resultbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := PushTagsResult{}
	err = json.Unmarshal(resultbytes, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// 计算签名
func (me *Request) computeSign(uri string) {

	keys := make([]string, 0, len(me.Values)+5)
	for k, _ := range me.Values {
		keys = append(keys, k)
	}

	//排序key
	sort.Strings(keys)

	buf := bytes.NewBufferString(Method + Host + uri)
	for _, k := range keys {
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(me.Get(k))
	}

	buf.WriteString(me.secret)

	src := buf.Bytes()
	log.Println("src=", string(src))
	result := md5.Sum(src)
	sign := hex.EncodeToString(result[0:])
	log.Println("sign=", sign)
	me.Set("sign", sign)
}

func (me *Request) SendTagsOr(tags ...string) {

}

func SendAccount(accessId int64, validTime int, secretKey string) {

}

func SendTags(accessId int64, validTime int, secretKey string) {

}

func SendAll() {

}
