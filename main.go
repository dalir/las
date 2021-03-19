package main

import (
	"log"
)

func main() {
	las := Las{}
	if err := las.Parse("./pointcloud_2.las"); err != nil {
		log.Fatalf("error in Parsing Las file. %v", err)
	}

	if err := las.Las2txt("./pointcloud_2.txt"); err != nil {
		log.Fatalf("error in converting to txt. %v", err)
	}

}
