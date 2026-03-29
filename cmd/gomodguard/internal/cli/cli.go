package cli

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/phayes/checkstyle"
	"go.yaml.in/yaml/v4"

	"github.com/ryancurrah/gomodguard/v2"
)

const (
	errFindingHomedir    = "unable to find home directory, %w"
	errReadingConfigFile = "could not read config file: %w"
	errParsingConfigFile = "could not parse config file: %w"
)

var (
	configFile           = ".gomodguard.yaml"
	logger               = log.New(os.Stderr, "", 0)
	errFindingConfigFile = errors.New("could not find config file")
)

// Run the gomodguard linter. Returns the exit code to use.
//
//nolint:funlen
func Run() int {
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		return MigrateConfig(configFile)
	}

	var (
		args           []string
		help           bool
		noTest         bool
		report         string
		reportFile     string
		issuesExitCode int
		printVersion   bool
		cwd, _         = os.Getwd()
	)

	flag.BoolVar(&printVersion, "version", false, "Print the version")
	flag.BoolVar(&help, "h", false, "Show this help text")
	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&noTest, "n", false, "Don't lint test files")
	flag.BoolVar(&noTest, "no-test", false, "")
	flag.StringVar(&report, "r", "", "Report results to one of the following formats: checkstyle. "+
		"A report file destination must also be specified")
	flag.StringVar(&report, "report", "", "")
	flag.StringVar(&reportFile, "f", "", "Report results to the specified file. A report type must also be specified")
	flag.StringVar(&reportFile, "file", "", "")
	flag.IntVar(&issuesExitCode, "i", 2, "Exit code when issues were found")
	flag.IntVar(&issuesExitCode, "issues-exit-code", 2, "")
	flag.Parse()

	if printVersion {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			fmt.Println("Failed to read build info")

			return 1
		}

		fmt.Printf("gomodguard version: %s\n", info.Main.Version)

		return 0
	}

	report = strings.TrimSpace(strings.ToLower(report))

	if help {
		showHelp()
		return 0
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

	config, err := getConfig(configFile)
	if err != nil {
		logger.Fatalf("error: %s", err)
	}

	filteredFiles := gomodguard.Find(cwd, noTest, args)

	processor, err := gomodguard.NewProcessor(config)
	if err != nil {
		logger.Fatalf("error: %s", err)
	}

	allowedModuleNames := make([]string, len(config.Allowed))
	for i, m := range config.Allowed {
		allowedModuleNames[i] = m.Module
	}

	blockedModuleNames := make([]string, len(config.Blocked))
	for i, m := range config.Blocked {
		blockedModuleNames[i] = m.Module
	}

	slices.Sort(allowedModuleNames)
	slices.Sort(blockedModuleNames)

	logger.Printf("info: allowed modules, %+v", allowedModuleNames)
	logger.Printf("info: blocked modules, %+v", blockedModuleNames)

	results := processor.ProcessFiles(filteredFiles)

	if report == "checkstyle" {
		err := WriteCheckstyle(reportFile, results)
		if err != nil {
			logger.Fatalf("error: %s", err)
		}
	}

	for _, r := range results {
		fmt.Println(r.String())
	}

	if len(results) > 0 {
		return issuesExitCode
	}

	return 0
}

// getConfig from YAML file.
func getConfig(configFile string) (*gomodguard.Configuration, error) {
	config := gomodguard.Configuration{}

	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf(errFindingHomedir, err)
	}

	homeDirCfgFile := filepath.Join(home, configFile)

	var cfgFile string

	switch {
	case fileExists(configFile):
		cfgFile = configFile
	case fileExists(homeDirCfgFile):
		cfgFile = homeDirCfgFile
	default:
		return nil, fmt.Errorf("%w: %s %s", errFindingConfigFile, configFile, homeDirCfgFile)
	}

	data, err := os.ReadFile(filepath.Clean(cfgFile))
	if err != nil {
		return nil, fmt.Errorf(errReadingConfigFile, err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf(errParsingConfigFile, err)
	}

	return &config, nil
}

// showHelp text for command line.
func showHelp() {
	helpText := `Usage: gomodguard <file> [files...]
Also supports package syntax but will use it in relative path, i.e. ./pkg/...

Commands:
  (default)  Lint Go module dependencies using the configuration file
  migrate    Convert a v1 .gomodguard.yaml config file to v2 format and print to stdout

Flags:`
	fmt.Println(helpText)
	flag.PrintDefaults()
}

// WriteCheckstyle takes the results and writes them to a checkstyle formated file.
func WriteCheckstyle(checkstyleFilePath string, results []gomodguard.Issue) error {
	check := checkstyle.New()

	for i := range results {
		file := check.EnsureFile(results[i].FileName)
		file.AddError(
			checkstyle.NewError(
				results[i].LineNumber, 1,
				checkstyle.SeverityError,
				results[i].Reason,
				"gomodguard",
			),
		)
	}

	body, err := xml.MarshalIndent(check, "", "  ")
	if err != nil {
		return err
	}

	header := []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	checkstyleXML := slices.Concat([]byte{'\n'}, header, []byte{'\n'}, body)

	err = os.WriteFile(checkstyleFilePath, checkstyleXML, 0644) //nolint:gosec
	if err != nil {
		return err
	}

	return nil
}

// fileExists returns true if the file path provided exists.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
