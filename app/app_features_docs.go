package app


type DocsFeature struct {
	Enabled bool
	DocsPath string
	DocsApiPath string
}

func Docs(docsPath string, docsApiPath string) DocsFeature {
	f := DocsFeature{
		Enabled: true,
		DocsPath: docsPath,
		DocsApiPath: docsApiPath,
	}
	
	return f
}