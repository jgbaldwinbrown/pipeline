package main

import (
	"os/exec"
	"os"
	"io"
)

func main() {
	d := exec.Command("echo", "-e", "apple\n", "banana\n", "carrot\n")
	g := exec.Command("grep", "banana\\|carrot")
	r, w := io.Pipe()
	g.Stdin = r
	g.Stdout = os.Stdout
	d.Stdout = w
	d.Start()
	g.Start()
	d.Wait()
	w.Close()
	g.Wait()
	r.Close()
}
