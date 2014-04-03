package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/garyburd/redigo/redis"
)

const CLUSTER_HASH_SLOTS = 16384

type Cluster struct {
	nodes [](*Node)
	addrs []string
}

func NewCluster(addr string) *Cluster {
	c := &Cluster{addrs: []string{addr}}
	return c
}

func (c *Cluster) Nodes() ([]*Node, error) {
	if len(c.nodes) == len(c.addrs) {
		return c.nodes, nil
	}

	nodes := make([](*Node), len(c.addrs))
	for i, addr := range c.addrs {
		n, err := NewNode(addr)
		if err != nil {
			return nil, err
		}
		nodes[i] = n
	}

	c.nodes = nodes
	return c.nodes, nil
}

func (c *Cluster) FetchNodes() ([]*Node, error) {
	i := rand.Intn(len(c.addrs))

	all_nodes, err := c.Nodes()
	if err != nil {
		return nil, err
	}

	node := all_nodes[i]
	c_nodes, err := redis.String(node.Call("CLUSTER", "NODES"))

	if err != nil {
		return nil, err
	}

	nodes := [](*Node){}
	n_nodes := strings.Split(c_nodes, "\n")
	for _, n := range n_nodes {
		parts := strings.Split(n, " ")
		if len(parts) <= 3 {
			continue
		}

		ip_port := parts[1]
		flags := parts[2]

		if strings.Contains(flags, "myself") {
			nodes = append(nodes, node)
		} else {
			new_node, err := NewNode(ip_port)
			if err != nil {
				new_node.SetInfo(parts)
			}
			nodes = append(nodes, new_node)
		}
	}

	c.nodes = nodes
	return c.nodes, nil
}

type InterfaceErrorCombo struct {
	result interface{}
	err    error
}

type EachFunction func(*Node, interface{}, error, string, []interface{})

func (c *Cluster) Each(f EachFunction, cmd string, args ...interface{}) ([]*InterfaceErrorCombo, error) {
	nodes, err := c.FetchNodes()

	if err != nil {
		return nil, err
	}

	ies := make([]*InterfaceErrorCombo, len(nodes))

	for i, node := range nodes {
		val, err := node.Call(cmd, args...)
		ie := &InterfaceErrorCombo{val, err}
		ies[i] = ie

		if f != nil {
			f(node, val, err, cmd, args)
		}
	}

	return ies, nil
}

func (c *Cluster) EachPrint(cmd string, args ...interface{}) ([]*InterfaceErrorCombo, error) {
	return c.Each(func(n *Node, result interface{}, err error, cmd string, args []interface{}) {
		val, _ := redis.String(result, err)

		if len(args) > 0 {
			string_args := ToStringArray(args)
			fmt.Printf("%s: %s %s\n", n.String(), cmd, strings.Join(string_args, " "))
		} else {
			fmt.Printf("%s: %s %s\n", n.String(), cmd)
		}

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(strings.Trim(val, " \n"))
		}
		fmt.Println("--")
	}, cmd, args...)
}
