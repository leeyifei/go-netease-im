package netease_im

import (
	"net/url"
)

const (
	USER_MESSAGE_HISTORY_PATH = "history/querySessionMsg.action"
)

type User_message_history_response struct {
	Size int `json:"size"`
	Msgs []Base_message_body `json:"msgs"`
	Netease_base_response
}

// 定义消息接口
type Base_message_interface interface {
	Get_message_type() int
}

// 消息基本结构
type Base_message_body struct {
	From string `json:"from"`
	Msgid int64 `json:"msgid"`
	Sendtime int64 `json:"sendtime"`
	Type int `json:"type"`
	Body map[string]interface{} `json:"body"`
}

// 获取消息的类型
func(this *Base_message_body) Get_message_type() int {
	return this.Type
}

// 用户之间的聊天记录
func (this *Netease_im) User_message_history (from string, to string, begintime string, endtime string, limit string, asc bool) (*User_message_history_response, error) {
	form_value := url.Values{}
	form_value.Set("from", from)
	form_value.Set("to", to)
	form_value.Set("begintime", begintime)
	form_value.Set("endtime", endtime)
	form_value.Set("limit", limit)
	if asc {
		form_value.Set("reverse", "1")
	} else {
		form_value.Set("reverse", "2")
	}
	response, err := this.request(USER_MESSAGE_HISTORY_PATH, form_value)
	if err != nil {
		return nil, err
	}

	dict, err := this.json_unserialize(new(User_message_history_response), response)
	if err != nil {
		return nil, err
	}
	return dict.(*User_message_history_response), nil
}
