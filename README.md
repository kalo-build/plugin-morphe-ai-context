# plugin-morphe-ai-context

Kalo plugin that reads a Morphe registry and emits structured AI context files for LLM agent consumption.

## Output Files

| File | Format | Contents |
|------|--------|----------|
| `domain_catalog.yaml` | YAML | All models with fields, types, constraints, relationships, and enums |
| `entity_graph.json` | JSON | Node-edge relationship graph for structured traversal |

## Input / Output Formats

- **Input:** `KA:MO1:YAML1` — Morphe YAML registry (`.mod`, `.enum`, `.str` files)
- **Output:** `KA:AC1:YAML1` — AI context files

## domain_catalog.yaml

Lists every model as an entity entry with:
- **Fields** — name, type, primary/required/unique constraints
- **Relationships** — target, type (ForOne/HasMany/HasOnePoly/etc.), alias, through, for
- **Enums** — name with sorted entry keys

```yaml
entities:
  Task:
    fields:
      ID: { type: UUID, primary: true, required: true }
      Title: { type: String, required: true }
      Status: { type: TaskStatus, required: true }
    relationships:
      - { target: Project, type: ForOne }
      - { target: Membership, type: ForOne, alias: Assignee }
enums:
  TaskStatus: [backlog, done, in_progress, todo]
```

## entity_graph.json

Provides a graph representation with nodes (models) and directed edges (relationships):

```json
{
  "nodes": [
    { "id": "Task", "type": "model", "fields": 5 }
  ],
  "edges": [
    { "from": "Task", "to": "Project", "type": "ForOne" },
    { "from": "Task", "to": "Membership", "type": "ForOne", "alias": "Assignee" }
  ]
}
```

## Relationship Types

| Type | Description |
|------|-------------|
| `ForOne` | Belongs-to (FK to one parent) |
| `HasOne` | Has one child |
| `HasMany` | Has many children |
| `ForOnePoly` | Polymorphic belongs-to (with `for` listing target models) |
| `HasOnePoly` | Polymorphic has-one (with `through` and optional `alias`) |
| `HasManyPoly` | Polymorphic has-many (with `through` and optional `alias`) |

## Configuration

```json
{
  "inputPath": "./morphe/registry",
  "outputPath": "./ai-context",
  "verbose": false
}
```

No additional config fields required. The plugin reads all models and enums from the registry and emits context files.

## Pipeline Usage

```yaml
stores:
  morphe-registry:
    type: filesystem
    path: morphe/registry

  ai-context:
    type: filesystem
    path: ai-context

pipelines:
  compile:
    stages:
      - name: "ai-context"
        steps:
          - "plugin: @kalo-build/plugin-morphe-ai-context"
```

## Build

```bash
# Native binary
go build -o plugin-morphe-ai-context ./cmd/plugin

# WASM (for Kalo CLI)
GOOS=wasip1 GOARCH=wasm go build -o plugin-morphe-ai-context.wasm ./cmd/plugin
```

## Test

```bash
go test ./... -v
```
