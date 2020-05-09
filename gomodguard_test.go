package gomodguard

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

var (
	recommendedModule = "github.com/ryancurrah/gomodguard"
	invalidModule     = "github.com/ryancurrah/gomodbump"
	blockedModule     = "github.com/someblocked/module"
	allowedModule     = "github.com/someallowed/module"
	allowedDomain     = "golang.org"
	recommendations   = Recommendations{Recommendations: []string{recommendedModule}, Reason: "test reason"}
	blockedModules    = BlockedModules{{blockedModule: recommendations}}
	allowed           = Allowed{Modules: []string{allowedModule}, Domains: []string{allowedDomain}}
)

func TestGomodguardIsRecommended(t *testing.T) {
	isRecommended := recommendations.IsRecommended(recommendedModule)
	if !isRecommended {
		t.Errorf("%s should have been recommended", recommendedModule)
	}

	isRecommended = recommendations.IsRecommended(invalidModule)
	if isRecommended {
		t.Errorf("%s should NOT have been recommended", invalidModule)
	}
}

func TestGomodguardRecommendationsString(t *testing.T) {
	recommendationsMsg := recommendations.String()
	if recommendationsMsg == "" {
		t.Error("recommendations string message should not be empty")
	}

	if strings.Contains(recommendationsMsg, "modules") {
		t.Errorf("recommendations string message should be singular: %s", recommendationsMsg)
	}

	multipleRecommendations := recommendations
	multipleRecommendations.Recommendations = append(multipleRecommendations.Recommendations, "github.com/some/thing", "github.com/some/otherthing")

	multipleRecommendationsMsg := multipleRecommendations.String()
	if !strings.Contains(multipleRecommendationsMsg, "modules") {
		t.Errorf("recommendations string message should be plural: %s", recommendationsMsg)
	}

	emptyRecommendations := Recommendations{}

	recommendationsMsg = emptyRecommendations.String()
	if recommendationsMsg != "" {
		t.Error("recommendations string message should be empty")
	}
}

func TestGomodguardHasRecommendations(t *testing.T) {
	hasRecommendations := recommendations.HasRecommendations()
	if !hasRecommendations {
		t.Error("should have recommendations when more than one recommended module in list")
	}

	emptyRecommendations := Recommendations{}

	hasRecommendations = emptyRecommendations.HasRecommendations()
	if hasRecommendations {
		t.Error("should not have recommendations when no recommended modules in list")
	}
}

func TestGomodguardBlockedModulesGet(t *testing.T) {
	blockedModulesList := blockedModules.Get()
	if len(blockedModulesList) == 0 {
		t.Error("blocked modules list should not be empty")
	}
}

func TestGomodguardRecommendedModules(t *testing.T) {
	recommendedModules := blockedModules.RecommendedModules(blockedModule)
	if len(recommendedModules.Recommendations) == 0 {
		t.Error("recommended modules list should not be empty")
	}

	recommendedModules = blockedModules.RecommendedModules(invalidModule)
	if recommendedModules != nil {
		t.Error("recommended modules should be nil when no recommendations for blocked module")
	}
}

func TestGomodguardIsBlockedPackage(t *testing.T) {
	blockedPkg := fmt.Sprintf("%s/util", blockedModule)

	isBlockedPackage := blockedModules.IsBlockedPackage(blockedPkg)
	if !isBlockedPackage {
		t.Errorf("package %s should be blocked when the module is in the blocked list", blockedPkg)
	}

	allowedPkg := "github.com/someallowed/module"

	isBlockedPackage = blockedModules.IsBlockedPackage(allowedPkg)
	if isBlockedPackage {
		t.Errorf("package %s should NOT be blocked when the module is NOT in the blocked list", allowedPkg)
	}
}

func TestGomodguardIsBlockedModule(t *testing.T) {
	isBlockedPackage := blockedModules.IsBlockedModule(blockedModule)
	if !isBlockedPackage {
		t.Errorf("module %s should be blocked when the module is in the blocked list", blockedModule)
	}

	isBlockedPackage = blockedModules.IsBlockedModule(allowedModule)
	if isBlockedPackage {
		t.Errorf("module %s should NOT be blocked when the module is NOT in the blocked list", allowedModule)
	}
}

