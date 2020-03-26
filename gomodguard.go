package gomodguard

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/mod/modfile"
)

var (
	blockedReason = "import of package `%s` is blocked because the module is not in the allowed modules list."
	goModFile     = "go.mod"
)

// Replacement is a list of blocked modules with a replacement module and reason why it should be replaced.
type Replacement struct {
	Modules     []string `yaml:"modules"`
	Replacement string   `yaml:"replacement"`
	Reason      string   `yaml:"reason"`
}

// String returns the replacement module and reason message.
func (r *Replacement) String() string {
	return fmt.Sprintf("`%s` should be used instead. reason: %s", r.Replacement, r.Reason)
}

// HasReplacement returns true if the blocked package has a replacement module.
func (r *Replacement) HasReplacement(pkg string) bool {
	for i := range r.Modules {
		if strings.HasPrefix(strings.ToLower(pkg), strings.ToLower(r.Modules[i])) {
			return true
		}
	}

	return false
}

// Replacements a list of replacement modules.
type Replacements []Replacement

// Get will return a replacement for the package provided. If there is no replacement nil will be returned.
func (r Replacements) Get(pkg string) *Replacement {
	for i := range r {
		if r[i].HasReplacement(pkg) {
			return &r[i]
		}
	}

	return nil
}

// Allow packages and domains.
type Allow struct {
	Modules []string `yaml:"modules"`
	Domains []string `yaml:"domains"`
}

// Configuration of gomodguard.
type Configuration struct {
	Allow        Allow        `yaml:"allow"`
	Replacements Replacements `yaml:"replacements"`
}

// Result represents the result of one error.
type Result struct {
	FileName   string
	LineNumber int
	Position   token.Position
	Reason     string
}

// String returns the filename, line number and reason of a Result.
func (r *Result) String() string {
	return fmt.Sprintf("%s:%d: %s", r.FileName, r.LineNumber, r.Reason)
}

// Processor processes Go files.
type Processor struct {
	config         Configuration
	logger         *log.Logger
	modfile        *modfile.File
	blockedModules []string
	result         []Result
}

// NewProcessorWithConfig will create a Processor to lint blocked packages.
func NewProcessorWithConfig(config Configuration, logger *log.Logger) *Processor {
	moddata, err := ioutil.ReadFile(goModFile)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	mfile, err := modfile.Parse(goModFile, moddata, nil)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	logger.Printf("info: allowed modules, %+v", config.Allow.Modules)
	logger.Printf("info: allowed module domains, %+v", config.Allow.Domains)

	p := &Processor{
		config:  config,
		logger:  logger,
		modfile: mfile,
		result:  []Result{},
	}

	p.setBlockedModules()

	return p
}

// ProcessFiles takes a string slice with file names (full paths) and lints them.
func (p *Processor) ProcessFiles(filenames []string) []Result {
	p.logger.Printf("info: go.mod file has '%d' blocked module(s), %+v", len(p.blockedModules), p.blockedModules)

	if len(p.blockedModules) == 0 {
		return p.result
	}

	for _, filename := range filenames {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			p.logger.Fatalf("error: %v", err)
		}

		p.process(filename, data)
	}

	return p.result
}

// process file imports and add lint error if blocked package is imported.
func (p *Processor) process(filename string, data []byte) {
	fileSet := token.NewFileSet()

	file, err := parser.ParseFile(fileSet, filename, data, parser.ParseComments)
	if err != nil {
		p.result = append(p.result, Result{
			FileName:   filename,
			LineNumber: 0,
			Reason:     fmt.Sprintf("invalid syntax, file cannot be linted (%s)", err.Error()),
		})

		return
	}

	imports := file.Imports
	for i := range imports {
		pkg := strings.Trim(imports[i].Path.Value, "\"")
		if p.isBlockedPackage(pkg) {
			reason := fmt.Sprintf(blockedReason, pkg)
			replacement := p.config.Replacements.Get(pkg)

			if replacement != nil {
				reason += fmt.Sprintf(" %s", replacement.String())
			}

			p.addError(fileSet, imports[i].Pos(), reason)
		}
	}
}

// Add an error for the file and line number for the current token.Pos with the given reason.
func (p *Processor) addError(fileset *token.FileSet, pos token.Pos, reason string) {
	position := fileset.Position(pos)

	p.result = append(p.result, Result{
		FileName:   position.Filename,
		LineNumber: position.Line,
		Position:   position,
		Reason:     reason,
	})
}

func (p *Processor) setBlockedModules() {
	blockedModules := make([]string, 0, len(p.modfile.Require))
	require := p.modfile.Require

	for i := range require {
		if !require[i].Indirect {
			if p.isAllowedModuleDomain(require[i].Mod.Path) {
				continue
			}

			if p.isAllowedModule(require[i].Mod.Path) {
				continue
			}

			blockedModules = append(blockedModules, require[i].Mod.Path)
		}
	}

	p.blockedModules = blockedModules
}

func (p *Processor) isAllowedModuleDomain(module string) bool {
	domains := p.config.Allow.Domains
	for n := range domains {
		if strings.HasPrefix(strings.ToLower(module), strings.ToLower(domains[n])) {
			return true
		}
	}

	return false
}

func (p *Processor) isAllowedModule(module string) bool {
	packages := p.config.Allow.Modules
	for n := range packages {
		if strings.EqualFold(module, packages[n]) {
			return true
		}
	}

	return false
}

func (p *Processor) isBlockedPackage(pkg string) bool {
	blockedModules := p.blockedModules
	for i := range blockedModules {
		if strings.HasPrefix(strings.ToLower(pkg), strings.ToLower(blockedModules[i])) {
			return true
		}
	}

	return false
}
