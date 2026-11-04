package project

import (
	"os"

	"github.com/ndx-au/MoBlend/internal/compile"
	"github.com/ndx-au/MoBlend/internal/spec"
	"github.com/ndx-au/MoBlend/internal/validate"
)

// LoadResult is either a single-file or project-resolved declaration.
type LoadResult struct {
	Decl      spec.Declaration
	SourceMap *SourceMap
	IsProject bool
	RootDir   string
	Files     []string
	Shot      string
	Resolved  *Resolved // set for project loads
}

// LoadSceneOrProject loads path as scene JSON or project manifest.
// shot is used only for projects with shots (Phase 8); empty selects default rules.
func LoadSceneOrProject(path, shot string) (LoadResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return LoadResult{}, err
	}
	kind, err := DetectKind(data)
	if err != nil {
		return LoadResult{}, err
	}
	if kind == KindScene {
		res, err := validate.ParseAndValidate(data)
		if err != nil {
			return LoadResult{}, err
		}
		return LoadResult{Decl: res.Decl, IsProject: false}, nil
	}
	resolved, err := LoadAndResolveShot(path, shot)
	if err != nil {
		return LoadResult{}, err
	}
	sm := resolved.SourceMap
	selected := shot
	if selected == "" {
		selected = resolved.Shot
	}
	return LoadResult{
		Decl:      resolved.Decl,
		SourceMap: &sm,
		IsProject: true,
		RootDir:   resolved.RootDir,
		Files:     resolved.Files,
		Shot:      selected,
		Resolved:  &resolved,
	}, nil
}

// CompileLoaded compiles a LoadResult and optionally writes a source map.
func CompileLoaded(lr LoadResult) (compile.Artifacts, error) {
	return compile.CompileDeclaration(lr.Decl)
}
