# Spec: Virgil

## Dogma rector

Virgil implementa la visión de Uncle Bob (julio 2026): no revisas código
del agente — verificas métricas, trazabilidad y cumplimiento de la
metodología desde un nivel superior. Seis principios:

1. Metodología e2e (idea → operación, operación es facade opcional)
2. Verificas trazabilidad Y fuerza (binding + mutation/CRAP/complexity)
3. Gestionas desde un nivel superior (dashboard de salud, no code review)
4. El agente opera bajo constraint, no bajo confianza (hooks + gates)
5. Un handoff, ejecución paralela con semántica de coordinación
6. Gates determinísticos, no advisory (validación mecánica en transiciones)

---

## Requisitos funcionales

### RF-01 — Inicialización del proyecto (`virgil init`)

El CLI configura un proyecto destino para usar Virgil.

**AC-01.1** — Genera gobernanza mínima
```
GIVEN un directorio con git inicializado
WHEN el usuario ejecuta `virgil init`
THEN Virgil genera un AGENTS.md de 40-80 líneas con axiomas,
     build/test commands, y bootstrap del SM mental model
```

**AC-01.2** — Instala skills como archivos locales
```
GIVEN `virgil init` ejecutándose
WHEN el compilador modular procesa los 31 docs de metodología
THEN genera un skill por fase/concepto en `.virgil/skills/`,
     cada uno con frontmatter (nombre, descripción, trigger)
     y contenido auto-contenido
```

**AC-01.3** — Registra git hooks
```
GIVEN `virgil init` ejecutándose
WHEN la fase de hooks se ejecuta
THEN instala hooks en `.git/hooks/` (o via husky/lefthook si existe):
     post-commit, post-merge, post-rewrite, post-checkout
     que delegan a `virgil refresh --diff`
```

**AC-01.4** — Inicializa el binding layer vacío
```
GIVEN `virgil init` ejecutándose
WHEN la fase de binding se ejecuta
THEN crea `.virgil/bindings.db` (SQLite) con schema vacío
     y lo agrega a .gitignore
```

**AC-01.5** — Genera configuración MCP
```
GIVEN `virgil init` ejecutándose
WHEN la fase de MCP se ejecuta
THEN genera `.virgil/mcp.json` con la configuración del MCP server
     para que el agente pueda consultar el binding layer
```

---

### RF-02 — Flujo greenfield (`/virgil-idea`)

El skill guía al agente fase por fase desde una idea hasta un handoff
aprobado. CERO código hasta que el handoff existe.

**AC-02.1** — Activación del skill
```
GIVEN un proyecto con Virgil inicializado
WHEN el MIM dice "/virgil-idea" o equivalente
THEN el agente carga el skill de idea, asume rol SM,
     y guía la conversación a producir idea.md
```

**AC-02.2** — Bloqueo de código prematuro
```
GIVEN el flujo /virgil-idea activo sin handoff aprobado
WHEN el agente o MIM intenta escribir código de implementación
THEN Virgil bloquea la acción (via hook o skill instruction)
     con mensaje: "handoff no existe o no está aprobado"
```

**AC-02.3** — Transición de fases
```
GIVEN idea.md aprobado por el MIM
WHEN el agente consulta el siguiente paso
THEN el skill de idea descarga y carga el skill de spec,
     pasando el contexto de idea.md como input
```

**AC-02.4** — Pipeline completo
```
GIVEN el flujo /virgil-idea activo
WHEN el MIM aprueba cada artefacto en secuencia
THEN el pipeline produce: idea.md → spec.md → design.md →
     tasks.md → handoff.md, cada uno siguiendo el schema
     definido en la metodología (29148, 42010, 21502, 15289)
```

**AC-02.5** — Binding layer se construye durante ejecución
```
GIVEN handoff.md aprobado y ejecución iniciada (Red/Green/Refactor)
WHEN el agente implementa un AC
THEN el agente declara el binding (AC → código) via MCP,
     y el post-commit hook infiere bindings adicionales del diff
```

---

### RF-03 — Flujo brownfield (`/virgil-takeover`)

El skill ejecuta arqueología profunda del codebase existente, deriva
docs funcionales equivalentes, y establece el binding layer baseline.

