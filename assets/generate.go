// +build ignore

package main

import (
	"github.com/andrey-yantsen/mattermost-talks-voting/assets/static"
	"github.com/andrey-yantsen/mattermost-talks-voting/assets/templates"
	"github.com/shurcooL/vfsgen"
	"log"
)

func main() {
	if err := vfsgen.Generate(static.FS, vfsgen.Options{
		Filename:     "assets/static/static_vfsdata.go",
		PackageName:  "static",
		BuildTags:    "deploy_build",
		VariableName: "FS",
	}); err != nil {
		log.Fatalln(err)
	}

	if err := vfsgen.Generate(templates.FS, vfsgen.Options{
		Filename:     "assets/templates/templates_vfsdata.go",
		PackageName:  "templates",
		BuildTags:    "deploy_build",
		VariableName: "FS",
	}); err != nil {
		log.Fatalln(err)
	}
}
