package main

import (
	"fmt"
	"bytes"
	"os/exec"
	"os"
	"io"
)

type Pipeline_entry struct {
	Cmd *exec.Cmd;
	Pipe_reader *io.PipeReader;
	Pipe_writer *io.PipeWriter;
}

type Pipeline [][]string

type Pipeline_internal []Pipeline_entry

func (p Pipeline) Run(input io.Reader, output io.Writer) {
	pi := Make_pipeline(p)
	pi.Run(input, output)
}

func Make_pipeline(cmds Pipeline) Pipeline_internal {
	out := make(Pipeline_internal, 0)
	for _, e := range cmds {
		out = append(out, Pipeline_entry{Cmd: exec.Command(e[0], e[1:]...), Pipe_reader: nil, Pipe_writer: nil})
	}
	return out
}

func (p Pipeline_internal) Run(input io.Reader, output io.Writer) {
	p.Start(input, output)
	p.Finish()
}

func (p Pipeline_internal) Start(input io.Reader, output io.Writer) {
	for i, e := range p {
		if i == 0 {
			e.Cmd.Stdin = input
		}
		if i != len(p) - 1 {
			r, w := io.Pipe()
			p[i+1].Cmd.Stdin = r
			p[i+1].Pipe_reader = r
			p[i].Cmd.Stdout = w
			p[i].Pipe_writer = w
		}
		if i == len(p) - 1 {
			e.Cmd.Stdout = output
		}
	}
	for _, e := range p {
		e.Cmd.Start()
	}
}

func (p Pipeline) Start(input io.Reader, output io.Writer) Pipeline_internal {
	out := Make_pipeline(p)
	out.Start(input, output)
	return out
}

func (p Pipeline_internal) Finish() {
	for _, e := range p {
		e.Cmd.Wait()
		if e.Pipe_reader != nil {
			e.Pipe_reader.Close()
		}
		if e.Pipe_writer != nil {
			e.Pipe_writer.Close()
		}
	}
}

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
