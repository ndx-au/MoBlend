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
	"github.com/ndx-au/MoBlend/internal/resources"
	"github.com/ndx-au/MoBlend/internal/runner"
	"github.com/ndx-au/MoBlend/internal/spec"
	"github.com/ndx-au/MoBlend/internal/validate"
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
		data, err := resources.SchemaJSON()
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
  blendctl schema
  blendctl assets list [--json]
  blendctl assets describe <asset-id> [--json]
  blendctl validate <scene.json> [--json]
  blendctl compile <scene.json> --out DIR
  blendctl doctor [--json] [--blender PATH]
  blendctl guide
  blendctl preview <scene.json> [--frame N] [--out DIR]
  blendctl render <scene.json> [--out DIR]
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
		fmt.Fprintln(os.Stderr, "usage: blendctl validate <scene.json> [--json]")
		return fmt.Errorf("missing path")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "validate: %v\n", err)
		return err
	}
	res, err := validate.ParseAndValidate(data)
	wantJSON := hasFlag(args, "--json")
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
		return encodeJSON(map[string]any{"ok": true, "schema_version": res.Decl.SchemaVersion})
	}
	fmt.Println("OK")
	return nil
}

func cmdCompile(args []string) error {
	path := firstNonFlag(args)
	outDir := flagValue(args, "--out")
	if path == "" || outDir == "" {
		return fmt.Errorf("usage: blendctl compile <scene.json> --out DIR")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	art, _, err := compile.CompileJSON(data)
	if err != nil {
		return err
	}
	if err := compile.WriteToDir(outDir, art); err != nil {
		return err
	}
	fmt.Printf("wrote %s\nwrote %s\n", filepath.Join(outDir, art.ScriptName), filepath.Join(outDir, art.DataName))
	return nil
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
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	art, validated, err := compile.CompileJSON(data)
	if err != nil {
		return err
	}
	outParent := flagValue(args, "--out")
	if outParent != "" {
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
		Decl:        validated.Decl,
		Mode:        mode,
		RenderFrame: frame,
	})
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
