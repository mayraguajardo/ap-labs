package main

import (
	"fmt"
	"time"
	"os"
)
func main() {
	channel := make(chan struct{})
	msg := make(chan int)
	var commsPerSecond int
	startTime := time.Now()

	go func(){
		msg <- 1
		for{
			commsPerSecond++
			msg<- <-msg
		}
		channel <- struct{}{}
	}()
	go func(){
		for{
			commsPerSecond++
			msg<- <-msg
		}
		channel <- struct{}{}
	}()
	for{
		select{
		case <- time.After(1 * time.Second):
			s := float64(commsPerSecond)/time.Since(startTime).Seconds()
			fmt.Println("Communications Per Second :", s)
		}
		if commsPerSecond > 50000000{
			fmt.Println("Communication is over")
			os.Exit(0)
		}
	}

	<- channel
}
