package main

import (
	"github.com/jgbaldwinbrown/pipeline"
	"fmt"
)

func main() {
	p := Pipeline {
		{"echo", "-e", "apple\nbanana\ncarrot\n"},
		{"grep", "banana\\|carrot"},
	}
	p2 := p.Start(nil, os.Stdout)
	p2.Finish()
	p.Run(nil, os.Stdout)
	p = Pipeline {
		{"awk", "{print($2, $3)}"},
		{"grep", "banana"},
	}
	var b bytes.Buffer
	var b2 bytes.Buffer
	b.WriteString("apple	banana	carrot\n")
	for i:=0; i<1000000; i++ {
		b.WriteString("apple	banana	banana apple\n")
	}
	b.WriteString("apple	grape	toucan apple\n")
	p.Run(&b, &b2)
	fmt.Printf("%v", b2.String())
}
