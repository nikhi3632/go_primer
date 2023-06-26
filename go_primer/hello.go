package main

import (
	"fmt"
	"time"
)

// make vs new, make is used only for channels, slices and maps

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

func run3() {
	channel1 := make(chan string)
	channel2 := make(chan string)
	channel3 := make(chan string)

	go func() {
		channel1 <- "cat"
	}()

	go func() {
		channel2 <- "dog"
	}()

	go func() {
		channel3 <- "cow"
	}()

	select { // select statement is going to block until it recieves msg from one of the channels
	case msgFromChannel1 := <-channel1:
		fmt.Println(msgFromChannel1)
	case msgFromChannel2 := <-channel2:
		fmt.Println(msgFromChannel2)
	case msgFromChannel3 := <-channel3:
		fmt.Println(msgFromChannel3)
	} // one of it gets printed randomly whichever is ready first
}

func doWork(done <-chan bool) { // done channel passed a read only
	for {
		select {
		case <-done:
			return
		default:
			fmt.Println("Doing WORK!")
		}
	}
}

func run4() {
	done := make(chan bool)
	go doWork(done)

	ticker := time.NewTicker(2 * time.Second)
	<-ticker.C
	close(done)
}

/*
    Pipeline: Start [Intial data] -> Stage1 [Intermediate data] ->
	Stage2 [Intermediate data] -> Stage3 [Output data]-> End.
*/

func sliceToChannel(sliceData []int) <-chan int {
	channelData := make(chan int)
	go func() {
		for _, sliceVal := range sliceData {
			channelData <- sliceVal // channelData is blocked until it's read below from channelSq.
		}
		close(channelData)
	}()
	return channelData
}

func sq(channelData <-chan int) chan int {
	channelSq := make(chan int)
	go func() {
		for channelVal := range channelData { // loop stops when channelData is closed above.
			channelSq <- channelVal * channelVal
		}
		close(channelSq)
	}()
	return channelSq
}

func run5() {
	dataSlice := []int{2, 3, 5, 7}
	dataChannel := sliceToChannel(dataSlice)
	sqChannel := sq(dataChannel)
	for n := range sqChannel {
		fmt.Println(n)
	}
}

/*
	Unbuffered channels have a capaticy of one. In the above example
	the go routines in sliceToChannel and sq are running at the same
	time and the communication between the go routines is synchronous.
*/

func main() {
	run1()
	fmt.Println("--------------------------------------")
	run2()
	fmt.Println("--------------------------------------")
	run3()
	fmt.Println("--------------------------------------")
	run4()
	fmt.Println("--------------------------------------")
	run5()
}
