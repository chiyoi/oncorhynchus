package main

import (
	"os"

	"github.com/chiyoi/oncorhynchus/internal/oncorhynchus"
)

func main() { oncorhynchus.Handler().Serve(os.Args[1:]) }
