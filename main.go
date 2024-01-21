package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var (

	//nameArg - Match only blocks with the given name
	nameArg *string

	// shellArg - Match only blocks with the given shell.
	shellArg *string

	// combineArg - If true then combine all the matching blocks into one.
	combineArg *bool

	// joinArg - If true write all blocks to one file.
	joinArg *bool

	// keepArg - If true don't delete the temporary file we create.
	keepArg *bool

	// runArg - If true run matching blocks
	runArg *bool
)

// CodeBlock holds the details for each code-block within a given file.
type CodeBlock struct {
	// Shell is the defined shell for this block
	Shell string

	// Name is the name of the block
	Name string

	// Content is the code within the block
	Content string
}

// parseBlocks reads the the given file and returns a structure
// containing each of the fenced code-blocks.
func parseBlocks(file string) ([]CodeBlock, error) {

	// Result we return to the caller
	var res []CodeBlock

	// Read the file
	content, err := os.ReadFile(file)
	if err != nil {
		return res, fmt.Errorf("error reading input %s: %s", file, err)
	}

	// Create the new helper
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	// Create a reader for processing the content.
	reader := text.NewReader(content)

	// Create the parser object
	doc := md.Parser().Parse(reader)

	// Walk the AST of the parser.
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		// We only care about code blocks
		switch n := n.(type) {
		case *ast.FencedCodeBlock:

			// If the block has no info-node then return
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
			if shll != "" && name != "" {

				res = append(res, CodeBlock{
					Name:    name,
					Shell:   shll,
					Content: buf.String(),
				})
			}
			return ast.WalkContinue, nil
		default:
			return ast.WalkContinue, nil
		}
	})

	return res, nil
}

// filterBlocks takes the list of fenced code blocks and
// returns the results of filtering those blocks
func filterBlocks(in []CodeBlock) []CodeBlock {

	// No filters?  Return them all
	if *nameArg == "" && *shellArg == "" {
		return in
	}

	// Otherwise we need to filter, so store the return value here
	var res []CodeBlock

	// Process each one
	for _, block := range in {

		// Matching name?
		if *nameArg != "" && block.Name == *nameArg {
			res = append(res, block)

			// Only add once, even if it matches name AND shell
			continue
		}

		// Matching shell?
		if *shellArg != "" && strings.Contains(block.Shell, *shellArg) {
			res = append(res, block)

			// Only add once, even if it matches name AND shell
			continue
		}
	}

	return res
}

func main() {

	// setup the flags - strings
	nameArg = flag.String("name", "", "Match only blocks with the specified name")
	shellArg = flag.String("shell", "", "Match only blocks with the specified shell")

	// setup the flags - bools
	joinArg = flag.Bool("join", false, "Join all matching blocks into one run")
	keepArg = flag.Bool("keep", false, "Keep the temporary files we created")
	runArg = flag.Bool("run", false, "Run the matching block(s)")

	// Parse the arguments
	flag.Parse()

	// Ensure we have a list of files.
	if len(flag.Args()) < 1 {
		fmt.Printf("Usage: runme [args] file1.md file2.md ..\n")
		return
	}

	// Process each file
	for _, file := range flag.Args() {

		// Get the blocks from within the file
		blocks, err := parseBlocks(file)

		// If there were errors we're done
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}

		// Now filter the blocks, based on --name/--shell
		blocks = filterBlocks(blocks)

		// And finally process them
		for _, block := range blocks {

			// If we're not running we just show the details
			if !*runArg {
				fmt.Printf("Shell:%s  Name:%s\n", block.Shell, block.Name)
				fmt.Printf("%s\n", block.Content)
			}

			//
			// Running here
			//

			// Create a temporary file
			file, err := os.CreateTemp(os.TempDir(), "rm")
			if err != nil {
				fmt.Printf("error writing temporary file %s\n", err.Error())
				return
			}

			// ensure we cleanup
			if *keepArg {
				fmt.Printf("wrote to %s\n", file.Name())
			} else {
				defer os.Remove(file.Name())
			}

			// Write the shebang + contents
			file.WriteString("#!" + block.Shell + "\n" + block.Content)
			file.Close()

			// Make it executable
			_ = os.Chmod(file.Name(), 0755)

			// Execute the newly created file.
			cmd, err := exec.Command("/bin/sh", "-c", file.Name()).Output()
			if err != nil {
				fmt.Printf("error executing temporary file %s [shell:%s block:%s]", err, block.Shell, block.Name)
				return
			}

			// Show the output
			if len(cmd) > 0 {
				fmt.Printf("%s", cmd)
			} else {
				fmt.Printf("[no output]\n")
			}
		}
	}
}
