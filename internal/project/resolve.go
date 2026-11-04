package project

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/ndx-au/MoBlend/internal/diag"
	"github.com/ndx-au/MoBlend/internal/resources"
	"github.com/ndx-au/MoBlend/internal/spec"
	"github.com/ndx-au/MoBlend/internal/validate"
)

type idMap map[string]string

var (
	nsRe    = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*$`)
	localRe = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)
)

type loadedModule struct {
	ID      string
	Path    string // absolute
	RelPath string // project-relative
	File    ModuleFile
	Depth   int
}

func resolveManifest(root, manRel string, manifestData []byte, man Manifest, shotID string, allListed bool) (Resolved, error) {
	if len(man.Modules) == 0 {
		return Resolved{}, diag.ErrorList{{
			Code: "PROJ_MODULES_EMPTY", Stage: "schema", Pointer: "/modules", Source: manRel,
			Message: "modules list is empty",
		}}
	}
	if len(man.Modules) > MaxModules {
		return Resolved{}, diag.ErrorList{{
			Code: "LIM_MODULES", Stage: "semantic", Pointer: "/modules", Source: manRel,
			Message: fmt.Sprintf("more than %d modules", MaxModules),
		}}
	}
	if man.Entry == "" {
		return Resolved{}, diag.ErrorList{{
			Code: "PROJ_ENTRY", Stage: "schema", Pointer: "/entry", Source: manRel,
			Message: "entry module id is required",
		}}
	}

	byID := map[string]ManifestMod{}
	seenID := map[string]struct{}{}
	for i, m := range man.Modules {
		p := fmt.Sprintf("/modules/%d", i)
		if m.ID == "" || m.Path == "" {
			return Resolved{}, diag.ErrorList{{
				Code: "PROJ_MODULE_FIELDS", Stage: "schema", Pointer: p, Source: manRel,
				Message: "module id and path are required",
			}}
		}
		if _, dup := seenID[m.ID]; dup {
			return Resolved{}, diag.ErrorList{{
				Code: "PROJ_DUP_MODULE_ID", Stage: "semantic", Pointer: p + "/id", Source: manRel,
				Message: fmt.Sprintf("duplicate module id %q", m.ID),
			}}
		}
		seenID[m.ID] = struct{}{}
		byID[m.ID] = m
	}
	if _, ok := byID[man.Entry]; !ok {
		return Resolved{}, diag.ErrorList{{
			Code: "PROJ_ENTRY_UNKNOWN", Stage: "semantic", Pointer: "/entry", Source: manRel,
			Message: fmt.Sprintf("entry %q is not in modules", man.Entry),
		}}
	}

	rootReal, err := filepath.EvalSymlinks(root)
	if err != nil {
		return Resolved{}, err
	}

	root = rootReal

	totalBytes := int64(len(manifestData))
	digests := []FileDigest{digest(manRel, manifestData)}
	loaded := map[string]*loadedModule{}
	var order []string

	var loadOne func(id string, depth int) error
	loadOne = func(id string, depth int) error {
		if depth > MaxImportDepth {
			return diag.ErrorList{{
				Code: "LIM_IMPORT_DEPTH", Stage: "semantic", Source: manRel,
				Message: fmt.Sprintf("import depth exceeds %d", MaxImportDepth),
			}}
		}
		if _, ok := loaded[id]; ok {
			return nil
		}
		mm, ok := byID[id]
		if !ok {
			return diag.ErrorList{{
				Code: "PROJ_MODULE_MISSING", Stage: "semantic", Source: manRel,
				Message: fmt.Sprintf("unknown module id %q", id),
			}}
		}
		abs, err := resolveContainedPath(root, rootReal, mm.Path)
		if err != nil {
			return err
		}
		data, err := readBounded(abs, MaxFileBytes)
		if err != nil {
			return withSource(err, rel(root, abs))
		}
		digests = append(digests, digest(rel(root, abs), data))
		totalBytes += int64(len(data))
		if totalBytes > int64(MaxProjectBytes) {
			return diag.ErrorList{{
				Code: "LIM_PROJECT_BYTES", Stage: "semantic", Source: manRel,
				Message: fmt.Sprintf("project exceeds %d bytes", MaxProjectBytes),
			}}
		}
		var mf ModuleFile
		if err := strictDecode(data, &mf); err != nil {
			return withSource(err, rel(root, abs))
		}
		if err := validateModuleContract(data, mf, rel(root, abs)); err != nil {
			return err
		}
		lm := &loadedModule{ID: id, Path: abs, RelPath: rel(root, abs), File: mf, Depth: depth}
		loaded[id] = lm
		order = append(order, id)
		for _, im := range mf.Imports {
			if err := loadOne(im.Module, depth+1); err != nil {
				return err
			}
		}
		return nil
	}

	// Load the entry closure first; whole-project mode additionally checks listed orphans.
	if err := loadOne(man.Entry, 0); err != nil {
		if el, ok := err.(diag.ErrorList); ok {
			return Resolved{}, el
		}
		return Resolved{}, err
	}
	// Check unused listed files as graph inputs, without merging their entities.
	closure := append([]string(nil), order...)
	if allListed {
		for _, m := range man.Modules {
			if err := loadOne(m.ID, 0); err != nil {
				return Resolved{}, err
			}
		}
	}

	// Cycle detection via DFS on import edges within loaded set.
	if err := detectCycles(loaded); err != nil {
		return Resolved{}, err
	}

	// Depth must be bounded along every path, even if a diamond loaded a
	// dependency first through a shorter path. Cycles were rejected above.
	depths := map[string]int{}
	var graphDepth func(string) int
	graphDepth = func(id string) int {
		if d, ok := depths[id]; ok {
			return d
		}
		d := 0
		for _, im := range loaded[id].File.Imports {
			if n := 1 + graphDepth(im.Module); n > d {
				d = n
			}
		}
		depths[id] = d
		return d
	}
	for id := range loaded {
		if graphDepth(id) > MaxImportDepth {
			return Resolved{}, diag.ErrorList{{Code: "LIM_IMPORT_DEPTH", Stage: "semantic", Source: manRel, Message: "import depth exceeds limit"}}
		}
	}
	// Only the selected entry closure participates in scene composition.
	order = closure
	entry := loaded[man.Entry]
	if entry.File.Kind != "entry" {
		return Resolved{}, diag.ErrorList{{
			Code: "PROJ_ENTRY_KIND", Stage: "semantic", Pointer: "/entry", Source: entry.RelPath,
			Message: "entry module must have kind \"entry\"",
		}}
	}
	if entry.File.Scene == nil || entry.File.Camera == nil {
		return Resolved{}, diag.ErrorList{{
			Code: "MOD_ENTRY_SCENE", Stage: "schema", Source: entry.RelPath,
			Message: "entry requires scene and camera",
		}}
	}
	if entry.File.SchemaVersion != spec.SchemaVersion || entry.File.AssetPackVersion != spec.AssetPackVersion {
		return Resolved{}, diag.ErrorList{{
			Code: "MOD_ENTRY_VERSIONS", Stage: "schema", Source: entry.RelPath,
			Message: "entry schema_version/asset_pack_version mismatch",
		}}
	}

	// Namespace uniqueness among loaded modules.
	nsOwner := map[string]string{}
	for _, id := range order {
		lm := loaded[id]
		if prev, ok := nsOwner[lm.File.Namespace]; ok {
			return Resolved{}, diag.ErrorList{{
				Code: "MOD_DUP_NAMESPACE", Stage: "semantic", Pointer: "/namespace", Source: lm.RelPath,
				Message: fmt.Sprintf("namespace %q already used by module %q", lm.File.Namespace, prev),
			}}
		}
		nsOwner[lm.File.Namespace] = id
	}

	sm := SourceMap{
		Objects: map[string]SourceRef{},
		Lights:  map[string]SourceRef{},
		Entry:   SourceRef{File: entry.RelPath, Module: man.Entry},
	}
	var decl spec.Declaration
	decl.SchemaVersion = entry.File.SchemaVersion
	decl.AssetPackVersion = entry.File.AssetPackVersion
	decl.Scene = *entry.File.Scene
	decl.Camera = *entry.File.Camera

	// Assemble in sorted module-id order for determinism (not import order).
	ids := append([]string(nil), order...)
	sort.Strings(ids)

	moduleMaps := map[string]idMap{}
	aliasMaps := map[string]map[string]string{} // module id -> alias -> imported module ID

	for _, id := range ids {
		lm := loaded[id]
		aliases := map[string]string{}
		for _, im := range lm.File.Imports {
			aliases[im.Alias] = im.Module
		}
		aliasMaps[id] = aliases
		local := idMap{}
		if id == man.Entry {
			local[entry.File.Camera.ID] = entry.File.Camera.ID
		}
		qualify := func(localID string) string {
			return lm.File.Namespace + "__" + localID
		}
		for i, o := range lm.File.Objects {
			if !localRe.MatchString(o.ID) {
				return Resolved{}, diag.ErrorList{{
					Code: "SEM_ID", Stage: "semantic",
					Pointer: fmt.Sprintf("/objects/%d/id", i), Source: lm.RelPath,
					Message: "invalid id",
				}}
			}
			if _, duplicate := local[o.ID]; duplicate {
				return Resolved{}, diag.ErrorList{{Code: "SEM_DUPLICATE_ID", Stage: "semantic", Pointer: fmt.Sprintf("/objects/%d/id", i), Source: lm.RelPath, Message: "ambiguous local id"}}
			}
			qid := qualify(o.ID)
			if _, exists := sm.Objects[qid]; exists {
				return Resolved{}, diag.ErrorList{{
					Code: "SEM_DUPLICATE_ID", Stage: "semantic",
					Pointer: fmt.Sprintf("/objects/%d/id", i), Source: lm.RelPath,
					Message: fmt.Sprintf("duplicate qualified id %q", qid),
				}}
			}
			obj := o
			obj.ID = qid
			decl.Objects = append(decl.Objects, obj)
			local[o.ID] = qid
			sm.Objects[qid] = SourceRef{File: lm.RelPath, Pointer: fmt.Sprintf("/objects/%d", i), Module: id}
		}
		for i, l := range lm.File.Lights {
			if !localRe.MatchString(l.ID) {
				return Resolved{}, diag.ErrorList{{
					Code: "SEM_ID", Stage: "semantic",
					Pointer: fmt.Sprintf("/lights/%d/id", i), Source: lm.RelPath,
					Message: "invalid id",
				}}
			}
			if _, duplicate := local[l.ID]; duplicate {
				return Resolved{}, diag.ErrorList{{Code: "SEM_DUPLICATE_ID", Stage: "semantic", Pointer: fmt.Sprintf("/lights/%d/id", i), Source: lm.RelPath, Message: "ambiguous local id"}}
			}
			qid := qualify(l.ID)
			if _, exists := sm.Lights[qid]; exists {
				return Resolved{}, diag.ErrorList{{
					Code: "SEM_DUPLICATE_ID", Stage: "semantic",
					Pointer: fmt.Sprintf("/lights/%d/id", i), Source: lm.RelPath,
					Message: fmt.Sprintf("duplicate qualified id %q", qid),
				}}
			}
			lt := l
			lt.ID = qid
			decl.Lights = append(decl.Lights, lt)
			local[l.ID] = qid
			sm.Lights[qid] = SourceRef{File: lm.RelPath, Pointer: fmt.Sprintf("/lights/%d", i), Module: id}
		}
		moduleMaps[id] = local
	}

	for _, id := range ids {
		lm := loaded[id]
		local := moduleMaps[id]
		aliasToNS := aliasMaps[id]
		for i, tr := range lm.File.Animation {
			target, err := resolveRef(tr.TargetID, lm.File.Namespace, local, aliasToNS, moduleMaps)
			if err != nil {
				return Resolved{}, diag.ErrorList{{
					Code: "SEM_TRACK_TARGET", Stage: "semantic",
					Pointer: fmt.Sprintf("/animation/%d/target_id", i), Source: lm.RelPath,
					Message: err.Error(),
				}}
			}
			tr2 := tr
			tr2.TargetID = target
			decl.Animation = append(decl.Animation, tr2)
			sm.Tracks = append(sm.Tracks, SourceRef{
				File: lm.RelPath, Pointer: fmt.Sprintf("/animation/%d", i), Module: id,
			})
		}
	}

	// Stable object/light ordering by qualified ID for determinism.
	sort.SliceStable(decl.Objects, func(i, j int) bool { return decl.Objects[i].ID < decl.Objects[j].ID })
	sort.SliceStable(decl.Lights, func(i, j int) bool { return decl.Lights[i].ID < decl.Lights[j].ID })

	cat, err := resources.LoadCatalogue()
	if err != nil {
		return Resolved{}, err
	}
	res, err := validate.Validate(decl, cat)
	if err != nil {
		if el, ok := err.(diag.ErrorList); ok {
			for i := range el {
				el[i] = annotateSource(el[i], sm, decl)
			}
			diag.Sort(el)
			return Resolved{}, el
		}
		return Resolved{}, err
	}

	files := make([]string, 0, len(digests))
	for _, d := range digests {
		files = append(files, d.Path)
	}
	sort.Strings(files)

	return Resolved{
		Decl:      res.Decl,
		SourceMap: sm,
		RootDir:   root,
		Files:     files,
		Digests:   digests,
		Shot:      shotID,
	}, nil
}

func resolveRef(ref, selfNS string, local idMap, aliases map[string]string, moduleMaps map[string]idMap) (string, error) {
	if strings.Contains(ref, ":") {
		parts := strings.Split(ref, ":")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return "", fmt.Errorf("invalid cross-module ref %q", ref)
		}
		moduleID, ok := aliases[parts[0]]
		if !ok {
			return "", fmt.Errorf("unknown import alias %q", parts[0])
		}
		if q, ok := moduleMaps[moduleID][parts[1]]; ok {
			return q, nil
		}
		return "", fmt.Errorf("unknown id %q in module alias %q", parts[1], parts[0])
	}
	if q, ok := local[ref]; ok {
		return q, nil
	}
	if strings.Contains(ref, "__") {
		return "", fmt.Errorf("qualified target %q requires an explicit alias:local reference", ref)
	}
	return "", fmt.Errorf("unknown target_id %q in namespace %s", ref, selfNS)
}

func detectCycles(loaded map[string]*loadedModule) error {
	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := map[string]int{}
	var visit func(id string) error
	visit = func(id string) error {
		color[id] = gray
		lm := loaded[id]
		for _, im := range lm.File.Imports {
			dep := im.Module
			if _, ok := loaded[dep]; !ok {
				continue
			}
			switch color[dep] {
			case gray:
				return diag.ErrorList{{
					Code: "PROJ_CYCLE", Stage: "semantic", Source: lm.RelPath,
					Message: fmt.Sprintf("import cycle involving %q -> %q", id, dep),
				}}
			case white:
				if err := visit(dep); err != nil {
					return err
				}
			}
		}
		color[id] = black
		return nil
	}
	for id := range loaded {
		if color[id] == white {
			if err := visit(id); err != nil {
				return err
			}
		}
	}
	return nil
}

func resolveContainedPath(root, rootReal, relPath string) (string, error) {
	if filepath.IsAbs(relPath) {
		return "", diag.ErrorList{{
			Code: "PROJ_PATH_ABS", Stage: "semantic",
			Message: "module path must be project-relative",
		}}
	}
	clean := filepath.Clean(relPath)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", diag.ErrorList{{
			Code: "PROJ_PATH_ESCAPE", Stage: "semantic",
			Message: fmt.Sprintf("path escapes project root: %q", relPath),
		}}
	}
	abs := filepath.Join(root, clean)
	if _, err := os.Lstat(abs); err != nil {
		if os.IsNotExist(err) {
			return "", diag.ErrorList{{
				Code: "PROJ_FILE_MISSING", Stage: "semantic",
				Message: fmt.Sprintf("missing module file %q", relPath),
			}}
		}
		return "", diag.ErrorList{{Code: "PROJ_IO", Stage: "parse", Message: err.Error()}}
	}
	real, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", diag.ErrorList{{
			Code: "PROJ_FILE_MISSING", Stage: "semantic",
			Message: fmt.Sprintf("missing module file %q", relPath),
		}}
	}
	rootPrefix := rootReal + string(filepath.Separator)
	if real != rootReal && !strings.HasPrefix(real, rootPrefix) {
		return "", diag.ErrorList{{
			Code: "PROJ_PATH_ESCAPE", Stage: "semantic",
			Message: fmt.Sprintf("path escapes project root: %q", relPath),
		}}
	}
	return real, nil
}

func annotateSource(d diag.Diagnostic, sm SourceMap, decl spec.Declaration) diag.Diagnostic {
	if d.Source != "" {
		return d
	}
	// Assembled pointers (/objects/N/...) map via qualified id on the resolved declaration.
	if n, ok := indexFromPointer(d.Pointer, "/objects/"); ok && n >= 0 && n < len(decl.Objects) {
		if ref, found := sm.Objects[decl.Objects[n].ID]; found {
			d.Source = ref.File
			d.Pointer = ref.Pointer + strings.TrimPrefix(d.Pointer, fmt.Sprintf("/objects/%d", n))
			return d
		}
	}
	if n, ok := indexFromPointer(d.Pointer, "/lights/"); ok && n >= 0 && n < len(decl.Lights) {
		if ref, found := sm.Lights[decl.Lights[n].ID]; found {
			d.Source = ref.File
			d.Pointer = ref.Pointer + strings.TrimPrefix(d.Pointer, fmt.Sprintf("/lights/%d", n))
			return d
		}
	}
	if n, ok := indexFromPointer(d.Pointer, "/animation/"); ok && n >= 0 && n < len(sm.Tracks) {
		d.Source = sm.Tracks[n].File
		d.Pointer = sm.Tracks[n].Pointer + strings.TrimPrefix(d.Pointer, fmt.Sprintf("/animation/%d", n))
		return d
	}
	if strings.HasPrefix(d.Pointer, "/scene") || strings.HasPrefix(d.Pointer, "/camera") || d.Pointer == "/schema_version" || d.Pointer == "/asset_pack_version" {
		d.Source = sm.Entry.File
	}
	return d
}

func indexFromPointer(pointer, prefix string) (int, bool) {
	if !strings.HasPrefix(pointer, prefix) {
		return 0, false
	}
	rest := strings.TrimPrefix(pointer, prefix)
	num := rest
	if i := strings.IndexByte(rest, '/'); i >= 0 {
		num = rest[:i]
	}
	var n int
	if _, err := fmt.Sscanf(num, "%d", &n); err != nil {
		return 0, false
	}
	return n, true
}
