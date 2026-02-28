package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	listFlag := flag.Bool("list", false, "List all discovered operations and their arguments")
	categoryFlag := flag.String("category", "", "Filter by category (e.g. arithmetic, resample)")
	enumsFlag := flag.Bool("enums", false, "List all discovered enum types")
	coverageFlag := flag.Bool("coverage", false, "Show coverage report of generated vs hand-written vs missing ops")
	generateFlag := flag.Bool("generate", false, "Generate C and Go bridge code")
	outputDir := flag.String("output", "", "Output directory for generated files (default: vips/)")
	flag.Parse()

	ops, err := Introspect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *listFlag {
		listOps(ops, *categoryFlag)
		return
	}

	if *enumsFlag {
		listEnums(ops)
		return
	}

	if *coverageFlag {
		showCoverage(ops)
		return
	}

	if *generateFlag {
		dir := *outputDir
		if dir == "" {
			dir = "vips"
		}
		if err := Generate(ops, dir); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating: %v\n", err)
			os.Exit(1)
		}
		return
	}

	flag.Usage()
}

func listOps(ops []OpDef, categoryFilter string) {
	// Group by category.
	categories := make(map[string][]OpDef)
	for _, op := range ops {
		cat := normalizeCategoryForOp(op.Name, op.Category)
		if categoryFilter != "" && cat != categoryFilter {
			continue
		}
		if excludeOps[op.Name] {
			continue
		}
		categories[cat] = append(categories[cat], op)
	}

	// Sort category names.
	catNames := make([]string, 0, len(categories))
	for name := range categories {
		catNames = append(catNames, name)
	}
	sortStrings(catNames)

	total := 0
	excluded := 0
	for _, op := range ops {
		if excludeOps[op.Name] {
			excluded++
		} else {
			total++
		}
	}

	fmt.Printf("Discovered %d operations (%d excluded, %d generatable)\n\n", len(ops), excluded, total)

	for _, cat := range catNames {
		catOps := categories[cat]
		fmt.Printf("=== %s (%d ops) ===\n", cat, len(catOps))
		for _, op := range catOps {
			fmt.Printf("  %s: %s\n", op.Name, op.Description)

			reqInputs := op.RequiredInputs()
			optInputs := op.OptionalInputs()
			outputs := op.Outputs()

			if len(reqInputs) > 0 {
				fmt.Printf("    required: %s\n", formatArgs(reqInputs))
			}
			if len(optInputs) > 0 {
				fmt.Printf("    optional: %s\n", formatArgs(optInputs))
			}
			if len(outputs) > 0 {
				fmt.Printf("    outputs:  %s\n", formatArgs(outputs))
			}
			fmt.Println()
		}
	}
}

func listEnums(ops []OpDef) {
	enumTypes := CollectEnumTypes(ops)
	fmt.Printf("Discovered %d enum types\n\n", len(enumTypes))

	enums, err := IntrospectEnums(enumTypes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error introspecting enums: %v\n", err)
		os.Exit(1)
	}

	for _, e := range enums {
		goName := goEnumName(e.CName)
		if goName == "" {
			goName = "(unmapped)"
		}
		fmt.Printf("%s -> %s\n", e.CName, goName)
		for _, v := range e.Values {
			fmt.Printf("  %s = %d (nick: %s)\n", v.CName, v.Value, v.Nick)
		}
		fmt.Println()
	}
}

func showCoverage(ops []OpDef) {
	catCounts := make(map[string]int)
	catExcluded := make(map[string]int)

	for _, op := range ops {
		cat := normalizeCategoryForOp(op.Name, op.Category)
		if excludeOps[op.Name] {
			catExcluded[cat]++
		} else {
			catCounts[cat]++
		}
	}

	fmt.Printf("Coverage Report\n")
	fmt.Printf("%-20s %8s %8s %8s\n", "Category", "Generate", "Excluded", "Total")
	fmt.Printf("%s\n", strings.Repeat("-", 50))

	catNames := make([]string, 0, len(catCounts))
	for name := range catCounts {
		catNames = append(catNames, name)
	}
	// Also add excluded-only categories.
	for name := range catExcluded {
		if catCounts[name] == 0 {
			catNames = append(catNames, name)
		}
	}
	sortStrings(catNames)

	totalGen := 0
	totalExcl := 0
	for _, cat := range catNames {
		gen := catCounts[cat]
		excl := catExcluded[cat]
		totalGen += gen
		totalExcl += excl
		fmt.Printf("%-20s %8d %8d %8d\n", cat, gen, excl, gen+excl)
	}
	fmt.Printf("%s\n", strings.Repeat("-", 50))
	fmt.Printf("%-20s %8d %8d %8d\n", "TOTAL", totalGen, totalExcl, totalGen+totalExcl)
}

func formatArgs(args []ArgDef) string {
	parts := make([]string, len(args))
	for i, a := range args {
		typeName := argTypeName(a.Type)
		if a.EnumType != "" {
			typeName = a.EnumType
		}
		parts[i] = fmt.Sprintf("%s:%s", a.Name, typeName)
	}
	return strings.Join(parts, ", ")
}

func argTypeName(t ArgType) string {
	switch t {
	case ArgTypeImage:
		return "image"
	case ArgTypeDouble:
		return "double"
	case ArgTypeInt:
		return "int"
	case ArgTypeBool:
		return "bool"
	case ArgTypeString:
		return "string"
	case ArgTypeEnum:
		return "enum"
	case ArgTypeFlags:
		return "flags"
	case ArgTypeArrayDouble:
		return "[]double"
	case ArgTypeArrayInt:
		return "[]int"
	case ArgTypeArrayImage:
		return "[]image"
	case ArgTypeBlob:
		return "blob"
	case ArgTypeInterpolate:
		return "interpolate"
	case ArgTypeSource:
		return "source"
	case ArgTypeTarget:
		return "target"
	default:
		return "unknown"
	}
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j-1] > s[j]; j-- {
			s[j-1], s[j] = s[j], s[j-1]
		}
	}
}