**AC-03.1** — Crawl exhaustivo
```
GIVEN un proyecto existente con código pero sin docs funcionales
WHEN el MIM dice "/virgil-takeover"
THEN Virgil ejecuta `virgil scan --full`:
     - Tree-Sitter AST parsing de todos los archivos fuente
     - Extracción de símbolos (funciones, clases, módulos, configs)
     - Construcción del grafo de dependencias internas
     - Almacenamiento en `.virgil/bindings.db`
```

**AC-03.2** — Derivación de docs funcionales
```
GIVEN el crawl exhaustivo completado
WHEN la fase de declaración funcional se ejecuta
THEN el agente produce docs equivalentes:
     - README/docs existentes → idea.md equiv
     - Tests existentes → spec.md parcial (ACs inferidos)
     - Arquitectura detectada → design.md equiv
     con gaps identificados (código sin doc, doc sin código)
```

**AC-03.3** — Scoring y plan
```
GIVEN docs funcionales derivados y gaps identificados
WHEN la fase de scoring se ejecuta
THEN aplica fastForward con overrides brownfield,
     determina tier de activación,
     y produce echo bootstrap plan
```

**AC-03.4** — Flujos de datos, no inventario
```
GIVEN el crawl ejecutándose
WHEN se analizan módulos del codebase
THEN el análisis produce flujos de datos end-to-end
     (qué dato entra → cómo se transforma → dónde sale),
     NO solo un listado de archivos y sus exports
```

---

### RF-04 — Binding layer (grafo de trazabilidad)

El vínculo lógico real entre requerimientos y código. Implementado como
SQLite knowledge graph con Tree-Sitter para parsing.

**AC-04.1** — Modelo de datos
```
GIVEN el binding layer inicializado
WHEN se consulta el schema
THEN contiene las entidades:
     REQUIREMENT (id, type, source_doc, description, status)
     CODE_ARTIFACT (id, file_path, symbol, type, last_commit)
     TEST_ARTIFACT (id, file_path, describes, status)
     BINDING (id, requirement_id, artifact_id, relation, confidence, created_by)
```

**AC-04.2** — Declaración de bindings por el agente
```
GIVEN el agente implementó AC-N
WHEN el agente ejecuta declare(AC-N → file:symbol)
THEN se crea un BINDING con confidence: declared, created_by: agent
```

**AC-04.3** — Inferencia de bindings por hook
```
GIVEN un commit realizado en el proyecto
WHEN el post-commit hook ejecuta `virgil refresh --diff`
THEN analiza el diff contra docs funcionales,
     infiere bindings nuevos (confidence: inferred, created_by: hook),
     y marca como STALE bindings cuyo código cambió
```

**AC-04.4** — Verificación on demand
```
GIVEN un AC con bindings existentes
WHEN el MIM o agente ejecuta `virgil verify AC-N`
THEN Virgil hace verificación estructural (Tree-Sitter es sintáctico,
     NO puede juzgar si el código cumple el given/when/then del AC —
     esa lectura semántica queda para el MIM o el agente):
     - Existencia de símbolo: los símbolos vinculados siguen presentes
       en el codebase
     - Staleness: los archivos/símbolos vinculados cambiaron desde el
       último scan
     - Test linkage: existe un test vinculado a este AC
     - Mutation strength (si hay herramientas de mutación disponibles):
       el test vinculado efectivamente guarda el AC
     - Retorna: bound | stale | unbound
     - Actualiza confidence a: verified (verificación estructural,
       no semántica)
```

**AC-04.5** — Detección de desincronización
```
GIVEN spec.md actualizado (un AC cambió su definición)
WHEN el post-commit hook detecta cambio en docs funcionales
THEN marca todos los bindings del AC modificado como STALE,
     y emite warning al agente en próxima consulta
```

**AC-04.6** — Granularidad a nivel símbolo
```
GIVEN un binding entre AC y código
WHEN se registra el binding
THEN la granularidad es a nivel símbolo:
     AC-3 → PdfService.generate() (no solo pdf.service.ts)
     usando Tree-Sitter para resolver funciones, clases y módulos
```

