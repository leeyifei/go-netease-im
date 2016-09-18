package netease_im

import (
	"testing"
	"fmt"
	"encoding/json"
)

var netEase *NeteaseIm

func init() {
	netEase = NeteaseImInstance("71b6da2c2393f0717b9f905d9ad0b9ac", "9443d5af1682", true)
}

func TestUploadFileMultipart(t *testing.T) {
	response, err := netEase.UploadImageMultipart("/Users/leeyifiei/Pictures/111.png")
	if err != nil {
		fmt.Printf("Test Upload_file_multipart occur error =====> %s\n", err.Error())
	} else {
		if response.IsSuccess() {
			fmt.Printf("test upload_file_multipart success \n")
			fmt.Printf("file url is =====> %s\n", response.Url)
			fmt.Printf("file name is =====> %s\n", response.Name)
			fmt.Printf("file ext is =====> %s\n", response.Ext[1:])
			fmt.Printf("file size is =====> %d\n", response.Size)
			fmt.Printf("file width is =====> %d\n", response.W)
			fmt.Printf("file height is =====> %d\n", response.H)
			fmt.Printf("file md5 is =====> %s\n", response.Md5)
		} else {
			fmt.Printf("test upload_file_multipart fail code =====> %d\n", response.FailCode())
			fmt.Printf("test upload_file_multipart fail reason ====> %s\n", response.FailReason())
		}
	}

}

func TestSendSingleTextMessage(t *testing.T) {
	body := make(map[string]interface{})
	body["msg"] = "test send message"
	option := make(map[string]bool)
	option["push"] = true
	response, err := netEase.SendSingleTextMessage("1", "2", "test send message", option, "test send message", "")
	if err != nil {
		fmt.Printf("test Send_single_message occur error =====> %s\n", err.Error() )
	}
	if response.IsSuccess() {
		fmt.Printf("test Send_single_message success \n")
	} else {
		fmt.Printf("test Send_single_message fail code ====> %d\n", response.FailCode())
		fmt.Printf("test Send_single_message fail reason ====> %s\n", response.FailReason())
	}
}

func TestSendSingleImageMessage(t *testing.T) {
	option := make(map[string]bool)
	option["push"] = true
	response, err := netEase.SendSingleImageMessage("1", "2", "/Users/leeyifiei/Pictures/111.png", option, "test send message", "")
	if err != nil {
		fmt.Printf("Test Upload_file_multipart occur error =====> %s\n", err.Error())
	} else {
		if response.IsSuccess() {
			fmt.Printf("test Send_single_image_message success \n")

		} else {
			fmt.Printf("test upload_file_multipart fail code =====> %d\n", response.FailCode())
			fmt.Printf("test upload_file_multipart fail reason ====> %s\n", response.FailReason())
		}
	}
}

func TestSendMultiTextMessage(t *testing.T) {
	option := make(map[string]bool)
	to := []string {"1", "2", "3"}
	response, err := netEase.SendMultiTextMessage("1", to, "test send multi message", option, "test send multi message", "")
	if err != nil {
		fmt.Printf("Test TestSend_multi_text_message occur error =====> %s\n", err.Error())
	} else {
		if response.IsSuccess() {
			fmt.Printf("test TestSend_multi_text_message success \n")
			if len(response.Unregister) > 0 {
				fmt.Printf("unregister users:\n")
				for _, user := range response.Unregister {
					fmt.Printf("\t%s\n", user)
				}
			}
		} else {
			fmt.Printf("test TestSend_multi_text_message fail code =====> %d\n", response.FailCode())
			fmt.Printf("test TestSend_multi_text_message fail reason ====> %s\n", response.FailReason())
		}
	}
}

func TestSendMultiImageMessage(t *testing.T) {
	option := make(map[string]bool)
	option["push"] = true
	to := []string{"1", "2", "3" , "111"}
	response, err := netEase.SendMultiImageMessage("1", to, "/Users/leeyifiei/Pictures/111.png", option, "test send multi image message", "")
	if err != nil {
		fmt.Printf("Test TestSend_multi_image_message occur error =====> %s\n", err.Error())
	} else {
		if response.IsSuccess() {
			fmt.Printf("test TestSend_multi_image_message success \n")
			if len(response.Unregister) > 0 {
				fmt.Printf("unregister users %s:\n", response.Unregister)
				var user_map []string
				err := json.Unmarshal([]byte(response.Unregister), &user_map)
				if err != nil {
					fmt.Printf("test TestSend_multi_image_message unregister json unserialize fail")
				} else {
					for _, user := range user_map {
						fmt.Printf("%s\n", user)
					}
				}
			}
		} else {
			fmt.Printf("test TestSend_multi_text_message fail code =====> %d\n", response.FailCode())
			fmt.Printf("test TestSend_multi_text_message fail reason ====> %s\n", response.FailReason())
		}
	}
}