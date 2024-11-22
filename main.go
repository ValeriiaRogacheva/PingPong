package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Game struct {
	playerScore    int
	computerScore  int
	ballX          int
	ballY          int
	ballDirection  int
	mu             sync.Mutex
	messageChannel chan string
}

func (g *Game) sendState(conn *websocket.Conn) {
	g.mu.Lock()
	defer g.mu.Unlock()

	state := fmt.Sprintf("state:ballX=%d,ballY=%d,playerScore=%d,computerScore=%d", g.ballX, g.ballY, g.playerScore, g.computerScore)
	conn.WriteMessage(websocket.TextMessage, []byte(state))
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	game := &Game{
		messageChannel: make(chan string),
		ballX:          250,
		ballY:          rand.Intn(180) + 10,
		ballDirection:  1,
	}

	// Фоновая задача: обновляет состояние мяча
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)

			game.mu.Lock()
			// Двигаем мяч в текущем направлении
			game.ballX += 50 * game.ballDirection

			// Проверяем границы: если мяч улетел за пределы
			if game.ballX < 0 {
				game.computerScore++ // Компьютер выигрывает
				game.resetBall(1)    // Мяч возвращается с направлением к игроку
			} else if game.ballX > 500 {
				game.playerScore++ // Игрок выигрывает
				game.resetBall(-1) // Мяч возвращается с направлением к компьютеру
			}
			game.mu.Unlock()

			game.messageChannel <- "update"
		}
	}()

	// Основной цикл чтения/обработки сообщений
	for {
		select {
		case msg := <-game.messageChannel:
			if msg == "update" {
				game.sendState(conn)
			}
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message:", err)
				return
			}

			if string(message) == "hit" {
				game.mu.Lock()
				game.ballDirection *= -1 // Меняем направление мяча
				game.ballY = rand.Intn(180) + 10
				game.mu.Unlock()
			}
		}
	}
}

func (g *Game) resetBall(direction int) {
	g.ballX = 250 // Центр стола
	g.ballY = rand.Intn(180) + 10
	g.ballDirection = direction
}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/game", gameHandler)

	fmt.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
