# Review Adversarial #002 — Modelo Completo de Planificación (Idea → Handoffs)

> Fecha: 2026-07-30
> Archivos revisados: `artifact-model.md`, `behavior-scrum-master-routing.md`,
> `role-profiles.md`, `operational-model.md`, `adversarial-review-001.md`
> Agentes: 3 (coherencia interna, completitud, implementabilidad)
> Total hallazgos brutos: 48. Consolidados (deduplicados): 38.
> Corregidos en sesión: 4. Restantes: 34.

---

## Estado del Review #001

| ID | Estado | Nota |
|----|--------|------|
| C1–C5 | RESUELTOS | Ver review-001 para resoluciones |
| H1 (TPM ausente en operational-model) | **VIGENTE** → H9 abajo |
| H2 (Gates vacíos) | Mejorado — gates ahora tienen criterios concretos |
| H3 (Sin protocolo de escalación) | Parcialmente resuelto — `escalate()` existe pero thin |
| H4 (Tablas de roles inconsistentes) | **VIGENTE, PEOR** → C6, H8 abajo |
| H5, H6 (Preguntas rígidas, carga cognitiva) | Sin cambio — concern de producto |
| H7 (Terminología SDD residual) | **VIGENTE** → C7 abajo |
| M1–M4 | Sin cambio |
| M5 (docs/ vs artifacts/) | **VIGENTE** → M7 abajo |
| M6 (Interfaz de adaptador) | Resuelto en artifact-model, NO en operational-model → M8 |

---

## CORREGIDOS EN ESTA SESIÓN

| ID | Hallazgo | Corrección |
|----|----------|-----------|
| FIX-1 | TPM read/write contradice Pattern B | Regla #3 reescrita: TPM = SOLO escritura, lectura libre via Pattern B |
| FIX-2 | Diagrama cadena con "territorio libre" stale | Etiquetas actualizadas a ISO 29148 §9.3 / ISO 21502 §7.6 |
| FIX-3 | Jerarquía de work items ausente | Nueva sección completa: 5 niveles (L0-L4), schema, dependencias, DAG, paralelismo |
| FIX-4 | Sin procedimiento de selección de metodología | Regla explícita: MIM elige → Scrum default → informa → cambiable en boundary |

---

## CRITICAL — Bloquean implementación coherente

### C1. SM "cero excepciones" vs "se convoca a sí mismo" en Fase 1

**Fuente**: coherencia + implementabilidad (convergencia 2/3 revisores)

La regla cardinal dice "El SM no lee archivos, no escribe, no produce.
**CERO excepciones.**" Fase 1 dice "el SM se convoca a sí mismo para
extraer reglas de proceso de un tech challenge." Extraer reglas de un
challenge requiere leer archivos — exactamente lo prohibido.

- **Docs afectados**: `behavior-scrum-master-routing.md` (Fase 1),
  `role-profiles.md` (diagrama Fase 1)
- **Impacto**: un agente implementador enfrenta dos directivas
  mutuamente excluyentes en el mismo archivo
- **Resolución**: → **CORREGIDO** — reescrito como delegación explícita
  a sub-agente con contrato SM-analista

### C2. Fast-forward sin algoritmo evaluable — solo adjetivos

**Fuente**: completitud + implementabilidad (convergencia 2/3 revisores)

El gradiente de certeza (baja/media/alta) es la decisión de mayor
impacto del sistema (determina si se consulta al humano o se salta
directo a ejecución) y no tiene scoring, thresholds, ni ejemplos
resueltos de frontera.

- **Docs afectados**: `behavior-scrum-master-routing.md` (Fast-Forward)
- **Impacto**: dos corridas del mismo modelo con el mismo input pueden
  clasificar diferente. Decisiones no reproducibles ni auditables.
- **Resolución**: → **CORREGIDO** — checklist de 4 factores con scoring
  0/1/2, thresholds numéricos, y 3 ejemplos resueltos de frontera

### C3. State reconstruction indefinida para RAG inconsistente

**Fuente**: implementabilidad

Solo hay un path feliz para derivar estado. No hay branch para:
spec.md completo pero idea.md ausente, dos artefactos en progreso
simultáneamente, o design.md completo pero inconsistente con spec.md
editado después.

- **Docs afectados**: `behavior-scrum-master-routing.md` (State Machine)
- **Impacto**: crash, compaction, o corrupción producen estado indefinido
- **Resolución**: → **CORREGIDO** — tabla de anomalías con branch
  CORRUPT → escalación a MIM + definición mecánica de "completo"

