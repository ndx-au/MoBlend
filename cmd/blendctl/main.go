// Command blendctl is the MoBlend CLI.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/ndx-au/MoBlend/internal/compile"
	"github.com/ndx-au/MoBlend/internal/diag"
	"github.com/ndx-au/MoBlend/internal/project"
	"github.com/ndx-au/MoBlend/internal/resources"
	"github.com/ndx-au/MoBlend/internal/runner"
	"github.com/ndx-au/MoBlend/internal/spec"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "version":
		fmt.Println(versionString())
	case "smoke":
		if err := cmdSmoke(); err != nil {
			fmt.Fprintf(os.Stderr, "smoke: %v\n", err)
			os.Exit(1)
		}
	case "pack-sha256":
		sum, err := resources.BasePackSHA256()
		if err != nil {
			fmt.Fprintf(os.Stderr, "pack-sha256: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(sum)
	case "schema":
		kind := "scene"
		if len(os.Args) > 2 && !hasFlag(os.Args[2:], "--json") {
			if k := firstNonFlag(os.Args[2:]); k != "" {
				kind = k
			}
		}
		data, err := resources.SchemaKindJSON(kind)
		if err != nil {
			fmt.Fprintf(os.Stderr, "schema: %v\n", err)
			os.Exit(1)
		}
		os.Stdout.Write(data)
		if len(data) == 0 || data[len(data)-1] != '\n' {
			fmt.Fprintln(os.Stdout)
		}
	case "assets":
		if err := cmdAssets(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "assets: %v\n", err)
			os.Exit(1)
		}
	case "validate":
		if err := cmdValidate(os.Args[2:]); err != nil {
			os.Exit(1)
		}
	case "compile":
		if err := cmdCompile(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "compile: %v\n", err)
			os.Exit(1)
		}
	case "doctor":
		code := cmdDoctor(os.Args[2:])
		os.Exit(code)
	case "guide":
		data, err := resources.GuideMarkdown()
		if err != nil {
			fmt.Fprintf(os.Stderr, "guide: %v\n", err)
			os.Exit(1)
		}
		os.Stdout.Write(data)
		if len(data) == 0 || data[len(data)-1] != '\n' {
			fmt.Fprintln(os.Stdout)
		}
	case "preview":
		if err := cmdPreview(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "preview: %v\n", err)
			os.Exit(1)
		}
	case "render":
		if err := cmdRender(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "render: %v\n", err)
			os.Exit(1)
		}
	case "inspect":
		if err := cmdInspect(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "inspect: %v\n", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `blendctl — MoBlend CLI (provisional)

Usage:
  blendctl version
  blendctl pack-sha256
  blendctl schema [scene|project|module]
  blendctl assets list [--json]
  blendctl assets describe <asset-id> [--json]
  blendctl validate <scene-or-project.json> [--json] [--shot ID]
  blendctl compile <scene-or-project.json> --out DIR [--shot ID]
  blendctl inspect <project.json> [--json]
  blendctl doctor [--json] [--blender PATH]
  blendctl guide
  blendctl preview <scene-or-project.json> [--frame N] [--out DIR] [--shot ID]
  blendctl render <scene-or-project.json> [--out DIR] [--shot ID]
  blendctl smoke [--out DIR]

Environment:
  MOBLEND_BLENDER   absolute path to Blender (else BLENDER, else PATH)

validate / schema / assets / compile / guide never launch Blender.
`)
}

// version is set via -ldflags "-X main.version=..." for releases.
var version = ""

func versionString() string {
	v := version
	if v == "" {
		v = "dev"
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			v = info.Main.Version
		}
	}
	return fmt.Sprintf("blendctl %s schema=%s pack=%s", v, spec.SchemaVersion, resources.PackVersion)
}

func cmdSmoke() error {
	outParent := ""
	args := os.Args[2:]
	for i := 0; i < len(args); i++ {
		if args[i] == "--out" && i+1 < len(args) {
			outParent = args[i+1]
			i++
		}
	}
	if outParent != "" {
		if err := os.MkdirAll(outParent, 0o755); err != nil {
			return fmt.Errorf("create --out parent: %w", err)
		}
	}
	res, err := runner.RunSmoke(context.Background(), runner.SmokeOptions{
		Parent: outParent,
	})
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if encErr := enc.Encode(res); encErr != nil && err == nil {
		return encErr
	}
	return err
}

