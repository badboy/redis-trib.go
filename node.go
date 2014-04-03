package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
)

const (
	UnusedHashSlot = iota
	NewHashSlot
)

type Node struct {
	id              string
	ip_port         string
	flags           []string
	master_id       string
	ping            string
	pong            string
	config          string
	state           string
	slots           map[int]int
	slots_importing map[int]string
	slots_migrating map[int]string
	client          redis.Conn
	dirty           bool
	last_error      error
}

func NewNode(ip_port string) (*Node, error) {
	n := &Node{
		ip_port:         ip_port,
		dirty:           false,
		slots:           make(map[int]int),
		slots_migrating: make(map[int]string),
		slots_importing: make(map[int]string),
	}
	_, err := n.LoadInfo()
	n.last_error = err
	return n, err
}

func (node *Node) LoadInfo() (*Node, error) {
	nodes, err := redis.String(node.Call("CLUSTER", "NODES"))

	if err != nil {
		//fmt.Println("Error in LoadInfo:", err)
		return node, err
	}

	n_nodes := strings.Split(nodes, "\n")
	for _, val := range n_nodes {
		parts := strings.Split(val, " ")
		if len(parts) <= 3 {
			continue
		}

		if strings.Contains(parts[2], "myself") {
			node.SetInfo(parts)
		}
	}

	return node, nil
}

func (node *Node) AddSlots(start, end int) {
	for i := start; i <= end; i++ {
		node.slots[i] = NewHashSlot
	}

	node.dirty = true
}

func (node *Node) SetInfo(parts []string) *Node {
	//id, ip_port, flags, master_id, ping, pong, config, state, *slots

	node.id = parts[0]
	node.flags = strings.Split(parts[2], ",")
	if !strings.Contains(parts[2], "myself") {
		node.ip_port = parts[1]
	}

	node.master_id = parts[3]
	node.ping = parts[4]
	node.pong = parts[5]
	node.config = parts[6]
	node.state = parts[7]

	if len(parts) > 7 {
		for i := 8; i < len(parts); i++ {
			slots := parts[i]
			if strings.Contains(slots, "<") {
				slot_id_str := strings.Split(slots, "-<-")
				slot_id, _ := strconv.Atoi(slot_id_str[0])
				node.slots_importing[slot_id] = slot_id_str[1]
			} else if strings.Contains(slots, ">") {
				slot_id_str := strings.Split(slots, "->-")
				slot_id, _ := strconv.Atoi(slot_id_str[0])
				node.slots_migrating[slot_id] = slot_id_str[1]
			} else if strings.Contains(slots, "-") {
				slot_id_str := strings.Split(slots, "-")
				first_id, _ := strconv.Atoi(slot_id_str[0])
				last_id, _ := strconv.Atoi(slot_id_str[1])
				node.AddSlots(first_id, last_id)
			} else {
				first_id, _ := strconv.Atoi(slots)
				node.AddSlots(first_id, first_id)
			}

		}
	}

	return node
}

func (node *Node) Call(cmd string, args ...interface{}) (interface{}, error) {
	c, err := node.Client()
	if err != nil {
		return c, err
	}

	return c.Do(cmd, args...)
}

func (node *Node) Client() (redis.Conn, error) {
	if node.client != nil {
		return node.client, nil
	}

	conn, err := redis.Dial("tcp", node.ip_port)
	node.client = conn
	return node.client, err
}

func (node *Node) IsDead() bool {
	for _, flag := range node.flags {
		switch flag {
		case "disconnected", "fail", "noaddr":
			return true
		}
	}

	return false
}

func (node *Node) IsAlive() bool {
	return !node.IsDead()
}

func (node *Node) IsClusterEnabled() bool {
	info, err := redis.String(node.Call("INFO", "cluster"))

	if err != nil {
		//fmt.Println("Error in IsClusterEnabled:", err)
		return false
	}

	return strings.Contains(info, "cluster_enabled:1")
}

func (node *Node) IsOnlyNode() bool {
	info, err := redis.String(node.Call("CLUSTER", "INFO"))

	if err != nil {
		//fmt.Println("Error in IsOnlyNode:", err)
		return false
	}

	return strings.Contains(info, "cluster_known_nodes:1")
}

func (node *Node) IsEmpty() bool {
	info, err := redis.String(node.Call("INFO", "keyspace"))

	if err != nil {
		//fmt.Println("Error in IsEmpty:", err)
		return false
	}

	return info == "# Keyspace\n"
}

func (node *Node) FlagIsSet(flag string) bool {
	for _, f := range node.flags {
		if strings.Contains(f, flag) {
			return true
		}
	}

	return false
}

func (node *Node) IsMaster() bool {
	return node.FlagIsSet("master")
}
func (node *Node) IsSlave() bool {
	return node.FlagIsSet("slave")
}
func (node *Node) IsFailed() bool {
	return node.FlagIsSet("fail")
}
func (node *Node) IsDisconnected() bool {
	return node.FlagIsSet("disconnected")
}
func (node *Node) HasNoAddress() bool {
	return node.FlagIsSet("noaddr")
}
func (node *Node) InFailState() bool {
	return node.IsFailed() ||
		node.IsDisconnected() ||
		node.HasNoAddress()
}

func (node *Node) String() string {
	return fmt.Sprintf("%s [%s]", node.id, node.ip_port)
}
