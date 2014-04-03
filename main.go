package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()

	app.Version = "0.0.1"
	app.Name = "redis-trib"
	app.Usage = "control your Redis Cluster"

	app.Flags = []cli.Flag{
		cli.BoolFlag{"verbose", "Verbose mode"},
	}

	app.Commands = []cli.Command{
		{
			Name:        "check",
			Usage:       "check the state of your Redis Cluster",
			Action:      runCheck,
			Description: "Checks for consistent config and that all slots are covered by one node.",
		},
		{
			Name:        "each",
			Usage:       "run command on each of the Cluster nodes",
			Action:      runEach,
			Description: "Runs the specified command against each of the Cluster nodes as reported by them and returns the result.",
		},
		{
			Name:        "create",
			Usage:       "create a cluster using the given nodes",
			Action:      runCreate,
			Description: "Creates the cluster and joins all given nodes and assigning slots to each. If --replicas <n> is given, it will assign <n> nodes as slaves.",
			Flags: []cli.Flag{
				cli.IntFlag{"replicas", 0, "Number of replicas to use"},
			},
		},
	}

	app.Run(os.Args)
}

func runEach(ctx *cli.Context) {
	args := ctx.Args()

	if len(args) < 2 {
		fmt.Println("Need the address of one node and a command to run. See 'help' for more.")
		os.Exit(2)
	}

	c := NewCluster(args[0])
	command := args[1]
	cmdArgs := ToInterfaceArray(args[2:])

	_, err := c.EachPrint(command, cmdArgs...)
	if err != nil {
		fmt.Println("Command failed:", err)
	}
}

func runCheck(ctx *cli.Context) {
	args := ctx.Args()
	if len(args) != 1 {
		fmt.Println("Need the address of one node. See 'help' for more.")
		os.Exit(2)
	}
	c := NewCluster(args[0])
	c.CheckCluster()
}

func runCreate(ctx *cli.Context) {
	args := ctx.Args()
	if len(args) < 1 {
		fmt.Println("Need the address of atleast one node. See 'help' for more.")
		os.Exit(2)
	}

	fmt.Println("[ERR] 'create' not implemented yet")
	os.Exit(2)
}