### C4. Circuit breaker sin hogar de persistencia

**Fuente**: implementabilidad

El contador de 3 fallos consecutivos es estado transitorio que NO es
un artefacto ni se deriva del RAG. Si hay compaction/crash entre fallo
#2 y #3, el counter se pierde y el circuit breaker se resetea
silenciosamente — contradice el claim de "context resilience."

- **Docs afectados**: `behavior-scrum-master-routing.md` (Supervisión)
- **Impacto**: circuit breaker no funciona cross-session como se
  documenta
- **Resolución**: → **CORREGIDO** — alcance explícito: counter es
  de sesión. Cross-session, el TPM trackea historial de fallos como
  metadata del artefacto.

### C5. Methodology swap: claim aspiracional, routing Scrum-only

**Fuente**: implementabilidad + coherencia (convergencia 2/3 revisores)

El routing file tiene una state machine de 8 fases con gates discretos
= Scrum puro. No hay routing alternativo para Kanban (flujo continuo,
WIP limits), Shape Up (betting table, appetite), ni PI Planning.
Además, la tabla de metodología dice Kanban tiene "sin roles fijos"
pero el SM persiste como "ÚNICO rol constante."

- **Docs afectados**: `artifact-model.md` (tabla de metodología),
  `behavior-scrum-master-routing.md` (state machine completa)
- **Impacto**: la abstracción funciona a nivel de artefactos pero es
  aspiracional a nivel de orquestación
- **Resolución**: → **CORREGIDO** — claim acotado: "artefactos son
  metodología-agnostic; routing default es Scrum; routing para Kanban/
  Shape Up/SAFe es extensible pero no implementado aún." Tabla de
  roles corregida: roles son funciones constantes con nombres
  methodology-específicos.

### C6. QA en Fase 3 (Design): afirmado y negado en el mismo archivo

**Fuente**: coherencia

3/4 fuentes (convocation map, role-profiles, operational-model) dicen
NO QA en Design. 1 tabla (Matriz Completa) dice SÍ. Un SM que lea la
Matriz Completa invocaría QA sin contrato definido.

- **Docs afectados**: `behavior-scrum-master-routing.md` (Matriz
  Completa vs Mapa de Convocatoria)
- **Impacto**: delegación sin contrato = output impredecible
- **Resolución**: → **CORREGIDO** — QA eliminado de Fase 3 en Matriz
  Completa, consistente con las otras 3 fuentes

### C7. operational-model.md usa vocabulario SDD huérfano

**Fuente**: coherencia (confirma review-001 H7)

"explorar, proponer" no existen como fases en el modelo de 8 etapas
ni como artefactos. Son residuos de los skills SDD (`/sdd-explore`,
`/sdd-propose`).

- **Docs afectados**: `operational-model.md` (Límites, problem statement)
- **Impacto**: un lector buscaría fases Explorar/Proponer que no existen
- **Resolución**: → **CORREGIDO** — reemplazado con nombres canónicos
  (Idea, Spec, Design, Tasks, Handoff)

### C8. ops-runbook.md sin fase ni trigger en la state machine

**Fuente**: completitud + coherencia (convergencia 2/3 revisores)

Definido como artefacto universal con ISO 20000 + ITIL 4 pero ninguna
fase del SM lo produce, ningún contrato de rol lo cubre.

- **Docs afectados**: `behavior-scrum-master-routing.md`,
  `role-profiles.md`, `artifact-model.md` (Pregunta Abierta #1)
- **Impacto**: artefacto estructuralmente huérfano
- **Resolución**: → **CORREGIDO** — ops-runbook explícitamente
  post-ejecución. Se produce en Verify/Accept cuando ya hay código.
  Pregunta abierta #1 resuelta.

### C9. Verify/Accept/Retro contradicen operational-model.md

**Fuente**: completitud

operational-model.md: "el scrum team está EN SILENCIO post-handoff."
role-profiles.md + behavior: contratos completos para 5 roles en
Fase 6/7/8. Dos documentos autoritativos dan instrucciones
contradictorias para el mismo punto del flujo.

- **Docs afectados**: `operational-model.md` (two-mode model)
- **Impacto**: un agente leyendo solo operational-model.md concluye que
  el scrum team muere después del handoff
- **Resolución**: → **CORREGIDO** — operational-model.md actualizado:
  V/A/R son fases del modo planificación que operan POST-ejecución.
  El scrum team se reactiva para review, no se mantiene activo
  durante ejecución.

---

## HIGH — Ambigüedad que un agente interpretaría mal

### H1. Fase 7 paralelo viola regla secuencial del PDC

**Fuente**: implementabilidad

PDC dice "no se puede lanzar otro sub-agente sin completar los 4 pasos
del anterior." Fase 7 lanza 3-5 delegaciones en paralelo. Ningún
cross-reference. Mecánica de merge para fallo parcial no definida.

- **Resolución**: → **CORREGIDO** — excepción de Fase 7 documentada
  dentro de la sección PDC. Regla de merge: voto faltante = BLOCK
  implícito, re-delegación solo del rol faltante.

### H2. Handoff self-containment: validado por quien lo produjo

**Fuente**: implementabilidad

TPM produce handoff.md. SM valida "autocontención" pero no puede leer
archivos (regla cardinal). Entonces confía en el reporte del TPM sobre
su propio trabajo. Viola regla #1 (quien produce nunca valida su
propio artefacto).

