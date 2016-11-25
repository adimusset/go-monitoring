package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Please give 2 arguments, not ", len(args))
		return
	}
	f, err := os.Create(args[1])
	if err != nil {
		fmt.Println("Error while creating log file - ", err.Error())
		return
	}
	fmt.Println("running")
	for {
		time.Sleep(time.Second)
		date := time.Now().Format("02/Jan/2006:15:04:05 -0700")
		log := fmt.Sprintf(`%s - - [%s] "GET %s HTTP/1.0" 200 50`, "127.0.0.1", date, "/section/page")
		f.WriteString(log + "\n")
	}
}