**AC-04.7** — Consulta del grafo via MCP
```
GIVEN el agente trabaja en el proyecto
WHEN el agente necesita contexto de un AC
THEN consulta el MCP server que expone:
     - trace(requirement_id) → bindings + código + tests
     - impact(file_path) → requirements afectados
     - coverage() → % de ACs con binding verified
     - stale() → bindings marcados como STALE
```

---

### RF-05 — Refresh incremental via git hooks

Después del scan inicial, el grafo se mantiene por diff incremental.
Costo proporcional al cambio, no al tamaño del codebase.

**AC-05.1** — post-merge hook
```
GIVEN git pull o git merge completado
WHEN el post-merge hook dispara
THEN `virgil refresh --diff` analiza solo los archivos cambiados,
     actualiza bindings afectados, marca STALE los rotos
```

**AC-05.2** — post-commit hook
```
GIVEN un commit realizado
WHEN el post-commit hook dispara
THEN analiza el diff del commit contra docs funcionales,
     infiere bindings nuevos (confidence: inferred)
```

**AC-05.3** — post-rewrite hook
```
GIVEN un rebase o amend completado
WHEN el post-rewrite hook dispara
THEN re-verifica bindings de commits reescritos,
     actualiza last_commit en CODE_ARTIFACT
```

**AC-05.4** — post-checkout hook
```
GIVEN cambio de branch (git switch/checkout)
WHEN el post-checkout hook dispara
THEN marca bindings como potentially stale
     si la branch difiere de la baseline
```

**AC-05.5** — Performance proporcional al delta
```
GIVEN un codebase de 100K líneas donde se cambian 20 líneas
WHEN cualquier hook dispara refresh
THEN el tiempo de procesamiento es proporcional a las 20 líneas,
     NO a las 100K líneas del codebase completo
```

---

### RF-06 — Delivery modular (skills + hooks + MCP)

Reemplaza el AGENTS.md monolítico de 1,061 líneas con mecanismos
nativos del agente.

**AC-06.1** — AGENTS.md mínimo (40-80 líneas)
```
GIVEN `virgil init` completado
WHEN el agente lee AGENTS.md
THEN contiene solo:
     - Axiomas no negociables (5-10 líneas)
     - Build/test commands
     - Bootstrap: "Soy el SM, mis herramientas están en .virgil/"
     - Lista de skills disponibles (nombre + descripción, 1 línea c/u)
```

**AC-06.2** — Progressive disclosure via skills
```
GIVEN el agente necesita guía para una fase específica
WHEN activa un skill (e.g., /virgil-idea, /virgil-spec)
THEN carga SOLO el contenido de esa fase (~50-150 líneas),
     no los 31 docs completos
```

**AC-06.3** — Hooks como guardrails determinísticos
```
GIVEN reglas críticas (commit por refactor, echo pre-commit)
WHEN se configuran como hooks
THEN el agente NO puede saltarlos — enforcement determinístico,
     no advisory
```

**AC-06.4** — MCP server para binding layer
```
GIVEN el binding layer con datos
WHEN el agente necesita trazar un AC a código
THEN consulta el MCP server en tiempo real,
     sin cargar todo el grafo al contexto
```

---

### RF-07 — SM mental model (activación del agente)

Cuando Virgil se activa, el agente reconoce su rol y herramientas
disponibles.

**AC-07.1** — Reconocimiento de rol
```
GIVEN un proyecto con Virgil inicializado
WHEN el agente lee AGENTS.md o recibe un /virgil-* command
THEN reconoce: "Soy el SM. Tengo acceso a:
     1. RAG de metodología (skills en .virgil/skills/)
     2. RAG del proyecto (docs funcionales en adapter)
     3. Binding layer (MCP consultable)
     4. Código (working tree)"
```

**AC-07.2** — Lazy loading de contexto
```
GIVEN el agente activado como SM
WHEN necesita información de metodología o binding
THEN consulta on demand (skill específico o MCP query),
     NO carga todo al contexto de una vez
```

