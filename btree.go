package main

import "encoding/binary"

const HEADER = 4

const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

type BNode []byte

type BTree struct {
	root uint64
	get  func(uint64) []byte
	new  func([]byte) uint64
	del  func(uint64)
}

const (
	BNODE_NODE = 1 // internal node
	BNODE_LEAF = 2 // leaf node
)

func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

func (node BNode) getPtr(idx uint16) uint64 {
	if idx >= node.nkeys() || idx < 0 {
		panic("idx out of range") // idx <- [0, 1, ..... 7] for 8 keys
	}

	pos := HEADER + 8*idx // header offset plus idx byte position
	return binary.LittleEndian.Uint64(node[pos:])
}

func (node BNode) setPtr(idx uint16, val uint64) {
	if idx >= node.nkeys() || idx < 0 {
		panic("idx out of range")
	}

	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node[pos:pos+8], val)
}

func offsetPos(node BNode, idx uint16) uint16 {
	// We do not consider idx 0 b.c. idx 0 always has offset 0
	if idx < 1 || idx >= node.nkeys() {
		panic("idx out of range")
	}
	return HEADER + node.nkeys()*8 + (idx-1)*2
}

func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16(node[offsetPos(node, idx):])
}

func (node BNode) setOffset(idx uint16, offset uint16) {
	binary.LittleEndian.PutUint16(node[offsetPos(node, idx):], offset)
}
