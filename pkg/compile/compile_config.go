package compile

import (
	"github.com/kalo-build/morphe-go/pkg/registry"
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
)

type CompileConfig struct {
	rcfg.MorpheLoadRegistryConfig

	RegistryHooks registry.LoadMorpheRegistryHooks

	OutputDirPath string
}

func DefaultCompileConfig(registryPath, outputDirPath string) CompileConfig {
	return CompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      registryPath + "/enums",
			RegistryModelsDirPath:     registryPath + "/models",
			RegistryStructuresDirPath: registryPath + "/structures",
			RegistryEntitiesDirPath:   registryPath + "/entities",
		},
		OutputDirPath: outputDirPath,
	}
}
