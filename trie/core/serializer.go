package core

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type Serializer struct {
	position int
}

const cachePath = "cache/"
const fileMode = 0777

func (s *Serializer) DeserializeNode(node *TNode) []*TNode {

	var successors []*TNode

	filePath := getFilePath(node.SerializationId)

	file, err := os.ReadFile(filePath)
	for err != nil {
		log.Println(err.Error())
		file, err = os.ReadFile(filePath)
	}

	err = json.Unmarshal(file, &successors)
	if err != nil {
		panic(err)
	}

	err = os.Remove(filePath)
	for err != nil {
		log.Println(err.Error())
		err = os.Remove(filePath)
	}

	return successors
}

func getFilePath(sId *int) string {
	return cachePath + strconv.Itoa(*sId) + ".bin"
}

func (s *Serializer) MarkNodeToBeSerialized(node *TNode) {
	if node.SerializationId != nil {
		return
	}

	node.SerializationId = s.nextUid()

	s.serializeNode(node)
}

func (s *Serializer) serializeNode(node *TNode) {
	node.prepareToSerialization()

	data, sErr := json.Marshal(&node.Successors)
	if sErr != nil {
		panic("Cannot serialize node")

	}
	err := os.WriteFile(getFilePath(node.SerializationId), data, fileMode)
	if err != nil {
		panic(err)
	}

	node.Successors = nil
}

func (s *Serializer) Init() {
	err := os.Mkdir(cachePath, fileMode)
	if err != nil {
		log.Println("Cannot create [cache] directory, probably this directory already exists")
	}
}

func (s *Serializer) nextUid() *int {
	position := s.position

	s.position = s.position + 1
	return &position
}
