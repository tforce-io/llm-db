# Project Context

Part of the project memory layer: human-readable memory of what this project is and what has been confirmed.

The reading and writing rules for this file and the rest of `.attunic/` are defined in [attunic-protocol.md](../.agents/attunic-protocol.md).

## Abstract

A newcomer should be able to tell what the software does without knowing what it is built with. One paragraph that describes the domain, not the technology, in simple terms:
<!-- Project-specific / Abstract -->

LLM DB is a curated source of truth about Large Language Models: their capabilities, costs, limits, modalities, and deployments across providers; kept in structured, machine-readable form. The project also ships LocalAI YAML configurations and a small CLI that exports the catalog to the formats other tools consume.

## Technical Specifications

A brief overview of the ingredients used for developing this project. These are things developers should know when developing the project.

Typical ingredients include, but are not limited to: primary language, frameworks, architecture styles, environment prerequisites, libraries, test runners, tools, databases...

<!-- Project-specific / Technical Specifications -->

- Primary language: Go.
- Data formats: JSON, YAML.
- CLI framework: `github.com/spf13/cobra`.
- Logging: `github.com/rs/zerolog`.
- Build system: `go`.
- Test runner: `go test`.
- Linter: `gofmt`.
- Supported output format: Bifrost Datasheet, OpenCode v1 Configuration.

## Layout

Layout of this project; [project-structure.md](../.agents/project-structure.md) holds the language's standard map, kept for reference.

List of important folders, files and their roles in the project:

- They visualize the actual file system where source files live.
- If a top-level file/folder is missing from the list, prompt the user to update it.
<!-- Project-specific / Layout -->

```
├── .agents/            # project AI agent instructions (attunic-protocol, coding-conventions, development-guideline, license-header, memory-template, project-structure, pull-request-guideline)
├── .attunic/           # project AI agent memory (context, memory)
├── json/               # data in JSON format
│   ├── refs/           # snapshot models data from popular providers
│   │   ├── bifrost/         # reference models data from Bifrost
│   │   ├── kilocode/        # reference models data from KiloCode
│   │   └── opencode/        # reference models data from OpenCode
│   ├── samples/        # sample data for testing
│   ├── models.json         # project models data
│   ├── models.schema.json  # project JSON schema for models data
│   ├── providers.json      # project providers data
│   └── providers.schema.json # project JSON schema for providers data
├── yaml/               # data in YAML format
│   └── LocalAI/        # sample LocalAI deployment configurations
├── schema/             # Go types mirroring each JSON format
│   ├── llmdb/              # project models and providers format
│   ├── bifrost/            # Bifrost datasheet format
│   ├── opencode/           # OpenCode v1 config format
│   ├── opencode_cloud/     # OpenCode cloud API format
│   └── dotenv/         # minimal .env parser
├── engine/             # cobra wiring, CLI modules (export, format)
├── common/             # shared scalar types
├── config/             # global logger setup
├── exported/           # generated exports
├── .copier-answers.yml # project Copier answers
├── .editorconfig       # EditorConfig configuration
├── .gitignore          # git ignore patterns
├── LICENSE             # MIT license text
├── README.md           # project overview
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
└── main.go             # program entrypoint
```

## Layer Dependencies

List of layers and their relationships:

- Declare the innermost layer first.
- A layer may depend on itself and on any layer declared before it, and on nothing declared after it.
- If a layer appears in the table below, the dependency graph must be enforced.

| Layer | Path | Depends on | Must never contain |
|-------|------|------------|--------------------|
<!-- Project-specific / Layer Dependencies -->

### Deliberate Exceptions

List of documented exceptions that cannot follow the architecture rules above:

- Any place the rules above are knowingly broken, with the reason and the condition for fixing it.
- An undocumented exception is a defect; a documented one is a decision.

| What | Why | Revisit-if |
|------|-----|------------|
<!-- Project-specific / Deliberate Exceptions -->

## Supporting Tools

All the supporting tools declared here are opinionated for the sake of consistent results across development environments.

List each required tool/library under its own heading, followed by the commands for each use case.
Common use cases include: Build, Test, Lint, Format, Install, ...

<!-- Project-specific / Supporting Tools -->

### Build

- `go build -o llm-db.exe .` — build the CLI from the root package (Windows).
- `go build -o llm-db .` — build the CLI from the root package (Unix).

### Test

- `go test ./...` — run the whole test suite.
- `go test -race ./...` — run the suite with the race detector; required for changes that touch concurrency (per `.agents/coding-conventions.md`).

### Lint

- `gofmt -l .` — check formatting; must report no files.
- `gofmt -w .` — apply canonical formatting.
- `go vet ./...` — go static analysis.

## Confirmed Facts

Facts about the project which have been confirmed; they outrank anything derived or guessed.

| Fact | Reason |
|------|--------|
<!-- Project-specific / Confirmed Facts -->

## Project-specific

<!-- Project-specific -->
