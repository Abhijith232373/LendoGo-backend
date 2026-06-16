package utils

import (
    "log"
    "github.com/gofiber/contrib/websocket"
)
func SafeWriteJSON(conn *websocket.Conn, payload any) (err error) {
    defer func() {
        if r := recover(); r != nil {
            log.Println("Recovered from websocket panic:", r)
        }
    }()
    return conn.WriteJSON(payload)
}

func IsAdminUser(userID string) bool {
    return userID == "0" || userID == "admin"
}