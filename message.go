package netease_im

import (
	"net/url"
	"encoding/json"
	"fmt"
	"os"
	"bytes"
	"mime/multipart"
	"path/filepath"
	"io"
	"net/http"
	"io/ioutil"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	_ "image/gif"
	"crypto/md5"
	"encoding/hex"
)

const (
	SEND_SINGLE_MESSAGE_PATH = "msg/sendMsg.action"
	SEND_MULTI_MESSAGE_PATH = "msg/sendBatchMsg.action"
	UPLOAD_FILE_PATH = "msg/upload.action"
	UPLOAD_FILE_MULTIPART_PATH = "msg/fileUpload.action"
)

// 上传文件响应结构
type Upload_response struct {
	Netease_base_response
	Url string `json:"url, omitempty"`
}

// 上传图片文件响应结构
type Upload_image_response struct {
	Upload_response
	Name string
	Ext string
	Md5 string
	W int
	H int
	Size int64
}

// 单条消息响应结构
type Send_singleessage_response struct {
	Netease_base_response
}

// 多条消息响应结构
type Send_multimessage_response struct {
	Unregsiter []string `json:"unregister, omitempty"`
	Netease_base_response
}



// 流形式上传文件,TODO
func (this *Netease_im) Upload_file() {

}


// 上传文件
// @param file_path 文件绝对路径
// @return1 Upload_response
// @return2 Error
func (this *Netease_im) Upload_image_multipart(file_path string) (*Upload_image_response, error) {
	var file_size int64 = 0
	var file_width, file_height int = 0, 0
	var file_name, file_ext, file_md5 string = "", "", ""
	file, err := os.Open(file_path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 获取文件内容 md5
	file_content, err := ioutil.ReadFile(file_path)
	if err == nil {
		encoder := md5.New()
		encoder.Write(file_content)
		file_md5 = hex.EncodeToString(encoder.Sum(nil))
	}

	// 获取文件大小
	fs, err := file.Stat()
	if err != nil {
		file_size = 0
	}
	// 文件大小
	file_size = fs.Size()
	// 文件名
	file_name = fs.Name()
	// 扩展名
	file_ext = filepath.Ext(file_path)
	// 获取图片宽,高等信息
	im, _, err := image.DecodeConfig(file)
	if err == nil {
		file_width = im.Width
		file_height = im.Height
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("content", filepath.Base(file_path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	request, _ := http.NewRequest("POST", EASE_IM_HOST + UPLOAD_FILE_MULTIPART_PATH, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	this.get_nonce(128)
	this.get_cur_time()
	this.get_checksum()

	request.Header.Add("AppKey", this.AppKey)
	request.Header.Add("Nonce", this.nonce)
	request.Header.Add("CurTime", this.curtime)
	request.Header.Add("CheckSum", this.checksum)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	dict, err := this.json_unserialize(new(Upload_image_response), string(data))
	if err != nil {
		return nil, err
	}
	image_info := dict.(*Upload_image_response)
	image_info.Size = file_size
	image_info.Name = file_name
	image_info.Ext = file_ext
	image_info.W = file_width
	image_info.H = file_height
	image_info.Md5 = file_md5
	return image_info, nil
}

// 发送单条文本信息
// @param from 发送者ID
// @param to 接受者ID
// @param msg 消息内容
// @param option 消息选项 {"push":false,"roam":true,"history":false,"sendersync":true,"route":false,"badge":false,"needPushNick":true}
// @param pushcontent 推送内容
// @param payload ios 推送对应的payload
func (this *Netease_im) Send_single_text_message(from string, to string, msg string, option map[string]bool, pushcontent string, payload string) (*Send_singleessage_response, error) {
	body := make(map[string]interface{})
	body["msg"] = msg
	dict, err := this.send_single_message(from, to, "0", body, option, pushcontent, payload)
	if err != nil {
		return nil, err
	}
	return dict, nil
}

// 发送单条图片消息
func (this *Netease_im) Send_single_image_message(from string, to string, file_path string, option map[string]bool, pushcontent string, payload string) (*Send_singleessage_response, error) {
	// 上传图片文件
	upload_resp, err := this.Upload_image_multipart(file_path)
	if err != nil {
		return nil, err
	}

	if upload_resp.Is_success() != true {
		return nil, errors.New("upload file fail")
	}
	body := make(map[string]interface{})
	body["name"] = upload_resp.Name
	body["ext"] = upload_resp.Ext[1:]
	body["md5"] = upload_resp.Md5
	body["url"] = upload_resp.Url
	body["w"] = upload_resp.W
	body["h"] = upload_resp.H
	body["size"] = upload_resp.Size
	dict, err := this.send_single_message(from, to, "1", body, option, pushcontent, payload)
	if err != nil {
		return nil, err
	}
	return dict, nil
}

// 发送单条信息
// @param from 发送者ID
// @param to 接受者ID
// @param msg_type 消息类型 0 文本 1 图片 2 语音 3 视频 4 地理位置 6 文件 100 自定义
// @param body 消息体
// @param payload ios 推送对应的payload
func (this *Netease_im) send_single_message(from string, to string, msg_type string, body map[string]interface{}, option map[string]bool, pushcontent string, payload string) (*Send_singleessage_response, error) {
	form_value := url.Values{}
	form_value.Add("from", from)
	form_value.Add("ope", "0")
	form_value.Add("to", to)
	form_value.Add("type", msg_type)
	json_body, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Printf("Send_single_message body =====> %s", string(json_body))
	}
	form_value.Add("body", string(json_body))
	json_option, err := json.Marshal(option)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Print("Send_single_message option =====> %s", string(json_option))
	}
	form_value.Add("option", string(json_option))
	form_value.Add("pushcontent", pushcontent)
	form_value.Add("payload", payload)

	response, err := this.request(SEND_SINGLE_MESSAGE_PATH, form_value)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Print(fmt.Sprintf("Send_single_message result =====> %s\n", response))
	}

	dict, err := this.json_unserialize(new(Send_singleessage_response), response)
	if err != nil {
		return nil, err
	}
	return dict.(*Send_singleessage_response), nil
}

// 发送多条消息
func (this *Netease_im) send_multi_message(from string, to []string, msg_type string, body map[string]interface{}, option map[string]bool, pushcontent string, payload string) (*Send_multimessage_response, error) {
	form_value := url.Values{}
	form_value.Add("fromAccid", from)
	json_to, err := json.Marshal(to)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Printf("send_multi_message to =====> %s\n", string(json_to))
	}
	form_value.Add("toAccids", string(json_to))
	form_value.Add("type", msg_type)
	json_body, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Printf("send_multi_message body =====> %s\n", string(json_body))
	}
	form_value.Add("body", string(json_body))
	json_option, err := json.Marshal(option)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Print("Send_single_message option =====> %s", string(json_option))
	}
	form_value.Add("option", string(json_option))
	form_value.Add("pushcontent", pushcontent)
	form_value.Add("payload", payload)

	response, err := this.request(SEND_MULTI_MESSAGE_PATH, form_value)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Print(fmt.Sprintf("Send_single_message result =====> %s\n", response))
	}
	dict, err := this.json_unserialize(new(Send_multimessage_response), response)
	if err != nil {
		return nil, err
	}
	return dict.(*Send_multimessage_response), nil
}