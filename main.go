package main

import (
	"flag"
	"fmt"
	"github.com/alecthomas/chroma/quick"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"unsafe"
)

var (
	help    bool
	version bool
	input   string
	output  string
	theme   string
)

func init() {
	flag.BoolVar(&help, "h", false, "help info")
	flag.BoolVar(&version, "v", false, "version info")
	flag.StringVar(&input, "i", "", "input source code file path")
	flag.StringVar(&output, "o", "code.svg", "output svg file path")
	flag.StringVar(&theme, "t", "dracula", "highlight theme")
	flag.Usage = usage
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `code2svg tool in golang to generate svg image contains source code
Version: 0.0.1
Usage: code2svg [-hvio] [-h help] [-v version] [-i input source code file path] [-o output svg file path] [-t highlight theme]
Options
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if help {
		flag.Usage()
	} else if version {
		fmt.Println("version: 0.0.1")
	} else if input == "" || output == "" {
		fmt.Println("too less arguments, use '-h' to see help info")
		os.Exit(1)
	} else {
		start := time.Now()

		ext := filepath.Ext(input)
		file, err := os.Open(input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if oExt := filepath.Ext(output); oExt != "svg" {
			output = output[:len(output)-len(oExt)] + ".svg"
		}

		err = toSVG(file, ext, output, theme)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		absPath, _ := filepath.Abs(output)

		fmt.Println("Successfully generated svg at: ", absPath)
		fmt.Println("time consumption: ", time.Now().Sub(start).Seconds())
	}
}

func toSVG(in io.Reader, ext, output, theme string) (err error) {
	src, err := ioutil.ReadAll(in)
	if err != nil {
		return
	}

	// file already exits
	//if _, err := os.Stat(output); !os.IsNotExist(err) {
	//	return fmt.Errorf("can NOT generate svg file at %s due to file existed", output)
	//}
	_, err = os.Stat(filepath.Dir(output))
	if os.IsNotExist(err) {
		return os.MkdirAll(filepath.Dir(output), 0755)
	}
	file, err := os.Create(output)
	defer func() {
		_ = file.Close()
	}()

	if err != nil {
		return
	}

	err = quick.Highlight(file, bytes2string(src), ext, "svg", theme)

	return
}

func bytes2string(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
