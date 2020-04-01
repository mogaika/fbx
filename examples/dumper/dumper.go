package main

import (
	"flag"
	"os"

	"github.com/mogaika/fbx"
)

func loadFbx(filename string) *fbx.FBX {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fbxData, err := fbx.Read(file)
	if err != nil {
		panic(err)
	}
	return fbxData
}

func saveFbx(filename string, f *fbx.FBX) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fbx.Write(file, f)
}

func main() {
	flag.Parse()

	f := loadFbx(flag.Arg(0))

	print(f.SPrint())

	saveFbx(flag.Arg(0)+".new.fbx", f)
}
