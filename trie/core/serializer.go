package core

import (
	"bytes"
	"encoding/gob"
	randomGenerator "github.com/serfom256/fuzzy-trie/trie/core/utils"
	"math/rand"
	"os"
	"sync"
)

type Serializer struct {
	serialized   int
	deserialized int
}

const path = "cache/"

func (s *Serializer) DeserializeNode(node *TNode) []*TNode {

	var successors []*TNode
	b := bytes.Buffer{}

	file, err := os.ReadFile(path + (*node.SerializationId))
	if err != nil {
		panic(err)
	}

	b.Write(file)
	d := gob.NewDecoder(&b)

	err = d.Decode(&successors)
	if err != nil {
		panic(err)
	}

	err = os.Remove(path + (*node.SerializationId))
	if err != nil {
		panic(err)
	}

	return successors
}

func (s *Serializer) MarkNodeToBeSerialized(node *TNode, lock *sync.RWMutex) {
	if rand.Intn(3) == 0 && node.SerializationId != nil {
		return
	}

	lock.Lock()

	if node.SerializationId == nil {
		uid := randomGenerator.GenerateUUID(6)
		node.SerializationId = &uid
	}
	s.serialized++

	s.serializeNode(node)

	lock.Unlock()
}

func (s *Serializer) serializeNode(node *TNode) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	node.prepareToSerialization()

	err := encoder.Encode(node.Successors)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(path+(*node.SerializationId), buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	node.Successors = nil
}

func (s *Serializer) Init() {
	randomGenerator.Init()
}
