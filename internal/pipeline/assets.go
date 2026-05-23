// Package pipeline holds the It's Important to Note vamp pipeline +
// its embedded prompts. The CLI binary (cmd/iitn) drives it.
package pipeline

import (
	"embed"
	"io/fs"
)

//go:embed prompts/*.md workflows/*.json
var assets embed.FS

// PromptsFS narrows the embed to prompts/ so PromptFS calls
// reference files by their bare name.
var PromptsFS fs.FS = mustSub(assets, "prompts")

// WorkflowsFS exposes ComfyUI workflow JSON files (e.g. sdxl_turbo.json
// for cover-art generation) to the ComfyUI stage via WorkflowFS.
var WorkflowsFS fs.FS = mustSub(assets, "workflows")

func mustSub(fsys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
