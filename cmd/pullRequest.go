// Copyright © 2018 Micke Lisinge <hi@micke.me>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
  "fmt"
  "regexp"

  "github.com/micke/devflow/git"
  "github.com/spf13/cobra"
)

var prTemplatePaths = []string{
  ".github/PULL_REQUEST_TEMPLATE.md",
}
var storyRegex, _ = regexp.Compile("/^([0-9]+)-(.*)/")
var fvRegex, _ = regexp.Compile("/^fv-(.*)/")

// pullRequestCmd represents the pullRequest command
var pullRequestCmd = &cobra.Command{
  Use:   "pullRequest",
  Short: "output a pr template for use with example git",
  Run: func(cmd *cobra.Command, args []string) {
    storyUrlBase := fmt.Sprintf("%s/entity", baseUrl)
    storyUrlPattern := fmt.Sprintf("%s/([0-9]+)", baseUrl)
    branch := git.GetCurrentBranch()

    fmt.Println(branch)
    fmt.Println(storyUrlBase)
    fmt.Println(storyUrlPattern)
  },
}

func init() {
  rootCmd.AddCommand(pullRequestCmd)
}
