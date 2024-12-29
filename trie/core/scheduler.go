package core

import (
	"runtime"
	"time"
)

type Scheduler struct {
	schedulerDelay time.Duration
	trie           *Trie
}

const (
	serializationCallDelayMinutes = 20
)

func (scheduler *Scheduler) persistsIfQueueIsFull() {

	for {
		time.Sleep(scheduler.schedulerDelay)

		scheduler.trie.lock.Lock()

		trieRoot := scheduler.trie.root

		iterations := 3
		distanceFromRoot := 3

		serializeNodesByLevelsRecursively(trieRoot, distanceFromRoot, iterations, scheduler.trie.serializer)
		runtime.GC()

		scheduler.trie.lock.Unlock()
	}
}

func getNodeCount(node *TNode) int {
	if node == nil {
		return 0
	}
	count := 1
	if node.Successors == nil {
		return count
	}
	for _, n := range node.Successors {
		count += getNodeCount(n)
	}
	return count
}

func serializeNodesByLevelsRecursively(node *TNode, distance int, iterations int, serializer *Serializer) {
	if iterations == 0 {
		return
	}
	serializationThreshold := 1000 * iterations

	nextNodes := getFarFromRootNodes(node, distance)

	for _, subNode := range nextNodes {
		if getNodeCount(subNode) > serializationThreshold {
			serializeNodesByLevelsRecursively(subNode, distance, iterations-1, serializer)
			serializer.MarkNodeToBeSerialized(subNode)
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
