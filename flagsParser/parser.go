package flagsParser

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type Flags struct {
	Directory   string `short:"d" long:"dir" description:"Specify the directory where rename is to take place." required:"true"`
	Interactive bool   `short:"i" long:"interactive" description:"Make the renaming interactive (Enter a new file name for each file)"`
	Recursive   bool   `short:"r" long:"recursive" description:"recursively traverse the subdirectories to find more files to rename."`
}

func New() *Flags {

	var f Flags

	if _, err := flags.Parse(&f); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	return &f
}
