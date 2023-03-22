package pages

import "fmt"


type Pages struct {
	// Name of the base template
	Base string

	// List of html files
	Templates []string	

	// path of the templates. Example : "./templates", "./static"
	TemplatesPath string
}

func NewPages(baseName, tmpPath string, templates []string) *Pages{

	var templateList []string

	for _, tmpl := range templates {
		templateList = append(templateList, fmt.Sprintf("%s/%s.html", tmpPath, tmpl))
	}


	return &Pages{
		Base: baseName,
		Templates: templateList,
		TemplatesPath: tmpPath,
	}

}
