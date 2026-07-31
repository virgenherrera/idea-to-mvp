# Diseño del Modelo Operativo — idea-to-mvp

> Documento de trabajo. Objetivo: definir CÓMO opera el framework antes de
> decidir CÓMO implementarlo (skills, agentes, paquetes, etc.).

---

## El Problema

El framework actualmente mezcla tres concerns en un solo repositorio:

1. **Reglas de gobernanza** — axiomas de AGENTS.md, compact rules, fases del
   pipeline. Estas SÍ PERTENECEN a cada repo que adopta el framework.
2. **Tooling de planificación** — fases del ciclo (idea, spec, diseño,
   tareas, handoff), roles del scrum team, persistencia de artefactos
   (engram, local, híbrido). Son OPERACIONALES y no deben contaminar los
   repos adoptantes.
3. **Tooling de ejecución** — patrón orquestador-minion, delegación a
   sub-agentes, inyección de personalidad/contexto, resolución de skills.
   Son de RUNTIME y no deben estar acoplados a las reglas de gobernanza.

Resultado: los repos adoptantes acumulan archivos `.tmp-*`, directorios
`openspec/`, documentos de feedback y estado del ciclo de planificacion que no tienen nada que ver
con su codebase.

Cada repo tiene derecho a su propio AGENTS.md. El tooling que ayuda a CREAR
y HACER CUMPLIR ese AGENTS.md debe vivir en otro lugar.

---

## Comportamiento Global

Ver [behavior-scrum-master-routing.md](behavior-scrum-master-routing.md) —
el SM actúa como router de fases, convoca roles, valida gates, y bloquea
avances prematuros.

---

## Dos Modos

### Modo 1 — Planificación (idea → handoffs)

**Propósito**: producir fuentes de verdad y planes. Sin ejecución de código.

**Quién participa**: el scrum team (PO, Dev Lead, SM, UX, QA, DevSecOps)
como lentes de revisión — no como agentes que escriben código.

#### Entradas aceptadas

El modo planificación arranca desde cualquier nivel de definición.
**En esta etapa NO se detecta stack, arquitectura ni tecnologías.** Lo
único que ocurre es: crear la entrada inicial del proyecto en el RAG.

El sistema detecta el tipo de entrada y elige qué rol del agile team
la procesa:

| Nivel de entrada | Ejemplo | Rol asignado | Acción |
|-----------------|---------|-------------|--------|
| Idea vaga | "Haz el Uber de las lanchas" | PO | Formula preguntas de negocio al MIM para acotar alcance y valor |
| Archivos de un challenge | README.md + seeds + schema de un tech challenge | PO + SM | PO extrae requisitos y constraints. SM delega a sub-agente SM-Process la extraccion de reglas del proceso (timebox, evaluacion, restricciones) |
| Ticket externo | Link a Jira, Linear, Confluence, GitHub Issue | PO | Lee, estructura, identifica ambigüedades (vía adaptador TBD) |
| Especificación parcial | "API REST con auth JWT y CRUD de productos" | PO | Identifica gaps en los requisitos y pregunta solo lo faltante |

```mermaid
flowchart TD
    INPUT["Entrada del usuario"]
    DETECT["Detectar tipo de entrada"]
    ROUTE["Elegir rol(es) del agile team"]

    DETECT --> ROUTE
    INPUT --> DETECT

    ROUTE -->|idea vaga| PO_ONLY["PO: preguntas de negocio"]
    ROUTE -->|challenge| PO_SM["PO: requisitos\nSM: reglas de proceso"]
    ROUTE -->|ticket| PO_TICKET["PO: estructurar y desambiguar"]
    ROUTE -->|spec parcial| PO_GAPS["PO: identificar gaps"]

    PO_ONLY --> CREATE["Crear entrada inicial\nen RAG (docs/)"]
    PO_SM --> CREATE
    PO_TICKET --> CREATE
    PO_GAPS --> CREATE
```

El punto es: **no importa qué tan vago o preciso sea el input**. El sistema
detecta el nivel de definición, elige el rol correcto, y produce UNA cosa:
la entrada inicial del proyecto en el RAG (`idea.md`). Nada más.

Las decisiones técnicas (stack, arquitectura, patrones) NO pertenecen a
esta etapa. Llegan después, cuando el Dev Lead y DevSecOps entran en las
fases de diseño.

#### Flujo: de idea a handoff

