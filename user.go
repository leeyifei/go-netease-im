package netease_im

import (
	"net/url"
	"fmt"
	//"encoding/json"
)

const (
	CREATE_USER_SUB_PATH = "user/create.action"
)

type Create_user_response struct {
	Info struct {
		Token string `json:"token"`
		Accid string `json:"accid"`
		Name string `json:"name"`
	} `json:"info"`
	Netease_base_response
}


//  create a new user
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