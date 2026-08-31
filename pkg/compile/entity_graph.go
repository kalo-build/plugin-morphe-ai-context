package compile

import (
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/morphe-go/pkg/yamlops"
)

func BuildEntityGraph(r *registry.Registry) (*EntityGraph, error) {
	graph := &EntityGraph{}

	allModels := r.GetAllModels()
	modelNames := sortedKeys(allModels)

	for _, modelName := range modelNames {
		model := allModels[modelName]
		nodeType := classifyModelType(model)
		graph.Nodes = append(graph.Nodes, GraphNode{
			ID:     modelName,
			Type:   nodeType,
			Fields: len(model.Fields),
		})
	}

	for _, modelName := range modelNames {
		model := allModels[modelName]
		edges := buildGraphEdges(modelName, model.Related)
		graph.Edges = append(graph.Edges, edges...)
	}

	return graph, nil
}

func classifyModelType(model yaml.Model) string {
	hasOnlyFKFields := true
	for _, field := range model.Fields {
		if field.Type != yaml.ModelFieldTypeAutoIncrement && field.Type != yaml.ModelFieldTypeUUID {
			hasOnlyFKFields = false
			break
		}
	}
	if hasOnlyFKFields && len(model.Fields) <= 2 && len(model.Related) >= 2 {
		return "join"
	}
	return "model"
}

func buildGraphEdges(fromModel string, relations map[string]yaml.ModelRelation) []GraphEdge {
	if len(relations) == 0 {
		return nil
	}

	relationNames := sortedKeys(relations)
	var edges []GraphEdge

	for _, name := range relationNames {
		rel := relations[name]
		edge := GraphEdge{
			From: fromModel,
			To:   yamlops.GetRelationTargetName(name, rel.Aliased),
			Type: rel.Type,
		}
		if yamlops.IsRelationAliased(rel.Aliased) {
			edge.Alias = name
		}
		if rel.Through != "" {
			edge.Through = rel.Through
		}
		if len(rel.For) > 0 {
			edge.For = rel.For
		}
		edges = append(edges, edge)
	}

	return edges
}
