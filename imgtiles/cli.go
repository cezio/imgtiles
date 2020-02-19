package imgtiles

import (
	"fmt"
)

func Main() {
	var opts = ParseArgs()
	var _, err = Run(opts)
	if err != nil {
		fmt.Println("Error while processing", opts, err)
	}
}
