package main

import (
	"flag"
	"io/ioutil"
	"log"
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

	for _, filename := range flag.Args() {
		log.Printf("Loading file %v", filename)
		f := loadFbx(filename)

		f.PrintConnectionsTree(0)

		log.Printf("Creating dump for %v", filename)
		ioutil.WriteFile(filename+".txt", []byte(f.SPrint()), 666)

		saveFbx(filename+".new.fbx", f)
	}
}
