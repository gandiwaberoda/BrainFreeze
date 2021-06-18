package diagnostic

import "time"

type FpsGauge struct {
	frameCount int
	fps        int
	lastCheck  time.Time
	ticker     *time.Ticker
	doneChan   chan bool
}

func NewFpsGauge() *FpsGauge {
	return &FpsGauge{
		frameCount: 0,
		fps:        -1,
		lastCheck:  time.Now(),
		doneChan:   make(chan bool),
	}
}

func (f *FpsGauge) Start() {
	f.ticker = time.NewTicker(time.Second * 2)

	go func() {
		for {
			select {
			case <-f.doneChan:
				return
			case <-f.ticker.C:
				elapsed := time.Since(f.lastCheck)
				f.fps = f.frameCount / int(elapsed.Seconds())
				f.lastCheck = time.Now()
				f.frameCount = 0
			}
		}
	}()
}

func (f *FpsGauge) Stop() {
	f.ticker.Stop()
}

func (f *FpsGauge) Tick() {
	f.frameCount++
}

func (f *FpsGauge) Read() int {
	return f.fps
}
