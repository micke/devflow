package githelpers

import (
  "fmt"
  "os"
  "regexp"

  "github.com/libgit2/git2go"
)

func GetCurrentBranch() string {
  repo := getRepo()

  head, err := repo.Head()
  checkFail("Failed to get head", err)

  branchName, err := head.Branch().Name()
  checkFail("Failed to get branch name", err)

  return branchName
}

func ExistingBranchForPattern(pattern *regexp.Regexp) *string {
  var matchingBranch *string
  repo := getRepo()

  branchIterator, err := repo.NewBranchIterator(git.BranchAll)
  checkFail("Failed to get branch iterator", err)

  f := func(branch *git.Branch, t git.BranchType) error {
    branchName, err := branch.Name()
    if err == nil && pattern.MatchString(branchName) {
      matchingBranch = &branchName
    }

    return nil
  }

  branchIterator.ForEach(f)

  return matchingBranch
}

func CheckoutBranch(branchName string) error {
  var commit *git.Commit
  repo := getRepo()
  checkoutOpts := &git.CheckoutOpts{
    Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing |
							git.CheckoutAllowConflicts | git.CheckoutUseTheirs,
  }

  branch, err := repo.LookupBranch(branchName, git.BranchLocal)
  if err == nil {
      fmt.Printf("Switched to branch '%s'\n", branchName)
  } else {
    remoteBranch, err := repo.LookupBranch("origin/" + branchName, git.BranchRemote)
    if err == nil {
      // Find commit for remote branch
      commit, err = repo.LookupCommit(remoteBranch.Target())
      checkFail("Failed to find remote commit", err)
      fmt.Printf("Found remote branch '%s'\n", branchName)
    } else {
      // If remote branch does not exist we create a new local branch
      head, err := repo.Head()
      checkFail("Failed to get head", err)

      commit, err = repo.LookupCommit(head.Target())
      checkFail("Failed to find local commit", err)
      fmt.Printf("Creating new branch '%s'\n", branchName)
    }

    branch, err = repo.CreateBranch(branchName, commit, false)
    checkFail("Failed to create local branch", err)

    if remoteBranch != nil {
      // Set upstream of local branch to remote branch
      err = branch.SetUpstream("origin/" + branchName)
      checkFail("Failed to set upstream to origin", err)
    }
  }

  commit, err = repo.LookupCommit(branch.Target())
  checkFail("Failed to find commit for branch", err)

  tree, err := repo.LookupTree(commit.TreeId())
  checkFail("Failed to lookup for tree", err)

  err = repo.CheckoutTree(tree, checkoutOpts)
  checkFail("Failed to checkout tree", err)

  repo.SetHead("refs/heads/" + branchName)
  return nil
}

func getWd() string{
  dir, err := os.Getwd()
  checkFail("Error getting current dir", err)

  return dir
}

func getRepo() *git.Repository{
  repo, err := git.OpenRepositoryExtended(getWd(), 0, "/")
  checkFail("Error opening repository", err)

  return repo
}

func checkFail(message string, err error) {
  if err != nil {
    fmt.Printf("%s: %s\n", message, err)
    os.Exit(1)
  }
}
