package painter

import (
	"os"

	"github.com/alecthomas/chroma/quick"
)

func PaintCpp(cppCode string) {
	// "monokai" is the theme, "terminal256" is the formatter
	err := quick.Highlight(os.Stdout, cppCode, "cpp", "terminal256", "monokai")
	if err != nil {
		panic(err)
	}
}