- **Resolución**: → **CORREGIDO** — se agrega un smoke test: sub-agente
  fresco recibe SOLO handoff.md y debe generar un plan de ejecución.
  Si necesita hacer preguntas → handoff no es autocontenido.

### H3. Personalidades de rol son prosa, no rúbricas verificables

**Fuente**: implementabilidad

"Escéptico" vs "Riguroso" no produce diferencia verificable en output.
Solo "conciso, 1-3 oraciones" (Fase 7) es realmente checkeable.

- **Resolución**: → **DIFERIDO** — se marca como mejora iterativa.
  Nota: las personalidades NO son vacías — orientan el tono y foco del
  agente. Pero convertirlas en checklists verificables es trabajo de
  segunda iteración. Se agrega nota en role-profiles.md.

### H4. TPM conflate completitud estructural y semántica

**Fuente**: implementabilidad

"completo" puede significar "5/5 secciones existen" (mecánico) o "sin
ambigüedades" (juicio semántico). El framework no distingue cuál check
aplica en cada gate.

- **Resolución**: → **CORREGIDO** — gates ahora explícitamente
  separados en: completitud estructural (TPM, mecánica: secciones
  requeridas presentes) + completitud semántica (rol validador,
  juicio: contenido coherente y sin ambigüedades).

### H5. Sin cascada/invalidación al editar artefacto upstream

**Fuente**: completitud

Si idea.md cambia después de que spec.md existe, nada invalida los
artefactos downstream. `verifyConsistency` solo se ejecuta pre-handoff.

- **Resolución**: → **CORREGIDO** — TPM Update operation ahora
  trigger obligatorio de `verifyConsistency` en artefactos downstream.
  Resultado: report de stale artifacts al SM.

### H6. Sin re-evaluación de roles condicionales mid-project

**Fuente**: completitud

UX/DevSecOps se activan en Fase 1 y solo se re-evalúan en retrospectiva.
Si el scope cambia mid-ciclo (CLI→UI), no hay trigger.

- **Resolución**: → **CORREGIDO** — cualquier rol o el SM puede flaggear
  "scope changed, re-evaluate activation" como escalación inmediata,
  sin esperar retro.

### H7. Handoff: instrucciones del SM omiten estrategia de pruebas

**Fuente**: coherencia

artifact-model.md requiere "Estrategia de pruebas" en handoff.md (ISO
29119-3). behavior-scrum-master-routing.md Fase 5 lista 6 items al TPM
pero NO incluye estrategia de pruebas.

- **Resolución**: → **CORREGIDO** — agregado como item #7 en
  instrucciones de Fase 5.

### H8. Fase 4: DevSecOps/QA condicionales sin anchor en convocation map

**Fuente**: coherencia

role-profiles.md tiene contratos completos para DevSecOps y QA en
Fase 4 (condicionales). behavior-scrum-master-routing.md Fase 4 solo
lista Dev Lead. El SM no sabría cuándo invocarlos.

- **Resolución**: → **CORREGIDO** — DevSecOps/QA agregados como
  participantes condicionales en Fase 4 del behavior file.

### H9. TPM ausente de operational-model.md (confirma review-001 H1)

**Fuente**: coherencia

operational-model.md describe adaptadores invocados directamente.
artifact-model.md y behavior insisten en TPM como intermediario
obligatorio.

- **Resolución**: → **CORREGIDO** — TPM agregado a operational-model.md
  como capa entre roles y adaptadores. Referencia cruzada a
  artifact-model.md para definición completa.

