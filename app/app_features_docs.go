package app


type DocsFeature struct {
	Enabled bool
	DocsPath string
	docsApiPath string
}

func Docs(docsPath string, docsApiPath string) DocsFeature {
	f := DocsFeature{
		Enabled: true,
		DocsPath: docsPath,
		docsApiPath: docsApiPath,
	}
	
	return f
}