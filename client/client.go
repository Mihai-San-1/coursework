
package main

import (
    "bufio"
    "context"
    "fmt"
    "log"
    "os"
    "strings"
    "sync"

    pb "grpc-chat/proto"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type ChatClient struct {
    username string
    stream   pb.ChatService_SendMessageClient
    wg       sync.WaitGroup
}

func NewChatClient(username string, conn *grpc.ClientConn) (*ChatClient, error) {
    client := pb.NewChatServiceClient(conn)
    stream, err := client.SendMessage(context.Background())
    if err != nil {
        return nil, err
    }

    return &ChatClient{
        username: username,
        stream:   stream,
    }, nil
}

func (c *ChatClient) SendMessage(text string) error {
    msg := &pb.ChatMessage{
        User: c.username,
        Text: text,
    }
    return c.stream.Send(msg)
}

func (c *ChatClient) ReceiveMessages() {
    c.wg.Add(1)
    go func() {
        defer c.wg.Done()
        for {
            msg, err := c.stream.Recv()
            if err != nil {
                log.Printf("Соединение закрыто: %v", err)
                return
            }
            fmt.Printf("\r[%s] %s\n> ", msg.User, msg.Text)
        }
    }()
}

func (c *ChatClient) Wait() {
    c.wg.Wait()
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Введите ваше имя: ")
    username, _ := reader.ReadString('\n')
    username = strings.TrimSpace(username)

    conn, err := grpc.Dial("0.0.0.0:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Ошибка подключения: %v", err)
    }
    defer conn.Close()

    client, err := NewChatClient(username, conn)
    if err != nil {
        log.Fatalf("Ошибка создания клиента: %v", err)
    }

    client.ReceiveMessages()
    fmt.Printf("Добро пожаловать в чат, %s!\n", username)
    fmt.Print("> ")

    for {
        text, _ := reader.ReadString('\n')
        text = strings.TrimSpace(text)
        
        if text == "/quit" {
            fmt.Println("Выход из чата...")
            return
        }
        
        if text != "" {
            if err := client.SendMessage(text); err != nil {
                log.Printf("Ошибка отправки: %v", err)
                return
            }
            fmt.Print("> ")
        }
    }
}
