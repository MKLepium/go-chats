package main

import (
	httpserver "github.com/mklepium/chats/internal/http"
)

func main() {
	go httpserver.StartServer()
	select {}

}