**AC-07.3** — Coexistencia con otros ecosistemas
```
GIVEN gentle-ai u otro sistema de agentes instalado
WHEN Virgil se activa
THEN opera sin interferir con otros sistemas,
     sin conflictos de hooks, skills, o AGENTS.md
```

---

### RF-08 — Adapter de storage para docs funcionales

Los docs funcionales (idea, spec, design, tasks, handoff) se almacenan
donde el proyecto configure.

**AC-08.1** — Adapter local (default)
```
GIVEN ninguna configuración de adapter especificada
WHEN Virgil necesita leer/escribir docs funcionales
THEN usa archivos locales en la raíz del proyecto:
     idea.md, spec.md, design.md, tasks.md, handoff.md
```

**AC-08.2** — Adapter configurable
```
GIVEN configuración en .virgil/config.yaml
WHEN se especifica adapter: engram | github | jira
THEN Virgil delega lectura/escritura al adapter correspondiente,
     manteniendo el mismo schema de contenido
```

**AC-08.3** — Interfaz de adapter
```
GIVEN un adapter custom
WHEN se implementa
THEN debe satisfacer la interfaz:
     read(artifact_type) → content
     write(artifact_type, content) → void
     list() → artifact_type[]
     exists(artifact_type) → bool
```

---

### RF-09 — Compilador modular

Reemplaza el compilador monolítico que producía 1 AGENTS.md de 1,061
líneas.

**AC-09.1** — Input: 31 docs de metodología
```
GIVEN los docs de metodología en docs/
WHEN el compilador ejecuta
THEN lee los 31 docs y los transforma en artefactos modulares
```

**AC-09.2** — Output: múltiples artefactos
```
GIVEN los docs procesados
WHEN la compilación completa
THEN produce:
     - 1 AGENTS.md mínimo (40-80 líneas)
     - N skills (1 por fase/concepto)
     - M sub-agent definitions (frontmatter con rol + modelo)
     - K hook scripts
     - 1 MCP server config
```

**AC-09.3** — Assets embebidos
```
GIVEN el binario de Virgil compilado
WHEN se distribuye
THEN los 31 docs de metodología y templates están embebidos
     via go:embed (o equivalente), sin dependencias externas
```

---

### RF-10 — Métricas de verificación (reemplazo de code review)

Virgil provee métricas que permiten gestionar desde un nivel superior
sin revisar código línea por línea (dogma Uncle Bob). Cubre dos capas:
trazabilidad (binding layer) y fuerza de tests (métricas externas
orquestadas).

**AC-10.1** — Coverage de binding
```
GIVEN un proyecto con binding layer poblado
WHEN se ejecuta `virgil coverage`
THEN reporta: total ACs, vinculados, verificados, stale, unbound,
     y porcentaje de cobertura con umbral configurable
```

**AC-10.2** — Reporte de staleness
```
GIVEN bindings marcados como STALE por hooks o cambios en docs
WHEN se ejecuta `virgil stale`
THEN lista todos los bindings desincronizados con:
     AC afectado, archivo, símbolo, razón del staleness
```

**AC-10.3** — Métricas de salud del proyecto
```
GIVEN un proyecto con Virgil activo
WHEN se ejecuta `virgil health`
THEN reporta 4 categorías:
     - Trazabilidad: binding coverage %, staleness %, gaps
     - Fuerza de tests: mutation score, CRAP score (si herramientas
       disponibles — Virgil reporta resultados, no los computa)
     - Estructura: complejidad ciclomática, module sizes
     - Salud documental: completitud idea→handoff, ACs sin binding
```

**AC-10.4** — Gate de CI configurable
```
GIVEN un pipeline de CI
WHEN incluye `virgil coverage --min 80`
THEN falla el build si el binding coverage está por debajo del 80%
AND el exit code y mensaje indican qué ACs no están cubiertos
```

**AC-10.5** — Orquestación de métricas externas
```
GIVEN un proyecto con herramientas de métricas instaladas
     (mutate4go, Stryker, pitest, gocyclo, etc.)
WHEN se ejecuta `virgil health` o `virgil metrics`
THEN Virgil ejecuta las herramientas detectadas,
     captura sus resultados estructurados,
     y los persiste junto al binding layer
AND reporta contra thresholds configurables por tier
     (strict | standard | relaxed | custom)
```

