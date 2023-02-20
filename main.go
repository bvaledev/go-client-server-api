package main

import (
	"DesafioClientServer/client"
	"DesafioClientServer/server"
)

func main() {
	go server.Main()
	client.Main()
}
