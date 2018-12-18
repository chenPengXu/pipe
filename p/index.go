package p
import (
	"pipe/utils"
)

var pipeContainer *pipe;

const (
	PIPE_NODE_STATUS_OFFLINE  = -1
	PIPE_NODE_STATUS_STOPPED  = -1
	PIPE_NODE_STATUS_CONTINUE  = 1
	default_pipe_sema_chan_len  = 50
)

type Pipes interface {
	StartReceiveData(startNode string, data interface{})
	SetNode(nodes... *PipeNode)
	SetSemaphoreLength(len int)
}

type PipeData struct {
	Status int8
	nextNodes map[string]int8
	SourceData interface{}
	Data interface{}
}

//pipes
type pipe struct {
	sema *utils.Semaphore
	nodes map[string]*PipeNode
	pipeDatas chan PipeData
}

func (p *pipe) StartReceiveData(startNode string, data interface{})  {
	if p.nodes[startNode] == nil {
		return
	}
	p.sema.P()
	datas := PipeData{
		SourceData: data,
		Data: data,
	}
	p.runPreAction(startNode, &datas)
	go p.runNode(startNode, datas)
	go func() {
		for {
			select {
			case input := <- p.pipeDatas:
			//当上一个noode通过data标记停止后，后面的node不再处理数据
				if PIPE_NODE_STATUS_STOPPED == input.Status {
					return
				}
			//当前node如果下线了也不再处理数据
				for nodeName, status := range input.nextNodes {
					if PIPE_NODE_STATUS_OFFLINE == status {
						continue
					}
					p.sema.P()
					p.runPreAction(startNode, &input)
					go p.runNode(nodeName, input)
				}
			}
		}
	}()
}

func (p *pipe) SetNode(nodes... *PipeNode) {
	for _, node := range nodes {
		nodeName := node.Name
		if nodeName == "" || p.nodes[nodeName] != nil {
			continue
		}
		node.nextNodes = make(map[string]int8)
		p.nodes[nodeName] = node
	}

	//init nextNodes Info
	for _, node := range nodes {
		nodeName := node.Name
		if nodeName == "" {
			continue
		}
		preNode := node.PreNode
		if preNode == "" {
			continue
		}
		if p.nodes[preNode] == nil {
			continue
		}
		p.nodes[preNode].nextNodes[nodeName] = node.Status
	}
}

func (p *pipe) runPreAction(nodeName string, datas *PipeData)  {
	p.nodes[nodeName].doDataProcessingBeforeRoutine(datas) //data action
}

func (p *pipe) runNode(nodeName string, datas PipeData)  {
	p.nodes[nodeName].doDataProcessing(&datas) //data action
	datas.nextNodes = p.nodes[nodeName].nextNodes
	p.pipeDatas <- datas
	p.sema.V()
}

func (p *pipe) SetSemaphoreLength(len int)  {
	p.sema = utils.NewSemaphore(len)
}

type PipeNode struct {
	Name string
	Status int8
	PreNode string
	PreHanderBeforeRoutine func(datas *PipeData)
	DataHander func(datas *PipeData) (status int8)
	nextNodes map[string]int8
}

func (node *PipeNode) doDataProcessing(datas *PipeData)  {
	if node.DataHander != nil {
		status := node.DataHander(datas)
		datas.Status = status
	}
}

func (node *PipeNode) doDataProcessingBeforeRoutine(datas *PipeData)  {
	if node.PreHanderBeforeRoutine != nil {
		node.PreHanderBeforeRoutine(datas)
	}
}

func (node *PipeNode) SetPipeNodeOffLine()  {
	node.Status = PIPE_NODE_STATUS_OFFLINE
}

func NewPipes() Pipes{
	semas := default_pipe_sema_chan_len
	nodes := make(map[string]*PipeNode)
	pipeContainer = &pipe{
		sema: utils.NewSemaphore(semas),
		nodes: nodes,
		pipeDatas:  make(chan PipeData, int(semas)),
	}
	return pipeContainer
}
