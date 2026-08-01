# idea-to-mvp --- Vista General del Framework

← [Índice](README.md)

> Mapa de navegación para entender el sistema completo antes de leer los
> documentos detallados. Los diagramas son la comunicación principal; el
> texto es tejido conectivo.

---

## Actores y Roles

El framework opera con tres capas de actores: el humano (MIM), la
infraestructura de orquestación (SM + TPM), y los roles productivos
(equipo default + ad-hoc).

```mermaid
%% Relación entre actores del framework
flowchart TD
    MIM["MIM\n(Humano)\nDecide, aprueba,\ndesbloquea"]

    subgraph INFRA["Infraestructura"]
        SM["SM\n(Session Manager)\nFacade / Orquestador"]
        TPM["TPM\n(Technical Program Manager)\nDBMS del artifact store"]
    end

    subgraph DEFAULT_TEAM["Equipo Default (5 roles)"]
        direction LR
        PO["PO\nValor de negocio"]
        DEV["Dev Lead\nArquitectura"]
        QA["QA\nVerificabilidad"]
        SEC["DevSecOps\nSeguridad + Infra"]
        UX["UX\nExperiencia de usuario"]
    end

    ADHOC["Roles Ad-Hoc\n(DBA, Performance Eng,\nDomain Expert, etc.)"]

    MIM -->|"toda interacción"| SM
    SM -->|"contratos de delegación"| DEFAULT_TEAM
    SM -->|"contratos ad-hoc"| ADHOC
    SM -->|"instrucciones CRUD"| TPM
    TPM -->|"estado de artefactos"| SM
    DEFAULT_TEAM -->|"contenido producido"| TPM
    ADHOC -->|"contenido producido"| TPM
```

Reglas clave:

- El SM **nunca produce contenido** --- solo orquesta, convoca y valida gates.
- El TPM **es el único que escribe** en el artifact store --- con criterio editorial.
- Los roles son sub-agentes con personalidad que **cambia por fase**.
- El SM puede **extender el equipo** con roles ad-hoc justificados.

> Detalle: [roles](planning/roles/README.md) (contratos por fase) y
> [comportamiento SM](planning/behavior/README.md)
> (reglas del SM).

---

## Modos de Funcionamiento

El framework separa planificación de ejecución con un contrato formal
entre ambos modos: el handoff.

```mermaid
%% Dos modos del framework y su interfaz
flowchart LR
    subgraph PLAN["MODO PLANIFICACIÓN"]
        direction TB
        P1["Entrada: idea, ticket,\nchallenge, spec parcial"]
        P2["Scrum team como\nlentes de revisión"]
        P3["Escribe en:\nartifact store\n(NUNCA en el repo)"]
    end

    HANDOFF["handoff.md\n(contrato entre modos)"]

    subgraph EXEC["MODO EJECUCIÓN"]
        direction TB
        E1["Entrada:\nhandoff + AGENTS.md"]
        E2["Orquestador\n+ sub-agentes"]
        E3["Escribe en:\nworking tree\n(SOLO el repo)"]
    end

    PLAN -->|"produce"| HANDOFF
    HANDOFF -->|"consume"| EXEC
    EXEC -->|"gap detectado"| PLAN
```

| Aspecto | Planificación | Ejecución |
|---------|---------------|-----------|
| Propósito | Producir fuentes de verdad | Implementar código |
| Participantes | Scrum team (lentes) | Orquestador + minions |
| Dónde escribe | Artifact store (fuera del repo) | Working tree del repo |
| Estado actual | **DEFINIDO** | **DEFINIDO** |

> Detalle: [modelo operativo](planning/operational-model.md).

---

## Pipeline Completo

El ciclo tiene 4 macro-fases, todas definidas. Las fases post-ejecución
están definidas como parte del Modo 1; la macro-fase de operación es el
Modo 3, opcional y reactivo.

