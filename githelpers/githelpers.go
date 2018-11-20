package githelpers

import (
  "fmt"
  "os"
  "os/exec"
  "strings"
  "log"

  "github.com/libgit2/git2go"
)

func GetCurrentBranch() string{
  branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
  branchArr, err := branchCmd.CombinedOutput()
  if err != nil {
    fmt.Printf("Error finding branch: %s %s\n", err, branchArr)
    os.Exit(1)
  }
  return strings.Trim(string(branchArr), "\n")
}

func CheckoutBranch(repo *git.Repository, branchName string) error {
  checkoutOpts := &git.CheckoutOpts{
    Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutAllowConflicts | git.CheckoutUseTheirs,
  }

  branch, err := repo.LookupBranch(branchName, git.BranchAll)
  if err != nil {
    // If branch does not exist we create it!
    head, err := repo.Head()
    if err != nil {
      log.Print("Failed to get head: " + branchName)
      return err
    }

    commit, err := repo.LookupCommit(head.Target())
    if err != nil {
      log.Print("Failed to find branch commit: " + branchName)
      return err
    }
    defer commit.Free()

    branch, err = repo.CreateBranch(branchName, commit, false)
    if err != nil {
      log.Print("Failed to create local branch: " + branchName)
      return err
    }
  }
  defer branch.Free()

  commit, err := repo.LookupCommit(branch.Target())
  if err != nil {
    log.Print("Failed to find branch commit: " + branchName)
    return err
  }
  defer commit.Free()

  tree, err := repo.LookupTree(commit.TreeId())
  if err != nil {
    log.Print("Failed to lookup for tree " + branchName)
    return err
  }
  defer tree.Free()

  err = repo.CheckoutTree(tree, checkoutOpts)
  if err != nil {
    log.Print("Failed to checkout tree " + branchName)
    return err
  }

  repo.SetHead("refs/heads/" + branchName)
  return nil
}
