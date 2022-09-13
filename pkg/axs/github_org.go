package axs

import (
	"fmt"
	"strings"

	"github.com/gocarina/gocsv"
)

var GithubOrgSteps = []string{
	"Open https://github.com/orgs/<org>/people",
	"Click Export",
	"Select 'CSV'",
	"Execute 'axsdump --github-org-members-csv=<path>'",
}

type githubMemberRecord struct {
	Login string `csv:"login"`
	Name  string `csv:"name"`
	Role  string `csv:"role"`
}

// GithubOrgMembers parses the CSV file generated by the Github Members page.
func GithubOrgMembers(path string) (*Artifact, error) {
	src, err := NewSource(path)
	if err != nil {
		return nil, fmt.Errorf("source: %w", err)
	}
	src.Kind = "github_org_members"
	src.Name = "Github Organization Members"
	src.Process = renderSteps(GithubOrgSteps, path)
	a := &Artifact{Metadata: src}

	records := []githubMemberRecord{}
	if err := gocsv.UnmarshalBytes(src.content, &records); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	for _, r := range records {
		u := User{
			Account: r.Login,
			Name:    strings.TrimSpace(r.Name),
			Role:    r.Role,
		}

		if strings.HasSuffix(u.Name, "Bot") || strings.HasSuffix(u.Account, "Bot") {
			a.Bots = append(a.Bots, u)
			continue
		}

		a.Users = append(a.Users, u)
	}

	return a, nil
}
