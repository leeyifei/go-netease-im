package netease_im

import (
	"testing"
	"fmt"
)

//var net_ease * Netease_im

func init() {
	net_ease = Netease_im_instance("71b6da2c2393f0717b9f905d9ad0b9ac", "9443d5af1682", true)
}


func TestCreate_user(t *testing.T) {
	response, err := net_ease.Create_user("3", "徐雯雯", "", "", "")
	if err != nil {
		fmt.Printf("test TestCreate_user occur error =====> %s\n", err.Error())
	}
	if response.Is_success() {
		fmt.Printf("new id ====> %s\n", response.Info.Accid)
		fmt.Printf("new name ====> %s\n", response.Info.Name)
		fmt.Printf("new token ====> %s\n", response.Info.Token)
	} else {
		fmt.Printf("fail code ====> %d\n", response.Fail_code())
		fmt.Printf("fail reason ====> %s\n", response.Fail_reason())
	}
}

func TestRefresh_token(t *testing.T) {
	response, err := net_ease.Refresh_token("1")
	if err != nil {
		fmt.Printf("test TestRefresh_token occur error =====> %s\n", err.Error())
	}
	if response.Is_success() {
		fmt.Printf("new id ====> %s\n", response.Info.Accid)
		fmt.Printf("new token ====> %s\n", response.Info.Token)
	} else {
		fmt.Printf("fail code ====> %d\n", response.Fail_code())
		fmt.Printf("fail reason ====> %s\n", response.Fail_reason())
	}
}

func TestBan_user(t *testing.T) {
	response, err := net_ease.Ban_user("1", false)
	if err != nil {
		fmt.Printf("test TestBan_user occur error =====> %s\n", err.Error())
	} else {
		if response.Is_success() {
			fmt.Printf("test TestBan_user success\n")
		} else {
			fmt.Printf("test TestBan_user fail\n")
			fmt.Printf("fail code ====> %d\n", response.Fail_code())
			fmt.Printf("fail reason ====> %s\n", response.Fail_reason())
		}
	}
}

func TestUnban_user(t *testing.T) {
	response, err := net_ease.Unban_user("1")
	if err != nil {
		fmt.Printf("test TestUnban_user occur error =====> %s\n", err.Error())
	} else {
		if response.Is_success() {
			fmt.Printf("test TestUnban_user success\n")
		} else {
			fmt.Printf("test TestUnban_user fail\n")
			fmt.Printf("fail code ====> %d\n", response.Fail_code())
			fmt.Printf("fail reason ====> %s\n", response.Fail_reason())
		}
	}
}