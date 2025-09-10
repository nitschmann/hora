package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nitschmann/hora/internal/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	horaCmd := cmd.NewRootCmd()
	err := doc.GenMarkdownTree(horaCmd, "docs/cli")
	if err != nil {
		fmt.Println("Error generating docs:", err)
		return
	}

	horaFile := filepath.Join("docs", "cli", "hora.md")
	readmeFile := filepath.Join("docs", "cli", "README.md")

	err = os.Rename(horaFile, readmeFile)
	if err != nil {
		fmt.Println("Error renaming hora.md to README.md:", err)
		return
	}

	err = updateCrossReferences("docs/cli")
	if err != nil {
		fmt.Println("Error updating cross-references:", err)
		return
	}

	fmt.Println("Documentation generated successfully!")
	fmt.Println("Main CLI documentation is available at docs/cli/README.md")
}

// updateCrossReferences updates all cross-references from hora.md to README.md
func updateCrossReferences(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file, "README.md") {
			continue
		}

		content, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		updatedContent := strings.ReplaceAll(string(content), "hora.md", "README.md")

		err = ioutil.WriteFile(file, []byte(updatedContent), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
