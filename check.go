package main

import (
	"fmt"
)

func (c *Cluster) PrintNodes() {
	nodes, _ := c.FetchNodes()

	for _, node := range nodes {
		node_type := ""
		if node.IsMaster() {
			node_type = "master"
		} else {
			node_type = "slave "
		}
		if node.IsFailed() {
			node_type = node_type + ",failed"
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
	nodes, _ := c.FetchNodes()

	clean := true
	old_sig := ""
	for _, node := range nodes {
		if len(old_sig) == 0 {
			old_sig = node.GetConfigSignature()
		} else {
			new_sig := node.GetConfigSignature()
			if old_sig != new_sig {
				fmt.Println("[ERR] Signatures don't match. Error in Config.")
				fmt.Println("      Error came up when checking node", node.String())
				clean = false
				break
			}
		}
	}

	if clean {
		fmt.Println("[OK] Config consistent.")
	}
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
