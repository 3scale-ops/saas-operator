package assets

import (
	"bytes"
	"text/template"
)

// Important: Run "make assets" to regenerate code after modifying/adding/removing any asset
//go:generate go-bindata -o bindata.go -pkg $GOPACKAGE dashboards

// SafeStringAsset Returns asset data as string
// panic if not found or any err is detected
func SafeStringAsset(name string) string {
	data, err := Asset(name)
	if err != nil {
		panic(err)
	}

	return string(data)
}

// TemplateAsset Executes one tamplate by applying it to a daata structure.
// panic if not found or any err is detected
func TemplateAsset(name string, data interface{}) string {
	tObj, err := template.New(name).Parse(SafeStringAsset(name))
	if err != nil {
		panic(err)
	}

	var tpl bytes.Buffer
	err = tObj.Execute(&tpl, data)
	if err != nil {
		panic(err)
	}

	return tpl.String()
}
