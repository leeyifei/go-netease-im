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

type Netease_response_interface interface {
	Is_success() bool
	Fail_code() int
	Fail_reason() string
}

type Netease_base_response struct {
	Code int `json:"code"`
	Desc string `json:"desc, omitempty"`
}

func (this *Netease_base_response) Is_success() bool {
	return this.Code == 200
}

func (this *Netease_base_response) Fail_code() int {
	return this.Code
}

func (this *Netease_base_response) Fail_reason() string {
	return this.Desc
}

type Netease_im struct {
	AppKey    string
	AppSecret string
	Debug bool

	nonce string
	curtime string
	checksum string
}

// construct function
func Netease_im_instance(app_key string, app_secret string, debug bool) *Netease_im {
	instance := new(Netease_im)
	instance.AppKey = app_key
	instance.AppSecret = app_secret
	instance.Debug = debug
	return instance
}


func (this *Netease_im) request(sub_path string, form_values url.Values) (string, error) {
	body := ioutil.NopCloser(strings.NewReader(form_values.Encode()))
	client := &http.Client{}
	req, _ := http.NewRequest("POST", EASE_IM_HOST + sub_path, body)
	req.Header.Add("Content-Type", REQUEST_CONTENT_TYPE)

	this.get_nonce( 128 )
	this.get_cur_time()
	this.get_checksum()

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


func (this *Netease_im) json_unserialize(model Netease_response_interface, json_str string) (Netease_response_interface, error) {
	err := json.Unmarshal([]byte(json_str), &model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

//generate nonce string
func (this *Netease_im) get_nonce(length int) {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	this.nonce = string(result)
	if this.Debug {
		fmt.Print(fmt.Sprintf("%s\n", this.nonce))
	}
}


//generate UTC time
func (this *Netease_im) get_cur_time() {
	var cur_str string
	cur := time.Now().Unix()
	cur_str = strconv.FormatInt(cur, 10)
	this.curtime = cur_str
	if this.Debug {
		fmt.Print(fmt.Sprintf("%s\n", this.curtime))
	}
}


//generate sha1 hex string
func (this *Netease_im) get_checksum() {
	init_string := this.AppSecret + this.nonce + this.curtime
	sha1_encoder := sha1.New()
	sha1_encoder.Write([]byte(init_string))
	binary_string := sha1_encoder.Sum(nil)
	this.checksum = fmt.Sprintf("%x", binary_string)
	if this.Debug {
		fmt.Print(fmt.Sprintf("%s\n", this.checksum))
	}

}