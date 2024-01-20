package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func main() {

	// Parse arguments
	nameArg := flag.String("name", "", "Match only blocks with the specified name")
	shellArg := flag.String("shell", "", "Match only blocks with the specified shell")
	keepArg := flag.Bool("keep", false, "Keep the temporary files we created")
	runArg := flag.Bool("run", false, "Run the matching block(s)")
	flag.Parse()

	content, err := os.ReadFile("README.md")
	if err != nil {
		fmt.Print(err)
		return
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	// Create a reader
	reader := text.NewReader(content)

	// and a parser
	doc := md.Parser().Parse(reader)

	// Walk the AST of the parser.
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		// We only care about code blocks
		switch n := n.(type) {
		case *ast.FencedCodeBlock:

			// See if the block has an info-node
			if n.Info == nil {
				return ast.WalkContinue, nil
			}

			// OK get the segment/info node.
			segment := n.Info.Segment
			info := segment.Value(reader.Source())

			// Now get the body of the block
			var buf bytes.Buffer
			lines := n.Lines()
			for i := 0; i < lines.Len(); i++ {
				line := lines.At(i)
				buf.Write(line.Value(reader.Source()))
			}

			// Save the shell and name here.
			shll := ""
			name := ""

			// Now we should have: bash "foo", or similar
			found := false
			for i := 0; i < len(info); i++ {
				if found {
					name += string(info[i])
				} else {
					if info[i] == ' ' {
						found = true
					} else {
						shll += string(info[i])
					}
				}
			}

			// Skip this block if we're missing either a shell, or a name
			if shll == "" || name == "" {
				return ast.WalkContinue, nil
			}

			// Should we skip this block?
			// Default to not skipping.
			skip := false

			if *shellArg != "" || *nameArg != "" {
				skip = true
			}

			// Show some debugging
			if shll == *shellArg {
				skip = false
			}
			if name == *nameArg {
				skip = false
			}

			if !skip {
				if *runArg {

					// Create a temporary file
					file, err := os.CreateTemp(os.TempDir(), "rm")
					if err != nil {
						fmt.Printf("error writing temporary file %s\n", err.Error())
						return ast.WalkContinue, nil
					}

					// ensure we cleanup
					if *keepArg {
						fmt.Printf("wrote to %s\n", file.Name())
					} else {
						defer os.Remove(file.Name())
					}

					// Write the shebang + contents
					file.WriteString("#!" + shll + "\n" + buf.String())
					file.Close()

					// Make it executable
					_ = os.Chmod(file.Name(), 0755)

					// Execute the newly created file.
					cmd, err := exec.Command("/bin/sh", "-c", file.Name()).Output()
					if err != nil {
						fmt.Printf("error executing temporary file %s [shell:%s block:%s]", err, shll, name)
						return ast.WalkContinue, nil
					}

					// Show the output
					fmt.Printf("%s\n", cmd)
				} else {
					fmt.Printf("Shell:%s  Name:%s\n", shll, name)
					fmt.Printf("%s\n", buf.String())
				}
			}
			return ast.WalkContinue, nil
		default:
			return ast.WalkContinue, nil
		}
	})

	// At this point we've got a list
}
