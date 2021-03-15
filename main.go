package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"io"
	"bufio"
	"gocv.io/x/gocv"
	"path/filepath"
)

func getFullPath(dir string) (path string) {
	// Get path to the user's home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}

	path = usr.HomeDir + dir
	return
}

func getImages(path string) []string {
	var images []string

	fmt.Println("Scanning images in directory...")
	err := filepath.Walk(path, func(pth string, info os.FileInfo, err error) error {
		if(filepath.Ext(pth) == ".jpg" || filepath.Ext(pth) == ".png") {
			images = append(images, info.Name())
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return images
}

func deleteImages(images []string) {
	for i, img := range images {
		fmt.Printf("%d %s\n", i, img)
		e := os.Remove(img)
		if e != nil {
			log.Fatal(e)
		}
	}
}

func simpleDeletion(path string) {
	var img gocv.Mat
	var marked []string
	var text string
	reader := bufio.NewReader(os.Stdin)
	window := gocv.NewWindow("IMG")
	deletion := 1

	fmt.Println("Keys:\n  [0] Skip\n  [1] Mark for deletion\n  [s] Delete images and exit\n  [x] Cancel")
	err := filepath.Walk(path, func(pth string, info os.FileInfo, err error) error {
		if(filepath.Ext(pth) == ".jpg" || filepath.Ext(pth) == ".png") {
			fullpath := path + info.Name()
			img = gocv.IMRead(fullpath, gocv.IMReadColor)
			loop:
				for {
					window.IMShow(img)
					switch k := window.WaitKey(0); k {
					case int('0'):
						fmt.Println("Skipped")
						break loop
					case int('1'):
						marked = append(marked, fullpath)
						fmt.Println("Marked")
						break loop
					case int('s'):
						fmt.Println("Delete marked images? [y/n]")
						text, _ = reader.ReadString('\n')
						if text == "y\n" {
							deletion = 1
							return io.EOF
						} else {
							deletion = 0
						}
					case int('x'):
						fmt.Println("Are you sure you want to cancel? [y/n]")
						text, _ = reader.ReadString('\n')
						if text == "y\n" {
							deletion = 0
							return io.EOF
						} else {
							break
						}
					default:
						fmt.Println("Not a valid key")
					}
			}
			img.Close()
		}
		return nil
	})
	if err != nil && err != io.EOF {
		panic(err)
	}

	if deletion == 1 {
		deleteImages(marked)
	}
}

func main() {
	path := getFullPath("/Pictures/test/")
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("What function would you like to use? \n[0] Simple deletion")
	switch text, _ := reader.ReadString('\n'); text {
	case "0\n":
		fmt.Println("Simple deletion")
		simpleDeletion(path)
	default:
		log.Fatal("Not a valid option")
	}
}
