package main

import (
    "io"
    "log"
    "sync"
    "time"

    pb "grpc-chat/proto"
)

type ChatServer struct {
    pb.UnimplementedChatServiceServer
    mu       sync.RWMutex
    streams  map[pb.ChatService_SendMessageServer]bool
    messages []*pb.ChatMessage
}

func NewChatServer() *ChatServer {
    return &ChatServer{
        streams: make(map[pb.ChatService_SendMessageServer]bool),
    }
}

func (s *ChatServer) SendMessage(stream pb.ChatService_SendMessageServer) error {
    s.mu.Lock()
    s.streams[stream] = true
    s.mu.Unlock()

    s.mu.RLock()
    for _, msg := range s.messages {
        if err := stream.Send(msg); err != nil {
            log.Printf("Ошибка отправки истории: %v", err)
            break
        }
    }
    s.mu.RUnlock()

    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("Ошибка получения сообщения: %v", err)
            break
        }

        msg.Timestamp = time.Now().Unix()
        log.Printf("[%s] %s", msg.User, msg.Text)

        s.mu.Lock()
        s.messages = append(s.messages, msg)
        if len(s.messages) > 100 {
            s.messages = s.messages[1:]
        }
        s.mu.Unlock()

        s.mu.RLock()
        for clientStream := range s.streams {
            if err := clientStream.Send(msg); err != nil {
                log.Printf("Ошибка отправки клиенту: %v", err)
            }
        }
        s.mu.RUnlock()
    }

    s.mu.Lock()
    delete(s.streams, stream)
    s.mu.Unlock()

    return nil
}