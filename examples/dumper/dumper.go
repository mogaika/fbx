package main

import (
	"flag"
	"os"

	"github.com/mogaika/fbx"
)

func main() {
	flag.Parse()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fbxData, err := fbx.Read(f)
	if err != nil {
		panic(err)
	}

	print(fbxData.SPrint())
}
