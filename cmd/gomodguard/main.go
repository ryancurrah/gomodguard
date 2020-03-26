package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-xmlfmt/xmlfmt"
	"github.com/mitchellh/go-homedir"
	"github.com/phayes/checkstyle"
	"github.com/ryancurrah/gomodguard"
	"gopkg.in/yaml.v2"
)

var (
	configFile     = ".gomodguard.yaml"
	checkstyleFile = "gomodguard-checkstyle.xml"
	logger         = log.New(os.Stderr, "", 0)
	lintErrorRC    = 2
)

func main() {
	var (
		args       []string
		help       bool
		noTest     bool
		cwd, _     = os.Getwd()
		files      = []string{}
		finalFiles = []string{}
		config     = gomodguard.Configuration{}
	)

	home, err := homedir.Dir()
	if err != nil {
		logger.Fatalf("error: unable to find home directory, %s", err)
	}

	cfgFile := ""
	homeDirCfgFile := filepath.Join(home, configFile)

	switch {
	case fileExists(configFile):
		cfgFile = configFile
	case fileExists(homeDirCfgFile):
		cfgFile = homeDirCfgFile
	default:
		logger.Fatalf("error: could not find config file in %s, %s", configFile, homeDirCfgFile)
	}

	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	flag.BoolVar(&help, "h", false, "Show this help text")
	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&noTest, "n", false, "Don't lint test files")
	flag.BoolVar(&noTest, "no-test", false, "")
	flag.Parse()

	if help {
		showHelp()
		return
	}

	args = flag.Args()
	if len(args) == 0 {
		args = []string{"./..."}
	}

	for _, f := range args {
		if strings.HasSuffix(f, "/...") {
			dir, _ := filepath.Split(f)

			files = append(files, expandGoWildcard(dir)...)

			continue
		}

		if _, err := os.Stat(f); err == nil {
			files = append(files, f)
		}
	}

	// Use relative path to print shorter names, sort out test files if chosen.
	for _, f := range files {
		if noTest {
			if strings.HasSuffix(f, "_test.go") {
				continue
			}
		}

		if relativePath, err := filepath.Rel(cwd, f); err == nil {
			finalFiles = append(finalFiles, relativePath)

			continue
		}

		finalFiles = append(finalFiles, f)
	}

	processor := gomodguard.NewProcessorWithConfig(config, logger)
	results := processor.ProcessFiles(finalFiles)

	writeCheckstyle(results)

	for _, r := range results {
		fmt.Println(r.String())
	}

	if len(results) > 0 {
		os.Exit(lintErrorRC)
	}
}

func expandGoWildcard(root string) []string {
	foundFiles := []string{}

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// Only append go files
		if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		foundFiles = append(foundFiles, path)

		return nil
	})

	return foundFiles
}

func showHelp() {
	helpText := `Usage: gomodguard <file> [files...]
Also supports package syntax but will use it in relative path, i.e. ./pkg/...
Flags:`
	fmt.Println(helpText)
	flag.PrintDefaults()
}

func writeCheckstyle(results []gomodguard.Result) {
	check := checkstyle.New()

	for i := range results {
		file := check.EnsureFile(results[i].FileName)
		file.AddError(checkstyle.NewError(results[i].LineNumber, 1, checkstyle.SeverityError, "import", results[i].Reason))
	}

	err := ioutil.WriteFile(checkstyleFile, []byte(xmlfmt.FormatXML(check.String(), "", "  ")), 0600)
	if err != nil {
		logger.Fatalf("error: %s", err)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
