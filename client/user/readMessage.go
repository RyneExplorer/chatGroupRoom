package user

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func ReadMessage(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("服务器连接断开:", err)
		}
		fmt.Println(strings.TrimSpace(msg))
	}
}
