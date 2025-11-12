package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fengqi/pve-pct/client"
	"github.com/fengqi/pve-pct/config"
)

func main() {
	flag.Parse()
	c := config.Init()
	s := client.Init(c)
	code, out := s.Exec(flag.Args())
	fmt.Println(out)
	os.Exit(code)
}
