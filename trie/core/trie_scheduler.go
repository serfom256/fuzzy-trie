package core

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

type Scheduler struct {
	schedulerDelay time.Duration
	trie           *Trie
}

const serializationNodeAmount = 1000

func (scheduler *Scheduler) persistsIfQueueIsFull() {
	log.Println("Start scheduling...")
	for {
		time.Sleep(scheduler.schedulerDelay)
		PrintMemUsage()
		println("Start")
		trieRoot := scheduler.trie.root
		scheduler.trie.lock.Lock()
		println(getNodeCount(trieRoot))
		serializeNodesByLevelsRecursively(trieRoot, 3, 2, scheduler.trie.serializer)
		runtime.GC()
		println("Done")
		PrintMemUsage()
		println(getNodeCount(trieRoot))
		scheduler.trie.lock.Unlock()
	}
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
func getNodeCount(node *TNode) int {
	count := 1
	if node == nil || node.Successors == nil {
		return count
	}
	for _, n := range node.Successors {
		count += getNodeCount(n)
	}
	return count
}

func serializeNodesByLevelsRecursively(node *TNode, distance int, rootDistance int, serializer *Serializer) {
	if rootDistance == 0 {
		return
	}
	nextNodes := getFarFromRootNodes(node, distance)
	for _, subNode := range nextNodes {
		if getNodeCount(subNode) > serializationNodeAmount {
			serializeNodesByLevelsRecursively(subNode, distance, rootDistance-1, serializer)
			if getNodeCount(subNode) > serializationNodeAmount {
				serializer.MarkNodeToBeSerialized(subNode)
			}
		}
	}
}

func getFarFromRootNodes(node *TNode, distance int) []*TNode {
	nodeList := []*TNode{}
	if distance == 0 {
		return append(nodeList, node)
	}

	if node == nil || node.SuccessorsCount == 0 {
		return nodeList
	}

	for _, n := range node.Successors {
		nextNodes := getFarFromRootNodes(n, distance-1)
		nodeList = append(nodeList, nextNodes...)
	}

	return nodeList
}

func NewScheduler(trie *Trie) *Scheduler {
	newScheduler := &Scheduler{schedulerDelay: time.Minute * 10, trie: trie}
	go newScheduler.persistsIfQueueIsFull()
	return newScheduler
}
