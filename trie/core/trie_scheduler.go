package core

import (
	"log"
	"runtime"
	"time"
)

type Scheduler struct {
	schedulerDelay time.Duration
	trie           *Trie
}

const (
	serializationNodeAmount       = 1000
	serializationCallDelayMinutes = 10
)

func (scheduler *Scheduler) persistsIfQueueIsFull() {
	log.Println("\nStart scheduling...")
	for {
		time.Sleep(scheduler.schedulerDelay)

		scheduler.trie.lock.Lock()

		trieRoot := scheduler.trie.root
		DisplayMemoryUsage(trieRoot)

		serializeNodesByLevelsRecursively(trieRoot, 3, 2, scheduler.trie.serializer)
		gc()

		DisplayMemoryUsage(trieRoot)
		scheduler.trie.lock.Unlock()
	}
}

func gc() {
	runtime.GC()
}

func DisplayMemoryUsage(trieRoot *TNode) {
	log.Printf("\nTrie node count: %v", getNodeCount(trieRoot))
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	memAlloc := mem.Alloc / 1024 / 1024
	log.Printf("\nMemory allocated: %v MiB", memAlloc)
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
	newScheduler := &Scheduler{schedulerDelay: time.Minute * serializationCallDelayMinutes, trie: trie}
	go newScheduler.persistsIfQueueIsFull()
	return newScheduler
}
