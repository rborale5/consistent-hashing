package main

import "crypto/md5"

type HashRing struct {
	numReplicas uint32
	nodes       *LinkedList
	items       []*HashRingItem
	numNodes    uint32
	numItems    uint32
	hashFn      HashFunction
}
type HashRingNode struct {
	name string
}
type HashRingItem struct {
	node *HashRingNode
	hash uint64
}

const (
	HASH_FUNCTION_MD5 HashFunction = iota
	HASH_FUNCTION_SHA1
)

type HashFunction int

type LinkedList struct {
	data interface{}
	next *LinkedList
}

func (ring *HashRing) Hash(data []byte) (uint64, error) {
	return uint64(md5.Sum(data)), nil
}

func (ring *HashRing) AddNode(name string) error {

}
