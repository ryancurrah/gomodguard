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
	configFile  = ".gomodguard.yaml"
	logger      = log.New(os.Stderr, "", 0)
	lintErrorRC = 2
)

func main() {
	var (
		args       []string
		help       bool
		noTest     bool
		report     string
		reportFile string
		cwd, _     = os.Getwd()
		files      = []string{}
		finalFiles = []string{}
	)

	flag.BoolVar(&help, "h", false, "Show this help text")
	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&noTest, "n", false, "Don't lint test files")
	flag.BoolVar(&noTest, "no-test", false, "")
	flag.StringVar(&report, "r", "", "Report results to one of the following formats: checkstyle. A report file destination must also be specified")
	flag.StringVar(&report, "report", "", "")
	flag.StringVar(&reportFile, "f", "", "Report results to the specified file. A report type must also be specified")
	flag.StringVar(&reportFile, "file", "", "")
	flag.Parse()

	report = strings.TrimSpace(strings.ToLower(report))

	if help {
		showHelp()
		return
	}

	if report != "" && report != "checkstyle" {
		logger.Fatalf("error: invalid report type '%s'", report)
	}

	if report != "" && reportFile == "" {
		logger.Fatalf("error: a report file must be specified when a report is enabled")
	}

	if report == "" && reportFile != "" {
		logger.Fatalf("error: a report type must be specified when a report file is enabled")
	}

	args = flag.Args()
	if len(args) == 0 {
		args = []string{"./..."}
	}

	config, err := getConfig()
	if err != nil {
		logger.Fatalf("error: %s", err)
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

	processor, err := gomodguard.NewProcessor(*config, logger)
	if err != nil {
		logger.Fatalf("error: %s", err)
	}

	results := processor.ProcessFiles(finalFiles)

	if report == "checkstyle" {
		err := writeCheckstyle(reportFile, results)
		if err != nil {
			logger.Fatalf("error: %s", err)
		}
	}

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
		// Only append go files.
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

func writeCheckstyle(checkstyleFilePath string, results []gomodguard.Result) error {
	check := checkstyle.New()

	for i := range results {
		file := check.EnsureFile(results[i].FileName)
		file.AddError(checkstyle.NewError(results[i].LineNumber, 1, checkstyle.SeverityError, results[i].Reason, "gomodguard"))
	}

	checkstyleXML := fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n%s", check.String())

	err := ioutil.WriteFile(checkstyleFilePath, []byte(xmlfmt.FormatXML(checkstyleXML, "", "  ")), 0644)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func getConfig() (*gomodguard.Configuration, error) {
	config := gomodguard.Configuration{}

	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("unable to find home directory, %s", err)
	}

	cfgFile := ""
	homeDirCfgFile := filepath.Join(home, configFile)

	switch {
	case fileExists(configFile):
		cfgFile = configFile
	case fileExists(homeDirCfgFile):
		cfgFile = homeDirCfgFile
	default:
		return nil, fmt.Errorf("could not find config file in %s, %s", configFile, homeDirCfgFile)
	}

	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %s", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("could not parse config file: %s", err)
	}

	return &config, nil
}
