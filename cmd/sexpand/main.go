package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lcrownover/sexpand/pkg/sexpand"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: provide a SLURM node range expression as the only arg")
		os.Exit(1)
	}
	nodes, err := sexpand.SExpand(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(strings.Join(nodes, ","))
}
