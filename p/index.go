package p
import (
	"pipe/utils"
)

var pipeContainer *pipe;

const PIPE_NODE_STATUS_OFFLINE  = -1
const default_pipe_sema_chan_len  = 50

type Pipes interface {
	StartReceiveData(startNode string, data interface{})
	SetNode(nodes... *PipeNode)
	SetSemaphoreLength(len int)
}

type PipeData struct {
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
	go p.runNode(startNode, datas)
	go func() {
		for {
			select {
			case input := <- p.pipeDatas:
				for nodeName, status := range input.nextNodes {
					if status == PIPE_NODE_STATUS_OFFLINE {
						continue
					}
					p.sema.P()
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
	DataHander func(datas *PipeData)
	nextNodes map[string]int8
}

func (node *PipeNode) doDataProcessing(datas *PipeData)  {
	if node.DataHander != nil {
		node.DataHander(datas)
	}
}

func (node *PipeNode) SetPipeNodeOffLine()  {
	node.Status = PIPE_NODE_STATUS_OFFLINE
}

func NewPipes() Pipes{
	semas := default_pipe_sema_chan_len
	if pipeContainer == nil {
		nodes := make(map[string]*PipeNode)
		pipeContainer = &pipe{
			sema: utils.NewSemaphore(semas),
			nodes: nodes,
			pipeDatas:  make(chan PipeData, int(semas)),
		}
		return pipeContainer
	}
	return pipeContainer
}


