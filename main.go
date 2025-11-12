package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fengqi/pve-pct/client"
	"github.com/fengqi/pve-pct/config"
)

func init() {
	flag.Parse()
}

func main() {
	c := config.Init()
	s := client.Init(c)
	code, out := s.Exec(flag.Args())
	fmt.Println(out)
	os.Exit(code)
}
