# Review Adversarial #003 — Planning Mode Completo (Idea → Handoff)

> Fecha: 2026-07-30
> Archivos revisados: `artifact-model.md`, `behavior-scrum-master-routing.md`,
> `role-profiles.md`, `operational-model.md`
> Agentes: 3 (coherencia, completitud, implementabilidad)
> Hallazgos brutos: 38. Deduplicados: 28.
> Incluye: correcciones de roles ad-hoc recien agregadas.

---

## Resumen Ejecutivo

| Severidad | Cantidad | Corregidos | Diferidos |
|-----------|----------|------------|-----------|
| CRITICAL | 5 | 5 | 0 |
| HIGH | 10 | 10 | 0 |
| MEDIUM | 9 | 9 | 0 |
| LOW | 4 | 4 | 0 |
| **Total** | **28** | **28** | **0** |

**Veredicto post-fix: 28 de 28 hallazgos corregidos.** H1 (write path)
fue el ultimo en resolverse — ver "Write Path: Acoplado vs Desacoplado"
en `behavior-scrum-master-routing.md`.

### Causas Raiz — RESUELTAS

| # | Causa | Estado |
|---|-------|--------|
| **R1** | `tasks.md` menos reconciliado | **RESUELTO** — contrato Fase 4, gate, ownership matrix, schema, paralelismo todos actualizados |
| **R2** | Pattern A/B inconsistencia | **RESUELTO** — seccion TPM en behavior.md reescrita con Pattern B como default |
| **R3** | Valvulas de seguridad faltantes | **RESUELTO** — PARTIAL cap, smoke test rubric, rollback FF, mid-flight edit, todos agregados |

---

## CRITICAL — ~~Bloquean implementacion correcta~~ TODOS CORREGIDOS

### ~~C1. Work Item Hierarchy desconectada del contrato de Fase 4~~ CORREGIDO
**Fuente**: COH-004 + COMP-001

**Correccion**: contrato Dev Lead Fase 4 en `role-profiles.md`
reescrito con schema universal (L3/L4, parent_id, depends_on FS/SS/FF,
traces_to, lane, complexity XS-XL). Gate en `behavior.md` actualizado
para exigir el mismo schema. Ambos documentos ahora referencian el
mismo modelo de work items.

### ~~C2. SM listado como co-productor de `tasks.md`~~ CORREGIDO
**Fuente**: COH-001

**Correccion**: SM removido de co-produce en la Ownership Matrix de
`artifact-model.md`. Ahora: `Produce: Dev Lead`, `Co-produce: QA
(verificabilidad por tarea)`.

### ~~C3. Productor valida su propio artefacto en 3 de 6~~ CORREGIDO
**Fuente**: COH-003 + IMPL-001

**Correccion**: Ownership Matrix y Rule #1 actualizados en
`artifact-model.md`:
- `idea.md`: validador = SM (estructural via TPM) + QA (verificabilidad
  de restricciones)
- `design.md`: validador = SM (estructural) + DevSecOps (seguridad) +
  UX (experiencia, condicional)
- `tasks.md`: validador = SM (estructural) + QA (verificabilidad por
  tarea) — QA ahora obligatorio en Fase 4 para gate semantico

### ~~C4. Seccion TPM en behavior.md aun modela Pattern A~~ CORREGIDO
**Fuente**: COH-002 + IMPL-002

**Correccion**: diagrama de secuencia TPM en `behavior.md` reescrito
mostrando Pattern B (sub-agente lee directo via topic_key). Tabla de
operaciones actualizada: "Leer" agregado, "Servir contexto" acotado
a Pattern A (8+ consumidores o busqueda fuzzy).

### ~~C5. Smoke test sin criterio de fallo confiable~~ CORREGIDO
**Fuente**: COMP-005 + IMPL-003

**Correccion**: contrato explicito del sub-agente de smoke test
agregado a `behavior.md` con criterio PASS/FAIL basado en preguntas
bloqueantes y asunciones criticas. Instruccion explicita: "NO asumas
— lista la pregunta."

---

## HIGH — ~~Comportamiento incorrecto o fragil~~ 10/10 CORREGIDOS

### ~~H1. Write path fuerza contenido a traves del SM~~ CORREGIDO
**Fuente**: IMPL-002

**Correccion**: seccion "Write Path: Acoplado vs Desacoplado" agregada
en `behavior-scrum-master-routing.md` (PDC, punto 2.1). Regla de
umbral (reutiliza el threshold de 500 tokens de M14): por debajo del
umbral se mantiene el write path acoplado (sub-agente devuelve
contenido, SM hace ECHO inline); por encima, el sub-agente persiste
directo via TPM y el SM solo recibe el Status Report. El ECHO se
desacopla en dos verdictos delegados — estructural (TPM, equivalente
al tier Local de la Estrategia de Delegacion Multi-Modelo) y semantico
(sub-agente fresco que lee el artefacto via Patron B, tier Cloud) — el
SM nunca ve el contenido completo por encima del umbral. Nota explicita
agregada: la seccion define COMPORTAMIENTO (que debe pasar, cuando, en
que orden), no IMPLEMENTACION (que runtime ejecuta el sub-agente fresco,
que protocolo de transporte usa la lectura — eso se resuelve al disenar
el Modo de Ejecucion o al elegir adaptador concreto).

