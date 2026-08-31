package compile_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-ai-context/pkg/compile"
)

type DomainCatalogTestSuite struct {
	suite.Suite
}

func TestDomainCatalogTestSuite(t *testing.T) {
	suite.Run(t, new(DomainCatalogTestSuite))
}

func (s *DomainCatalogTestSuite) newRegistry() *registry.Registry {
	r := &registry.Registry{}
	return r
}

func (s *DomainCatalogTestSuite) TestBuildDomainCatalog_EmptyRegistry() {
	r := s.newRegistry()

	catalog, err := compile.BuildDomainCatalog(r)

	s.NoError(err)
	s.NotNil(catalog)
	s.Empty(catalog.Entities)
	s.Empty(catalog.Enums)
}

func (s *DomainCatalogTestSuite) TestBuildDomainCatalog_SingleModel() {
	r := s.newRegistry()
	r.SetModel("User", yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID":    {Type: yaml.ModelFieldTypeUUID},
			"Name":  {Type: yaml.ModelFieldTypeString},
			"Email": {Type: yaml.ModelFieldTypeString, Attributes: []string{"unique"}},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	})

	catalog, err := compile.BuildDomainCatalog(r)

	s.NoError(err)
	s.Len(catalog.Entities, 1)

	user := catalog.Entities["User"]
	s.Len(user.Fields, 3)

	s.Equal("UUID", user.Fields["ID"].Type)
	s.True(user.Fields["ID"].Primary)
	s.True(user.Fields["ID"].Required)

	s.Equal("String", user.Fields["Name"].Type)
	s.False(user.Fields["Name"].Primary)
	s.True(user.Fields["Name"].Required)

	s.Equal("String", user.Fields["Email"].Type)
	s.True(user.Fields["Email"].Unique)
	s.True(user.Fields["Email"].Required)
}

func (s *DomainCatalogTestSuite) TestBuildDomainCatalog_OptionalField() {
	r := s.newRegistry()
	r.SetModel("Task", yaml.Model{
		Name: "Task",
		Fields: map[string]yaml.ModelField{
			"ID":      {Type: yaml.ModelFieldTypeUUID},
			"DueDate": {Type: yaml.ModelFieldTypeDate, Attributes: []string{"optional"}},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	})

	catalog, err := compile.BuildDomainCatalog(r)

	s.NoError(err)
	task := catalog.Entities["Task"]
	s.True(task.Fields["ID"].Required)
	s.False(task.Fields["DueDate"].Required)
}

func (s *DomainCatalogTestSuite) TestBuildDomainCatalog_Relationships() {
	r := s.newRegistry()
	r.SetModel("Organization", yaml.Model{
		Name: "Organization",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeUUID},
			"Name": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Membership": {Type: "HasMany"},
			"Project":    {Type: "HasMany"},
		},
	})

	catalog, err := compile.BuildDomainCatalog(r)

	s.NoError(err)
	org := catalog.Entities["Organization"]
	s.Len(org.Relationships, 2)

	s.Equal("Membership", org.Relationships[0].Target)
	s.Equal("HasMany", org.Relationships[0].Type)
	s.Equal("Project", org.Relationships[1].Target)
	s.Equal("HasMany", org.Relationships[1].Type)
}

func (s *DomainCatalogTestSuite) TestBuildDomainCatalog_AliasedRelationship() {
	r := s.newRegistry()
	r.SetModel("Task", yaml.Model{
		Name: "Task",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeUUID},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Assignee": {Type: "ForOne", Aliased: "Membership"},
			"Creator":  {Type: "ForOne", Aliased: "Membership"},
		},
	})

	catalog, err := compile.BuildDomainCatalog(r)

	s.NoError(err)
	task := catalog.Entities["Task"]
	s.Len(task.Relationships, 2)

	s.Equal("Membership", task.Relationships[0].Target)
	s.Equal("ForOne", task.Relationships[0].Type)
	s.Equal("Assignee", task.Relationships[0].Alias)

	s.Equal("Membership", task.Relationships[1].Target)
	s.Equal("ForOne", task.Relationships[1].Type)
	s.Equal("Creator", task.Relationships[1].Alias)
}

func (s *DomainCatalogTestSuite) TestBuildDomainCatalog_PolymorphicRelationship() {
	r := s.newRegistry()
	r.SetModel("Comment", yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeUUID},
			"Text": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Task", "Project"},
			},
		},
	})

	catalog, err := compile.BuildDomainCatalog(r)

	s.NoError(err)
	comment := catalog.Entities["Comment"]
	s.Len(comment.Relationships, 1)

	s.Equal("Commentable", comment.Relationships[0].Target)
	s.Equal("ForOnePoly", comment.Relationships[0].Type)
	s.Equal([]string{"Task", "Project"}, comment.Relationships[0].For)
}

func (s *DomainCatalogTestSuite) TestBuildDomainCatalog_HasManyPolyWithThrough() {
	r := s.newRegistry()
	r.SetModel("Person", yaml.Model{
		Name: "Person",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Note": {
				Type:    "HasManyPoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
		},
	})

	catalog, err := compile.BuildDomainCatalog(r)

	s.NoError(err)
	person := catalog.Entities["Person"]
	s.Len(person.Relationships, 1)

	s.Equal("Comment", person.Relationships[0].Target)
	s.Equal("HasManyPoly", person.Relationships[0].Type)
	s.Equal("Note", person.Relationships[0].Alias)
	s.Equal("Commentable", person.Relationships[0].Through)
}

func (s *DomainCatalogTestSuite) TestBuildDomainCatalog_Enums() {
	r := s.newRegistry()
	r.SetEnum("TaskStatus", yaml.Enum{
		Name: "TaskStatus",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"backlog":     "Backlog",
			"todo":        "To Do",
			"in_progress": "In Progress",
			"done":        "Done",
		},
	})
	r.SetEnum("Priority", yaml.Enum{
		Name: "Priority",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"low":    "Low",
			"normal": "Normal",
			"high":   "High",
		},
	})

	catalog, err := compile.BuildDomainCatalog(r)

	s.NoError(err)
	s.Len(catalog.Enums, 2)

	s.Equal([]string{"high", "low", "normal"}, catalog.Enums["Priority"])
	s.Len(catalog.Enums["TaskStatus"], 4)
}

func (s *DomainCatalogTestSuite) TestBuildDomainCatalog_EnumFieldType() {
	r := s.newRegistry()
	r.SetEnum("Status", yaml.Enum{
		Name:    "Status",
		Type:    yaml.EnumTypeString,
		Entries: map[string]any{"active": "Active", "inactive": "Inactive"},
	})
	r.SetModel("User", yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID":     {Type: yaml.ModelFieldTypeUUID},
			"Status": {Type: "Status"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	})

	catalog, err := compile.BuildDomainCatalog(r)

	s.NoError(err)
	user := catalog.Entities["User"]
	s.Equal("Status", user.Fields["Status"].Type)
}

func (s *DomainCatalogTestSuite) TestBuildDomainCatalog_MultipleModels() {
	r := s.newRegistry()
	r.SetModel("Alpha", yaml.Model{
		Name:   "Alpha",
		Fields: map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	})
	r.SetModel("Beta", yaml.Model{
		Name:   "Beta",
		Fields: map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	})

	catalog, err := compile.BuildDomainCatalog(r)

	s.NoError(err)
	s.Len(catalog.Entities, 2)
	s.Contains(catalog.Entities, "Alpha")
	s.Contains(catalog.Entities, "Beta")
}

