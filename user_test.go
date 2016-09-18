package netease_im

import (
	"testing"
	"fmt"
)

//var netEase * NeteaseIm

func init() {
	netEase = NeteaseImInstance("71b6da2c2393f0717b9f905d9ad0b9ac", "9443d5af1682", true)
}


func TestCreateUser(t *testing.T) {
	response, err := netEase.CreateUser("3", "徐雯雯", "", "", "")
	if err != nil {
		fmt.Printf("test TestCreate_user occur error =====> %s\n", err.Error())
	}
	if response.IsSuccess() {
		fmt.Printf("new id ====> %s\n", response.Info.Accid)
		fmt.Printf("new name ====> %s\n", response.Info.Name)
		fmt.Printf("new token ====> %s\n", response.Info.Token)
	} else {
		fmt.Printf("fail code ====> %d\n", response.FailCode())
		fmt.Printf("fail reason ====> %s\n", response.FailReason())
	}
}

func TestRefreshToken(t *testing.T) {
	response, err := netEase.RefreshToken("1")
	if err != nil {
		fmt.Printf("test TestRefresh_token occur error =====> %s\n", err.Error())
	}
	if response.IsSuccess() {
		fmt.Printf("new id ====> %s\n", response.Info.Accid)
		fmt.Printf("new token ====> %s\n", response.Info.Token)
	} else {
		fmt.Printf("fail code ====> %d\n", response.FailCode())
		fmt.Printf("fail reason ====> %s\n", response.FailReason())
	}
}

func TestBanUser(t *testing.T) {
	response, err := netEase.BanUser("1", false)
	if err != nil {
		fmt.Printf("test TestBan_user occur error =====> %s\n", err.Error())
	} else {
		if response.IsSuccess() {
			fmt.Printf("test TestBan_user success\n")
		} else {
			fmt.Printf("test TestBan_user fail\n")
			fmt.Printf("fail code ====> %d\n", response.FailCode())
			fmt.Printf("fail reason ====> %s\n", response.FailReason())
		}
	}
}

func TestUnbanUser(t *testing.T) {
	response, err := netEase.UnbanUser("1")
	if err != nil {
		fmt.Printf("test TestUnban_user occur error =====> %s\n", err.Error())
	} else {
		if response.IsSuccess() {
			fmt.Printf("test TestUnban_user success\n")
		} else {
			fmt.Printf("test TestUnban_user fail\n")
			fmt.Printf("fail code ====> %d\n", response.FailCode())
			fmt.Printf("fail reason ====> %s\n", response.FailReason())
		}
	}
}