### ~~H2. Sin cap para loops infinitos de PARTIAL~~ CORREGIDO
**Correccion**: regla agregada en `behavior.md`: 3 PARTIAL consecutivos
con mismo progreso (X/Y sin cambio) se tratan como FAILED y disparan
el circuit breaker.

### ~~H3. Ownership del paralelismo~~ CORREGIDO
**Correccion**: `artifact-model.md` clarificado — Dev Lead produce el
grafo de dependencias (Fase 4), el orquestador de ejecucion lo consume.
Campo `lane` agregado al schema de work items. Paso 3 del algoritmo
ahora referencia `lane` en vez de campo inexistente.

### ~~H4. Sin protocolo para reabrir artefacto mid-planning~~ CORREGIDO
**Correccion**: fila agregada a tabla de anomalias en `behavior.md`:
SM marca artefacto como "en revision", re-convoca productor original,
cascada via verifyConsistency, fase actual se pausa.

### ~~H5. Fast-forward F1 sin score para parciales~~ CORREGIDO
**Correccion**: nota agregada — artefacto incompleto = 0.5 puntos,
redondeo al entero mas cercano. "Incompleto" = TPM reporta secciones
faltantes.

### ~~H6. Sin protocolo para requests bundled~~ CORREGIDO
**Correccion**: regla agregada en Fase 1 de `behavior.md`: SM
descompone requests compuestos en L1 features independientes, informa
al MIM, ejecuta en secuencia o paralelo segun dependencias.

### ~~H7. TPM criterio editorial no operacionalizado~~ CORREGIDO
**Correccion**: checklist concreto de 4 reglas agregado a
`artifact-model.md`: oraciones unicas, sin TODOs en completos,
referencias con IDs trazables, formato Markdown consistente.

### ~~H8. Contratos condensados sin template~~ CORREGIDO
**Correccion**: ejemplo completo de PO condensado (Fase 1+2) agregado
a `role-profiles.md` con regla de transicion de personalidad explicita.

### ~~H9. Score fast-forward no se persiste~~ CORREGIDO
**Correccion**: instruccion agregada en `behavior.md` — SM instruye al
TPM para persistir score F1-F4 en idea.md seccion "Decisiones tomadas"
con formato auditable.

### ~~H10. operational-model.md "SM extrae"~~ CORREGIDO
**Correccion**: cambiado a "SM delega a sub-agente SM-Process" en
`operational-model.md`.

---

## MEDIUM — ~~TODOS CORREGIDOS~~

### ~~M1. operational-model.md omite roles ad-hoc~~ CORREGIDO
**Correccion**: fila Ad-hoc agregada a tabla de scrum team + nota sobre
extensibilidad en `operational-model.md`.

### ~~M2. Sin protocolo para desactivar rol mid-cycle~~ CORREGIDO
**Correccion**: punto 4 agregado a reglas de activacion en
`role-profiles.md`: protocolo de 4 pasos (documentar, conservar
artefactos, remover de Fase 7, notificar MIM).

### ~~M3. Sin resolucion de conflicto pre-Fase 7~~ CORREGIDO
**Correccion**: seccion agregada en `role-profiles.md`: conflicto
tecnico → MIM decide, alcance → PO prioridad, seguridad → DevSecOps
prioridad, otros → MIM.

### ~~M4. Ad-hoc sin clausula de escalacion upstream~~ CORREGIDO
**Correccion**: campo "Escalacion upstream" agregado como obligatorio
en template de contrato ad-hoc en `role-profiles.md`.

### ~~M5. Sin recovery si MIM rechaza fast-forward~~ CORREGIDO
**Correccion**: protocolo de rollback agregado en `behavior.md`:
identificar artefactos con asunciones, marcar en revision, re-evaluar
F1-F4, retomar desde el primer artefacto afectado.

### ~~M6. Ad-hoc no en flowchart de Fase 1~~ CORREGIDO
**Correccion**: rama Q_ADHOC agregada al flowchart de activacion en
`role-profiles.md`.

### ~~M7. Edicion mid-flight no cubierta~~ CORREGIDO
**Correccion**: fila agregada a tabla de anomalias en `behavior.md`:
SM encola el edit, aplica PDC al retorno, evalua invalidacion.

