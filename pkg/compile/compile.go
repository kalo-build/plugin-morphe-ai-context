package compile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"gopkg.in/yaml.v3"
)

func MorpheToAIContext(config CompileConfig) error {
	r, err := registry.LoadMorpheRegistry(config.RegistryHooks, config.MorpheLoadRegistryConfig)
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	catalog, err := BuildDomainCatalog(r)
	if err != nil {
		return fmt.Errorf("failed to build domain catalog: %w", err)
	}

	graph, err := BuildEntityGraph(r)
	if err != nil {
		return fmt.Errorf("failed to build entity graph: %w", err)
	}

	if err := os.MkdirAll(config.OutputDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := writeDomainCatalog(config.OutputDirPath, catalog); err != nil {
		return err
	}

	if err := writeEntityGraph(config.OutputDirPath, graph); err != nil {
		return err
	}

	return nil
}

func writeDomainCatalog(outputDir string, catalog *DomainCatalog) error {
	data, err := yaml.Marshal(catalog)
	if err != nil {
		return fmt.Errorf("failed to marshal domain catalog: %w", err)
	}
	path := filepath.Join(outputDir, "domain_catalog.yaml")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write domain catalog: %w", err)
	}
	return nil
}

func writeEntityGraph(outputDir string, graph *EntityGraph) error {
	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal entity graph: %w", err)
	}
	data = append(data, '\n')
	path := filepath.Join(outputDir, "entity_graph.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write entity graph: %w", err)
	}
	return nil
}
