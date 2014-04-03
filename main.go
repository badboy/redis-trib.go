package main

import (
	"fmt"
	"os"
)

func main() {
	//n := NewNode("localhost:7001")
	//fmt.Println("hello", n.slots)
	//fmt.Println("dead?", n.IsDead())
	//fmt.Println("alive?", n.IsAlive())
	//fmt.Println("cluster enabled?", n.IsClusterEnabled())
	//fmt.Println("empty?", n.IsEmpty())

	//c := NewCluster("localhost:7001")
	//fmt.Println("hello", c)
	//fmt.Println("nodes", c.FetchNodes())

	//c.Each(func(n *Node, result interface{}, err error, cmd string, args []interface{}) {
	//val, _ := redis.String(result, err)
	//fmt.Printf("%s: %s %s\n", n.String(), cmd, args)
	//fmt.Print(val)
	//fmt.Println("--")
	//}, "INFO", "memory")

	if len(os.Args) < 3 {
		fmt.Printf("usage: %s [command] [ip:port] [arguments]\n", os.Args[0])
		os.Exit(1)
	}

	c := NewCluster(os.Args[2])

	switch os.Args[1] {
	case "each":
		runEach(c)
	case "check":
		runCheck(c)
	default:
		fmt.Printf("command '%s' not implemented\n", os.Args[1])
	}

}

func runEach(c *Cluster) {
	if len(os.Args) < 4 {
		fmt.Printf("usage: %s each [ip:port] [arguments]\n", os.Args[0])
		os.Exit(1)
	}
	command := os.Args[3]

	args := ToInterfaceArray(os.Args[4:])

	_, err := c.EachPrint(command, args...)
	if err != nil {
		fmt.Println("Command failed:", err)
	}
}

func runCheck(c *Cluster) {
	c.CheckCluster()
}