**AC-10.6** — Degradación elegante sin herramientas externas
```
GIVEN un proyecto SIN herramientas de mutation/CRAP instaladas
WHEN se ejecuta `virgil health`
THEN las categorías sin herramienta reportan "no disponible"
     en lugar de fallar, y sugiere qué instalar por lenguaje
```

---

### RF-11 — Fase de operación (facade opcional)

La metodología cubre hasta la operación del producto, pero cada proyecto
decide si activa esta fase.

**AC-11.1** — Detección de necesidad operativa
```
GIVEN un proyecto completando la fase de handoff
WHEN el handoff declara "documentación operativa requerida"
THEN Virgil activa el skill de operación que guía la producción de
     ops-runbook.md (o equivalente según el tipo de proyecto)
```

**AC-11.2** — Skip explícito
```
GIVEN un proyecto sin superficie operativa (librería, CLI, paquete)
WHEN el handoff declara "sin documentación operativa requerida"
THEN Virgil omite la fase de operación sin warnings,
     y Accept no verifica ops-runbook
```

**AC-11.3** — Adapter de operación
```
GIVEN un proyecto que activa operación
WHEN configura el tipo de operación (servicio, CLI, librería)
THEN Virgil carga el template correspondiente:
     servicio → ops-runbook.md (SLAs, monitoreo, runbooks)
     CLI → usage-guide.md (flags, exit codes, ejemplos)
     librería → api-reference.md (API, migración, changelog)
```

---

### RF-12 — Gate determinístico del handoff (`virgil handoff lint`)

Validación mecánica del contrato de handoff antes de habilitar código.
El gate "ZERO código hasta handoff aprobado" se convierte en enforcement
real, no solo advisory.

**AC-12.1** — Validación de completitud
```
GIVEN un handoff.md generado por el pipeline
WHEN se ejecuta `virgil handoff lint`
THEN valida:
     - Todo AC tiene ID único (AC-{n})
     - Toda task tiene dependencias válidas (sin refs rotas)
     - Todo AC referencia un RF del spec (sin ACs huérfanos)
     - Toda decisión arquitectónica tiene ADR
     - DAG de tasks sin ciclos
     - Toda task tiene estimación (S/M/L)
AND emite errores con guía de reparación por cada violación
```

**AC-12.2** — Hook pre-ejecución
```
GIVEN un handoff.md que no pasa lint
WHEN el agente intenta iniciar ejecución
THEN Virgil bloquea con lista de errores reparables,
     al estilo linter: qué falló, dónde, cómo arreglarlo
```

**AC-12.3** — Lint como CI gate
```
GIVEN un pipeline de CI
WHEN incluye `virgil handoff lint --strict`
THEN falla si el handoff tiene errores de completitud
AND el exit code indica el tipo de error más severo
```

---

### RF-13 — Coordinación de lanes paralelos

Semántica operativa para que múltiples lanes ejecuten el handoff en
paralelo sin pisarse. Principio 5 del dogma.

**AC-13.1** — Claiming de tasks
```
GIVEN un handoff en ejecución con lanes paralelos
WHEN un lane reclama una task
THEN el estado cambia a claimed con owner y timestamp
AND ningún otro lane puede reclamar la misma task
AND las precondiciones se verifican: todas las deps en "done"
```

**AC-13.2** — Execution state persistente
```
GIVEN un handoff en ejecución
WHEN se consulta `virgil status`
THEN muestra el estado de cada task:
     pending | claimed (owner, since) | done (commit, timestamp)
AND el estado sobrevive a crash o compaction del agente
```

**AC-13.3** — Reanudación determinística
```
GIVEN un agente que crasheó o perdió contexto (compaction)
WHEN el siguiente agente consulta el handoff
THEN sabe exactamente qué tasks están done, claimed o pending,
     y puede retomar sin duplicar trabajo
```

**AC-13.4** — Concurrencia de store
```
GIVEN múltiples hooks y lanes escribiendo al binding layer
WHEN las escrituras son simultáneas
THEN el store serializa escrituras sin corrupción
AND las lecturas son no-bloqueantes
```

---

