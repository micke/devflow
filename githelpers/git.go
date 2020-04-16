package git

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func GetCurrentBranch() string {
	output, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()

	if err != nil {
		fmt.Println("Getting current branch failed:", string(output))
		os.Exit(1)
	}

	return string(output[:len(output)-1])
}

func BranchExistsForPattern(pattern *regexp.Regexp) bool {
	for _, branch := range branches() {
		if pattern.MatchString(branch) {
			return true
		}
	}

	return false
}

func CheckoutBranch(branchName string) error {
	output, err := exec.Command("git", "checkout", "-b", branchName).CombinedOutput()

	if err != nil {
		fmt.Println("Getting current branch failed:", string(output))
		os.Exit(1)
	}

	return nil
}

func branches() []string {
	output, err := exec.Command("git", "branch", "--all", "--format=%(refname:lstrip=2)", "--no-merged").CombinedOutput()

	if err != nil {
		fmt.Println("Listing existing branches failed:", string(output))
		os.Exit(1)
	} else if len(output) == 0 {
		return []string{}
	}

	return strings.Split(string(output[:len(output)-1]), "\n")
}
