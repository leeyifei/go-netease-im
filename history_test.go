package netease_im

import (
	"testing"
	"fmt"
)

/*
var net_ease *Netease_im

func init() {
	net_ease = Netease_im_instance("71b6da2c2393f0717b9f905d9ad0b9ac", "9443d5af1682", true)
}
*/

func TestUserMessageHistory(t *testing.T) {
	response, err := netEase.UserMessageHistory("1", "2", "1473134400000", "1473350399000", "100", true)
	if err != nil {
		fmt.Printf("test TestUser_message_history occur error =====> %v\n", err)
	} else {
		if response.IsSuccess() {
			fmt.Printf("test TestUser_message_history success \n")
			msg_lists := response.Msgs
			if length := len(msg_lists) ; length > 0 {
				for idx, msg := range msg_lists {
					switch msg.GetMessageType() {
					case 0:
						fmt.Printf("the %d is a text message\n", idx)
						fmt.Printf("    message : %s\n", msg.Body["msg"])
					case 1:
						fmt.Printf("the %d is a image message\n", idx)
						fmt.Printf("    name : %s\n", msg.Body["name"])
						fmt.Printf("    md5 : %s\n", msg.Body["md5"])
						fmt.Printf("    url : %s\n", msg.Body["url"])
						fmt.Printf("    ext : %s\n", msg.Body["ext"])
						fmt.Printf("    w : %d\n", int(msg.Body["w"].(float64)))
						fmt.Printf("    h : %d\n", int(msg.Body["h"].(float64)))
						fmt.Printf("    size : %d\n", int(msg.Body["size"].(float64)))
					}
				}
			}
		} else {
			fmt.Printf("test TestUser_message_history fail code =====> %d\n", response.FailCode())
			fmt.Printf("test TestUser_message_history fail reason ====> %s\n", response.FailReason())
		}
	}
}
