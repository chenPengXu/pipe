package main

import (
	"fmt"
	"time"
	"sync"

	"pipe/p"
)

const NUM  = 20
var waitChan sync.WaitGroup

func main()  {
	pipeMana := p.NewPipes()
	pipe1 := p.PipeNode {
		Name: "pipe1",
		DataHander: pipeHandFn,
	}
	pipe2 := p.PipeNode{
		Name: "pipe2",
		PreNode: "pipe1",
		DataHander: pipeHandFn,
	}
	pipe3 := p.PipeNode{
		Name: "pipe3",
		PreNode: "pipe2",
		DataHander: pipeHandFn3,
	}

	pipeMana.SetNode(&pipe1, &pipe2, &pipe3)
	//pipeMana.SetSemaphoreLength(20)

	//test time Calculation
	fmt.Println("Job Starting ..@time:", time.Now().Unix())
	for i := 0; i < NUM; i++ {
		waitChan.Add(1)
		func (n int){
			pipeMana.StartReceiveData("pipe1", n)
		}(i)

	}
	waitChan.Wait()
	fmt.Println("Job Finished ..@time:", time.Now().Unix())
}

func pipeHandFn(pipeData *p.PipeData)  {
	data := pipeData.Data
	switch v := data.(type) {
	case int:
		pipeData.Data = v + 1
		//fmt.Println("After hander the Data:", pipeData)
	}
}

func pipeHandFn3(pipeData *p.PipeData)  {
	data := pipeData.Data
	switch v := data.(type) {
	case int:
		pipeData.Data = v + 1
		time.Sleep(time.Duration(2) * time.Second)
		// if line 32 SetSemaphoreLength(1), you will see the total time will be 2*NUM seconds;
		//fmt.Println("After Last hander the Data:", pipeData.Data)
		waitChan.Done()
	}
}