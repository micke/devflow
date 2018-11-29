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
  "strings"

  "github.com/manifoldco/promptui"
  "github.com/micke/devflow/githelpers"
  "github.com/micke/devflow/targetprocess"
  "github.com/spf13/cobra"
)

var cleanerRegex, _ = regexp.Compile("[^\\w-]")

// checkoutCmd represents the branch command
var checkoutCmd = &cobra.Command{
  Use:   "checkout",
  Short: "Check out a branch for the story you have in progress",
  Run: func(cmd *cobra.Command, args []string) {
    tp := targetprocess.TargetProcess{
      AccessToken: accessToken,
      BaseUrl: baseUrl,
    }

    assignments := tp.GetAssignments(userId)

    branchName := branchify(assignments.SelectAssignment())

    prompt := promptui.Prompt{
      Label:     "Branch name",
      Default:   branchName,
      AllowEdit: true,
    }

    branchName, err := prompt.Run()

    if err == promptui.ErrInterrupt {
      return
    } else if err != nil {
      fmt.Printf("Prompt failed %v\n", err)
      return
    }

    githelpers.CheckoutBranch(branchName)
  },
}

func init() {
  rootCmd.AddCommand(checkoutCmd)

  // Here you will define your flags and configuration settings.

  // Cobra supports Persistent Flags which will work for this command
  // and all subcommands, e.g.:
  // checkoutCmd.PersistentFlags().String("foo", "", "A help for foo")

  // Cobra supports local flags which will only run when this command
  // is called directly, e.g.:
  // checkoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
