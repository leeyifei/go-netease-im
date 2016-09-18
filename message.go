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
type UploadResponse struct {
	NeteaseBaseResponse
	Url string `json:"url, omitempty"`
}

// 上传图片文件响应结构
type UploadImageResponse struct {
	UploadResponse
	Name string
	Ext string
	Md5 string
	W int
	H int
	Size int64
}

// 单条消息响应结构
type SendSingleMessageResponse struct {
	NeteaseBaseResponse
}

type Unregister struct {
	User string
}

// 多条消息响应结构
type SendMultimessageResponse struct {
	Unregister string `json:"unregister, omitempty"`
	NeteaseBaseResponse
}



// 流形式上传文件,TODO
func (this *NeteaseIm) UploadFile() {

}


// 表单形式上传文件
// @param string file_path 文件绝对路径
// @return1 *Upload_response
// @return2 error
func (this *NeteaseIm) UploadImageMultipart(filePath string) (*UploadImageResponse, error) {
	var fileSize int64 = 0
	var fileWidth, fileHeight int = 0, 0
	var fileName, fileExt, fileMd5 string = "", "", ""
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 获取文件内容 md5
	fileContent, err := ioutil.ReadFile(filePath)
	if err == nil {
		encoder := md5.New()
		encoder.Write(fileContent)
		fileMd5 = hex.EncodeToString(encoder.Sum(nil))
	}

	// 获取文件大小
	fs, err := file.Stat()
	if err != nil {
		fileSize = 0
	}
	// 文件大小
	fileSize = fs.Size()
	// 文件名
	fileName = fs.Name()
	// 扩展名
	fileExt = filepath.Ext(filePath)
	// 获取图片宽,高等信息
	im, _, err := image.DecodeConfig(file)
	if err == nil {
		fileWidth = im.Width
		fileHeight = im.Height
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("content", filepath.Base(filePath))
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

	this.getNonce(128)
	this.getCurTime()
	this.getChecksum()

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
	dict, err := this.json_unserialize(new(UploadImageResponse), string(data))
	if err != nil {
		return nil, err
	}
	image_info := dict.(*UploadImageResponse)
	image_info.Size = fileSize
	image_info.Name = fileName
	image_info.Ext = fileExt
	image_info.W = fileWidth
	image_info.H = fileHeight
	image_info.Md5 = fileMd5
	return image_info, nil
}

// 发送单条文本信息
// @param string from 发送者ID
// @param string to 接受者ID
// @param string msg 消息内容
// @param map[string]bool option 消息选项 {"push":false,"roam":true,"history":false,"sendersync":true,"route":false,"badge":false,"needPushNick":true}
// @param string pushcontent 推送内容
// @param string payload ios 推送对应的payload
// @return *Send_singleessage_response
// @return error
func (this *NeteaseIm) SendSingleTextMessage(from string, to string, msg string, option map[string]bool, pushcontent string, payload string) (*SendSingleMessageResponse, error) {
	body := make(map[string]interface{})
	body["msg"] = msg
	dict, err := this.sendSingleMessage(from, to, "0", body, option, pushcontent, payload)
	if err != nil {
		return nil, err
	}
	return dict, nil
}

// 发送单条图片消息
// @param string from
// @param string to
// @param string file_path
// @param map[string]bool option
// @param string pushcontent
// @param string payload
// @return *Send_singleessage_response
// @return error
func (this *NeteaseIm) SendSingleImageMessage(from string, to string, filePath string, option map[string]bool, pushcontent string, payload string) (*SendSingleMessageResponse, error) {
	// 上传图片文件
	uploadResp, err := this.UploadImageMultipart(filePath)
	if err != nil {
		return nil, err
	}

	if uploadResp.IsSuccess() != true {
		return nil, errors.New("upload file fail")
	}
	body := make(map[string]interface{})
	body["name"] = uploadResp.Name
	body["ext"] = uploadResp.Ext[1:]
	body["md5"] = uploadResp.Md5
	body["url"] = uploadResp.Url
	body["w"] = uploadResp.W
	body["h"] = uploadResp.H
	body["size"] = uploadResp.Size
	dict, err := this.sendSingleMessage(from, to, "1", body, option, pushcontent, payload)
	if err != nil {
		return nil, err
	}
	return dict, nil
}

// 批量发送多人文本消息
// @param string from
// @param []string to
// @param string msg
// @param map[string]bool option
// @param string pushcontent
// @param string payload
// @return *sendMultimessageResponse
// @return error
func (this *NeteaseIm) SendMultiTextMessage(from string, to []string, msg string, option map[string]bool, pushcontent string, payload string) (*SendMultimessageResponse, error) {
	body := make(map[string]interface{})
	body["msg"] = msg
	dict, err := this.sendMultiMessage(from, to, "0" ,body, option, pushcontent, payload)
	if err != nil {
		return nil, err
	}
	return dict, nil
}

// 批量发送多人图片消息
// @param string from
// @param []string to
// @param string filePath
// @param map[string]bool option
// @param string pushcontent
// @param string payload
// @return SendMultimessageResponse
// @return error
func (this *NeteaseIm) SendMultiImageMessage(from string, to []string, filePath string, option map[string]bool, pushcontent string, payload string) (*SendMultimessageResponse, error) {
	uploadResp, err := this.UploadImageMultipart(filePath)
	if err != nil {
		return nil, err
	}

	if uploadResp.IsSuccess() != true {
		return nil, errors.New("upload file fail")
	}

	body := make(map[string]interface{})
	body["name"] = uploadResp.Name
	body["ext"] = uploadResp.Ext[1:]
	body["md5"] = uploadResp.Md5
	body["url"] = uploadResp.Url
	body["w"] = uploadResp.W
	body["h"] = uploadResp.H
	body["size"] = uploadResp.Size
	dict, err := this.sendMultiMessage(from, to, "1", body, option, pushcontent, payload)
	if err != nil {
		return nil, err
	}
	return dict, nil
}

// 发送单条信息
// @param string from 发送者ID
// @param string to 接受者ID
// @param string msgType 消息类型 0 文本 1 图片 2 语音 3 视频 4 地理位置 6 文件 100 自定义
// @param map[string]interface{} body 消息体
// @param map[string]bool option
// @param string pushcontent
// @param payload ios 推送对应的payload
// @return SendSingleMessageResponse
// @return error
func (this *NeteaseIm) sendSingleMessage(from string, to string, msgType string, body map[string]interface{}, option map[string]bool, pushcontent string, payload string) (*SendSingleMessageResponse, error) {
	formValue := url.Values{}
	formValue.Add("from", from)
	formValue.Add("ope", "0")
	formValue.Add("to", to)
	formValue.Add("type", msgType)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Printf("Send_single_message body =====> %s", string(jsonBody))
	}
	formValue.Add("body", string(jsonBody))
	jsonOption, err := json.Marshal(option)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Print("Send_single_message option =====> %s", string(jsonOption))
	}
	formValue.Add("option", string(jsonOption))
	formValue.Add("pushcontent", pushcontent)
	formValue.Add("payload", payload)

	response, err := this.request(SEND_SINGLE_MESSAGE_PATH, formValue)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Print(fmt.Sprintf("Send_single_message result =====> %s\n", response))
	}

	dict, err := this.json_unserialize(new(SendSingleMessageResponse), response)
	if err != nil {
		return nil, err
	}
	return dict.(*SendSingleMessageResponse), nil
}

