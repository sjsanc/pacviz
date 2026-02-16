package themes

import (
	"embed"
	"fmt"
)

//go:embed themes/*.toml
var Embedded embed.FS

func ReadTheme(name string) ([]byte, error) {
	path := fmt.Sprintf("themes/%s.toml", name)
	return Embedded.ReadFile(path)
}
