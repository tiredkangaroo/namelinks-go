package cli

import (
	"flag"
	"fmt"
)

var NewShort string
var NewLong string

func init() {
	flag.StringVar(&NewShort, "-ns", "", "Add a new short link.")
	flag.StringVar(&NewLong, "-nl", "", "Add a new long link.")
	flag.Parse()
}

func main() {
	if NewShort != "" && NewLong == "" {
		fmt.Println("Must specify new long (specified new short, but not new long).")
		return
	}
	if NewShort == "" && NewLong != "" {
		fmt.Println("Must specify new short (specified new long, but not new short).")
		return
	}
}
