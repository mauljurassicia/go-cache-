package main

import "fmt"

type Node struct {
	Key   interface{}
	Val   interface{}
	Left  *Node
	Right *Node
}

type Queue struct {
	Head   *Node
	Tail   *Node
	Length int
}

type Cache struct {
	Queue    *Queue
	Hash     *Hash
	Capacity int
}

type Hash map[interface{}]*Node

func NewCache(size int) *Cache {
	return &Cache{Queue: NewQueue(), Hash: &Hash{}, Capacity: size}
}

func NewQueue() *Queue {
	return &Queue{Length: 0}
}

func (c *Cache) Put(key interface{}, value interface{}) {
	if c.Queue.Length == 0 {
		node := &Node{Val: value, Key: key}
		(*c.Hash)[key] = node
		c.Queue.Head = node
		c.Queue.Tail = node
		c.Queue.Length++
		return
	}

	if val, exist := (*c.Hash)[key]; exist {
		val.Val = value
		c.moveToResent(val)
		return
	}

	newNode := &Node{Val: value, Right: c.Queue.Head}
	(*c.Hash)[key] = newNode
	c.Queue.Head.Left = newNode
	c.Queue.Head = newNode
	c.Queue.Length++

	if c.Queue.Length > c.Capacity {
		c.removeTail()
		c.Queue.Length--
	}

}

func (c *Cache) Get(key interface{}) interface{} {
	if val, exist := (*c.Hash)[key]; exist {
		if c.Queue.Length == 1 {
			return val.Val
		}

		c.moveToResent(val)
		return val.Val
	}
	return nil
}

func (c *Cache) Delete(key interface{}) {
	c.remove(key)
}

func (c *Cache) Length() int {
	return c.Queue.Length
}

func (c *Cache) remove(key interface{}) {
	node := (*c.Hash)[key]
	if c.Queue.Head == node {
		c.removeHead()
		c.Queue.Length--
		return
	}
	if c.Queue.Tail == node {
		c.removeTail()
		c.Queue.Length--
		return
	}
	c.takeOut(node)
	delete(*c.Hash, key)
	c.Queue.Length--
}

func (c *Cache) takeOut(node *Node) {
	if node.Left != nil {
		node.Left.Right = node.Right
	}
	if node.Right != nil {
		node.Right.Left = node.Left
	}
}

func (c *Cache) moveToResent(node *Node) {
	c.takeOut(node)
	node.Left = nil
	if c.Queue.Head != node {
		c.Queue.Head.Left = node
		node.Right = c.Queue.Head
	}
	c.Queue.Head = node
}

func (c *Cache) removeTail() {
	removedKey := c.Queue.Tail.Key
	delete(*c.Hash, removedKey)
	c.Queue.Tail = c.Queue.Tail.Left
	if c.Queue.Tail != nil {
		c.Queue.Tail.Right = nil
	}

}

func (c *Cache) removeHead() {
	removedKey := c.Queue.Head.Key
	delete(*c.Hash, removedKey)
	c.Queue.Head = c.Queue.Head.Right
	if c.Queue.Head != nil {
		c.Queue.Head.Left = nil
	}
}

func main() {
	fmt.Println("Start Cache")

	cache := NewCache(15)

	cache.Put("makan", "nasi")
	fmt.Println(cache.Get("makan"), cache.Queue.Head.Right)
	cache.Put("makan", "tempe")

	fmt.Println(cache.Get("makan"), cache.Queue.Length, cache.Queue.Head.Right)
	cache.Put("Baso", "nasi")
	fmt.Println(cache.Queue.Length)
	cache.Delete("Baso")
	fmt.Println(cache.Length(), cache.Capacity, *cache.Queue.Head, *cache.Queue.Tail)

}