---

## MEDIUM — Debería corregirse pero no bloquea

### M1. MIM dice "no sé" — sin protocolo

**Fuente**: completitud

Si el MIM no puede responder una pregunta de gate, no hay protocolo.
El SM se bloquearía indefinidamente.

- **Resolución**: → **CORREGIDO** — regla: PO registra asunción
  explícita en "Decisiones tomadas" como riesgo flaggeado. Gate se
  satisface con la asunción documentada. No se bloquea indefinidamente.

### M2. Sin gate humano entre handoff y ejecución

**Fuente**: completitud

La transición Handoff → Execution es automática. No hay "MIM confirma:
empezar ejecución ahora."

- **Resolución**: → **CORREGIDO** — gate explícito de confirmación MIM
  antes de transición a Execution Mode.

### M3. TPM editorial vs intención del productor

**Fuente**: completitud

TPM aplica "criterio editorial" sin que el rol productor revise la
versión editada. Posible drift semántico silencioso.

- **Resolución**: → **DIFERIDO** — bajo riesgo en la práctica (TPM
  edita formato, no semántica). Se agrega nota: "ediciones TPM que
  cambien estructura/contenido requieren confirmación del productor."

### M4. Cadena de dependencias no enforced mecánicamente

**Fuente**: completitud

Nada impide al TPM crear design.md si spec.md no está completo. La
cadena es descriptiva, no enforced.

- **Resolución**: → **CORREGIDO** — TPM Create valida completitud de
  artefactos upstream como precondición. Rechaza/advierte si se viola.

### M5. Sin modelo de concurrencia/locking para RAG

**Fuente**: completitud

Nada aborda dos sesiones concurrentes escribiendo al mismo artefacto.

- **Resolución**: → **CORREGIDO** — explícitamente fuera de scope para
  adaptador local ("sesión activa única asumida"). Concurrency safety
  diferida a adaptadores futuros (Jira, DBMS).

### M6. MIM como entidad única vs aprobador distinto

**Fuente**: completitud

Fase 1 pregunta "¿quién aprueba el resultado final?" pero luego MIM
es un interlocutor único para todo. Si requester ≠ approver, no hay
routing.

- **Resolución**: → **DIFERIDO** — scope actual asume MIM = persona
  única. Multi-stakeholder es extensión futura. Se documenta como
  limitación conocida.

### M7. Path default: docs/ vs artifacts/

**Fuente**: coherencia (confirma review-001 M5)

operational-model.md usa `artifacts/` en un lugar, `docs/` en otro.
artifact-model.md canonicaliza `docs/`.

- **Resolución**: → **CORREGIDO** — `artifacts/` reemplazado por
  `docs/` en operational-model.md.

### M8. Interfaz de adaptador: 4 ops vs 7 ops

**Fuente**: coherencia

operational-model.md lista 4 operaciones. artifact-model.md define 7.

- **Resolución**: → **CORREGIDO** — operational-model.md referencia la
  interfaz completa de 7 ops en artifact-model.md.

### M9. TPM Delete sin primitiva en interfaz de adaptador

**Fuente**: coherencia

TPM Operations tiene "Delete" pero la interfaz universal no tiene
delete.

- **Resolución**: → **CORREGIDO** — `delete(artifact, reason)` agregado
  a la interfaz de adaptador.

### M10. Breakeven N=7 no coincide con roster real de Fase 7

**Fuente**: coherencia + implementabilidad

El breakeven dice 7+ para Pattern A, pero Fase 7 tiene 5 roles, no 7.
El ejemplo propio del framework no alcanza su propio threshold.

- **Resolución**: → **CORREGIDO** — aclarado: Fase 7 con 5 roles usa
  Pattern B (bajo el breakeven). Pattern A solo para escenarios reales
  de fan-out alto (custom roles, multi-team reviews).

### M11. idea.md Owner lista SM antes que PO

**Fuente**: coherencia

"Owner: SM (orquesta) → PO (formula)" sugiere que SM produce, violando
regla cardinal. La Matriz Ownership dice "Produce: PO, Co-produce: SM."

- **Resolución**: → **CORREGIDO** — Owner field alineado con la Matriz:
  "Produce: PO. Co-produce: SM (formula preguntas)."

### M12. "SM no mantiene estado" ambiguo: ¿sesión o cross-session?

**Fuente**: implementabilidad

Leído literalmente, prohíbe caching dentro de una sesión continua,
forzando dispatch al TPM para cada pregunta de estado trivial.

