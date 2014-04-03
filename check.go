package main

import (
	"fmt"
)

//func (c *Cluster) EachPrint(cmd string, args ...interface{}) ([]*InterfaceErrorCombo, error) {

func (c *Cluster) PrintNodes() {
	nodes, _ := c.FetchNodes()

	for _, node := range nodes {
		node_type := ""
		if node.IsMaster() {
			node_type = "master"
		} else {
			node_type = "slave "
		}
		fmt.Println(node, node_type, "#slots:", len(node.slots))
	}
}

func (c *Cluster) CheckCluster() {
	c.PrintNodes()
	c.CheckConfigConsistency()
	c.CheckOpenSlots()
	c.CheckSlotsCoverage()
}

func (c *Cluster) CheckConfigConsistency() {
}

func (c *Cluster) CheckOpenSlots() {
}

func (c *Cluster) CheckSlotsCoverage() {
	all_slots := c.coveredSlots()

	if len(all_slots) == CLUSTER_HASH_SLOTS {
		fmt.Printf("[OK] All %d slots covered.\n", CLUSTER_HASH_SLOTS)
	} else {
		fmt.Printf("[ERR] Not all %d slots covered are covered by nodes.\n", CLUSTER_HASH_SLOTS)
	}
}

func (c *Cluster) coveredSlots() map[int]int {
	nodes, _ := c.FetchNodes()

	slots := make(map[int]int)

	for _, node := range nodes {
		for key, val := range node.slots {
			slots[key] = val
		}
	}

	return slots
}
