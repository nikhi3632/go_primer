package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

// Sum numbers from channel `nums` and output sum to `out`.
// You should only output to `out` once.
// Do NOT modify function signature.
func sumWorker(nums chan int, out chan int) {
	// TODO: implement me
	// HINT: use for loop over `nums`
	var sum int

	for num := range nums {
		sum = sum + num
	}

	out <- sum
	close(out)
}

// Read integers from the file `fileName` and return sum of all values.
// This function must launch `num` go routines running
// `sumWorker` to find the sum of the values concurrently.
// You should use `checkError` to handle potential errors.
// Do NOT modify function signature.
func sum(num int, fileName string) int {
	// TODO: implement me
	// HINT: use `readInts` and `sumWorkers`
	// HINT: used buffered channels for splitting numbers between workers
	file, err := os.Open(fileName)
	checkError(err)
	defer file.Close()

	nums, err := readInts(file)
	checkError(err)

	numChunks := len(nums) / num
	var outChannels []chan int

	for chunk := 0; chunk < num; chunk++ {
		channel := make(chan int, numChunks)
		out := make(chan int, num)
		outChannels = append(outChannels, out)
		lowerBound := chunk * numChunks
		upperBound := (chunk + 1) * numChunks
		for _, numVal := range nums[lowerBound:upperBound] {
			channel <- numVal
		}
		close(channel)
		go sumWorker(channel, out)
	}

	result := 0
	for _, channel := range outChannels {
		partialSum := <-channel
		result = result + partialSum
	}

	return result
}

// Read a list of integers separated by whitespace from `r`.
// Return the integers successfully read with no error, or
// an empty slice of integers and the error that occurred.
// Do NOT modify this function.
func readInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var elems []int
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return elems, err
		}
		elems = append(elems, val)
	}
	return elems, nil
}
