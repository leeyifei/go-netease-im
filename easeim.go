package netease_im

import (
	"time"
	"math/rand"
	"crypto/sha1"
	"fmt"
	"net/http"
	"strings"
	"net/url"
	"io/ioutil"
	"strconv"
	"encoding/json"
)

const (
	EASE_IM_HOST = "https://api.netease.im/nimserver/"
	REQUEST_CONTENT_TYPE = "application/x-www-form-urlencoded"
)

type NeteaseResponseInterface interface {
	IsSuccess() bool
	FailCode() int
	FailReason() string
}

type NeteaseBaseResponse struct {
	Code int `json:"code"`
	Desc string `json:"desc, omitempty"`
}

// 请求是否成功
func (this *NeteaseBaseResponse) IsSuccess() bool {
	return this.Code == 200
}

// 错误码
func (this *NeteaseBaseResponse) FailCode() int {
	return this.Code
}

// 错误内容
func (this *NeteaseBaseResponse) FailReason() string {
	return this.Desc
}

type NeteaseIm struct {
	AppKey    string
	AppSecret string
	Debug bool

	nonce string
	curtime string
	checksum string
}

// 构造函数
func NeteaseImInstance(appKey string, appSecret string, debug bool) *NeteaseIm {
	instance := new(NeteaseIm)
	instance.AppKey = appKey
	instance.AppSecret = appSecret
	instance.Debug = debug
	return instance
}

// 请求云信API
// @param string subPath
// @param url.Values formValues
// @return string
// @return error
func (this *NeteaseIm) request(subPath string, formValues url.Values) (string, error) {
	body := ioutil.NopCloser(strings.NewReader(formValues.Encode()))
	client := &http.Client{}
	req, _ := http.NewRequest("POST", EASE_IM_HOST + subPath, body)
	req.Header.Add("Content-Type", REQUEST_CONTENT_TYPE)

	this.getNonce( 128 )
	this.getCurTime()
	this.getChecksum()

	req.Header.Add("AppKey", this.AppKey)
	req.Header.Add("Nonce", this.nonce)
	req.Header.Add("CurTime", this.curtime)
	req.Header.Add("CheckSum", this.checksum)

	response, err := client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if this.Debug {
		fmt.Println(string(data), err)
	}
	return string(data), nil
}


// 云信JSON解析
// @param NeteaseResponseInterface model
// @param string jsonStr
// @return *NeteaseResponseInterface
// @return error
func (this *NeteaseIm) json_unserialize(model NeteaseResponseInterface, jsonStr string) (NeteaseResponseInterface, error) {
	err := json.Unmarshal([]byte(jsonStr), &model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

// 生成指定长度的随机数
// @param int length
func (this *NeteaseIm) getNonce(length int) {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	this.nonce = string(result)
}

// 获取当前时间戳
func (this *NeteaseIm) getCurTime() {
	var curStr string
	cur := time.Now().Unix()
	curStr = strconv.FormatInt(cur, 10)
	this.curtime = curStr
}


// 生成签名
func (this *NeteaseIm) getChecksum() {
	initString := this.AppSecret + this.nonce + this.curtime
	sha1Encoder := sha1.New()
	sha1Encoder.Write([]byte(initString))
	binaryString := sha1Encoder.Sum(nil)
	this.checksum = fmt.Sprintf("%x", binaryString)
}