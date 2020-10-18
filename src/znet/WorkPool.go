package znet

//限制
//1.worker工作池的任务队列的最大值
//2.任务队列中任务的最大数量
//协程池
type WorkerPool struct {
	cap       int
	tasksSize int
	TaskQueue []chan IConnection //信道集合

}

//启动一个worker工作池，开启工作池只能发生一次
func (W *WorkerPool) StartWorkPool() {
	//根据任务队列的大小，分别开启worker，每个worker用go来承载，每一个worker对应一个任务队列
	for i := 0; i < W.cap; i++ {
		//为每个worker开辟缓冲信道（任务队列）
		W.TaskQueue[i] = make(chan IConnection, W.tasksSize)
		//启动worker，阻塞等待任务从channel中到来

		go W.StartOneWorker(i, W.TaskQueue[i])
	}
}

func (W *WorkerPool) StartOneWorker(id int, taskqueue chan IConnection) {

	for {
		select {
		case request := <-taskqueue:
			//如果有消息过来，则处理业务
			request.Start()
		default:
			continue
		}
	}
}

func (W *WorkerPool) Put(Connection IConnection) {
	index := Connection.GetConnID() % uint32(W.cap)
	W.TaskQueue[index] <- Connection
}

func NewWorkerPool(cap int, len int) *WorkerPool {
	return &WorkerPool{
		cap:       cap,
		tasksSize: len,
		TaskQueue: make([]chan IConnection, cap),
	}
}
