package compile

// DomainCatalog is the top-level structure for domain_catalog.yaml.
type DomainCatalog struct {
	Entities map[string]CatalogEntity `yaml:"entities" json:"entities"`
	Enums    map[string][]string      `yaml:"enums" json:"enums"`
}

// CatalogEntity describes a single domain entity in the catalog.
type CatalogEntity struct {
	Fields        map[string]CatalogField        `yaml:"fields" json:"fields"`
	Relationships []CatalogRelationship          `yaml:"relationships,omitempty" json:"relationships,omitempty"`
}

// CatalogField describes a single field on an entity.
type CatalogField struct {
	Type     string `yaml:"type" json:"type"`
	Primary  bool   `yaml:"primary,omitempty" json:"primary,omitempty"`
	Required bool   `yaml:"required,omitempty" json:"required,omitempty"`
	Unique   bool   `yaml:"unique,omitempty" json:"unique,omitempty"`
}

// CatalogRelationship describes one relationship edge from an entity.
type CatalogRelationship struct {
	Target   string `yaml:"target" json:"target"`
	Type     string `yaml:"type" json:"type"`
	Alias    string `yaml:"alias,omitempty" json:"alias,omitempty"`
	Optional bool   `yaml:"optional,omitempty" json:"optional,omitempty"`
	Through  string `yaml:"through,omitempty" json:"through,omitempty"`
	For      []string `yaml:"for,omitempty" json:"for,omitempty"`
}

// EntityGraph is the top-level structure for entity_graph.json.
type EntityGraph struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

// GraphNode is a single node in the entity relationship graph.
type GraphNode struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Fields int    `json:"fields"`
}

// GraphEdge is a single directed edge in the entity relationship graph.
type GraphEdge struct {
	From     string   `json:"from"`
	To       string   `json:"to"`
	Type     string   `json:"type"`
	Alias    string   `json:"alias,omitempty"`
	Optional bool     `json:"optional,omitempty"`
	Through  string   `json:"through,omitempty"`
	For      []string `json:"for,omitempty"`
}
