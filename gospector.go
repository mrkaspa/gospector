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
	extToTrailing map[string](bool)
}

func createGospector(dir string, config *gospectorConf) *gospector {
	g := &gospector{
		dir:        dir,
		config:     config,
		extToWords: make(map[string]([]string), len(config.Rules)),
		extToTrailing: make(map[string](bool), len(config.Rules)),
	}

	for _, rule := range g.config.Rules {
		for _, ext := range rule.Extensions {
			if words, ok := g.extToWords[ext]; ok {
				g.extToWords[ext] = append(words, rule.Words...)
			} else {
				g.extToWords[ext] = rule.Words
			}
			g.extToTrailing[ext] = rule.Trailing
		}
	}

	return g
}

func (g *gospector) execute() []error {
	return g.executeDir(g.dir, true)
}

func (g *gospector) executeDir(dir string, checkFiles bool) []error {
	errArr := []error{}
	files, err := filepath.Glob(dir + "/*")
	if err != nil {
		return append(errArr, err)
	}

	for _, fileName := range files {
		file, _ := os.Stat(fileName)
		if shouldExec, checkFiles := g.shouldExecuteDir(fileName); file.IsDir() &&  shouldExec {
			errRet := g.executeDir(fileName, checkFiles)
			errArr = append(errArr, errRet...)
		} else if checkFiles && g.shouldExecuteFile(fileName) {
			errRet := g.executeFile(fileName)
			errArr = append(errArr, errRet...)
		}
	}

	return errArr
}

func (g *gospector) executeFile(file string) []error {
	errArr := []error{}
	fileOpened, err := os.Open(file)
	defer fileOpened.Close()
	if err != nil {
		return append(errArr, err)
	}

	fileExt := filepath.Ext(file)
	words := g.extToWords[fileExt]
	trailing := g.extToTrailing[fileExt]
	reader := bufio.NewReader(fileOpened)
	lineNumber := 0
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		lineNumber++
		if trailing && strings.HasSuffix(line, " \n"){
			errorMessage := fmt.Sprintf("%s => trailing space found %d", file, lineNumber)
			errArr = append(errArr, errors.New(errorMessage))
		}
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

func (g *gospector) shouldExecuteDir(dir string) (bool, bool) {
	if len(g.config.Subdirs) == 0 && len(g.config.Excluded) == 0 {
		return true, true
	}

	for _, ex := range g.config.Excluded {
		fullSubdir := g.dir + "/" + ex
		if strings.LastIndex(dir, fullSubdir) == 0 {
			return false, false
		}
	}

	for _, subdir := range g.config.Subdirs {
		fullSubdir := g.dir + "/" + subdir
		if strings.LastIndex(dir, fullSubdir) == 0 {
			return true, true
		} else if strings.LastIndex(fullSubdir, dir) == 0 {
			return true, false
		}
	}
	return false, false
}
