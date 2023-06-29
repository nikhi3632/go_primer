package mapreduce

import (
	"container/list"
	"fmt"
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
		ok := call(w.address, "Worker.Shutdown", args, &reply)
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
			worker := <-mr.registerChannel
			jobArgs := DoJobArgs{
				File:          mr.file,
				Operation:     JobType(operation),
				JobNumber:     jobNumber,
				NumOtherPhase: num,
			}
			jobReply := DoJobReply{}
			mapreduceRpc := call(worker, "Worker.DoJob", jobArgs, &jobReply)
			if operation == Map {
				if mapreduceRpc {
					mapper <- jobArgs
					mr.registerChannel <- worker
					break
				}
			} else {
				if mapreduceRpc {
					reducer <- jobArgs
					mr.registerChannel <- worker
					break
				}
			}
		}
	}

	for i := 0; i < mr.nMap; i++ {
		go jobAllotment(i, Map, mr.nReduce)
		<-mapper
	}
	close(mapper)

	for i := 0; i < mr.nReduce; i++ {
		go jobAllotment(i, Reduce, mr.nMap)
		<-reducer
	}
	close(reducer)

	return mr.KillWorkers()
}
