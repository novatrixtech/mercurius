package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

//copy flies from src to dest
func mustCopyFile(destFilename, srcFilename string) {
	destFile, err := os.Create(destFilename)
	if err != nil {
		panic(err)
	}

	srcFile, err := os.Open(srcFilename)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		panic(err)
	}

	err = destFile.Close()
	if err != nil {
		panic(err)
	}

	err = srcFile.Close()
	if err != nil {
		panic(err)
	}
}

// parse template and change to new app data
func mustRenderTemplate(destPath, srcPath string, data map[string]interface{}) {
	tmpl, err := template.ParseFiles(srcPath)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(destPath)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}
}

// copyDir copies a directory tree over to a new directory
// Also, dot files and dot directories are skipped.
func mustCopyDir(destDir, srcDir string, data map[string]interface{}) error {
	return filepath.Walk(srcDir, func(srcPath string, info os.FileInfo, err error) error {
		// Get the relative path from the source base, and the corresponding path in
		// the dest directory.
		relSrcPath := strings.TrimLeft(srcPath[len(srcDir):], string(os.PathSeparator))
		destPath := path.Join(destDir, relSrcPath)

		// Skip dot files and dot directories.
		if strings.HasPrefix(relSrcPath, ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Create a subdirectory if necessary.
		if info.IsDir() {
			err := os.MkdirAll(path.Join(destDir, relSrcPath), 0777)
			if !os.IsExist(err) {
				if err != nil {
					panic(err)
				}
			}
			return nil
		}

		mustRenderTemplate(destPath, srcPath, data)
		return nil
	})
}

// help to manipulate terminal input and output
func terminal(question, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(fmt.Sprintf("%s (%s):", question, defaultValue))
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	if text == "" {
		return defaultValue
	}
	return text
}
