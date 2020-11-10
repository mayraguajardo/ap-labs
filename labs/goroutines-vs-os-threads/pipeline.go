package main

import (
	"fmt"
	"time"
)

var channel = make(chan struct{})
var transitTime = time.Now()

func pipeline(first chan struct{}, second chan struct{}, cuStages int, maxStages int){
	if cuStages <= maxStages{
		fmt.Println("Maximum number of pipeline stages   : ", maxStages)
		fmt.Println("Time to transit trough the pipeline : ", time.Since(transitTime))
		transitTime = time.Now()
		go pipeline(second,first,cuStages+1,maxStages)
		second<-<- first
	}else{
		fmt.Println("Maximum number of stages, need to wait")
		channel <- struct{}{}
	}
}
func main() {
	var firstPipeline chan(struct{})
	var secondPipeline chan(struct{})
	go pipeline(firstPipeline, secondPipeline, 0, 100)
	<-channel
}
