package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// wsMessage 与服务端相同的消息结构
type wsMessage struct {
	Type  string `json:"type"`
	Data  string `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

// AIChatRequest AI聊天请求结构
type AIChatRequest struct {
	Question string `json:"question"`
	Context  any    `json:"context,omitempty"`
}

func main() {
	// 命令行参数
	url := flag.String("url", "ws://localhost:8080/gdesign/ai/ws", "WebSocket服务器URL")
	token := flag.String("token", "test_user", "认证token")
	question := flag.String("q", "什么是有限自动机？", "测试问题")
	flag.Parse()

	// 构建带token的WebSocket URL
	wsURL := fmt.Sprintf("%s?token=%s", *url, *token)

	// 连接到WebSocket服务器
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer conn.Close()

	// 设置读写超时
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// 发送初始请求
	req := AIChatRequest{
		Question: *question,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.Fatal("序列化请求失败:", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, reqBytes); err != nil {
		log.Fatal("发送请求失败:", err)
	}

	fmt.Printf("已连接到 %s\n", wsURL)
	fmt.Printf("已发送问题: %s\n", *question)
	fmt.Println("正在接收响应...")
	fmt.Println("----------------------------------------")

	// 处理信号，允许用户通过Ctrl+C关闭连接
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// 创建一个通道用于接收消息
	messages := make(chan []byte)
	errors := make(chan error)

	// 启动一个goroutine读取消息
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				errors <- err
				return
			}
			messages <- message
		}
	}()

	// 主循环
	var fullResponse string
	for {
		select {
		case message := <-messages:
			var wsMsg wsMessage
			if err := json.Unmarshal(message, &wsMsg); err != nil {
				fmt.Printf("解析消息失败: %v\n", err)
				continue
			}

			switch wsMsg.Type {
			case "chunk":
				fmt.Print(wsMsg.Data)
				fullResponse += wsMsg.Data
			case "error":
				fmt.Printf("\n错误: %s\n", wsMsg.Error)
			case "done":
				fmt.Println("\n----------------------------------------")
				fmt.Println("响应完成")
				return
			default:
				fmt.Printf("未知消息类型: %s\n", wsMsg.Type)
			}

		case err := <-errors:
			fmt.Printf("\n连接错误: %v\n", err)
			return

		case <-interrupt:
			fmt.Println("\n\n正在关闭连接...")
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			time.Sleep(500 * time.Millisecond)
			return
		}
	}
}