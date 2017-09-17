package watchdog

import (
	"time"
	"sync"
)

var watchdog Watchdog

type Watchdog struct {
	updateFreq   int64
	remindBefore int64
	remindAfter  int64

	onTick   func ()
	onUpdate func()

	ticks int64
}

type WatchdogConfig struct {
	UpdateFreq   int64 `json:"updateFreq"`
	RemindBefore int64 `json:"remindBefore"`
	RemindAfter  int64 `json:"remindAfter"`
}

func Initialize(config WatchdogConfig) {
	watchdog = Watchdog{
		updateFreq:   config.UpdateFreq,
		remindBefore: config.RemindBefore,
		remindAfter:  config.RemindAfter,
	}
}

func Get() *Watchdog {
	return &watchdog
}

func (w *Watchdog) OnUpdate(callback func()) {
	w.onUpdate = callback
}

func (w *Watchdog) OnTick(callback func()) {
	w.onTick = callback
}

func (w *Watchdog) Listen(wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		w.tick()
	}
	wg.Done()
}

func (w *Watchdog) tick() {
	w.ticks++
	if w.ticks > w.updateFreq {
		w.onUpdate()
	}
	w.onTick()
}
