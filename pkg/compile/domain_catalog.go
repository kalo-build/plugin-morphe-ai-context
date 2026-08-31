package compile

import (
	"sort"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/morphe-go/pkg/yamlops"
)

func BuildDomainCatalog(r *registry.Registry) (*DomainCatalog, error) {
	catalog := &DomainCatalog{
		Entities: make(map[string]CatalogEntity),
		Enums:    make(map[string][]string),
	}

	allModels := r.GetAllModels()
	modelNames := sortedKeys(allModels)
	for _, modelName := range modelNames {
		model := allModels[modelName]
		entity, err := buildCatalogEntity(model)
		if err != nil {
			return nil, err
		}
		catalog.Entities[modelName] = entity
	}

	allEnums := r.GetAllEnums()
	enumNames := sortedKeys(allEnums)
	for _, enumName := range enumNames {
		enum := allEnums[enumName]
		catalog.Enums[enumName] = buildCatalogEnumValues(enum)
	}

	return catalog, nil
}

func buildCatalogEntity(model yaml.Model) (CatalogEntity, error) {
	entity := CatalogEntity{
		Fields: make(map[string]CatalogField),
	}

	primaryFields := getPrimaryIdentifierFields(model)

	fieldNames := sortedKeys(model.Fields)
	for _, fieldName := range fieldNames {
		field := model.Fields[fieldName]
		catalogField := CatalogField{
			Type:     string(field.Type),
			Required: !hasAttribute(field.Attributes, "optional"),
		}
		if _, isPrimary := primaryFields[fieldName]; isPrimary {
			catalogField.Primary = true
		}
		if hasAttribute(field.Attributes, "unique") {
			catalogField.Unique = true
		}
		entity.Fields[fieldName] = catalogField
	}

	entity.Relationships = buildCatalogRelationships(model.Related)

	return entity, nil
}

func buildCatalogRelationships(relations map[string]yaml.ModelRelation) []CatalogRelationship {
	if len(relations) == 0 {
		return nil
	}

	relationNames := sortedKeys(relations)
	var result []CatalogRelationship

	for _, name := range relationNames {
		rel := relations[name]
		cr := CatalogRelationship{
			Target: yamlops.GetRelationTargetName(name, rel.Aliased),
			Type:   rel.Type,
		}
		if yamlops.IsRelationAliased(rel.Aliased) {
			cr.Alias = name
		}
		if hasAttribute([]string{}, "optional") {
			cr.Optional = true
		}
		if rel.Through != "" {
			cr.Through = rel.Through
		}
		if len(rel.For) > 0 {
			cr.For = rel.For
		}
		result = append(result, cr)
	}

	return result
}

func buildCatalogEnumValues(enum yaml.Enum) []string {
	keys := make([]string, 0, len(enum.Entries))
	for k := range enum.Entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func getPrimaryIdentifierFields(model yaml.Model) map[string]struct{} {
	result := make(map[string]struct{})
	primary, exists := model.Identifiers["primary"]
	if !exists {
		return result
	}
	for _, f := range primary.Fields {
		result[f] = struct{}{}
	}
	return result
}

func hasAttribute(attributes []string, attr string) bool {
	for _, a := range attributes {
		if a == attr {
			return true
		}
	}
	return false
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
