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
var storyRegex, _ = regexp.Compile("^([0-9]+)-(.*)")
var fvRegex, _ = regexp.Compile("^fv-(.*)")
var tpHeaderRegex, _ = regexp.Compile("\n*#+\\s*TP")

// pullRequestCmd represents the pullRequest command
var pullRequestCmd = &cobra.Command{
	Use:     "pullRequest",
	Aliases: []string{"pr", "pull-request", "pullrequest"},
	Short:   "output a pr template for use with example git (Aliased as pr, pull-request and pullrequest)",
	Run: func(cmd *cobra.Command, args []string) {
		storyUrlBase := fmt.Sprintf("%s/entity", baseUrl)
		storyUrlPattern, _ := regexp.Compile(fmt.Sprintf("%s/([0-9]+)", storyUrlBase))
		branch := githelpers.GetCurrentBranch()
		var storyId string
		var storyUrl string
		var title string

		// Match branch to find the type of PR to create
		if storyRegex.MatchString(branch) {
			matches := storyRegex.FindStringSubmatch(branch)
			storyId = matches[1]
			subject := matches[2]
			storyUrl = fmt.Sprintf("%s/%s", storyUrlBase, storyId)
			title = fmt.Sprintf("[#%s] %s", storyId, titleize(subject))
		} else if fvRegex.MatchString(branch) {
			matches := fvRegex.FindStringSubmatch(branch)
			subject := matches[1]
			storyUrl = ""
			title = fmt.Sprintf("[FV] %s", titleize(subject))
		} else {
			storyUrl = ""
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

		if storyUrl == "" {
			// Remove any storyUrl and TP header from the template
			template = storyUrlPattern.ReplaceAllString(template, "")
			template = tpHeaderRegex.ReplaceAllString(template, "")
		} else {
			if storyUrlPattern.MatchString(template) {
				// Replace placeholder URL if it exists
				template = storyUrlPattern.ReplaceAllString(template, storyUrl)
			} else {
				// Otherwise just append it after some artistic newlines
				template = fmt.Sprintf("%s\n\n%s\n", template, storyUrl)
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
