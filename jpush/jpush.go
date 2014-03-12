package jpush

import (
	"github.com/wuyongzhi/gopush/utils"
	"time"
	"net/url"
)


const JPushServerUrl  string = "http://api.jpush.cn:8800/v2/push"

type Message url.Values




var defaultHttpClient *utils.HttpClient

func init() {
	timeout, _ := time.ParseDuration("10s")


	defaultHttpClient = utils.NewHttpClient(20, timeout, timeout, false)
	defaultHttpClient.PostForm(JPushServerUrl, nil)
}




