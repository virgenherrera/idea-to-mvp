# idea-to-mvp --- Vista General del Framework

> Mapa de navegacion para entender el sistema completo antes de leer los
> documentos detallados. Los diagramas son la comunicacion principal; el
> texto es tejido conectivo.

---

## Actores y Roles

El framework opera con tres capas de actores: el humano (MIM), la
infraestructura de orquestacion (SM + TPM), y los roles productivos
(equipo default + ad-hoc).

```mermaid
%% Relacion entre actores del framework
flowchart TD
    MIM["MIM\n(Humano)\nDecide, aprueba,\ndesbloquea"]

    subgraph INFRA["Infraestructura"]
        SM["SM\n(Scrum Master)\nFacade / Orquestador"]
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

    MIM -->|"toda interaccion"| SM
    SM -->|"contratos de delegacion"| DEFAULT_TEAM
    SM -->|"contratos ad-hoc"| ADHOC
    SM -->|"instrucciones CRUD"| TPM
    TPM -->|"estado de artefactos"| SM
    DEFAULT_TEAM -->|"contenido producido"| TPM
    ADHOC -->|"contenido producido"| TPM
```

Reglas clave:

- El SM **nunca produce contenido** --- solo orquesta, convoca y valida gates.
- El TPM **es el unico que escribe** en el artifact store --- con criterio editorial.
- Los roles son sub-agentes con personalidad que **cambia por fase**.
- El SM puede **extender el equipo** con roles ad-hoc justificados.

> Detalle: [role-profiles.md](role-profiles.md) (contratos por fase) y
> [behavior-scrum-master-routing.md](behavior-scrum-master-routing.md)
> (reglas del SM).

---

## Modos de Funcionamiento

El framework separa planificacion de ejecucion con un contrato formal
entre ambos modos: el handoff.

```mermaid
%% Dos modos del framework y su interfaz
flowchart LR
    subgraph PLAN["MODO PLANIFICACION"]
        direction TB
        P1["Entrada: idea, ticket,\nchallenge, spec parcial"]
        P2["Scrum team como\nlentes de revision"]
        P3["Escribe en:\nartifact store\n(NUNCA en el repo)"]
    end

    HANDOFF["handoff.md\n(contrato entre modos)"]

    subgraph EXEC["MODO EJECUCION"]
        direction TB
        E1["Entrada:\nhandoff + AGENTS.md"]
        E2["Orquestador\n+ sub-agentes"]
        E3["Escribe en:\nworking tree\n(SOLO el repo)"]
    end

    PLAN -->|"produce"| HANDOFF
    HANDOFF -->|"consume"| EXEC
    EXEC -->|"gap detectado"| PLAN
```

| Aspecto | Planificacion | Ejecucion |
|---------|---------------|-----------|
| Proposito | Producir fuentes de verdad | Implementar codigo |
| Participantes | Scrum team (lentes) | Orquestador + minions |
| Donde escribe | Artifact store (fuera del repo) | Working tree del repo |
| Estado actual | **DEFINIDO** | **TBD** |

> Detalle: [operational-model.md](operational-model.md).

---

## Pipeline Completo

El ciclo tiene 4 macro-fases. Solo la primera esta completamente
definida en el framework actual.

```mermaid
%% Pipeline completo con macro-fases
flowchart TD
    subgraph MACRO_1["Idea a Handoff (DEFINIDO)"]
        direction LR
        F1["Fase 1\nDefinir Idea"]
        F2["Fase 2\nEspecificar"]
        F3["Fase 3\nDisenar"]
        F4["Fase 4\nDesglosar\nTareas"]
        F5["Fase 5\nGenerar\nHandoff"]
        F1 --> F2 --> F3 --> F4 --> F5
    end

    subgraph MACRO_2["Handoff a Ejecucion (TBD)"]
        direction LR
        EX["Ejecucion\nOrquestador\n+ minions\nimplementan"]
    end

    subgraph POST["Post-Ejecucion (DEFINIDO)"]
        direction LR
        F6["Fase 6\nVerificar"]
        F7["Fase 7\nAceptar"]
        F8["Fase 8\nRetrospectiva"]
        F6 --> F7 --> F8
    end

    subgraph MACRO_3["Ejecucion a Operacion (TBD)"]
        direction LR
        OPS["Operacion\nMonitoreo,\nrunbooks,\nSRE"]
    end

    F5 -->|"handoff.md"| EX
    EX -->|"codigo\nimplementado"| F6
    F8 -->|"siguiente ciclo"| F1
    F7 -->|"ops-runbook.md"| OPS

    style MACRO_2 stroke-dasharray: 5 5
    style MACRO_3 stroke-dasharray: 5 5
```

