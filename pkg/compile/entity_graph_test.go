package compile_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-ai-context/pkg/compile"
)

type EntityGraphTestSuite struct {
	suite.Suite
}

func TestEntityGraphTestSuite(t *testing.T) {
	suite.Run(t, new(EntityGraphTestSuite))
}

func (s *EntityGraphTestSuite) newRegistry() *registry.Registry {
	r := &registry.Registry{}
	return r
}

func (s *EntityGraphTestSuite) TestBuildEntityGraph_EmptyRegistry() {
	r := s.newRegistry()

	graph, err := compile.BuildEntityGraph(r)

	s.NoError(err)
	s.NotNil(graph)
	s.Empty(graph.Nodes)
	s.Empty(graph.Edges)
}

func (s *EntityGraphTestSuite) TestBuildEntityGraph_SingleModelNode() {
	r := s.newRegistry()
	r.SetModel("User", yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID":    {Type: yaml.ModelFieldTypeUUID},
			"Name":  {Type: yaml.ModelFieldTypeString},
			"Email": {Type: yaml.ModelFieldTypeString},
		},
	})

	graph, err := compile.BuildEntityGraph(r)

	s.NoError(err)
	s.Len(graph.Nodes, 1)
	s.Equal("User", graph.Nodes[0].ID)
	s.Equal("model", graph.Nodes[0].Type)
	s.Equal(3, graph.Nodes[0].Fields)
}

func (s *EntityGraphTestSuite) TestBuildEntityGraph_ForOneEdge() {
	r := s.newRegistry()
	r.SetModel("Task", yaml.Model{
		Name:   "Task",
		Fields: map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Related: map[string]yaml.ModelRelation{
			"Project": {Type: "ForOne"},
		},
	})
	r.SetModel("Project", yaml.Model{
		Name:   "Project",
		Fields: map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
	})

	graph, err := compile.BuildEntityGraph(r)

	s.NoError(err)
	s.Len(graph.Nodes, 2)
	s.Len(graph.Edges, 1)

	s.Equal("Task", graph.Edges[0].From)
	s.Equal("Project", graph.Edges[0].To)
	s.Equal("ForOne", graph.Edges[0].Type)
	s.Empty(graph.Edges[0].Alias)
}

func (s *EntityGraphTestSuite) TestBuildEntityGraph_HasManyEdge() {
	r := s.newRegistry()
	r.SetModel("Project", yaml.Model{
		Name:   "Project",
		Fields: map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Related: map[string]yaml.ModelRelation{
			"Task": {Type: "HasMany"},
		},
	})
	r.SetModel("Task", yaml.Model{
		Name:   "Task",
		Fields: map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
	})

	graph, err := compile.BuildEntityGraph(r)

	s.NoError(err)

	var projectEdge *compile.GraphEdge
	for i := range graph.Edges {
		if graph.Edges[i].From == "Project" {
			projectEdge = &graph.Edges[i]
		}
	}
	s.NotNil(projectEdge)
	s.Equal("Task", projectEdge.To)
	s.Equal("HasMany", projectEdge.Type)
}

func (s *EntityGraphTestSuite) TestBuildEntityGraph_AliasedEdge() {
	r := s.newRegistry()
	r.SetModel("Task", yaml.Model{
		Name:   "Task",
		Fields: map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Related: map[string]yaml.ModelRelation{
			"Assignee": {Type: "ForOne", Aliased: "Membership"},
			"Creator":  {Type: "ForOne", Aliased: "Membership"},
		},
	})
	r.SetModel("Membership", yaml.Model{
		Name:   "Membership",
		Fields: map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
	})

	graph, err := compile.BuildEntityGraph(r)

	s.NoError(err)
	s.Len(graph.Edges, 2)

	s.Equal("Membership", graph.Edges[0].To)
	s.Equal("Assignee", graph.Edges[0].Alias)
	s.Equal("Membership", graph.Edges[1].To)
	s.Equal("Creator", graph.Edges[1].Alias)
}

func (s *EntityGraphTestSuite) TestBuildEntityGraph_PolymorphicForEdge() {
	r := s.newRegistry()
	r.SetModel("Comment", yaml.Model{
		Name:   "Comment",
		Fields: map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Task", "Project"},
			},
		},
	})

	graph, err := compile.BuildEntityGraph(r)

	s.NoError(err)
	s.Len(graph.Edges, 1)

	edge := graph.Edges[0]
	s.Equal("Comment", edge.From)
	s.Equal("Commentable", edge.To)
	s.Equal("ForOnePoly", edge.Type)
	s.Equal([]string{"Task", "Project"}, edge.For)
}

func (s *EntityGraphTestSuite) TestBuildEntityGraph_PolymorphicHasWithThrough() {
	r := s.newRegistry()
	r.SetModel("Person", yaml.Model{
		Name:   "Person",
		Fields: map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeAutoIncrement}},
		Related: map[string]yaml.ModelRelation{
			"Note": {
				Type:    "HasManyPoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
		},
	})

	graph, err := compile.BuildEntityGraph(r)

	s.NoError(err)
	s.Len(graph.Edges, 1)

	edge := graph.Edges[0]
	s.Equal("Person", edge.From)
	s.Equal("Comment", edge.To)
	s.Equal("HasManyPoly", edge.Type)
	s.Equal("Note", edge.Alias)
	s.Equal("Commentable", edge.Through)
}

func (s *EntityGraphTestSuite) TestBuildEntityGraph_NodesSorted() {
	r := s.newRegistry()
	r.SetModel("Zebra", yaml.Model{
		Name:   "Zebra",
		Fields: map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
	})
	r.SetModel("Alpha", yaml.Model{
		Name:   "Alpha",
		Fields: map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
	})
	r.SetModel("Middle", yaml.Model{
		Name:   "Middle",
		Fields: map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
	})

	graph, err := compile.BuildEntityGraph(r)

	s.NoError(err)
	s.Len(graph.Nodes, 3)
	s.Equal("Alpha", graph.Nodes[0].ID)
	s.Equal("Middle", graph.Nodes[1].ID)
	s.Equal("Zebra", graph.Nodes[2].ID)
}

func (s *EntityGraphTestSuite) TestBuildEntityGraph_EdgesSortedByModel() {
	r := s.newRegistry()
	r.SetModel("B", yaml.Model{
		Name:   "B",
		Fields: map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Related: map[string]yaml.ModelRelation{
			"A": {Type: "ForOne"},
		},
	})
	r.SetModel("A", yaml.Model{
		Name:   "A",
		Fields: map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Related: map[string]yaml.ModelRelation{
			"B": {Type: "HasMany"},
		},
	})

	graph, err := compile.BuildEntityGraph(r)

	s.NoError(err)
	s.Len(graph.Edges, 2)
	s.Equal("A", graph.Edges[0].From)
	s.Equal("B", graph.Edges[1].From)
}

func (s *EntityGraphTestSuite) TestBuildEntityGraph_NoRelationsNoEdges() {
	r := s.newRegistry()
	r.SetModel("Standalone", yaml.Model{
		Name: "Standalone",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeUUID},
			"Name": {Type: yaml.ModelFieldTypeString},
		},
	})

	graph, err := compile.BuildEntityGraph(r)

	s.NoError(err)
	s.Len(graph.Nodes, 1)
	s.Empty(graph.Edges)
}
