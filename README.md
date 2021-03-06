# Futil 
<!-- TOC -->

- [Futil](#futil)
    - [Download](#download)
    - [TODO](#todo)
    - [Examples](#examples)

<!-- /TOC -->


A small utility library for managing files. Created primarily for managing my music library.

## Download
`go get -u github.com/Necroforger/futil`

## TODO
- Add progress tracking via channels

## Examples

```go
package main

import (
	"log"
	"os"

	"github.com/Necroforger/futil"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
  // Zip
  must(futil.ZipDir("folder", "zipped.zip"))
  must(futil.Unzip("zipped.zip", "unzippedfolder"))
  
  must(os.Remove("zipped.zip"))
  must(os.RemoveAll("unzippedfolder"))

  // Copy
  must(futil.CpDir("from", "to"))
  must(os.RemoveAll("to"))

  // Move
  must(futil.MvDir("from", "to"))
  must(futil.MvDir("to", "from"))
}
```