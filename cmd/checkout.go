package cmd

import (
	"fmt"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/micke/devflow/githelpers"
	"github.com/micke/devflow/targetprocess"
	"github.com/spf13/cobra"
)

var cleanerRegex = regexp.MustCompile("[^\\w-]")

// checkoutCmd represents the branch command
var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Check out a branch for the story you have in progress",
	Run: func(cmd *cobra.Command, args []string) {
		tp := targetprocess.TargetProcess{
			AccessToken: accessToken,
			BaseURL:     baseURL,
		}

		assignments := tp.GetAssignments(userID)

		var branchName string
		assignment := assignments.SelectAssignment()
		branchPrefix := fmt.Sprintf("tp%d-", assignment.ID)
		branchPattern := regexp.MustCompile("^" + branchPrefix)
		existingBranch := githelpers.ExistingBranchForPattern(branchPattern)

		if existingBranch != nil {
			fmt.Printf("Found existing branch matching story: %s\n", *existingBranch)
			prompt := promptui.Prompt{
				Label:     "Do you want to switch to this branch?",
				IsConfirm: true,
				Default:   "y",
			}

			confirmed, err := prompt.Run()

			if err == promptui.ErrInterrupt {
				return
			}

			if confirmed != "n" {
				branchName = *existingBranch
			} else {
				existingBranch = nil
			}
		}

		if existingBranch == nil {
			prompt := promptui.Prompt{
				Label:     "Branch name",
				Default:   branchPrefix,
				AllowEdit: true,
			}

			var err error
			branchName, err = prompt.Run()

			if err == promptui.ErrInterrupt {
				return
			} else if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
		}

		githelpers.CheckoutBranch(branchName)
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)

	checkoutCmd.Flags().BoolP("master", "m", false, "Create this feature branch from master")
}
