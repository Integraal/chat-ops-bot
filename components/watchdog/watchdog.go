package watchdog

import (
	"time"
	"sync"
)

var watchdog Watchdog

type Watchdog struct {
	updateFreq   int64
	RemindBefore int64
	RemindAfter  int64
	DontRemindAfter  int64

	onTick   func ()
	onUpdate func()

	ticks int64
}

type WatchdogConfig struct {
	UpdateFreq   int64 `json:"updateFreq"`
	RemindBefore int64 `json:"remindBefore"`
	RemindAfter  int64 `json:"remindAfter"`
	DontRemindAfter  int64 `json:"dontRemindAfter"`
}

func Initialize(config WatchdogConfig) {
	watchdog = Watchdog{
		updateFreq:   config.UpdateFreq,
		RemindBefore: config.RemindBefore,
		RemindAfter:  config.RemindAfter,
		DontRemindAfter: config.DontRemindAfter,
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
	w.ticks = w.updateFreq
	w.tick()
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		w.tick()
	}
	wg.Done()
}

func (w *Watchdog) tick() {
	w.ticks++
	if w.ticks >= w.updateFreq {
		w.onUpdate()
		w.ticks = 0
	}
	w.onTick()
}
