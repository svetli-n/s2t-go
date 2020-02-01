package storage

type Transcript struct {
	Raw      string   `json:"raw"`
	Edited   string   `json:"edited"`
	OrgEmail string   `json:"org_email"`
	Emails   []string `json:"emails"`
}