// 发送多条消息
// @param string from
// @param []string to
// @param string msgType
// @param map[string]interface{} body
// @param map[string]bool option
// @param string pushcontent
// @param string payload
// @return SendMultimessageResponse
// @return error
func (this *NeteaseIm) sendMultiMessage(from string, to []string, msgType string, body map[string]interface{}, option map[string]bool, pushcontent string, payload string) (*SendMultimessageResponse, error) {
	formValue := url.Values{}
	formValue.Add("fromAccid", from)
	jsonTo, err := json.Marshal(to)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Printf("send_multi_message to =====> %s\n", string(jsonTo))
	}
	formValue.Add("toAccids", string(jsonTo))
	formValue.Add("type", msgType)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Printf("send_multi_message body =====> %s\n", string(jsonBody))
	}
	formValue.Add("body", string(jsonBody))
	json_option, err := json.Marshal(option)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Print("Send_single_message option =====> %s", string(json_option))
	}
	formValue.Add("option", string(json_option))
	formValue.Add("pushcontent", pushcontent)
	formValue.Add("payload", payload)

	response, err := this.request(SEND_MULTI_MESSAGE_PATH, formValue)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Print(fmt.Sprintf("Send_single_message result =====> %s\n", response))
	}
	dict, err := this.json_unserialize(new(SendMultimessageResponse), response)
	if err != nil {
		return nil, err
	}
	return dict.(*SendMultimessageResponse), nil
}