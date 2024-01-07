package main

import (
	chiserver "lesson15_16/4"
	grpcclient "lesson15_16/5/client"
	grpcserver "lesson15_16/5/server"
)

func main() {
	chiserver.StartServer()

	grpcclient.StartClient()
	grpcserver.StartServer()
}
