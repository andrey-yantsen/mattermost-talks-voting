// +build !deploy_build

package static

import (
	"net/http"
)

// Assets is not used in development and is always nil.
var FS http.FileSystem = http.Dir("static")
