package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Позволяет всем источникам, измените по необходимости
    },
}

type Client struct {
	conn *websocket.Conn
}

// Глобальный список клиентов и мьютекс для безопасного доступа
var clients = make(map[*Client]bool)
var mu sync.Mutex

func handleHTTP(w http.ResponseWriter, r *http.Request) {
    //fmt.Fprintf(w, "Hello, this is an HTTP response!")
	if r.URL.Path != "/settings_local.json" {
		http.ServeFile(w, r, "." + r.URL.Path)
	} else {
		http.ServeFile(w, r, "err403.html")
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()

	client := &Client{conn: conn}

	// Добавляем клиента в список
	mu.Lock()
	clients[client] = true
	mu.Unlock()

	// Отправляем текущий чат при подключении
	chat, err := os.ReadFile("chat.txt")
	if err != nil {
		chat = []byte(err.Error())
	}

	err = conn.WriteMessage(1, chat)
	if err != nil {
		fmt.Println("Error while writing message:", err)
	}

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			break
		}
		fmt.Printf("Received message: %s\n", msg)
		appendToFile("chat.txt", string(msg) + "\n")

		// Отправляем сообщение всем клиентам
		mu.Lock()
		for c := range clients {
			err := c.conn.WriteMessage(msgType, msg)
			if err != nil {
				fmt.Println("Error while writing message:", err)
				c.conn.Close()
				delete(clients, c) // Удаляем клиента из списка в случае ошибки
			}
		}
		mu.Unlock()
	}
}

func main() {
	port, err := getValueFromJson("settings_local.json", "port")
	if err != nil {
		os.Exit(1)
	}
	str_port := fmt.Sprintf("%s", port)

	address, err := getValueFromJson("settings_local.json", "address")
	if err != nil {
		os.Exit(2)
	}
	str_address := fmt.Sprintf("%s", address)

    http.HandleFunc("/", handleHTTP) // Обработка HTTP-запросов
    http.HandleFunc("/ws", handleWebSocket) // Обработка WebSocket-запросов

    fmt.Println("Server is running on " + str_address + ":" + str_port)
    err = http.ListenAndServe(str_address + ":" + str_port, nil) // Запуск сервера
    if err != nil {
        fmt.Println("Error starting server:", err)
    }
}

func getValueFromJson(filePath, keyPath string) (interface{}, error) {
	// Чтение содержимого файла
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %v", err)
	}

	// Парсинг JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	// Разбиваем путь на массив ключей
	keys := strings.Split(keyPath, ".")

	// Получаем значение по ключам
	value := jsonData
	for _, key := range keys {
		if v, ok := value[key]; ok {
			// Если значение - это вложенная структура, продолжаем
			if nestedValue, ok := v.(map[string]interface{}); ok {
				value = nestedValue
			} else {
				return v, nil // Возвращаем найденное значение
			}
		} else {
			return nil, fmt.Errorf("ключ \"%s\" не найден", key)
		}
	}

	return nil, fmt.Errorf("значение не найдено")
}

func appendToFile(fileName, textToAdd string) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Ошибка при открытии файла:", err)
        return
    }
    defer file.Close() // Закрываем файл в конце

    // Создаем новый записыватель
    writer := bufio.NewWriter(file)

    // Записываем текст в файл
    _, err = writer.WriteString(textToAdd)
    if err != nil {
        fmt.Println("Ошибка при записи в файл:", err)
        return
    }

    // Сбрасываем буфер
    err = writer.Flush()
    if err != nil {
        fmt.Println("Ошибка при сбросе буфера:", err)
        return
    }
}