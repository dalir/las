## usage

```go
package main

import (
	"github.com/dalir/las"
	"log"
)

func main() {
	las := las.Las{}
	if err := las.Parse("./pointcloud.las"); err != nil {
		log.Fatalf("error in Parsing Las file. %v", err)
	}

	if err := las.Las2txt("./pointcloud.txt"); err != nil {
		log.Fatalf("error in converting to txt. %v", err)
	}

}
```