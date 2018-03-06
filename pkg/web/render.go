package web

// Render describes a web render interface used to render a web page.
type Render interface {
	Run() error
	Stop()
	ExtractLinksFromSelector(pageURL string, selector string) ([]interface{}, error)
}
