package models

import (
	"flag"
)

var (
	Port     = flag.String("port", "8080", "Port number")
	Prt_help = flag.Bool("help", false, "Show help screen")
)