```mermaid
flowchart TD
    INPUT["Entrada del usuario\n(idea, archivos, ticket, spec parcial)"]
    DETECT["Detectar nivel de definición"]
    QUESTIONS["Generar preguntas para el MIM\n(lo que falta para acotar)"]
    MIM["MIM responde\n(stakeholder o creador)"]
    DOC["Crear/actualizar entrada en RAG\n(docs/ por defecto)"]
    ENOUGH{{"¿Suficiente para\nla siguiente fase?"}}
    NEXT["Siguiente fase de planificación\n(spec, diseño, tareas...)"]

    INPUT --> DETECT
    DETECT --> QUESTIONS
    QUESTIONS --> MIM
    MIM --> DOC
    DOC --> ENOUGH
    ENOUGH -->|No| QUESTIONS
    ENOUGH -->|Sí| NEXT
    NEXT -->|produce required params\npara la siguiente fase| NEXT
```

#### El RAG como fuente de verdad progresiva

El artifact store NO es solo persistencia — es un **RAG** que los agentes
consultan para obtener contexto ACOTADO sin crawlear el codebase.

Principio fundamental: **cada fase consume el output de la anterior y
produce los required params de la siguiente**. Ningún agente necesita leer
"todo" — solo el slice que le corresponde.

```mermaid
flowchart LR
    subgraph RAG["RAG — docs/ (default adapter)"]
        direction TB
        R1["idea.md\n(inputs del usuario)"]
        R2["spec.md\n(ACs, contratos, constraints)"]
        R3["design.md\n(arquitectura, decisiones)"]
        R4["tasks.md\n(desglose ordenado)"]
        R5["handoff.md\n(contrato para ejecución)"]
    end

    R1 -->|"required params"| R2
    R2 -->|"required params"| R3
    R3 -->|"required params"| R4
    R4 -->|"required params"| R5
```

Evolución del contexto en cada fase:

| Fase | Consume del RAG | Produce al RAG | Quién consulta después |
|------|----------------|---------------|----------------------|
| Definir idea | Nada (input fresco del usuario) | `idea.md` — el problema, alcance, restricciones | Fase de spec |
| Especificar | `idea.md` | `spec.md` — ACs, contratos, constraints | Fase de diseño |
| Diseñar | `idea.md` + `spec.md` | `design.md` — arquitectura, patrones, tradeoffs | Fase de tareas |
| Desglosar tareas | `spec.md` + `design.md` | `tasks.md` — tareas ordenadas con dependencias | Fase de handoff |
| Generar handoff | `spec.md` + `design.md` + `tasks.md` | `handoff.md` — contrato autocontenido | Modo ejecución |

**Clave**: cuando un agente de ejecución necesita contexto, no lee 15
archivos del repo — hace fetch al RAG y obtiene exactamente el slice que
necesita. Al principio (definir idea) el RAG solo contiene los inputs del
usuario. Al final (handoff) contiene toda la cadena de decisiones.

#### Guía al MIM (generación de preguntas)

El sistema no espera que el MIM sepa qué preguntar. Según el nivel de
entrada detectado, genera preguntas dirigidas:

Para una **idea vaga** ("Uber de lanchas"):
- ¿Quién es el usuario final? (pasajeros, lancheros, ambos)
- ¿Cuál es el flujo core? (reservar, pagar, rastrear)
- ¿Qué plataforma? (web, mobile, ambos)
- ¿Hay restricciones técnicas? (stack, hosting, presupuesto)
- ¿MVP o producto completo? ¿Deadline?

Para un **tech challenge** (archivos del repo):
- ¿Cuál es el timebox?
- ¿Hay restricciones de stack no documentadas?
- ¿Qué se evalúa? (código, proceso, arquitectura, todo)
- ¿Se puede usar tooling de AI? ¿Con qué restricciones?

Para un **ticket externo**:
- ¿Los ACs están completos o hay ambigüedad?
- ¿Hay dependencias bloqueantes?
- ¿Quién aprueba el resultado?

Las preguntas se adaptan: si el MIM ES el stakeholder/creador, las responde
directamente. Si no lo es, las usa como guía para obtener las respuestas.

#### Adaptador por defecto: `docs/` como RAG local

