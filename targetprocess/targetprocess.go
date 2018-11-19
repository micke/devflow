package targetprocess

import (
  "fmt"
  "os"

  "gopkg.in/resty.v1"
)

type TargetProcess struct {
  AccessToken string
  BaseUrl string
}

type TPAssignments struct {
  Items []TPAssignment
}

type TPAssignment struct {
  TPAssignable `json:"Assignable"`
}

type TPAssignable struct {
  Id int
  Name string
}

func (tp TargetProcess) client() *resty.Client{
  return resty.New().
    SetQueryParams(map[string]string{
      "format": "json",
      "access_token": tp.AccessToken,
    }).
    SetHeader("Content-Type", "application/json").
    SetHeader("Accept", "application/json").
    SetHostURL(tp.BaseUrl)
}

func (tp TargetProcess) request() *resty.Request{
  return tp.client().R()
}

func (tp TargetProcess) GetAssignments(userId string) *TPAssignments{
  resp, err := tp.request().
    SetQueryParams(map[string]string{
      "where": fmt.Sprintf("(GeneralUser.Id eq %s)and(Assignable.EntityState.Name eq 'In Progress')", userId),
      "include": "[Assignable[Id,Name]]",
    }).
    SetResult(TPAssignments{}).
    Get("/api/v1/assignments")

  if err != nil {
    fmt.Printf("error getting assignments, %s\n", err)
    os.Exit(1)
  }

  return resp.Result().(*TPAssignments)
}