- **Resolución**: → **CORREGIDO** — clarificado: "El SM no persiste
  estado cross-session. Dentro de una sesión, puede cachear el último
  status report del TPM y re-consultar solo cuando el estado puede
  haber cambiado."

### M13. Overhead de TPM / budget de delegaciones por tier sin cuantificar

**Fuente**: implementabilidad

Un ciclo completo realista puede requerir 30-50+ dispatches. Para un
challenge con timebox, eso podría consumir el presupuesto antes de
escribir código.

- **Resolución**: → **DIFERIDO** — requiere benchmarking empírico.
  Se agrega nota: "comprimir interacción TPM en tiers bajos es una
  optimización pendiente (batch writes por fase)."

### M14. Token economics omite overhead de reasoning para Pattern B

**Fuente**: implementabilidad

El modelo de costos solo cuenta tokens movidos, no reasoning del agente
para decidir qué queries hacer. Para artefactos pequeños, el overhead
podría dominar.

- **Resolución**: → **DIFERIDO** — requiere medición empírica. Se
  agrega nota sobre threshold mínimo de tamaño de artefacto.

---

## LOW — Cosmético

### L1. Diagrama "Planificación" incluye Ejecución como stage

**Fuente**: coherencia — sección titulada "Ciclo de Planificación"
contiene nodo Ejecución dentro del loop. Renombrar a "Ciclo Completo."

### L2. Interfaz adaptador: verbos en español vs inglés

**Fuente**: coherencia — operational-model.md usa "guardar, leer,
buscar, listar", artifact-model.md usa "save, read, search, list."
Estandarizar a inglés.

### L3. Slugificación de nombre de proyecto no definida

**Fuente**: completitud — topic_keys asumen slug limpio. Sin regla de
normalización ni manejo de colisiones.

### L4. Phase naming drift

**Fuente**: completitud — operational-model.md usa "Propuesta" donde
los otros docs usan "Fase 1 — Definir Idea."

### L5. Sources field es self-reported sin verificación

**Fuente**: implementabilidad — sub-agentes reportan qué consultaron
en Sources pero el SM no puede cross-check.

### L6. DevSecOps "no triviales" sin threshold operacional

**Fuente**: implementabilidad — "se activa si infra no trivial" sin
definir qué es "no trivial."

---

## Patrones observados

### Patrón 1: Documentos divergentes (PEOR que en review-001)

operational-model.md fue el primer documento y nunca se reconcilió
completamente con los otros tres. Es la fuente de la mayoría de
contradicciones (C7, C9, H9, M7, M8, L1, L2, L4). Recomendación:
reescribir operational-model.md como documento de referencia rápida
que LINKE a los otros tres en vez de duplicar definiciones.

### Patrón 2: Prose-based decisions (RECURRENTE)

Fast-forward gradient, personality shifts, "no triviales", "criterio
editorial" — el framework delega decisiones a "juicio del agente" sin
dar rúbricas verificables. Esto funciona cuando el agente es potente
(Opus, Sonnet) pero es frágil y no reproducible.

### Patrón 3: Methodology abstraction leak

La abstracción de metodología es sólida a nivel de ARTEFACTOS (los 6
artifacts son realmente universales) pero ASPIRA a nivel de
ORQUESTACIÓN. El routing, gates, convocatoria, y ceremonies son
Scrum-shaped. No bloquea el uso inmediato (Scrum default) pero la
claim "intercambiable" es prematura para routing.

### Patrón 4: Happy-path bias

Los docs son exhaustivos para el flujo feliz pero thin en: errores del
MIM, artefactos inconsistentes, fallos parciales en votación, drift
editorial del TPM. Los exception paths necesitan la misma rigurosidad
que el happy path.

---

## Resumen ejecutivo

| Severidad | Total | Corregidos | Diferidos | Pendientes |
|-----------|-------|------------|-----------|------------|
| Ya corregidos (pre-review) | 4 | 4 | — | — |
| CRITICAL | 9 | 9 | — | 0 |
| HIGH | 9 | 8 | 1 (H3) | 0 |
| MEDIUM | 14 | 10 | 3 (M3, M13, M14) | 1 (M6) |
| LOW | 6 | 0 | 0 | 6 |
| **TOTAL** | **42** | **31** | **4** | **7** |

Los 9 CRITICAL y 8/9 HIGH se corrigen en esta sesión.
Los LOW quedan como polish para la próxima iteración.
