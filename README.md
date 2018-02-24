# go-renamer
A simple renaming program written in go

A quick tutorial for its use is as follows:

```
Usage:
  renamer [OPTIONS]

Application Options:
  -d, --dir=         Specify the directory where rename is to take place.
  -i, --interactive  Make the renaming interactive (Enter a new file name for
                     each file)
  -r, --recursive    recursively traverse the subdirectories to find more files
                     to rename.

Help Options:
  -h, --help         Show this help message
 ```
  
Incase the file contains any extensions, they will be preserved.

Eg. 
if the file's name is `test.jpg` and you rename it to `test1`, then the file will be saved as `test1.jpg`
