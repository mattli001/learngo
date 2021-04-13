package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mattli001/learngo/worker_pool/dispatcher"
	"github.com/mattli001/learngo/worker_pool/worker"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	start := time.Now()
	dd := dispatcher.New(2).Start()

	// terms := []int{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
	// 	4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
	// 	4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
	// 	4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
	// 	4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
	// 	4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
	// 	4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
	// 	4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
	// 	4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
	// 	4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
	// 	4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
	// 	4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
	// 	4, 1, 2, 3, 4}
	terms := []int{1, 2, 3, 4, 5}
	for i := range terms {
		log.Printf("Submit JobID: %d ... start\n", i)
		dd.Submit(worker.Job{
			ID:        i,
			Name:      fmt.Sprintf("JobID::%d", i),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		log.Printf("Submit JobID: %d ... done\n", i)
	}
	end := time.Now()
	log.Print(end.Sub(start).Seconds())
	log.Print("exit")
}
