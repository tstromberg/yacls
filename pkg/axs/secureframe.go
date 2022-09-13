package axs

import (
	"fmt"
	"strings"

	"github.com/gocarina/gocsv"
)

var SecureframeSteps = []string{
	"Open https://app.secureframe.com/personnel",
	"Deselect any active filters",
	"Click Export...",
	"Select 'Direct Download'",
	"Download resulting CSV file for analysis",
	"Execute 'axsdump --secureframe-personnel-csv=<path>'",
}

type secureframePersonnelRecord struct {
	Email string `csv:"Name (email)"`
	Role  string `csv:"Access role"`
}

// SecureframePersonnel parses the CSV file generated by the Secureframe Personnel page.
func SecureframePersonnel(path string) (*Artifact, error) {
	src, err := NewSource(path)
	if err != nil {
		return nil, fmt.Errorf("source: %w", err)
	}
	src.Kind = "secureframe_personnel"
	src.Name = "Secureframe Personnel"
	src.Process = renderSteps(SecureframeSteps, path)
	a := &Artifact{Metadata: src}

	records := []secureframePersonnelRecord{}
	if err := gocsv.UnmarshalBytes(src.content, &records); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	for _, r := range records {
		if r.Role == "" {
			continue
		}

		id, _, _ := strings.Cut(r.Email, "@")
		u := User{
			Account: id,
			Role:    r.Role,
		}

		a.Users = append(a.Users, u)
	}

	return a, nil
}
