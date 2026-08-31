package compile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/go-util/assertfile"
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-ai-context/internal/testutils"
	"github.com/kalo-build/plugin-morphe-ai-context/pkg/compile"
)

type CompileTestSuite struct {
	assertfile.FileSuite

	TestDirPath            string
	TestGroundTruthDirPath string

	ModelsDirPath     string
	EnumsDirPath      string
	StructuresDirPath string
	EntitiesDirPath   string
}

func TestCompileTestSuite(t *testing.T) {
	suite.Run(t, new(CompileTestSuite))
}

func (s *CompileTestSuite) SetupTest() {
	s.TestDirPath = testutils.GetTestDirPath()
	s.TestGroundTruthDirPath = filepath.Join(s.TestDirPath, "ground-truth", "compile-minimal")

	s.ModelsDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "models")
	s.EnumsDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "enums")
	s.StructuresDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "structures")
	s.EntitiesDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "entities")
}

func (s *CompileTestSuite) TearDownTest() {
	s.TestDirPath = ""
}

func (s *CompileTestSuite) TestMorpheToAIContext() {
	workingDirPath := filepath.Join(s.TestDirPath, "working")
	s.Nil(os.Mkdir(workingDirPath, 0755))
	defer os.RemoveAll(workingDirPath)

	config := compile.CompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      s.EnumsDirPath,
			RegistryModelsDirPath:     s.ModelsDirPath,
			RegistryStructuresDirPath: s.StructuresDirPath,
			RegistryEntitiesDirPath:   s.EntitiesDirPath,
		},
		OutputDirPath: workingDirPath,
	}

	compileErr := compile.MorpheToAIContext(config)

	s.NoError(compileErr)

	catalogPath := filepath.Join(workingDirPath, "domain_catalog.yaml")
	gtCatalogPath := filepath.Join(s.TestGroundTruthDirPath, "domain_catalog.yaml")
	s.FileExists(catalogPath)
	s.FileEquals(catalogPath, gtCatalogPath)

	graphPath := filepath.Join(workingDirPath, "entity_graph.json")
	gtGraphPath := filepath.Join(s.TestGroundTruthDirPath, "entity_graph.json")
	s.FileExists(graphPath)
	s.FileEquals(graphPath, gtGraphPath)
}
