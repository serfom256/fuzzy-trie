package core

import (
	"encoding/json"
	"log"
	"os"

	randomGenerator "github.com/serfom256/fuzzy-trie/trie/core/utils"
)

type Serializer struct {
}

const path = "cache/"

func (s *Serializer) DeserializeNode(node *TNode) []*TNode {

	var successors []*TNode

	file, err := os.ReadFile(path + (*node.SerializationId))
	for err != nil {
		log.Println(err.Error())
		file, err = os.ReadFile(path + (*node.SerializationId))
	}

	err = json.Unmarshal(file, &successors)
	if err != nil {
		panic(err)
	}

	err = os.Remove(path + (*node.SerializationId))
	for err != nil {
		log.Println(err.Error())
		err = os.Remove(path + (*node.SerializationId))
	}

	return successors
}

func (s *Serializer) MarkNodeToBeSerialized(node *TNode) {
	if node.SerializationId != nil {
		return
	}

	uid := randomGenerator.GenerateUUID(6)
	node.SerializationId = &uid

	s.serializeNode(node)
}

func (s *Serializer) serializeNode(node *TNode) {
	node.prepareToSerialization()

	data, sErr := json.Marshal(&node.Successors)
	if sErr != nil {
		panic("Cannot serialize node")

	}
	err := os.WriteFile(path+(*node.SerializationId), data, 0644)
	if err != nil {
		panic(err)
	}

	node.Successors = nil
}

func (s *Serializer) Init() {
	randomGenerator.Init()
}