## Requisitos no funcionales

### RNF-01 — Performance

| Métrica | Target |
|---------|--------|
| `virgil scan --full` en codebase de 50K LOC | < 30 segundos |
| `virgil refresh --diff` en delta de 100 líneas | < 2 segundos |
| Query MCP (trace, impact, coverage) | < 500 ms |
| `virgil init` | < 5 segundos |
| Memoria máxima durante scan completo | < 512 MB |

### RNF-02 — Portabilidad

- Binario standalone 100% pure Go (`CGO_ENABLED=0`), sin toolchain C
- Compatible: macOS (arm64, amd64), Linux (amd64, arm64), Windows (amd64)
- Cross-compile trivial: `GOOS/GOARCH` estándar, sin zig/musl
- `go install github.com/.../virgil@latest` funciona sin dependencias
- Sin dependencias externas en runtime (SQLite pure Go, Tree-Sitter pure Go)

### RNF-03 — Seguridad

- El binding DB (`.virgil/bindings.db`) contiene solo metadata
  (paths, símbolos, relaciones), NO contenido de código fuente
- Los hooks NO ejecutan código arbitrario — solo invocan `virgil refresh`
- El MCP server solo expone operaciones de lectura del grafo

### RNF-04 — Ergonomía

- Zero-config para el caso más común (local adapter, default hooks)
- Mensajes de error accionables (qué falló + cómo resolverlo)
- Progreso visible durante scan (`virgil scan --full` muestra porcentaje)

### RNF-05 — Precisión de inferencia

- La inferencia de bindings de `virgil refresh --diff` debe alcanzar
  ≥ 90% de precisión (≤ 10% de tasa de false-stale), medida contra los
  resultados de `virgil verify` sobre el corpus de referencia
- **False-stale**: `refresh` marca un binding como STALE, pero `verify`
  confirma que el símbolo sigue existiendo y no cambió
- Esta métrica se mide en la misma suite de performance de RNF-01
  (corpus fijo, CI) y se reporta junto a los targets de tiempo/memoria

---

## Contratos de interfaz

### CLI Commands

| Command | Input | Output | Descripción |
|---------|-------|--------|-------------|
| `virgil init` | directorio con .git | AGENTS.md + .virgil/ | Inicializa proyecto |
| `virgil scan --full` | codebase | .virgil/bindings.db populated | Crawl exhaustivo |
| `virgil refresh --diff` | git diff | bindings actualizados | Refresh incremental |
| `virgil verify <AC-ID>` | AC identifier | bound\|stale\|unbound | Verifica un AC (estructural: símbolo, staleness, test linkage, mutation strength) |
| `virgil trace <AC-ID>` | AC identifier | bindings + archivos + tests | Muestra trazabilidad |
| `virgil impact <file>` | file path | ACs afectados | Análisis de impacto |
| `virgil coverage` | — | % ACs verified + reporte | Cobertura de bindings |
| `virgil status` | — | resumen del binding layer + execution state | Estado general |
| `virgil claim <task-id>` | task ID + `--lane` | claimed \| rejected (ya claimed / deps no done) | Reclama una task de un handoff en ejecución (RF-13) |
| `virgil health` | — | reporte 4 categorías | Salud del proyecto |
| `virgil metrics` | — | resultados métricas externas | Mutation/CRAP/complexity |
| `virgil handoff lint` | handoff.md | errores + guía | Validación del contrato |

### MCP Server Tools

| Tool | Input | Output |
|------|-------|--------|
| `virgil_trace` | `{requirement_id: string}` | `{bindings: Binding[], code: CodeArtifact[], tests: TestArtifact[]}` |
| `virgil_impact` | `{file_path: string}` | `{requirements: Requirement[], bindings: Binding[]}` |
| `virgil_coverage` | `{}` | `{total: int, verified: int, stale: int, unbound: int, pct: float}` |
| `virgil_stale` | `{}` | `{stale_bindings: Binding[]}` |
| `virgil_declare` | `{requirement_id: string, artifact: {path, symbol}}` | `{binding_id: string}` |

### Adapter Interface (Go)

