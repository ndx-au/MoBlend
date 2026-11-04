package project

import (
	"encoding/json"
	"fmt"
	"github.com/ndx-au/MoBlend/internal/diag"
	"github.com/ndx-au/MoBlend/internal/spec"
)

func validateModuleContract(raw []byte, mf ModuleFile, source string) error {
	fail := func(code, pointer, message string) error {
		return diag.ErrorList{{Code: code, Stage: "schema", Pointer: pointer, Source: source, Message: message}}
	}
	if mf.ModuleVersion != ModuleVersion {
		return fail("MOD_VERSION", "/module_version", "unsupported module_version")
	}
	if mf.Kind != "module" && mf.Kind != "entry" {
		return fail("MOD_KIND", "/kind", "kind must be module or entry")
	}
	if !nsRe.MatchString(mf.Namespace) {
		return fail("MOD_NAMESPACE", "/namespace", "invalid namespace")
	}
	if mf.Kind == "module" {
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(raw, &fields); err != nil {
			return err
		}
		for _, key := range []string{"scene", "camera", "schema_version", "asset_pack_version"} {
			if _, present := fields[key]; present {
				return fail("MOD_ENTRY_ONLY", "/"+key, "field is allowed only on an entry")
			}
		}
	} else {
		if mf.Scene == nil || mf.Camera == nil {
			return fail("MOD_ENTRY_SCENE", "", "entry requires scene and camera")
		}
		if mf.SchemaVersion != spec.SchemaVersion || mf.AssetPackVersion != spec.AssetPackVersion {
			return fail("MOD_ENTRY_VERSIONS", "", "entry schema_version/asset_pack_version mismatch")
		}
	}
	if len(mf.Imports) > MaxModules {
		return fail("LIM_IMPORTS", "/imports", "too many imports")
	}
	seen := map[string]bool{}
	for i, im := range mf.Imports {
		pointer := fmt.Sprintf("/imports/%d/alias", i)
		if !nsRe.MatchString(im.Alias) {
			return fail("MOD_IMPORT_ALIAS", pointer, "invalid import alias")
		}
		if seen[im.Alias] {
			return fail("MOD_DUP_IMPORT_ALIAS", pointer, "duplicate import alias")
		}
		seen[im.Alias] = true
	}
	return nil
}
