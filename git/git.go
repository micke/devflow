package git

import (
  "fmt"
  "os"
  "os/exec"
  "strings"
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