```mermaid
%% Pipeline completo con macro-fases
flowchart TD
    subgraph MACRO_1["Idea a Handoff (DEFINIDO)"]
        direction LR
        F1["Fase 1\nDefinir Idea"]
        F2["Fase 2\nEspecificar"]
        F3["Fase 3\nDiseñar"]
        F4["Fase 4\nDesglosar\nTareas"]
        F5["Fase 5\nGenerar\nHandoff"]
        F1 --> F2 --> F3 --> F4 --> F5
    end

    subgraph MACRO_2["Handoff a Ejecución (DEFINIDO)"]
        direction LR
        EX_C["Pre-Fase\nContratos"]
        EX_R["Fase Red\nTests"]
        EX_G["Fase Green\nImplementación"]
        EX_RF["Fase Refactor\nCalidad"]
        EX_A["Fase Accept\nCertificación QA"]
        EX_C --> EX_R --> EX_G --> EX_RF --> EX_A
    end

    subgraph POST["Post-Ejecución (DEFINIDO)"]
        direction LR
        F6["Fase 6\nVerificar"]
        F7["Fase 7\nAceptar"]
        F8["Fase 8\nRetrospectiva"]
        F6 --> F7 --> F8
    end

    subgraph MACRO_3["Ejecución a Operación (DEFINIDO)"]
        direction LR
        OPS["Modo 3\nOperación\nUsuario + agente\nasistente"]
    end

    F5 -->|"handoff.md"| EX_C
    EX_A -->|"código\ncertificado"| F6
    F8 -->|"siguiente ciclo"| F1
    F7 -->|"ops-runbook.md"| OPS
    OPS -->|"gap detectado"| F1
```

Vista alternativa: la línea de tiempo pone el foco en el **orden temporal**
de las fases dentro de cada macro-etapa, en vez de las dependencias entre
sub-fases.

```mermaid
timeline
    title Pipeline del Proyecto
    section Planificación
        Fase 1 - Idea : idea.md aprobado
        Fase 2 - Spec : spec.md aprobado
        Fase 3 - Diseño : design.md aprobado
        Fase 4 - Tareas : tasks.md aprobado
        Fase 5 - Handoff : handoff.md aprobado
    section Ejecución
        Pre-fase - Contratos : interfaces definidas
        Fase Red : suite de tests completa
        Fase Green : tests pasando
        Fase Refactor : calidad aprobada
        Fase Accept : certificación QA
    section Cierre
        Verificar : QA aprueba
        Aceptar : panel vota
        Retrospectiva : ciclo cerrado
```

### Roles convocados por fase

| Fase | Roles activos |
|------|---------------|
| 1. Idea | PO |
| 2. Spec | PO + QA + UX (condicional) |
| 3. Diseño | Dev Lead + DevSecOps + UX (condicional) |
| 4. Tareas | Dev Lead + DevSecOps (cond) + QA (cond) |
| 5. Handoff | TPM (compila bajo instrucción del SM) |
| 6. Verificar | QA + Dev Lead + DevSecOps (condicional) |
| 7. Aceptar | Todos los roles activos (votación paralela) |
| 8. Retro | Todos los roles activos |

> Detalle: [comportamiento SM](planning/behavior/README.md)
> (fases 1-8) y [roles](planning/roles/README.md) (contratos de cada
> rol por fase).

---

## Artefactos

El framework produce 6 artefactos universales respaldados por estándares
ISO/IEC/IEEE. Cada fase consume el output de la anterior y produce los
parámetros requeridos para la siguiente.

```mermaid
%% Cadena de artefactos con respaldo ISO
flowchart LR
    IDEA["idea.md\nISO 29148\nsec 9.3 BRS"]
    SPEC["spec.md\nISO 29148\nStRS/SRS"]
    DESIGN["design.md\nISO 42010\n+ IEEE 1016"]
    TASKS["tasks.md\nISO 21502\nsec 7.6"]
    HANDOFF["handoff.md\nISO 15289\ntransition"]
    OPS["ops-runbook.md\nISO 20000\n+ ITIL 4"]

    IDEA -->|"problema\nvalor\nrestricciones"| SPEC
    SPEC -->|"ACs\ncontratos\nconstraints"| DESIGN
    DESIGN -->|"stack\narquitectura\npatrones"| TASKS
    TASKS -->|"tareas\ndeps\nACs"| HANDOFF
    HANDOFF -->|"post-ejecución"| OPS
```

### Quién produce y quién valida

| Artefacto | Produce | Valida (gate) |
|-----------|---------|---------------|
| `idea.md` | PO | SM (estructural via TPM) |
| `spec.md` | PO | QA (testeabilidad) + UX (experiencia) |
| `design.md` | Dev Lead + DevSecOps | SM (via TPM) + DevSecOps + UX |
| `tasks.md` | Dev Lead | QA (verificabilidad) + SM (via TPM) |
| `handoff.md` | TPM (compila) | SM (autocontención) |
| `ops-runbook.md` | DevSecOps + Dev Lead | SM (gate) |

> Regla cardinal: **quién produce nunca valida su propio artefacto**.
>
> Detalle: [artefactos](planning/artifacts/README.md) (schemas, contenido
> mínimo, jerarquía de work items, adaptadores de persistencia).

