package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	enumsPkg = "github.com/traust-security/traust-sdk/go/v1/enums"
)

func main() {
	enumsDir := flag.String("enums", "../../enums", "path to enums/ directory containing JSON files")
	outDir := flag.String("out", "./enums", "output directory for generated Go files")
	schemasDir := flag.String("schemas", "", "path to schemas/ directory (optional; enables type generation)")
	typesDir := flag.String("types", "", "output directory for generated type files")
	validateFile := flag.String("validate", "", "output path for generated schemas_gen.go (e.g. ./validate/schemas_gen.go)")
	storageSourceDir := flag.String("storage", "", "path to the contracts storage/v1 source directory")
	storageOutDir := flag.String("storage-out", "", "output directory for the storage package")
	fixtureDir := flag.String("storage-fixtures", "", "path to storage JSON fixtures")
	storageSamples := flag.String("storage-samples", "", "path to authored storage_samples.py")
	contractsRef := flag.String("contracts-ref", "unknown", "pinned traust-contracts git revision")
	flag.Parse()

	entries, err := os.ReadDir(*enumsDir)
	if err != nil {
		log.Fatalf("reading enums directory %s: %v", *enumsDir, err)
	}

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatalf("creating output directory %s: %v", *outDir, err)
	}

	var generated []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		src := filepath.Join(*enumsDir, entry.Name())
		if err := generateEnum(src, *outDir, *contractsRef); err != nil {
			log.Fatalf("generating from %s: %v", entry.Name(), err)
		}
		generated = append(generated, entry.Name())
	}

	fmt.Printf("generated %d enum files in %s\n", len(generated), *outDir)

	if err := generateNormalization(*enumsDir, *outDir, *contractsRef); err != nil {
		log.Fatalf("generating enum normalization metadata: %v", err)
	}

	if *schemasDir != "" && *typesDir != "" {
		schemas, err := LoadSchemas(*schemasDir)
		if err != nil {
			log.Fatalf("loading schemas: %v", err)
		}
		if err := ResolveRefs(schemas); err != nil {
			log.Fatalf("resolving schema refs: %v", err)
		}
		enumMappings, err := LoadEnumMappings(*enumsDir)
		if err != nil {
			log.Fatalf("loading enum mappings: %v", err)
		}
		if err := GenerateTypes(schemas, enumMappings, enumsPkg, *typesDir, *contractsRef); err != nil {
			log.Fatalf("generating types: %v", err)
		}
		if err := generateVersionFile(*typesDir, *contractsRef); err != nil {
			log.Fatalf("generating version.go: %v", err)
		}
		fmt.Printf("generated %d type files in %s\n", len(schemas), *typesDir)
	}

	if *schemasDir != "" && *validateFile != "" {
		if err := GenerateSchemasEmbed(*schemasDir, *validateFile, *contractsRef); err != nil {
			log.Fatalf("generating schemas embed: %v", err)
		}
		fmt.Printf("generated schemas map in %s\n", *validateFile)
	}

	if *storageSourceDir != "" && *storageOutDir != "" && *schemasDir != "" {
		if err := GenerateStorage(*storageSourceDir, *schemasDir, *storageOutDir, *fixtureDir, *storageSamples, *contractsRef); err != nil {
			log.Fatalf("generating storage assets: %v", err)
		}
		fmt.Printf("generated storage assets in %s\n", *storageOutDir)
	}
}

// generateVersionFile records the pinned traust-contracts commit the SDK was generated from.
func generateVersionFile(typesDir, contractsRef string) error {
	content := fmt.Sprintf(
		"// Code generated from traust-contracts %s. DO NOT EDIT.\n\npackage types\n\n"+
			"// ContractsRef is the pinned traust-contracts commit.\nconst ContractsRef = %q\n",
		contractsRef, contractsRef,
	)
	path := filepath.Join(typesDir, "version.go")
	return formatAndWrite(path, []byte(content))
}