### Roles convocados por fase

| Fase | Roles activos |
|------|---------------|
| 1. Idea | PO |
| 2. Spec | PO + QA + UX (condicional) |
| 3. Diseno | Dev Lead + DevSecOps + UX (condicional) |
| 4. Tareas | Dev Lead + DevSecOps (cond) + QA (cond) |
| 5. Handoff | TPM (compila bajo instruccion del SM) |
| 6. Verificar | QA + Dev Lead + DevSecOps (condicional) |
| 7. Aceptar | Todos los roles activos (votacion paralela) |
| 8. Retro | Todos los roles activos |

> Detalle: [behavior-scrum-master-routing.md](behavior-scrum-master-routing.md)
> (fases 1-8) y [role-profiles.md](role-profiles.md) (contratos de cada
> rol por fase).

---

## Artefactos

El framework produce 6 artefactos universales respaldados por estandares
ISO/IEC/IEEE. Cada fase consume el output de la anterior y produce los
parametros requeridos para la siguiente.

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
    HANDOFF -->|"post-ejecucion"| OPS
```

### Quien produce y quien valida

| Artefacto | Produce | Valida (gate) |
|-----------|---------|---------------|
| `idea.md` | PO | SM (estructural via TPM) |
| `spec.md` | PO | QA (testeabilidad) + UX (experiencia) |
| `design.md` | Dev Lead + DevSecOps | SM (via TPM) + DevSecOps + UX |
| `tasks.md` | Dev Lead | QA (verificabilidad) + SM (via TPM) |
| `handoff.md` | TPM (compila) | SM (autocontencion) |
| `ops-runbook.md` | DevSecOps + Dev Lead | SM (gate) |

> Regla cardinal: **quien produce nunca valida su propio artefacto**.
>
> Detalle: [artifact-model.md](artifact-model.md) (schemas, contenido
> minimo, jerarquia de work items, adaptadores de persistencia).

---

## Maquina de Estados de Artefactos

Cada artefacto transiciona a traves de una state machine configurable.
El estado `approved` es el que habilita la siguiente fase.

```mermaid
%% State machine default de artefactos
stateDiagram-v2
    [*] --> draft: Artefacto creado
    draft --> review: Productor solicita revision
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
    TPM_Q["SM pregunta al TPM:\n¿que artefactos existen?"]
    TPM_Q --> D1{{"idea\napproved?"}}
    D1 -->|No| PH1["Fase 1:\nDefinir Idea"]
    D1 -->|Si| D2{{"spec\napproved?"}}
    D2 -->|No| PH2["Fase 2:\nEspecificar"]
    D2 -->|Si| D3{{"design\napproved?"}}
    D3 -->|No| PH3["Fase 3:\nDisenar"]
    D3 -->|Si| D4{{"tasks\napproved?"}}
    D4 -->|No| PH4["Fase 4:\nTareas"]
    D4 -->|Si| PH5["Fase 5:\nHandoff"]
```

> Detalle: [artifact-model.md](artifact-model.md) (seccion `transition`)
> y [behavior-scrum-master-routing.md](behavior-scrum-master-routing.md)
> (state machine del proyecto).

---

## Modelo de Delegacion

El SM delega trabajo via **contratos de delegacion** con campos
obligatorios. Despues de cada retorno, ejecuta el **PDC** (Post-Delegation
Checkpoint).

```mermaid
%% Ciclo de delegacion SM -> sub-agente -> PDC
sequenceDiagram
    participant SM as SM
    participant SUB as Sub-agente
    participant TPM as TPM

    SM ->> SUB: Contrato (rol, personalidad,<br/>contexto, input, output, restricciones)
    activate SUB
    SUB ->> SUB: Lee del artifact store<br/>via Pattern B (topic_keys)
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
%% Dos patrones de retrieval y cuando usar cada uno
flowchart TD
    NEED["Sub-agente necesita contexto"]
    NEED --> Q{{"¿Target conocido\ny determinista?"}}

    Q -->|"Si"| PB["Pattern B\nSM pasa topic_key\nSub-agente lee directo\n(6x mas barato)"]
    Q -->|"No (busqueda\nfuzzy o fan-out 8+)"| PA["Pattern A\nSM busca, cura, inyecta\n(calidad sobre costo)"]
