package core

import (
	"log"
	"time"

	"github.com/dddddai/go-utils/linkedhashmap"
)

type Scheduler struct {
	linkedHashMap  linkedhashmap.LinkedHashMap
	serializer     Serializer
	schedulerDelay time.Duration
	queueCapacity  int
}

func (scheduler *Scheduler) persistsIfQueueIsFull() {
	log.Println("Start scheduling...")
	for {
		if scheduler.linkedHashMap.Size() < scheduler.queueCapacity {
			time.Sleep(scheduler.schedulerDelay)
			continue
		}

		pair, err := scheduler.linkedHashMap.PollLast()
		if err != nil {
			log.Panic(err.Error())
		}
		scheduler.serializer.MarkNodeToBeSerialized(pair.Value.(*TNode))
	}
}

func NewScheduler() *Scheduler {
	newScheduler := &Scheduler{linkedHashMap: *linkedhashmap.New(), queueCapacity: 1000, schedulerDelay: time.Minute / 3}
	go newScheduler.persistsIfQueueIsFull()
	return newScheduler
}
