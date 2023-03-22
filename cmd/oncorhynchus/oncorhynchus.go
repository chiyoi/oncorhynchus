package main

import (
	"os"

	"github.com/chiyoi/oncorhynchus/internal/app/oncorhynchus"
)

func main() { oncorhynchus.Handler().Serve(os.Args[1:]) }
