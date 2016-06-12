package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"fmt"

	"github.com/syndtr/goleveldb/leveldb/errors"
)

type gospector struct {
	dir        string
	config     *gospectorConf
	extToWords map[string]([]string)
}

func createGospector(dir string, config *gospectorConf) *gospector {
	g := &gospector{
		dir:        dir,
		config:     config,
		extToWords: make(map[string]([]string)),
	}

	for _, rule := range g.config.Rules {
		for _, ext := range rule.Extensions {
			if words, ok := g.extToWords[ext]; ok {
				g.extToWords[ext] = append(words, rule.Words...)
			} else {
				g.extToWords[ext] = rule.Words
			}
		}
	}

	return g
}

func (g *gospector) execute() []error {
	return g.executeDir(g.dir)
}

func (g *gospector) executeDir(dir string) []error {
	errArr := []error{}
	files, err := filepath.Glob(dir + "/*")
	if err != nil {
		return append(errArr, err)
	}

	for _, fileName := range files {
		file, _ := os.Stat(fileName)
		if file.IsDir() && g.shouldExecuteDir(fileName) {
			errRet := g.executeDir(fileName)
			errArr = append(errArr, errRet...)
		} else if g.shouldExecuteFile(fileName) {
			errRet := g.executeFile(fileName)
			errArr = append(errArr, errRet...)
		}
	}

	return errArr
}

func (g *gospector) executeFile(file string) []error {
	errArr := []error{}
	fileOpened, err := os.Open(file)
	if err != nil {
		return append(errArr, err)
	}

	fileExt := filepath.Ext(file)
	words := g.extToWords[fileExt]
	reader := bufio.NewReader(fileOpened)
	lineNumber := 0
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		lineNumber++
		for _, word := range words {
			if strings.Contains(line, word) {
				errorMessage := fmt.Sprintf("%s => %s found on line %d", file, word, lineNumber)
				errArr = append(errArr, errors.New(errorMessage))
			}
		}
	}
	return errArr
}

func (g *gospector) shouldExecuteFile(file string) bool {
	fileExt := filepath.Ext(file)
	for _, rule := range g.config.Rules {
		for _, ext := range rule.Extensions {
			if fileExt == ext {
				return true
			}
		}
	}
	return false
}

func (g *gospector) shouldExecuteDir(dir string) bool {
	if len(g.config.Subdirs) == 0 {
		return true
	}
	for _, subdir := range g.config.Subdirs {
		if strings.LastIndex(dir, g.dir+"/"+subdir) == 0 {
			return true
		}
	}
	return false
}