func cmdAssets(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: blendctl assets list|describe ...")
	}
	raw, err := resources.CatalogueJSON()
	if err != nil {
		return err
	}
	var cat struct {
		AssetPackVersion string           `json:"asset_pack_version"`
		Assets           []map[string]any `json:"assets"`
	}
	if err := json.Unmarshal(raw, &cat); err != nil {
		return err
	}
	asJSON := hasFlag(args, "--json")
	switch args[0] {
	case "list":
		if asJSON {
			return encodeJSON(cat.Assets)
		}
		for _, a := range cat.Assets {
			fmt.Printf("%s\t%v\n", a["id"], a["type"])
		}
		return nil
	case "describe":
		rest := args[1:]
		id := firstNonFlag(rest)
		if id == "" {
			return fmt.Errorf("usage: blendctl assets describe <asset-id>")
		}
		for _, a := range cat.Assets {
			if a["id"] == id {
				return encodeJSON(a)
			}
		}
		return fmt.Errorf("unknown asset_id %q", id)
	default:
		return fmt.Errorf("unknown assets subcommand %q", args[0])
	}
}

func cmdValidate(args []string) error {
	path := firstNonFlag(args)
	if path == "" {
		fmt.Fprintln(os.Stderr, "usage: blendctl validate <scene-or-project.json> [--json] [--shot ID]")
		return fmt.Errorf("missing path")
	}
	wantJSON := hasFlag(args, "--json")
	shot := flagValue(args, "--shot")
	if hasFlag(args, "--all-shots") {
		return cmdValidateAllShots(path, wantJSON)
	}
	lr, err := project.LoadSceneOrProject(path, shot)
	if err != nil {
		elist, ok := err.(diag.ErrorList)
		if !ok {
			elist = diag.ErrorList{{Code: "SPEC_JSON", Stage: "parse", Message: err.Error()}}
		}
		if wantJSON {
			_ = encodeJSON(map[string]any{"ok": false, "diagnostics": elist})
		} else {
			for _, d := range elist {
				fmt.Fprintln(os.Stderr, d.Error())
			}
		}
		return err
	}
	if wantJSON {
		return encodeJSON(map[string]any{
			"ok":             true,
			"schema_version": lr.Decl.SchemaVersion,
			"project":        lr.IsProject,
			"shot":           lr.Shot,
			"files":          lr.Files,
		})
	}
	fmt.Println("OK")
	return nil
}

func cmdValidateAllShots(path string, wantJSON bool) error {
	rep, err := project.Inspect(path)
	if err != nil {
		return err
	}
	if len(rep.Shots) == 0 {
		_, err := project.LoadSceneOrProject(path, "")
		return err
	}
	var all diag.ErrorList
	for _, s := range rep.Shots {
		if _, err := project.LoadSceneOrProject(path, s.ID); err != nil {
			if el, ok := err.(diag.ErrorList); ok {
				for i := range el {
					if el[i].Hint == "" {
						el[i].Hint = "shot " + s.ID
					}
				}
				all = append(all, el...)
			} else {
				all = append(all, diag.Diagnostic{Code: "SHOT_VALIDATE", Stage: "semantic", Message: err.Error(), Hint: "shot " + s.ID})
			}
		}
	}
	if len(all) > 0 {
		if wantJSON {
			_ = encodeJSON(map[string]any{"ok": false, "diagnostics": all})
		} else {
			for _, d := range all {
				fmt.Fprintln(os.Stderr, d.Error())
			}
		}
		return all
	}
	if wantJSON {
		return encodeJSON(map[string]any{"ok": true, "shots": len(rep.Shots)})
	}
	fmt.Println("OK")
	return nil
}

