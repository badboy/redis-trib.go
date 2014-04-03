package main

import (
	"fmt"
	"os"
)

func main() {
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
