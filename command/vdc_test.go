package command

import (
	"fmt"
)

func ExampleextractBuildID() {
	fmt.Println(extractBuildID("/api/vdc-builds/3686"))
	// Output: 3686
}
