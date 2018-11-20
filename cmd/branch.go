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
  "os"
  "regexp"
  "strings"

  "github.com/micke/devflow/targetprocess"
  "github.com/spf13/cobra"
  "gopkg.in/src-d/go-git.v4"
  "gopkg.in/src-d/go-git.v4/plumbing"
)

var cleanerRegex, _ = regexp.Compile("[^\\w-]")

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
  Use:   "branch",
  Short: "Check out a branch for the story you have in progress",
  Run: func(cmd *cobra.Command, args []string) {
    resp := targetprocess.TargetProcess{
      AccessToken: accessToken,
      BaseUrl: baseUrl,
    }.GetAssignments(userId)

    branchName := branchify(resp.SelectAssignment())

    dir, err := os.Getwd()
    if err != nil {
      fmt.Printf("Error getting current dir: %s\n", err)
      os.Exit(1)
    }

    repo, err := git.PlainOpen(dir)
    if err != nil {
      fmt.Printf("Error opening repository: %s\n", err)
      os.Exit(1)
    }

    worktree, err := repo.Worktree()
    if err != nil {
      fmt.Printf("Error getting worktree: %s\n", err)
      os.Exit(1)
    }

    err = worktree.Checkout(&git.CheckoutOptions{
      Branch: plumbing.NewBranchReferenceName(branchName),
      Create: true,
    })

    if err != nil {
      fmt.Printf("Error creating branch: %s\n", err)
    }
    fmt.Println("SUCCESS!")
  },
}

func init() {
  rootCmd.AddCommand(branchCmd)

  // Here you will define your flags and configuration settings.

  // Cobra supports Persistent Flags which will work for this command
  // and all subcommands, e.g.:
  // branchCmd.PersistentFlags().String("foo", "", "A help for foo")

  // Cobra supports local flags which will only run when this command
  // is called directly, e.g.:
  // branchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func branchify(assignable targetprocess.TPAssignable) string{
  return fmt.Sprintf("%d-%s", assignable.Id, clean(assignable.Name))
}

func clean(str string) string{
  return cleanerRegex.ReplaceAllString(
    strings.ToLower(strings.Replace(str, " ", "-", -1)),
    "",
  )
}
