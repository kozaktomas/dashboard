package gui

import (
	"github.com/xanzy/go-gitlab"
	"strconv"
)

func formatGitlabMergeRequestTitle(project *gitlab.Project, mr *gitlab.MergeRequest) string {
	icon := "🔵"
	if mr.State == "opened" {
		icon = "\U0001F7E2"
	}

	if mr.State == "Closed" {
		icon = "🔴"
	}
	ret := icon + " "
	ret += "[" + project.NameWithNamespace + "!" + strconv.Itoa(mr.IID) + "](fg:cyan)"
	ret += " - " + mr.Title

	return ret
}
