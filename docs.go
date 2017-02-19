package gocqrs

type APPDocs struct {
	Name     string              `json:"name"`
	Version  string              `json:"version"`
	Entities map[string][]string `json:"entities"`
}

func GenerateDocs(app *App) APPDocs {
	var docs APPDocs
	docs.Entities = make(map[string][]string)
	docs.Name = app.Name
	docs.Version = app.Version

	for e, c := range app.Entities {
		docs.Entities[e] = []string{}
		for event, _ := range c.EventHandlers {
			docs.Entities[e] = append(docs.Entities[e], event)
		}
	}

	return docs
}
