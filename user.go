package netease_im

import (
	"net/url"
	"fmt"
	//"encoding/json"
)

const (
	CREATE_USER_SUB_PATH = "user/create.action"
	BAN_USER_PATH = "user/block.action"
	UNBAN_USER_PATH = "user/unblock.action"
	REFRESH_TOKEN_PATH = "user/refreshToken.action"
)

// 创建用户响应结构
type Create_user_response struct {
	Info struct {
		Token string `json:"token"`
		Accid string `json:"accid"`
		Name string `json:"name"`
	} `json:"info"`
	Netease_base_response
}

// 刷新用户TOKEN响应结构
type Refresh_token_response struct {
	Info struct {
		Token string `json:"token"`
		Accid string `json:"accid"`
	} `json:"info"`
	Netease_base_response
}

// 封禁用户响应结构
type Ban_user_response struct {
	Netease_base_response
}

// 创建一个新用户
// @param string id
// @param string nick
// @param string props
// @param string icon
// @param string token
// @return *Create_user_response
// @return error
func (this *Netease_im) Create_user (id string, nick string, props string, icon string, token string) (*Create_user_response, error) {
	form_value := url.Values{}
	form_value.Set("accid", id)
	form_value.Set("name", nick)
	form_value.Set("props", props)
	form_value.Set("icon", icon)
	form_value.Set("token", token)

	response, err := this.request(CREATE_USER_SUB_PATH, form_value)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Print(fmt.Sprintf("%s\n", response))
	}
	dict, err := this.json_unserialize(new(Create_user_response), response)
	if err != nil {
		return nil, err
	}
	return dict.(*Create_user_response), nil
}

// 更新并获取用户新TOKEN
// @param string accid
// @return *Refresh_token_response
// @return error
func (this *Netease_im) Refresh_token (accid string) (*Refresh_token_response, error) {
	form_value := url.Values{}
	form_value.Set("accid", accid)

	response, err := this.request(REFRESH_TOKEN_PATH, form_value)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Print(fmt.Sprintf("%s\n", response))
	}
	dict, err := this.json_unserialize(new(Refresh_token_response), response)
	if err != nil {
		return nil, err
	}
	return dict.(*Refresh_token_response), nil
}

// 封禁一个用户
// @param string accid
// @param string needkick
// @return *Ban_user_response
// @return error
func (this *Netease_im) Ban_user (accid string, needkick bool) (*Ban_user_response, error) {
	form_value := url.Values{}
	form_value.Set("accid", accid)
	if needkick {
		form_value.Set("needkick", "true")
	} else {
		form_value.Set("needkick", "false")
	}

	response, err := this.request(BAN_USER_PATH, form_value)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Printf("%s\n", response)
	}
	dict, err := this.json_unserialize(new(Ban_user_response), response)
	if err != nil {
		return nil, err
	}
	return dict.(*Ban_user_response), nil
}

// 解封一个用户
// @param string accid
// @return *Ban_user_response
// @return error
func (this *Netease_im) Unban_user (accid string) (*Ban_user_response, error) {
	form_value := url.Values{}
	form_value.Set("accid", accid)

	response, err := this.request(UNBAN_USER_PATH, form_value)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Printf("%s\n", response)
	}
	dict, err := this.json_unserialize(new(Ban_user_response), response)
	if err != nil {
		return nil, err
	}
	return dict.(*Ban_user_response), nil
}