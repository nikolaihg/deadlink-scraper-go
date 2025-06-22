package linktype

type LinkType int

const (
	InternalLink LinkType = iota
	ExternalLink
	PageLink
)

type Link struct {
	URL  string
	Type LinkType
}
