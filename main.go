package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"renamer/flagsParser"
	"strconv"
	"strings"
	"syscall"

	input "github.com/tcnksm/go-input"
)

func main() {

	flags := flagsParser.New()
	Start(flags)
}

type Renamer struct {
	directoryPath string  // the directory path where the renaming will take place
	totalFiles    int     // holds the total files in the directory
	opt           options // options for the cmd invoking
	index         int
	baseFileName  string   // name of the original file
	ui            input.UI // input receiver from the tcknsm package
}

type options struct {
	using_r bool //using recursive flag
	using_i bool //using interactive flag
}

func (o *options) check_options(flags *flagsParser.Flags) {
	o.using_i = flags.Interactive

	if flags.Recursive {
		o.using_r = true
	}

}

func Start(flags *flagsParser.Flags) *Renamer {
	var err error
	var r Renamer
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	r.directoryPath, err = filepath.Abs(flags.Directory)
	if err != nil {
		log.Fatal("Error in making absolute path")
	}
	r.opt.check_options(flags)
	go func() {
		<-sigs
		done <- true
		os.Exit(1)

	}()
	r.index = 0
	if !r.opt.using_i {
		r.baseFileName, err = r.ui.Ask("Please enter a file name", &input.Options{
			Required:  true,
			Loop:      true,
			HideOrder: true,
		})

	}
	r.rename(r.directoryPath)
	return &r
}

func (r *Renamer) rename(directoryPath string) {

	r.ui = input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}
	rootDir := true
	err := filepath.Walk(directoryPath, func(path string, file os.FileInfo, err error) error {

		if file.IsDir() {
			if r.opt.using_r {
				log.Println("Entering directory : " + file.Name())
			} else {
				if !rootDir { // to ensure that that the root directory is traversed
					rootDir = false
					return filepath.SkipDir
				}
			}
		} else {

			basePath, fileName := filepath.Split(path)

			if r.opt.using_i {
				// interactive renaming
				query := "Current File name is: " + file.Name() + "\nEnter new name (enter to skip) "
				name, err := r.ui.Ask(query, &input.Options{
					Default:     file.Name(),
					Required:    true,
					Loop:        true,
					HideDefault: true,
					HideOrder:   true,
				})
				if err != nil {
					log.Println("error")
				}
				if name != fileName {
					//user wants to change the name
					fileName = name
				}
			} else {
				// batch renaming
				fileName = r.getNewFileName(fileName)
				log.Println("new name set to: " + fileName)
			}
			os.Rename(path, basePath+fileName)
		}
		return nil
	})

	if err != nil {
		log.Println("error")
	}

}

//returns a new name if the rename is not done interactively
func (r *Renamer) getNewFileName(fileName string) string {
	r.index++
	fileData := strings.Split(fileName, ".")
	if len(fileData) > 1 {
		return r.baseFileName + strconv.Itoa(r.index) + "." + fileData[len(fileData)-1] //adds extension if found
	}
	//creates a file name by adding the base file name plus (String)index
	return r.baseFileName + strconv.Itoa(r.index)
}