func TestGomodguardIsAllowedModule(t *testing.T) {
	isAllowedModule := allowed.IsAllowedModule(allowedModule)
	if !isAllowedModule {
		t.Errorf("module %s should be allowed when the module is in the allowed modules list", allowedModule)
	}

	isAllowedModule = allowed.IsAllowedModule(blockedModule)
	if isAllowedModule {
		t.Errorf("module %s should NOT be allowed when the module is NOT in the allowed modules list", blockedModule)
	}
}

func TestGomodguardIsAllowedModuleDomain(t *testing.T) {
	isAllowedModuleDomain := allowed.IsAllowedModuleDomain(allowedDomain)
	if !isAllowedModuleDomain {
		t.Errorf("module domain %s should be allowed when the module domain is in the allowed domains list", allowedDomain)
	}

	isAllowedModuleDomain = allowed.IsAllowedModuleDomain("blocked.domain")
	if isAllowedModuleDomain {
		t.Errorf("module domain %s should NOT be allowed when the module domain is NOT in the allowed domains list", "blocked.domain")
	}
}

func TestGomodguardProcessFilesWithAllowed(t *testing.T) {
	config, logger, cwd, err := setup()
	if err != nil {
		t.Error(err)
	}

	// Test that setting skip files to true does NOT return test files
	filteredFilesNoTests := GetFilteredFiles(cwd, true, []string{"./..."})

	testFileFound := false

	for _, finalFile := range filteredFilesNoTests {
		if strings.HasSuffix(finalFile, "_test.go") {
			testFileFound = true
		}
	}

	if testFileFound {
		t.Errorf("should NOT have returned files found that end with _test.go")
	}

	// Test that setting skip files to false DOES return test files
	filteredFiles := GetFilteredFiles(cwd, false, []string{"./..."})
	if len(filteredFiles) == 0 {
		t.Errorf("should have found a file to lint")
	}

	testFileFound = false

	for _, finalFile := range filteredFiles {
		if strings.HasSuffix(finalFile, "_test.go") {
			testFileFound = true
		}
	}

	if !testFileFound {
		t.Errorf("should have been able to find files that end with _test.go")
	}

	processor, err := NewProcessor(*config, logger)
	if err != nil {
		t.Errorf("should have been able to init a new processor without an error")
	}

	results := processor.ProcessFiles(filteredFiles)
	if len(results) > 0 {
		t.Errorf("should not have found a lint error")
	}
}

func TestGomodguardProcessFilesAllAllowed(t *testing.T) {
	config, logger, cwd, err := setup()
	if err != nil {
		t.Error(err)
	}

	config.Allowed.Modules = []string{}
	config.Allowed.Domains = []string{}
	config.Blocked.Modules = BlockedModules{}

	filteredFiles := GetFilteredFiles(cwd, false, []string{"./..."})
	if len(filteredFiles) == 0 {
		t.Errorf("should have found a file to lint")
	}

	processor, err := NewProcessor(*config, logger)
	if err != nil {
		t.Errorf("should have been able to init a new processor without an error")
	}

	results := processor.ProcessFiles(filteredFiles)
	if len(results) > 0 {
		t.Errorf("should not have found a lint error")
	}
}

func TestGomodguardProcessFilesWithBlockedModules(t *testing.T) {
	config, logger, cwd, err := setup()
	if err != nil {
		t.Error(err)
	}

	config.Allowed.Modules = []string{"github.com/someallowed/module"}
	config.Allowed.Domains = []string{}
	config.Blocked.Modules = BlockedModules{
		BlockedModule{"golang.org/x/mod": Recommendations{}},
		BlockedModule{"gopkg.in/yaml.v2": Recommendations{Recommendations: []string{"github.com/something/else"}, Reason: "test reason"}},
	}

	filteredFiles := GetFilteredFiles(cwd, false, []string{"./..."})
	if len(filteredFiles) == 0 {
		t.Errorf("should have found a file to lint")
	}

	processor, err := NewProcessor(*config, logger)
	if err != nil {
		t.Errorf("should have been able to init a new processor without an error")
	}

	results := processor.ProcessFiles(filteredFiles)
	if len(results) == 0 {
		t.Errorf("should have found at least one lint error")
	}
}

func setup() (*Configuration, *log.Logger, string, error) {
	config, err := GetConfig(".gomodguard.yaml")
	if err != nil {
		return nil, nil, "", err
	}

	logger := log.New(os.Stderr, "", 0)

	cwd, _ := os.Getwd()

	return config, logger, cwd, nil
}
