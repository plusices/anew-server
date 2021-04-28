package _examples

import (
	"context"
	"fmt"
	"time"
	"ts-go-server/pkg/sshx"
)

func main() {
	//
	conf, _ := sshx.NewAuthConfig("root", "xxxx", "", "")
	ssha := sshx.New("192.168.56.100:22", conf)
	ctxa, cancela := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancela()
	whoami, err := ssha.Command(ctxa, "whoami")
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	fmt.Printf("whoami: %s", whoami)
	_, err = ssha.SendFile("/root/test/test.txt", "D://VIP视频教程账号.txt", true, true)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
}