---

## Máquina de Estados de Artefactos

Cada artefacto transiciona a través de una state machine configurable.
El estado `approved` es el que habilita la siguiente fase.

```mermaid
%% State machine default de artefactos
stateDiagram-v2
    [*] --> draft: Artefacto creado
    draft --> review: Productor solicita revisión
    draft --> cancelled: SM o MIM cancela

    review --> approved: Gate aprobado
    review --> rejected: Gate rechazado
    review --> draft: Devuelto para correcciones

    rejected --> draft: Productor corrige y reintenta

    approved --> draft: SM reabre (mid-planning edit)
```

La state machine del **proyecto** se deriva del estado de los artefactos
en el RAG. El SM no persiste estado --- lo reconstruye consultando al TPM:

```mermaid
%% SM deriva fase actual del estado de artefactos
flowchart LR
    TPM_Q["SM pregunta al TPM:\n¿qué artefactos existen?"]
    TPM_Q --> D1{{"idea\napproved?"}}
    D1 -->|No| PH1["Fase 1:\nDefinir Idea"]
    D1 -->|Sí| D2{{"spec\napproved?"}}
    D2 -->|No| PH2["Fase 2:\nEspecificar"]
    D2 -->|Sí| D3{{"design\napproved?"}}
    D3 -->|No| PH3["Fase 3:\nDiseñar"]
    D3 -->|Sí| D4{{"tasks\napproved?"}}
    D4 -->|No| PH4["Fase 4:\nTareas"]
    D4 -->|Sí| PH5["Fase 5:\nHandoff"]
```

> Detalle: [artefactos](planning/artifacts/README.md) (sección `transition`)
> y [comportamiento SM](planning/behavior/README.md)
> (state machine del proyecto).

---

## Modelo de Delegación

El SM delega trabajo vía **contratos de delegación** con campos
obligatorios. Después de cada retorno, ejecuta el **PDC** (Post-Delegation
Checkpoint).

```mermaid
%% Ciclo de delegación SM -> sub-agente -> PDC
sequenceDiagram
    participant SM as SM
    participant SUB as Sub-agente
    participant TPM as TPM

    SM ->> SUB: Contrato (rol, personalidad,<br/>contexto, input, output, restricciones)
    activate SUB
    SUB ->> SUB: Lee del artifact store<br/>vía Pattern B (topic_keys)
    SUB -->> SM: Resultado + Status Report
    deactivate SUB

    Note over SM: PDC (4 pasos obligatorios)

    SM ->> SM: 1. ECHO: ¿coherente con contrato?
    SM ->> SM: 2. VERIFY: ¿cubre todo el scope?
    SM ->> TPM: 3. MARK: persistir resultado
    SM ->> SM: 4. DECIDE: ¿avanzar, re-delegar, escalar?
```

### Pattern A vs Pattern B (retrieval)

```mermaid
%% Dos patrones de retrieval y cuándo usar cada uno
flowchart TD
    NEED["Sub-agente necesita contexto"]
    NEED --> Q{{"¿Target conocido\ny determinista?"}}

    Q -->|"Sí"| PB["Pattern B\nSM pasa topic_key\nSub-agente lee directo\n(6x más barato)"]
    Q -->|"No (búsqueda\nfuzzy o fan-out 8+)"| PA["Pattern A\nSM busca, cura, inyecta\n(calidad sobre costo)"]
```

**Circuit breaker**: si 3 delegaciones consecutivas al mismo rol fallan, el SM detiene la cadena y escala al MIM.

> Detalle: [comportamiento SM](planning/behavior/README.md)
> (PDC, circuit breaker, context resilience) y
> [roles](planning/roles/README.md) (contratos por fase).

---

## Fast-Forward

El SM no avanza siempre una fase a la vez. Evalúa un **gradiente de
certeza** con 4 factores (F1-F4) y avanza proporcionalmente.

```mermaid
%% Scoring de fast-forward
flowchart LR
    subgraph SCORE["Checklist de certeza (F1-F4, 0-2 pts c/u)"]
        direction TB
        F1["F1: Artefactos existentes\n0=RAG vacío\n2=spec+design+tasks approved"]
        F2["F2: Estandarización\n0=dominio custom\n2=estándar abierto puro"]
        F3["F3: Ambigüedad\n0=infinitas interpretaciones\n2=determinista"]
        F4["F4: Referencia existente\n0=sin codebase\n2=codebase con patrones"]
    end

    subgraph RESULT["Resultado"]
        direction TB
        LOW["0-2 pts: Baja\nIdea + preguntas"]
        MED["3-5 pts: Media\nIdea + spec parcial"]
        HIGH["6-8 pts: Alta\nHasta handoff\no ejecución"]
    end

    SCORE --> RESULT
```

