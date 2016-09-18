package netease_im

import (
	"net/url"
)

const (
	USER_MESSAGE_HISTORY_PATH = "history/querySessionMsg.action"
)

type UserMessageHistoryResponse struct {
	Size int `json:"size"`
	Msgs []BaseMessageBody `json:"msgs"`
	NeteaseBaseResponse
}

// 定义消息接口
type BaseMessageInterface interface {
	GetMessageType() int
}

// 消息基本结构
type BaseMessageBody struct {
	From string `json:"from"`
	Msgid int64 `json:"msgid"`
	Sendtime int64 `json:"sendtime"`
	Type int `json:"type"`
	Body map[string]interface{} `json:"body"`
}

// 获取消息的类型
func(this *BaseMessageBody) GetMessageType() int {
	return this.Type
}

// 用户之间的聊天记录
func (this *NeteaseIm) UserMessageHistory (from string, to string, begintime string, endtime string, limit string, asc bool) (*UserMessageHistoryResponse, error) {
	formValue := url.Values{}
	formValue.Set("from", from)
	formValue.Set("to", to)
	formValue.Set("begintime", begintime)
	formValue.Set("endtime", endtime)
	formValue.Set("limit", limit)
	if asc {
		formValue.Set("reverse", "1")
	} else {
		formValue.Set("reverse", "2")
	}
	response, err := this.request(USER_MESSAGE_HISTORY_PATH, formValue)
	if err != nil {
		return nil, err
	}

	dict, err := this.json_unserialize(new(UserMessageHistoryResponse), response)
	if err != nil {
		return nil, err
	}
	return dict.(*UserMessageHistoryResponse), nil
}