```go
type DocAdapter interface {
    Read(artifactType string) ([]byte, error)
    Write(artifactType string, content []byte) error
    List() ([]string, error)
    Exists(artifactType string) (bool, error)
}
```

---

## Restricciones y asunciones

### Restricciones

1. **Go 100% pure (CGO_ENABLED=0)** — core, binding engine, MCP server, CLI. No Python. No CGo.
2. **Sin dependencias de runtime** — todo embebido (docs, GoTreeSitter, modernc.org/sqlite)
3. **No interfiere con otros ecosistemas** — gentle-ai, Cursor rules, etc.
4. **Binding DB local** — `.virgil/bindings.db`, gitignoreado por default
5. **Skills como archivos de texto** — no binarios, no compilados, editables
6. **Multiplataforma sin fricción** — cross-compile con GOOS/GOARCH estándar

### Asunciones

1. El proyecto destino usa git como VCS
2. El agente soporta skills/slash commands (Claude Code, Cursor, etc.)
3. Tree-Sitter tiene grammar para los lenguajes del proyecto destino
4. El MIM aprueba cada transición de fase explícitamente
5. Los docs funcionales siguen el schema de la metodología (29148, 42010, etc.)

---

## Priorización (MoSCoW)

### Must Have (v2.0)

- RF-01: `virgil init` (gobernanza + skills + hooks + binding DB)
- RF-02: `/virgil-idea` flow completo
- RF-03: `/virgil-takeover` con crawl y derivación
- RF-04: Binding layer con CRUD y queries básicas
- RF-05: Git hooks para refresh incremental
- RF-06: Delivery modular (skills reemplazan AGENTS.md monolítico)
- RF-07: SM mental model activation
- RF-09: Compilador modular con go:embed
- RF-10: Métricas de verificación (coverage, staleness, health, metrics externas)
- RF-11: Fase de operación como facade opcional
- RF-12: Gate determinístico del handoff (virgil handoff lint)
- RF-13: Coordinación de lanes paralelos (claiming, execution state)

### Should Have (v2.1)

- RF-04.7: MCP server completo (trace, impact, coverage, stale, declare)
- RF-08: Adapter configurable (engram, GitHub Issues)
- Semántica de embeddings para inferencia de bindings más precisa
- `virgil dashboard` — visualización del grafo en terminal

### Could Have (v2.2+)

- Adapter para Jira, Linear
- Cross-repo binding (monorepos con múltiples proyectos)
- Integración con CI/CD para verificación automática de coverage
- Plugin para IDE (VS Code, JetBrains) — binding layer visual

### Won't Have (v2.x)

- UI web para gestión de bindings
- Generación automática de código desde specs
- Ejecución de tests (Virgil traza, no ejecuta)
- Gestión de branches/PRs (responsabilidad de git/gh)

---

## Trazabilidad

| Requisito | Traza a idea.md |
|-----------|----------------|
| RF-01 | Contexto del sistema — "El proyecto destino instala Virgil" |
| RF-02 | Dos flujos — "Flujo 1: Idea → Producto (greenfield)" |
| RF-03 | Dos flujos — "Flujo 2: Takeover → Ownership (brownfield)" |
| RF-04 | Binding Layer — "EL GAP CRÍTICO" |
| RF-05 | Binding Layer — "Estrategia de crawl: exhaustivo una vez, incremental siempre" |
| RF-06 | Delivery — "Cómo la Metodología Llega al Agente" |
| RF-07 | Contexto del sistema — "Modelo mental del agente al activarse" |
| RF-08 | Capa 2 — "Almacenados en el adapter configurado" |
| RF-09 | Delivery — "Ya NO se compila en un archivo monolítico" |
| RF-10 | Dogma — "No revisas código, verificas métricas y trazabilidad" |
| RF-11 | Capa 2 — operación como facade opcional por tipo de proyecto |
| RF-12 | Dogma principio 6 — gates determinísticos |
| RF-13 | Dogma principio 5 — coordinación de lanes paralelos |

---

## Metadata

- Fecha de creación: 2026-08-05
- Estado: borrador
- Fuente: idea-virgil.md + research (RTM + code intelligence + GoReleaser)
- Revisores pendientes: MIM (aprobación)