- Path por defecto: `{repo}/docs/` — DENTRO del repositorio destino
  (corregido en commit 07eafad; ver también sección "Adaptadores del
  Artifact Store" más abajo para el detalle completo del adaptador)
- Formato: archivos markdown, uno por fase
- Legible por humanos, versionable con git
- Los agentes hacen fetch de archivos específicos, no crawl completo
- Adaptadores adicionales (engram, Jira, Confluence, etc.): ver sección
  "Adaptadores del Artifact Store"

**Restricción clave**: el modo planificación NUNCA toca el working tree del
repo destino. Lee el codebase para informar decisiones, pero toda la salida
va al artifact store — no a archivos `.tmp-*` dispersos en el repo.

#### Scrum team en este modo

Cada rol es un LENTE que revisa artefactos de planificación desde su
perspectiva. PO valida alcance contra valor de usuario. QA valida
testeabilidad. DevSecOps valida superficie de seguridad. Producen
veredictos de revisión, no código.

Los lentes se activan DESPUÉS de que una fase produce su artefacto —
revisan lo producido, no participan en la generación. Si un lente
encuentra un gap, el sistema regresa al ciclo de preguntas para esa fase.

---

### Modo 2 — Ejecución (handoffs → código funcional)

**Propósito**: implementar lo que la planificación produjo. Se escribe código.

**Quién participa**: el orquestador y los sub-agentes (minions).

**Qué consume**:
- Documentos de handoff del Modo 1
- El AGENTS.md del repo destino (reglas de gobernanza)
- Compact rules resueltas y estándares del proyecto

**Qué produce**:
- Cambios de código en el repo destino
- Resultados de pruebas
- Reportes de verificación
- Commits (siguiendo las convenciones del repo)

**Cómo funciona**:
- El orquestador lee el handoff y el AGENTS.md del repo
- El orquestador delega a sub-agentes con:
  - Alcance específico de tarea (del handoff)
  - Personalidad/contexto (apropiado para el rol de la tarea)
  - Compact rules inyectadas (de AGENTS.md, resueltas una vez)
  - Axiomas del proyecto inyectados (de AGENTS.md, si están definidos)
- Los sub-agentes escriben código, ejecutan pruebas, reportan
- El orquestador verifica, hace commit, pasa a la siguiente tarea

**Restricción clave**: el modo ejecución SOLO escribe en el working tree del
repo destino. NO crea artefactos de planificación, documentos de feedback ni
archivos de estado del proceso en el repo.

---

## Límites

```mermaid
flowchart TD
    subgraph PLAN["MODO PLANIFICACIÓN"]
        direction TB
        P_IN["Entrada: idea, problema, feature request"]
        P_TOOLS["Herramientas: idea, spec, diseño, tareas, handoff"]
        P_WHO["Participantes: scrum team como lentes de revisión"]
        P_OUT["Escribe en: artifact store — NUNCA en el repo destino"]
    end

    PLAN -->|handoff| EXEC

    subgraph EXEC["MODO EJECUCIÓN"]
        direction TB
        E_IN["Entrada: handoff + AGENTS.md del repo"]
        E_TOOLS["Herramientas: orquestador, sub-agentes, verificar, aceptar"]
        E_WHO["Participantes: orquestador + minions"]
        E_OUT["Escribe en: working tree del repo ÚNICAMENTE"]
    end

    EXEC -->|"gap detectado"| PLAN
```

---

## El Handoff como Contrato

El documento de handoff es la interfaz entre modos. Debe ser:

- **Autocontenido**: un ejecutor que nunca vio la conversación de
  planificación puede actuar sin hacer preguntas.
- **Portable**: funciona independientemente del adaptador que lo produjo.
- **Acotado**: dice exactamente qué hacer, qué NO hacer, y cómo se ve el
  éxito.

El handoff NO es un archivo en el repo destino. Vive en el artifact store
y es LEÍDO por el modo ejecución.

---

## Qué Vive DÓNDE

| Artefacto | Dónde vive | Por qué |
|-----------|-----------|---------|
| AGENTS.md | Repo destino (raíz) | La gobernanza es por repo. Cada proyecto posee sus reglas. |
| Artefactos de planificación (propuestas, specs, diseños, tareas) | Artifact store (depende del adaptador) | Informan el trabajo, no son el trabajo. |
| Documentos de handoff | Artifact store | Contrato entre planificación y ejecución. |
| Estado del ciclo (tracking de fases, DAG) | Artifact store | Estado operacional, no estado del proyecto. |
| Feedback de adoptantes | Artifact store (etiquetado al framework fuente) | Input para evolución del framework, no contenido del repo. |
| Código, pruebas, configs | Repo destino | El entregable real. |
| Archivos `.tmp-*` | EN NINGÚN LUGAR del repo destino | Eliminados. Los artefactos de planificación van al store. |

---

## Adaptadores del Artifact Store

El framework necesita una capa de persistencia pluggable. Cada adaptador
implementa la misma interfaz universal de 8 operaciones (ver
`artifact-model.md` → "Adaptadores de Persistencia" para la definición
completa: `save`, `read`, `search`, `list`, `markComplete`,
`verifyConsistency`, `delete`, `history`). Todas las operaciones de
escritura son mediadas por el TPM (ver `artifact-model.md` → "TPM como
DBMS"); las lecturas pueden ser directas vía Pattern B.

### Adaptador local (por defecto)
- Almacena artefactos como archivos markdown en `{repo}/docs/`
- Dentro del repositorio destino — agentes y usuarios acceden al mismo filesystem
- Ventajas: cero dependencias, legible por humanos, versionable con git
- Desventaja: sin acceso cross-machine, sin búsqueda semántica

### Adaptador engram
- Almacena artefactos como observaciones engram con topic keys estructurados.
- Ventajas: cross-session, buscable, sobrevive compaction.
- Desventaja: requiere servidor MCP de engram, contenido puede truncarse en
  resultados de búsqueda (se necesita `mem_get_observation` para contenido
  completo).

### Adaptador híbrido
- Escribe en ambos: local y engram.
- Ventajas: lo mejor de ambos — legibilidad local + persistencia
  cross-session.
- Desventaja: mayor costo de tokens por operación.

---

## Scrum Team — Cuándo y Cómo

El scrum team es una herramienta de PLANIFICACIÓN, no de ejecución.

| Rol | Cuándo se activa | Qué hace | Qué NO hace |
|-----|-----------------|----------|-------------|
| PO | Idea, Spec, Verificar, Aceptar, Retro | Valida alcance, prioriza, define ACs, acepta entregables | Escribir código, revisar PRs |
| Dev Lead | Diseño, Tareas, Verificar, Aceptar, Retro | Valida arquitectura, estima, secuencia, revisa calidad técnica | Ejecutar tareas (eso es modo ejecución) |
| SM | Todas las fases | Facilita, remueve bloqueos, valida proceso, orquesta gates | Producir contenido, leer archivos |
| UX | Spec, Diseño, Verificar, Aceptar, Retro | Valida decisiones que afectan al usuario | Implementar UI |
| QA | Spec, Tareas(cond), Verificar, Aceptar, Retro | Valida testeabilidad, define estrategia de pruebas, verifica cobertura | Escribir código de producción |
| DevSecOps | Diseño, Tareas(cond), Verificar, Aceptar, Retro | Valida superficie de seguridad, decisiones de infra, postura de seguridad | Desplegar |
| *Ad-hoc* | Cualquier fase, segun contrato | Expertise especializado fuera del equipo default (DBA, Performance Engineer, Domain Expert, etc.). El SM los define y convoca con contrato completo. | Depende del contrato |

> **Nota**: los 6 roles de arriba son el equipo **default**. El SM puede
> extender el equipo con roles ad-hoc cuando el proyecto requiere expertise
> que ningun rol default cubre. Ver `role-profiles.md` seccion "Roles Ad-Hoc".

**Durante ejecucion**, el scrum team esta en silencio. El orquestador y los
sub-agentes hacen el trabajo. Si la ejecución revela un gap de planificación,
el orquestador puede escalar DE VUELTA al modo planificación.

**Post-ejecución** (Verificar, Aceptar, Retrospectiva), el scrum team se
RE-ACTIVA como panel de revisión. Estas fases son parte del modo
planificación — operan sobre los resultados de la ejecución, no sobre
código directamente. Ver `behavior-scrum-master-routing.md` Fases 6-8 y
`role-profiles.md` para los contratos de delegación de cada rol en estas
fases.

---

## Qué Permanece en AGENTS.md (gobernanza por repo)

Estas son las cosas que cada repo adoptante recibe:

- Axiomas (principios no negociables)
- Fases del pipeline (la secuencia de trabajo)
- Gates de fase (DOR, DOD, checkpoints MIM)
- Compact rules (estándares de código específicos del proyecto)
- Definiciones de roles (qué valida cada rol en los gates)
- Tiers de activación (cuánta ceremonia según la madurez del proyecto)

Estas son REGLAS, no HERRAMIENTAS. Dicen qué debe ser verdad, no cómo
hacerlo verdad.

---

## Qué NO Va en AGENTS.md

- Definiciones de fases del ciclo de planificación (idea, spec, diseño, etc.)
- Configuración del artifact store
- Patrones de delegación del orquestador
- Templates de personalidad de sub-agentes
- Formatos de topic key de engram
- Protocolos de resolución de skills
- Tablas de asignación de modelos

Estos son OPERACIONALES. Pertenecen a la capa de tooling, no a la capa de
gobernanza.

---

## Estrategia de Delegación Multi-Modelo

El SM selecciona el tier de modelo por tarea usando un criterio simple:
**¿La salida correcta es derivable de reglas/templates, o requiere juicio?**

| Tier | Runtime | Criterio de selección | Costo |
|------|---------|----------------------|-------|
| **Local** (Docker / Ollama) | Modelo local, cero costo por token | La salida es determinista o template-driven. No requiere razonamiento complejo. | Cero (solo compute local) |
| **Cloud** (Claude / Codex / equivalente) | API remota, costo por token | Requiere síntesis, juicio, creatividad, o razonamiento sobre contexto ambiguo. | Proporcional al uso |

### Asignación por componente

| Componente | Tier | Justificación |
|-----------|------|---------------|
| **TPM** (validar formato, verificar schema, generar markdown, batch writes, slug) | Local | Operaciones mecánicas con reglas bien definidas |
| **Echo Protocol** — checks estructurales (completitud, formato, campos requeridos) | Local | Verificable con reglas |
| **Echo Protocol** — checks semánticos (coherencia, contradicciones, calidad) | Cloud | Requiere comprensión del contenido |
| **SM** (coordinación, decisiones de routing, gate evaluation) | Cloud | Requiere juicio sobre contexto |
| **PO** (spec desde input ambiguo, priorización, ACs) | Cloud | Síntesis y juicio |
| **Dev Lead** (diseño arquitectónico, estimación, secuenciación) | Cloud | Razonamiento técnico profundo |
| **QA** (adversarial review, estrategia de pruebas, verificación) | Cloud | Juicio y creatividad adversarial |
| **DevSecOps** (threat model, surface analysis) | Cloud | Razonamiento de seguridad |
| **UX** (validación de decisiones de usuario) | Cloud | Empatía y juicio de producto |
| **Retro** (síntesis stop/start/continue, agreements) | Cloud | Síntesis de múltiples perspectivas |

### Regla de decisión

```
if (output == template_con_slots && sin_ambiguedad)
  → Local
else
  → Cloud
```

El SM no necesita un scoring complejo. Si puede escribir el template y los
slots de antemano, la tarea es mecánica. Si necesita que el agente **piense**,
es cloud.

### Nota de implementación

La selección de modelo es una decisión de **tooling**, no de gobernanza.
Cada proyecto puede configurar qué modelo local usar (llama3, mistral,
phi, etc.) y qué proveedor cloud preferir. El framework define el CRITERIO
de selección, no el modelo específico.

---

## Preguntas Abiertas

1. **¿Dónde vive el tooling?**
   Opciones: skills de Claude Code (instalables), un paquete npm separado,
   una convención de dotfiles (`~/.idea-to-mvp/`), o una combinación.

2. **¿Cómo hace un repo "opt in" al framework?**
   Actualmente: copiar AGENTS.md. ¿Debería existir un comando bootstrap
   (`/sdd-init` o equivalente) que configure la capa de tooling sin
   contaminar el repo?

3. **¿Cómo fluye el feedback de vuelta al framework?**
   fullstack-base produjo feedback para idea-to-mvp. ¿Dónde vive ese
   feedback? ¿Cómo se rastrea? Actualmente es un archivo `.tmp-*` en el
   repo del framework — que es la misma contaminación que queremos eliminar.

4. **¿Debería estandarizarse el formato del handoff?**
   Si el handoff es el contrato entre modos, su estructura importa.
   ¿Un schema? ¿Un template? ¿Campos mínimos requeridos?

5. **¿Cómo afectan los tiers de activación a la separación de modos?**
   Un desarrollador solo en un proyecto de fin de semana no necesita una
   revisión de 6 roles. ¿Cómo escala hacia abajo el modo planificación de
   forma elegante?

6. ~~**¿Verificación en modo ejecución — quién la hace?**~~
   **RESUELTO**: Verify (Fase 6) y Accept (Fase 7) son fases
   POST-ejecución del modo planificación. El scrum team se reactiva
   como panel de revisión. Retro (Fase 8) cierra el ciclo y alimenta
   el siguiente. Ver `behavior-scrum-master-routing.md` Fases 6-8.
