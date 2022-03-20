package main

import (
	"fmt"
	"go-websocket/server/connection"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		wsconn *websocket.Conn
		err    error
		data   []byte
		conn   *connection.Connection
	)

	if wsconn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}

	//初始化Connection
	if conn, err = connection.InitConnection(wsconn); err != nil {
		goto ERR
	}

	// 设置检测心跳
	go func() {
		var (
			err error
		)
		for {
			if err = conn.WriteMessge([]byte("heartbeat")); err != nil {
				return
			}
			time.Sleep(5 * time.Second)
		}

	}()

	//读取消息
	for {
		if data, err = conn.ReadMessge(); err != nil {
			goto ERR
		}

		if err = conn.WriteMessge(data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()

}

func main() {
	fmt.Println("testing websocket ... ")

	http.HandleFunc("/test_websocket", wsHandler)

	http.ListenAndServe("0.0.0.0:6066", nil)
}