```

**Circuit breaker**: si 3 delegaciones consecutivas al mismo rol fallan, el SM detiene la cadena y escala al MIM.

> Detalle: [behavior-scrum-master-routing.md](behavior-scrum-master-routing.md)
> (PDC, circuit breaker, context resilience) y
> [role-profiles.md](role-profiles.md) (contratos por fase).

---

## Fast-Forward

El SM no avanza siempre una fase a la vez. Evalua un **gradiente de
certeza** con 4 factores (F1-F4) y avanza proporcionalmente.

```mermaid
%% Scoring de fast-forward
flowchart LR
    subgraph SCORE["Checklist de certeza (F1-F4, 0-2 pts c/u)"]
        direction TB
        F1["F1: Artefactos existentes\n0=RAG vacio\n2=spec+design+tasks approved"]
        F2["F2: Estandarizacion\n0=dominio custom\n2=estandar abierto puro"]
        F3["F3: Ambiguedad\n0=infinitas interpretaciones\n2=determinista"]
        F4["F4: Referencia existente\n0=sin codebase\n2=codebase con patrones"]
    end

    subgraph RESULT["Resultado"]
        direction TB
        LOW["0-2 pts: Baja\nIdea + preguntas"]
        MED["3-5 pts: Media\nIdea + spec parcial"]
        HIGH["6-8 pts: Alta\nHasta handoff\no ejecucion"]
    end

    SCORE --> RESULT
```

Ejemplos:

| Input | Score | Certeza | Accion |
|-------|-------|---------|--------|
| "Hazme el uber de lanchas" | 0 | Baja | Idea + preguntas |
| "Agrega auth JWT" (codebase Express) | 4 | Media | Idea + spec parcial |
| "Crea modulo OTEL" (codebase NestJS) | 6 | Alta | Hasta handoff |
| Epic ya groomeado (todo en RAG) | 8 | Alta | Fast-forward a ejecucion |

El SM registra el score F1-F4 en `idea.md` para auditabilidad. El
fast-forward tambien aplica **mid-cycle** (bugs en produccion, epics
ya groomeados).

> Detalle: [behavior-scrum-master-routing.md](behavior-scrum-master-routing.md)
> (seccion fast-forward contextual).

---

## Artifact Store y Adaptadores

Los artefactos se persisten via una **interfaz universal de 9
operaciones**. El adaptador es pluggable --- el framework define la
interfaz, no la implementacion.

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

El TPM actua como DBMS: no decide que datos crear, pero decide como se
almacenan, valida integridad, y sirve consultas con criterio editorial.

> Detalle: [artifact-model.md](artifact-model.md) (interfaz universal,
> contrato de comportamiento, garantias ACID, adaptadores).

---

## Que Sigue (areas TBD)

El framework actualmente cubre la **macro-fase de planificacion** en
detalle. Las siguientes areas estan identificadas pero no definidas:

| Area | Estado | Descripcion |
|------|--------|-------------|
| Modo Ejecucion | TBD | Como el orquestador consume el handoff, delega a sub-agentes, y produce codigo. Patron orquestador-minion. |
| Modo Operacion | TBD | Para proyectos con servicios vivos: monitoreo, runbooks, SRE, alertas. Consumidor de `ops-runbook.md`. |
| Adaptadores avanzados | TBD | Jira, DBMS, Git repo, MS Project como adaptadores del artifact store. |
| Routing no-Scrum | TBD | Routing tables para Kanban (WIP limits), Shape Up (bets), SAFe (PIs). Los artefactos son universales; la orquestacion no. |
| Tiers de activacion | TBD | Como escala hacia abajo el modo planificacion para proyectos simples o challenges con timebox. |
| Transacciones del adaptador | TBD | Primitivas `begin`/`commit`/`rollback` para adaptadores sin soporte nativo. |

---

## Indice de Documentos Detallados

| Documento | Que define |
|-----------|-----------|
| [operational-model.md](operational-model.md) | Dos modos, ownership, limites, adaptador por defecto |
| [artifact-model.md](artifact-model.md) | 6 artefactos, TPM, interfaz de adaptadores, state machine, jerarquia de work items |
| [behavior-scrum-master-routing.md](behavior-scrum-master-routing.md) | SM como facade, state machine del proyecto, fast-forward, PDC, circuit breaker |
| [role-profiles.md](role-profiles.md) | Contratos de delegacion por fase, personalidades, activacion condicional, roles ad-hoc |
