package targetprocess

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
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

func (assignments TPAssignments) SelectAssignment() TPAssignable{
	if (len(assignments.Items) == 1) {
		return assignments.Items[0].TPAssignable
	}

	templates := &promptui.SelectTemplates{
		Label:		"{{ . }}?",
		Active:		"{{ .Id | cyan }} {{ .Name | yellow }}",
		Inactive: "{{ .Id | cyan }} {{ .Name }}",
		Selected: "{{ .Id | green }} {{ .Name | green }}",
	}

	searcher := func(input string, index int) bool {
		item := assignments.Items[index]
		name := strings.Replace(strings.ToLower(item.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:		 "Multiple tasks in progress. Which one are you working on",
		Items:		 assignments.Items,
		Templates: templates,
		Size:			 4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		os.Exit(1)
	}

	return assignments.Items[i].TPAssignable
}
