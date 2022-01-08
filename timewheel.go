package helper

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

const (
	modeIsCircle  = true
	modeNotCircle = false

	modeIsAsync  = true
	modeNotAsync = false
)

type taskID int64

type Task struct {
	delay    time.Duration
	id       taskID
	round    int
	callback func()

	async  bool
	stop   bool
	circle bool
}

type TimeWheel struct {
	randomID int64

	tick      time.Duration
	ticker    *time.Ticker
	tickQueue chan time.Time

	bucketsNum    int
	buckets       []map[taskID]*Task
	bucketIndexes map[taskID]int

	currentIndex int

	onceStart sync.Once

	addC    chan *Task
	removeC chan *Task
	stopC   chan struct{}

	exited bool
}

// NewTimeWheel 创建一个时间轮
func NewTimeWheel(tick time.Duration, bucketsNum int, safeMode bool) (*TimeWheel, error) {
	if tick.Seconds() < 0.1 {
		return nil, errors.New("invalid params, must tick >= 100 ms")
	}
	if bucketsNum <= 0 {
		return nil, errors.New("invalid params, must bucketsNum > 0")
	}

	tw := &TimeWheel{
		// tick
		tick: tick,

		// store
		bucketsNum:    bucketsNum,
		bucketIndexes: make(map[taskID]int, 1024*100),
		buckets:       make([]map[taskID]*Task, bucketsNum),
		currentIndex:  0,

		// signal
		addC:    make(chan *Task, 1024*5),
		removeC: make(chan *Task, 1024*2),
		stopC:   make(chan struct{}),
	}

	for i := 0; i < bucketsNum; i++ {
		tw.buckets[i] = make(map[taskID]*Task, 16)
	}

	if safeMode {
		tw.tickQueue = make(chan time.Time, 10)
	}

	return tw, nil
}

// Start 启动时间轮
func (tw *TimeWheel) Start() {
	// 只允许启动一次
	tw.onceStart.Do(
		func() {
			tw.ticker = time.NewTicker(tw.tick)
			go tw.schduler()
			go tw.tickGenerator()
		},
	)
}

// Add 添加任务
func (tw *TimeWheel) Add(delay time.Duration, callback func()) *Task {
	return tw.addAny(delay, callback, modeNotCircle, modeIsAsync)
}

// AddCron 添加循环定时任务
func (tw *TimeWheel) AddCron(delay time.Duration, callback func()) *Task {
	return tw.addAny(delay, callback, modeIsCircle, modeIsAsync)
}

//删除任务
func (tw *TimeWheel) Remove(task *Task) error {
	tw.removeC <- task
	return nil
}

// Sleep 睡眠
func (tw *TimeWheel) Sleep(delay time.Duration) {
	queue := make(chan bool, 1)
	tw.addAny(delay,
		func() {
			queue <- true
		},
		modeNotCircle, modeNotAsync,
	)
	<-queue
}

func (tw *TimeWheel) After(delay time.Duration) <-chan time.Time {
	queue := make(chan time.Time, 1)
	tw.addAny(delay,
		func() {
			queue <- time.Now()
		},
		modeNotCircle, modeNotAsync,
	)
	return queue
}

// Stop 停止时间轮
func (tw *TimeWheel) Stop() {
	tw.stopC <- struct{}{}
}

//生成时间间隔
func (tw *TimeWheel) tickGenerator() {
	if tw.tickQueue == nil {
		return
	}

	for !tw.exited {
		select {
		case <-tw.ticker.C:
			select {
			case tw.tickQueue <- time.Now():
			default:
				panic("raise long time blocking")
			}
		}
	}
}

//定时执行任务
func (tw *TimeWheel) schduler() {
	queue := tw.ticker.C
	if tw.tickQueue != nil {
		queue = tw.tickQueue
	}

	for {
		select {
		case <-queue:
			tw.handleTick()

		case task := <-tw.addC:
			tw.put(task)

		case key := <-tw.removeC:
			tw.remove(key)

		case <-tw.stopC: //停止
			tw.exited = true
			tw.ticker.Stop()
			return
		}
	}
}

//清理任务
func (tw *TimeWheel) cleanTask(task *Task) {
	index := tw.bucketIndexes[task.id]
	delete(tw.bucketIndexes, task.id)
	delete(tw.buckets[index], task.id)
}

func (tw *TimeWheel) handleTick() {
	bucket := tw.buckets[tw.currentIndex]
	for k, task := range bucket {
		if task.stop {
			tw.cleanTask(task)
			continue
		}

		if bucket[k].round > 0 {
			bucket[k].round--
			continue
		}

		if task.async {
			go task.callback()
		} else {
			// optimize gopool
			task.callback()
		}

		// 循环执行
		if task.circle {
			tw.cleanTask(task)
			tw.putCircle(task, modeIsCircle)
			continue
		}

		// gc
		tw.cleanTask(task)
	}

	if tw.currentIndex == tw.bucketsNum-1 {
		tw.currentIndex = 0
		return
	}

	tw.currentIndex++
}

func (tw *TimeWheel) addAny(delay time.Duration, callback func(), circle, async bool) *Task {
	if delay <= 0 {
		delay = tw.tick
	}

	id := tw.genUniqueID()

	var task = new(Task)

	task.delay = delay
	task.id = id
	task.callback = callback
	task.circle = circle
	task.async = async

	tw.addC <- task
	return task
}

func (tw *TimeWheel) put(task *Task) {
	tw.store(task, false)
}

func (tw *TimeWheel) putCircle(task *Task, circleMode bool) {
	tw.store(task, circleMode)
}

func (tw *TimeWheel) store(task *Task, circleMode bool) {
	round := tw.calculateRound(task.delay)
	index := tw.calculateIndex(task.delay)

	if round > 0 && circleMode {
		task.round = round - 1
	} else {
		task.round = round
	}

	tw.bucketIndexes[task.id] = index
	tw.buckets[index][task.id] = task
}

func (tw *TimeWheel) calculateRound(delay time.Duration) (round int) {
	delaySeconds := delay.Seconds()
	tickSeconds := tw.tick.Seconds()
	round = int(delaySeconds / tickSeconds / float64(tw.bucketsNum))
	return
}

func (tw *TimeWheel) calculateIndex(delay time.Duration) (index int) {
	delaySeconds := delay.Seconds()
	tickSeconds := tw.tick.Seconds()
	index = (int(float64(tw.currentIndex) + delaySeconds/tickSeconds)) % tw.bucketsNum
	return
}

func (tw *TimeWheel) remove(task *Task) {
	tw.cleanTask(task)
}

//生成唯一任务ID
func (tw *TimeWheel) genUniqueID() taskID {
	id := atomic.AddInt64(&tw.randomID, 1)
	return taskID(id)
}