### ~~M8. Fase 7 sin desempate~~ CORREGIDO
**Correccion**: regla agregada en `behavior.md`: empate → escalar al
MIM, default conservador = REQUEST CHANGES.

### ~~M9. F2/F3 scoring subjetivo~~ CORREGIDO (parcial)
**Correccion**: 3 ejemplos de frontera F2/F3 agregados (Stripe, Winston,
REST→GraphQL) con explicacion de scoring por factor. Subjetividad
reducida pero no eliminada — documentada como limitacion conocida.

---

## LOW — ~~TODOS CORREGIDOS~~

### ~~L1. TPM ops mismatch~~ CORREGIDO
**Correccion**: "Leer" agregado a tabla de operaciones TPM en
`behavior.md`. Ahora 7 ops en ambos documentos.

### ~~L2. "SM-analista" vs "SM-Process"~~ CORREGIDO
**Correccion**: unificado a "SM-Process" en `behavior.md`.

### ~~L3. SM-Process sin campos obligatorios~~ CORREGIDO
**Correccion**: "NO hace" y "Status Report" agregados al contrato
SM-Process en `behavior.md`.

### ~~L4. Artefacto vacio no distinguido~~ CORREGIDO
**Correccion**: fila agregada a tabla de anomalias: artefacto vacio
(0 secciones) se trata como "no existe" para la state machine.

---

## Patrones Observados

### Patron 1 (recurrente): Injerto sin propagacion — CORREGIDO
La Work Item Hierarchy, Pattern B, y "SM extrae" ahora estan propagados
a los 4 documentos.

**Regla adoptada**: todo cambio a un artefacto en `artifact-model.md`
DEBE incluir actualizacion correspondiente en (1) el contrato de
delegacion en `role-profiles.md`, (2) el gate en `behavior.md`, y
(3) la referencia en `operational-model.md`.

### Patron 2: operational-model.md como rezagado cronico — CORREGIDO
Actualizado con SM-Process, roles ad-hoc. Tres reviews consecutivos
lo detectaron — la regla del Patron 1 previene recurrencia.

### Patron 3: gates semanticos sin validador independiente — CORREGIDO
Los 6 artefactos ahora tienen validadores independientes asignados en
la Ownership Matrix.

---

---

## Pendientes Unificados — Reviews 001 + 002 + 003

Lista consolidada de TODO lo que queda realmente pendiente o diferido
a traves de los 3 adversarial reviews. Items resueltos por trabajo
posterior fueron actualizados en sus reviews respectivos.

### DIFERIDOS — TODOS RESUELTOS

| # | Origen | Hallazgo | Tipo | Nota |
|---|--------|----------|------|------|
| ~~D1~~ | ~~R003-H1~~ | ~~Write path fuerza contenido a traves del SM~~ | ~~RESUELTO~~ | ~~Write path desacoplado con umbral de 500 tokens (reutiliza M14): ECHO estructural (TPM, tier Local) + ECHO semantico (sub-agente fresco via Patron B, tier Cloud). SM nunca ve el contenido completo por encima del umbral. Ver `behavior-scrum-master-routing.md` seccion 2.1.~~ |
| ~~D2~~ | ~~R002-H3~~ | ~~Personalidades no verificables~~ | ~~RESUELTO~~ | ~~Output constraints por personalidad en role-profiles.md~~ |
| ~~D3~~ | ~~R002-M3~~ | ~~TPM editorial drift~~ | ~~RESUELTO~~ | ~~Ediciones nivel 1 (formato) vs nivel 2 (estructura) con notificacion~~ |
| ~~D4~~ | ~~R002-M6~~ | ~~Multi-stakeholder MIM~~ | ~~RESUELTO~~ | ~~Routing requester/approver en Fase 1~~ |
| ~~D5~~ | ~~R002-M13~~ | ~~TPM dispatch budget~~ | ~~RESUELTO~~ | ~~3 tiers de batch writes (normal/comprimido/ultra-comprimido)~~ |
| ~~D6~~ | ~~R002-M14~~ | ~~Token overhead Pattern B~~ | ~~RESUELTO~~ | ~~Threshold 500 tokens: < 500 → Pattern A, >= 500 → Pattern B~~ |
| ~~D7~~ | ~~R001-M1~~ | ~~RAG multi-dev~~ | ~~RESUELTO~~ | ~~Capa de adaptadores con 6 implementaciones~~ |

### ABIERTOS — TODOS DIFERIDOS O RESUELTOS

