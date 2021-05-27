package main

import (
	"chatRoom/src/server/apps"
)

func main() {
	apps.OnInit()
	apps.Run(4)
}
