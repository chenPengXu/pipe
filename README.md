# pipe
data processing pipe with golang

This is a service to process data hand by hand in backgroud using golang goroutine. This pipes can be organized like a tree system to process source data or processed data from parent pipe.

### Example
```go
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

pipeMana.SetNode(&pipe1, &pipe2)


func pipeHandFn(pipeData *p.PipeData)  {
    data := pipeData.Data
    switch v := data.(type) {
    case int:
        pipeData.Data = v + 1
        //fmt.Println("After hander the Data:", pipeData)
    }
}
```

### The data Flow will like this
<img src="https://gw.alipayobjects.com/os/q/cms/images/jogzdyrw/34e81543-e4f5-4994-ac51-9d40d8e6dcb8_w593_h689.png" width=590>