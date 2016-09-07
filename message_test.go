package netease_im

import (
	"testing"
	"fmt"
)

var net_ease *Netease_im

func init() {
	net_ease = Netease_im_instance("71b6da2c2393f0717b9f905d9ad0b9ac", "9443d5af1682", true)
}

func TestSend_single_text_message(t *testing.T) {
	body := make(map[string]interface{})
	body["msg"] = "test send message"
	option := make(map[string]bool)
	option["push"] = true
	response, err := net_ease.Send_single_text_message("1", "2", "test send message", option, "test send message", "")
	if err != nil {
		fmt.Printf("test Send_single_message occur error =====> %s\n", err.Error() )
	}
	if response.Is_success() {
		fmt.Printf("test Send_single_message success \n")
	} else {
		fmt.Printf("test Send_single_message fail code ====> %d\n", response.Fail_code())
		fmt.Printf("test Send_single_message fail reason ====> %s\n", response.Fail_reason())
	}
}


func TestUpload_file_multipart(t *testing.T) {
	response, err := net_ease.Upload_image_multipart("/Users/leeyifiei/Pictures/111.png")
	if err != nil {
		fmt.Printf("Test Upload_file_multipart occur error =====> %s\n", err.Error())
	} else {
		if response.Is_success() {
			fmt.Printf("test upload_file_multipart success \n")
			fmt.Printf("file url is =====> %s\n", response.Url)
			fmt.Printf("file name is =====> %s\n", response.Name)
			fmt.Printf("file ext is =====> %s\n", response.Ext[1:])
			fmt.Printf("file size is =====> %d\n", response.Size)
			fmt.Printf("file width is =====> %d\n", response.W)
			fmt.Printf("file height is =====> %d\n", response.H)
			fmt.Printf("file md5 is =====> %s\n", response.Md5)
		} else {
			fmt.Printf("test upload_file_multipart fail code =====> %d\n", response.Fail_code())
			fmt.Printf("test upload_file_multipart fail reason ====> %s\n", response.Fail_reason())
		}
	}

}

func TestSend_single_image_message(t *testing.T) {
	option := make(map[string]bool)
	option["push"] = true
	response, err := net_ease.Send_single_image_message("1", "2", "/Users/leeyifiei/Pictures/111.png", option, "test send message", "")
	if err != nil {
		fmt.Printf("Test Upload_file_multipart occur error =====> %s\n", err.Error())
	} else {
		if response.Is_success() {
			fmt.Printf("test Send_single_image_message success \n")

		} else {
			fmt.Printf("test upload_file_multipart fail code =====> %d\n", response.Fail_code())
			fmt.Printf("test upload_file_multipart fail reason ====> %s\n", response.Fail_reason())
		}
	}
}