// +build ignore

package main

import (
	"log"
	"github.com/shurcooL/vfsgen"
	static "github.com/andrey-yantsen/mattermost-talks-voting/assets/static"
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
}
