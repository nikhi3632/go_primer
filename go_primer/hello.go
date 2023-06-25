package main

import (
	"fmt"
	"time"
)

func run1() {
	channel := make(chan int, 2)

	go func() {
		for j := 0; j < 3; j++ {
			fmt.Println(time.Now(), j, "sending")
			channel <- j
			fmt.Println(time.Now(), j, "sent")
		}

		// XXX: There could be cases where this message is not completed,
		// this is solved in futured examples
		fmt.Println(time.Now(), "all completed") // race condition
	}()

	time.Sleep(2 * time.Second)

	fmt.Println(time.Now(), "waiting for messages")

	fmt.Println(time.Now(), "received", <-channel)
	fmt.Println(time.Now(), "received", <-channel)
	fmt.Println(time.Now(), "received", <-channel)

	fmt.Println(time.Now(), "exiting")
}

func run2() {
	ch := make(chan int, 2)
	exit := make(chan struct{})

	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println(time.Now(), i, "sending")
			ch <- i
			fmt.Println(time.Now(), i, "sent")

			time.Sleep(1 * time.Second)
		}

		fmt.Println(time.Now(), "all completed, leaving")

		close(ch)
	}()

	go func() {
		// XXX: This is overcomplicated because is only channel only, "select"
		// shines when using multiple channels.
		for {
			select {
			case v, open := <-ch:
				if !open {
					close(exit)
					return
				}

				fmt.Println(time.Now(), "received", v)
			}
		}

		// XXX: In cases where only one channel is used
		// for v := range ch {
		// 	fmt.Println(time.Now(), "received", v)
		// }

		// close(exit)
	}()

	fmt.Println(time.Now(), "waiting for everything to complete")

	<-exit

	fmt.Println(time.Now(), "exiting")
}

func main() {
	run1()
	fmt.Println("--------------------------------------")
	run2()
}
