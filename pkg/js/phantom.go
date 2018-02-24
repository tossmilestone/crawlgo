package js

import (
	"fmt"

	"github.com/benbjohnson/phantomjs"
)

type Phantom struct {
	userAgent  string
	pageEncode string
}

func NewPhantom() *Phantom {
	return &Phantom{
		userAgent:  "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36",
		pageEncode: "utf-8",
	}
}

func (p *Phantom) Run() error {
	if err := phantomjs.DefaultProcess.Open(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (p *Phantom) Stop() {
	phantomjs.DefaultProcess.Close()
}

func (p *Phantom) ExtractLinksFromSelector(pageUrl string, selector string) ([]interface{}, error) {
	page, err := phantomjs.CreateWebPage()
	if err != nil {
		return nil, err
	}
	defer page.Close()

	if err := page.Open(pageUrl); err != nil {
		return nil, err
	}

	info, err := page.Evaluate(fmt.Sprintf(`function() {
        var links = document.body.querySelectorAll('%s');
        var result = [];
        for (var i = 0; i < links.length; ++i) {
            result.push(links[i].href);
        }
        return result;
    }`, selector))
	if err != nil {
		return nil, err
	}

	if info == nil {
		return nil, nil
	}

	return info.([]interface{}), nil
}