| # | Origen | Hallazgo | Tipo | Nota |
|---|--------|----------|------|------|
| A1 | R001-H3 | Escalacion ejecucion → planificacion sin protocolo completo | DIFERIDO (depende de Modo de Ejecucion) | `escalate(gap)` existe con 3 targets. Protocolo completo requiere disenar el Execution Mode. |
| A2 | R001-H5 | Preguntas pre-definidas de Fase 1 son rigidas | DIFERIDO (concern de producto) | Las 6 preguntas del PO son una lista estatica. No se adaptan al tipo de proyecto. |
| A3 | R001-H6 | Carga cognitiva insostenible para nuevos usuarios | DIFERIDO (concern de producto) | 4 docs, ~5000 lineas, sin onboarding guide ni quick-start. |
| ~~A4~~ | ~~R001-M4~~ | ~~Feedback del proceso~~ | ~~RESUELTO~~ | ~~Fase 8 Retrospectiva~~ |
| ~~A5~~ | ~~R001-C5p~~ | ~~Estrategia de delegacion multi-modelo no documentada~~ | ~~RESUELTO~~ | ~~Documentada en `operational-model.md`: tier Local (Docker/Ollama) para tareas mecanicas, tier Cloud (Claude/Codex) para tareas de razonamiento. Criterio: determinismo vs juicio (commit 500b79f).~~ |
| A6 | R002-C9 | Verify/Accept/Retro contradicen operational-model.md | DIFERIDO (depende de Modo de Ejecucion) | TBD hasta disenar el modo ejecucion. Misma razon que A1. |

### COSMETICOS — ~~Review-002 LOW~~ TODOS CORREGIDOS

| # | Origen | Hallazgo | Estado |
|---|--------|----------|--------|
| ~~L1~~ | ~~R002-L1~~ | ~~Diagrama "Planificacion"~~ | CORREGIDO — renombrado "Ciclo Completo" |
| ~~L2~~ | ~~R002-L2~~ | ~~Verbos es/en~~ | CORREGIDO — ya estandarizado |
| ~~L3~~ | ~~R002-L3~~ | ~~Slugificacion~~ | CORREGIDO — regla definida (lowercase, max 64, colision -N) |
| ~~L4~~ | ~~R002-L4~~ | ~~Phase naming drift~~ | CORREGIDO — ya alineado |
| ~~L5~~ | ~~R002-L5~~ | ~~Sources self-reported~~ | CORREGIDO — `read` en history() para cross-check |
| ~~L6~~ | ~~R002-L6~~ | ~~DevSecOps threshold~~ | CORREGIDO — checklist de 6 criterios concretos |

### Totales

| Categoria | Cantidad | Resueltos | Diferidos |
|-----------|----------|-----------|-----------|
| Diferidos (arquitectura) | 7 | 7 | 0 |
| Abiertos (producto/diseno) | 6 | 2 (A4, A5) | 4 (A1, A2, A3, A6) |
| Cosmeticos | 6 | 6 | 0 |
| **Total** | **19** | **15** | **4** |

**Resumen**: de ~81 hallazgos brutos a traves de 3 reviews, quedan **4**
diferidos a fases futuras — ninguno permanece abierto sin clasificar.
D1 (write path) se resolvio (write path desacoplado con umbral, ver
`behavior-scrum-master-routing.md` seccion 2.1). A1 y A6 requieren
decision arquitectonica sobre el Modo de Ejecucion. A2 y A3 son
concerns de producto para iteracion futura (preguntas adaptativas,
onboarding guide).

**Todos los hallazgos han sido resueltos o diferidos a fases futuras
(Modo Ejecución, decisiones de producto). Las revisiones adversariales
del Modo de Planificación están cerradas.**

---

## Estado Final del Planning Mode (Idea → Handoff)

| Aspecto | Estado |
|---------|--------|
| Artefactos (6) con schema ISO | Completo |
| Ownership sin contradicciones | Completo |
| Gates con validador independiente | Completo |
| Work Item Hierarchy integrada end-to-end | Completo |
| Pattern A/B consistente | Completo |
| Roles default (5) + ad-hoc extensible | Completo |
| Fast-forward via state machine (no heuristicas) | Completo |
| Circuit breaker (FAILED + PARTIAL) | Completo |
| Smoke test con rubric confiable | Completo |
| Mid-planning edit y mid-flight edit | Completo |
| Bundled requests | Completo |
| Role deactivation | Completo |
| Conflict resolution pre-Fase 7 | Completo |
| Condensed contracts | Completo (1 ejemplo) |
| Adapter ACID + transactions | Completo |
| State machine configurable (draft→approved) | Completo |
| TPM ingest (material del MIM) | Completo |
| Batch writes por tier (normal/comprimido/ultra) | Completo |
| Personalidades verificables (output constraints) | Completo |
| TPM anti-drift (nivel 1 vs nivel 2) | Completo |
| Multi-stakeholder routing | Completo |
| Retrospectiva con meta-config | Completo |
| Estrategia de delegacion multi-modelo (Local/Cloud) | Completo |
| Write path desacoplado (umbral 500 tokens, ECHO estructural + semantico) | Completo |
