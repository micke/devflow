package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/micke/devflow/githelpers"
	"github.com/spf13/cobra"
)

var prTemplatePaths = []string{
	".github/PULL_REQUEST_TEMPLATE.md",
}
var storyRegex = regexp.MustCompile("^tp([0-9]+)-(.*)")
var fvRegex = regexp.MustCompile("^fv-(.*)")
var tpHeaderRegex = regexp.MustCompile("\n*#+\\s*TP")

// pullRequestCmd represents the pullRequest command
var pullRequestCmd = &cobra.Command{
	Use:     "pullRequest",
	Aliases: []string{"pr", "pull-request", "pullrequest"},
	Short:   "output a pr template for use with example git (Aliased as pr, pull-request and pullrequest)",
	Run: func(cmd *cobra.Command, args []string) {
		storyURLBase := fmt.Sprintf("%s/entity", baseURL)
		storyURLPattern := regexp.MustCompile(fmt.Sprintf("%s/([0-9]+)", storyURLBase))
		branch := githelpers.GetCurrentBranch()
		var storyID string
		var storyURL string
		var title string

		// Match branch to find the type of PR to create
		if storyRegex.MatchString(branch) {
			matches := storyRegex.FindStringSubmatch(branch)
			storyID = matches[1]
			subject := matches[2]
			storyURL = fmt.Sprintf("%s/%s", storyURLBase, storyID)
			title = fmt.Sprintf("[#%s] %s", storyID, titleize(subject))
		} else if fvRegex.MatchString(branch) {
			matches := fvRegex.FindStringSubmatch(branch)
			subject := matches[1]
			storyURL = ""
			title = fmt.Sprintf("[FV] %s", titleize(subject))
		} else {
			storyURL = ""
			title = titleize(branch)
		}

		// Find and load template file
		var template string
		for _, file := range prTemplatePaths {
			if _, err := os.Stat(file); !os.IsNotExist(err) {
				templateBytes, _ := ioutil.ReadFile(file)
				template = string(templateBytes)
				break
			}
		}

		if storyURL == "" {
			// Remove any storyURL and TP header from the template
			template = storyURLPattern.ReplaceAllString(template, "")
			template = tpHeaderRegex.ReplaceAllString(template, "")
		} else {
			if storyURLPattern.MatchString(template) {
				// Replace placeholder URL if it exists
				template = storyURLPattern.ReplaceAllString(template, storyURL)
			} else {
				// Otherwise just append it after some artistic newlines
				template = fmt.Sprintf("%s\n\n%s\n", template, storyURL)
			}
		}

		fmt.Printf("%s\n\n", title)
		fmt.Printf(template)
	},
}

func init() {
	rootCmd.AddCommand(pullRequestCmd)
}

func titleize(str string) string {
	splitted := strings.Split(str, "-")
	splitted[0] = strings.Title(splitted[0])
	return strings.Join(splitted, " ")
}
