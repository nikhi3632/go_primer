package mapreduce

import (
	"container/list"
	"fmt"
	"log"
)

type WorkerInfo struct {
	address string
	// You can add definitions here.
}

// Clean up all workers by sending a Shutdown RPC to each one of them Collect
// the number of jobs each work has performed.
func (mr *MapReduce) KillWorkers() *list.List {
	l := list.New()
	for _, w := range mr.Workers {
		DPrintf("DoWork: shutdown %s\n", w.address)
		args := &ShutdownArgs{}
		var reply ShutdownReply
		ok := call(w.address, "Worker.Shutdown", args, &reply) // shutdown each worker
		if ok == false {
			fmt.Printf("DoWork: RPC %s shutdown error\n", w.address)
		} else {
			l.PushBack(reply.Njobs)
		}
	}
	return l
}

func (mr *MapReduce) RunMaster() *list.List {
	// Your code here
	mapper := make(chan interface{})
	reducer := make(chan interface{})
	jobAllotment := func(jobNumber int, operation string, num int) {
		for {
			worker := <-mr.registerChannel // take out the worker from registerChannel
			jobArgs := DoJobArgs{
				File:          mr.file,
				Operation:     JobType(operation),
				JobNumber:     jobNumber,
				NumOtherPhase: num,
			}
			jobReply := DoJobReply{}
			call(worker, "Worker.DoJob", jobArgs, &jobReply)
			// fmt.Println("Mr-Rpc:", mapreduceRpc, ", Worker:", worker, ", JobArgs:", jobArgs, ", JobReply:", jobReply)
			// if mapreduceRpc {
			if operation == Map {
				mapper <- jobArgs
				mr.registerChannel <- worker
				break
			} else if operation == Reduce {
				reducer <- jobArgs
				mr.registerChannel <- worker
				break
			} else {
				log.Panic("unknown operation for map reduce\n")
			}
			// }
		}
	}

	for i := 0; i < mr.nMap; i++ {
		go jobAllotment(i, Map, mr.nReduce)
		<-mapper // if mapreduceRpc is true finish the map job assigned
	}
	close(mapper)

	for i := 0; i < mr.nReduce; i++ {
		go jobAllotment(i, Reduce, mr.nMap)
		<-reducer // if mapreduceRpc is true finish the reduce job assigned
	}
	close(reducer)

	return mr.KillWorkers() // killing the workers sets mapreduceRpc to false
}
