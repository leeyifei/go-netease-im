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
type CreateUserResponse struct {
	Info struct {
		Token string `json:"token"`
		Accid string `json:"accid"`
		Name string `json:"name"`
	} `json:"info"`
	NeteaseBaseResponse
}

// 刷新用户TOKEN响应结构
type RefreshTokenResponse struct {
	Info struct {
		Token string `json:"token"`
		Accid string `json:"accid"`
	} `json:"info"`
	NeteaseBaseResponse
}

// 封禁用户响应结构
type BanUserResponse struct {
	NeteaseBaseResponse
}

// 创建一个新用户
// @param string id
// @param string nick
// @param string props
// @param string icon
// @param string token
// @return *CreateUserResponse
// @return error
func (this *NeteaseIm) CreateUser (id string, nick string, props string, icon string, token string) (*CreateUserResponse, error) {
	formValue := url.Values{}
	formValue.Set("accid", id)
	formValue.Set("name", nick)
	formValue.Set("props", props)
	formValue.Set("icon", icon)
	formValue.Set("token", token)

	response, err := this.request(CREATE_USER_SUB_PATH, formValue)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Print(fmt.Sprintf("%s\n", response))
	}
	dict, err := this.json_unserialize(new(CreateUserResponse), response)
	if err != nil {
		return nil, err
	}
	return dict.(*CreateUserResponse), nil
}

// 更新并获取用户新TOKEN
// @param string accid
// @return *RefreshTokenResponse
// @return error
func (this *NeteaseIm) RefreshToken (accid string) (*RefreshTokenResponse, error) {
	formValue := url.Values{}
	formValue.Set("accid", accid)

	response, err := this.request(REFRESH_TOKEN_PATH, formValue)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Print(fmt.Sprintf("%s\n", response))
	}
	dict, err := this.json_unserialize(new(RefreshTokenResponse), response)
	if err != nil {
		return nil, err
	}
	return dict.(*RefreshTokenResponse), nil
}

// 封禁一个用户
// @param string accid
// @param string needkick
// @return *BanUserResponse
// @return error
func (this *NeteaseIm) BanUser (accid string, needkick bool) (*BanUserResponse, error) {
	formValue := url.Values{}
	formValue.Set("accid", accid)
	if needkick {
		formValue.Set("needkick", "true")
	} else {
		formValue.Set("needkick", "false")
	}

	response, err := this.request(BAN_USER_PATH, formValue)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Printf("%s\n", response)
	}
	dict, err := this.json_unserialize(new(BanUserResponse), response)
	if err != nil {
		return nil, err
	}
	return dict.(*BanUserResponse), nil
}

// 解封一个用户
// @param string accid
// @return *BanUserResponse
// @return error
func (this *NeteaseIm) UnbanUser (accid string) (*BanUserResponse, error) {
	formValue := url.Values{}
	formValue.Set("accid", accid)

	response, err := this.request(UNBAN_USER_PATH, formValue)
	if err != nil {
		return nil, err
	}
	if this.Debug {
		fmt.Printf("%s\n", response)
	}
	dict, err := this.json_unserialize(new(BanUserResponse), response)
	if err != nil {
		return nil, err
	}
	return dict.(*BanUserResponse), nil
}