func cmdCompile(args []string) error {
	path := firstNonFlag(args)
	outDir := flagValue(args, "--out")
	if path == "" || outDir == "" {
		return fmt.Errorf("usage: blendctl compile <scene-or-project.json> --out DIR [--shot ID]")
	}
	lr, err := project.LoadSceneOrProject(path, flagValue(args, "--shot"))
	if err != nil {
		return err
	}
	art, err := compile.CompileDeclaration(lr.Decl)
	if err != nil {
		return err
	}
	if err := compile.WriteToDir(outDir, art); err != nil {
		return err
	}
	if lr.SourceMap != nil {
		data, _ := json.MarshalIndent(lr.SourceMap, "", "  ")
		_ = os.WriteFile(filepath.Join(outDir, "source_map.json"), append(data, '\n'), 0o644)
	}
	if lr.Resolved != nil {
		if _, err := project.WriteDependencyManifest(outDir, *lr.Resolved); err != nil {
			return err
		}
	}
	fmt.Printf("wrote %s\nwrote %s\n", filepath.Join(outDir, art.ScriptName), filepath.Join(outDir, art.DataName))
	return nil
}

func cmdInspect(args []string) error {
	path := firstNonFlag(args)
	if path == "" {
		return fmt.Errorf("usage: blendctl inspect <project.json> [--json]")
	}
	rep, err := project.Inspect(path)
	if err != nil {
		return err
	}
	return encodeJSON(rep)
}

func cmdDoctor(args []string) int {
	blender := flagValue(args, "--blender")
	rep, err := runner.Doctor(context.Background(), blender)
	if hasFlag(args, "--json") {
		os.Stdout.Write(runner.DoctorJSON(rep))
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		if rep.Hint != "" {
			fmt.Fprintln(os.Stderr, rep.Hint)
		}
	} else {
		fmt.Printf("OK %s (%s)\n", rep.BlenderPath, rep.VersionLine)
	}
	if err != nil {
		return 3
	}
	return 0
}

func cmdPreview(args []string) error {
	return cmdApply(args, "preview", "usage: blendctl preview <scene.json> [--frame N] [--out DIR]")
}

func cmdRender(args []string) error {
	return cmdApply(args, "render", "usage: blendctl render <scene.json> [--out DIR]")
}

func cmdApply(args []string, mode, usage string) error {
	path := firstNonFlag(args)
	if path == "" {
		return fmt.Errorf("%s", usage)
	}
	lr, err := project.LoadSceneOrProject(path, flagValue(args, "--shot"))
	if err != nil {
		return err
	}
	art, err := compile.CompileDeclaration(lr.Decl)
	if err != nil {
		return err
	}
	outParent := flagValue(args, "--out")
	if outParent != "" {
		if err := os.MkdirAll(outParent, 0o755); err != nil {
			return err
		}
	}
	// Per-shot isolation: nest under shot id when selected.
	if shot := flagValue(args, "--shot"); shot != "" && outParent != "" {
		outParent = filepath.Join(outParent, shot)
		if err := os.MkdirAll(outParent, 0o755); err != nil {
			return err
		}
	}
	frame := 0
	if v := flagValue(args, "--frame"); v != "" {
		if _, err := fmt.Sscanf(v, "%d", &frame); err != nil {
			return fmt.Errorf("bad --frame: %w", err)
		}
	}
	res, err := runner.RunApply(context.Background(), runner.ApplyOptions{
		Parent:      outParent,
		Artifacts:   art,
		Decl:        lr.Decl,
		Mode:        mode,
		RenderFrame: frame,
	})
	if err == nil && lr.Resolved != nil {
		_, _ = project.WriteDependencyManifest(filepath.Join(res.JobDir, "out"), *lr.Resolved)
	}
	if encErr := encodeJSON(res); encErr != nil && err == nil {
		return encErr
	}
	return err
}

func flagValue(args []string, name string) string {
	for i := 0; i < len(args); i++ {
		if args[i] == name && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func hasFlag(args []string, name string) bool {
	for _, a := range args {
		if a == name {
			return true
		}
	}
	return false
}

func firstNonFlag(args []string) string {
	for _, a := range args {
		if a != "" && a[0] != '-' {
			return a
		}
	}
	return ""
}

func encodeJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
