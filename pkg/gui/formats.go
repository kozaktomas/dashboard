package gui

import (
	"github.com/kozaktomas/dashboard/pkg/gitlab"
	"strconv"
)

func formatGitlabMergeRequestTitle(mr gitlab.MergeRequest) string {
	icon := "🔀"
	if mr.IsOpen() {
		icon = "👷"
	}

	if mr.IsClosed() {
		icon = "🔒"
	}
	ret := icon + " "
	ret += "[" + mr.ProjectName + "](fg:cyan)"
	ret += "!" + strconv.Itoa(mr.Id)
	ret += " - " + mr.Title

	return ret
}