Ejemplos:

| Input | Score | Certeza | Acción |
|-------|-------|---------|--------|
| "Hazme el uber de lanchas" | 0 | Baja | Idea + preguntas |
| "Agrega auth JWT" (codebase Express) | 4 | Media | Idea + spec parcial |
| "Crea módulo OTEL" (codebase NestJS) | 6 | Alta | Hasta handoff |
| Epic ya groomeado (todo en RAG) | 8 | Alta | Fast-forward a ejecución |

El SM registra el score F1-F4 en `idea.md` para auditabilidad. El
fast-forward también aplica **mid-cycle** (bugs en producción, epics
ya groomeados).

> Detalle: [comportamiento SM](planning/behavior/README.md)
> (sección fast-forward contextual).

---

## Artifact Store y Adaptadores

Los artefactos se persisten vía una **interfaz universal de 9
operaciones**. El adaptador es pluggable --- el framework define la
interfaz, no la implementación.

```mermaid
%% Interfaz universal y adaptadores
flowchart TD
    subgraph INTERFACE["Interfaz del Adaptador"]
        direction LR
        OPS_I["ingest | save | read\nsearch | list | delete\nverifyConsistency\nhistory | transition"]
    end

    subgraph ADAPTERS["Implementaciones"]
        direction LR
        LOCAL["Local (DEFAULT)\nArchivos .md en\n~/.idea-to-mvp/"]
        ENGRAM["Engram\nCross-session\nbuscable"]
        FUTURE["Jira | DBMS | Git\n(TBD)"]
    end

    TPM_W["TPM media\ntodas las escrituras"] --> INTERFACE
    INTERFACE --> LOCAL
    INTERFACE --> ENGRAM
    INTERFACE --> FUTURE

    style FUTURE stroke-dasharray: 5 5
```

El TPM actúa como DBMS: no decide qué datos crear, pero decide cómo se
almacenan, valida integridad, y sirve consultas con criterio editorial.

> Detalle: [artefactos](planning/artifacts/README.md) (interfaz universal,
> contrato de comportamiento, garantías ACID, adaptadores).

---

## Qué Sigue (áreas TBD)

El framework cubre la macro-fase de planificación y la macro-fase de
ejecución en detalle. Las siguientes áreas están identificadas pero no
definidas:

| Área | Estado | Descripción |
|------|--------|-------------|
| Modo Ejecución | **DEFINIDO** | 5 fases (Contratos → Red → Green → Refactor → Accept). Contract-first, modelo de boundaries (App + E2E), revisión multi-dimensional. Ver [modelo de ejecución](execution/README.md). |
| Modo Operación | **DEFINIDO** | Opcional. Para proyectos con superficie operativa: el usuario consume el producto con asistencia del agente. Reactivo, sin fases. Ver [modelo de operación](operation/README.md). |
| Adaptadores avanzados | TBD | Jira, DBMS, Git repo, MS Project como adaptadores del artifact store. |
| Routing no-Scrum | TBD | Routing tables para Kanban (WIP limits), Shape Up (bets), SAFe (PIs). Los artefactos son universales; la orquestación no. |
| Tiers de activación | TBD | Cómo escala hacia abajo el modo planificación para proyectos simples o challenges con timebox. |
| Transacciones del adaptador | TBD | Primitivas `begin`/`commit`/`rollback` para adaptadores sin soporte nativo. |

---

## Índice de Documentos Detallados

| Documento | Qué define |
|-----------|-----------|
| [Modelo operativo](planning/operational-model.md) | Dos modos, ownership, límites, adaptador por defecto |
| [Artefactos](planning/artifacts/README.md) | 6 artefactos, TPM, interfaz de adaptadores, state machine, jerarquía de work items |
| [Comportamiento SM](planning/behavior/README.md) | SM como facade, state machine del proyecto, fast-forward, PDC, circuit breaker |
| [Roles](planning/roles/README.md) | Contratos de delegación por fase, personalidades, activación condicional, roles ad-hoc |
| [Modo Ejecución](execution/README.md) | Modo 2: Contract-first, Red-Green-Refactor macro, roles de ejecución, conexión con Modo 1 |
| [Modo Operación](operation/README.md) | Modo 3: opcional y reactivo, sin fases, usuario + agente asistente, conexión con Modo 1 y Modo 2 |
