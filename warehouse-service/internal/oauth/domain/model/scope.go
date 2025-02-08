package oauthDomain

type Scope struct {
	Common
	Scope       string `db:"scope" json:"scope"`
	Description string `db:"description" json:"desc"`
	IsDefault   bool   `db:"is_default"`
}
