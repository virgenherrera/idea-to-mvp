# Review Adversarial #001 — Modelo Operativo

> Fecha: 2026-07-30
> Archivos revisados: `operational-model.md`, `behavior-scrum-master-routing.md`
> Agentes: 3 (viabilidad técnica, producto/UX, coherencia interna)
> Total hallazgos brutos: ~47. Consolidados (deduplicados): 15.

---

## CRITICAL — Bloquean la viabilidad del diseño

### ~~C1. La state machine del SM no sobrevive entre sesiones~~ → RESUELTO

**Resolución**: el estado de la state machine se DERIVA del RAG, no se
persiste internamente. El SM le pregunta al TPM "¿qué artefactos existen?"
y reconstruye su fase actual. Igual que un SM humano abre Jira para saber
en qué va el proyecto.

- `idea.md` completo + nada más → estamos en Fase 2
- Todo hasta `tasks.md` → estamos en Fase 5
- Nueva sesión, compaction, crash → el SM se recupera del RAG

Las decisiones condicionales ("UX desactivado") se persisten como metadatos
dentro de `idea.md` (sección "roles activos para este proyecto").

*Actualizado en: `behavior-scrum-master-routing.md`, sección State Machine*

### ~~C2. El heartbeat es imposible con la API actual~~ → RESUELTO

**Resolución**: reemplazado por **supervisión post-hoc** (patrón validado
en `nest-base`, `virgenherrera`, `fullstack-base`):

1. **Status Report obligatorio** en el output del sub-agente (Status /
   Progress / Blocker / Artifacts). Sin este bloque → FAILED.
2. **PDC de 4 pasos** (ECHO/VERIFY/MARK/DECIDE) después de cada retorno.
3. **Circuit Breaker** (3 fallos → escalar al MIM).
4. **Context Resilience** (artefactos son la memoria, reglas viajan como
   texto, skill resolution feedback detecta pérdida de contexto).

No es heartbeat real-time — es evaluación batch post-retorno. Funciona
con fire-and-forget.

*Actualizado en: `behavior-scrum-master-routing.md`, sección Supervisión*

### ~~C3. No hay escape hatch~~ → RESUELTO (fast-forward contextual)

**Resolución**: no es un "skip" global — es **fast-forward contextual**.
El MIM puede inyectar trabajo en cualquier punto del ciclo según lo que
ya existe:

- **Bug en producción** → el MIM dice "esto tronó" → SM orquesta:
  reproducir → diagnosticar → fix → promover al ambiente apropiado.
  No pasa por Idea → Spec → Design.
- **Epic ya groomeado** → todo está en el RAG (spec, design, tasks,
  handoff) → SM detecta artefactos completos → fast-forward directo
  a ejecución.
- **Idea vaga** → pasa por todas las fases, sin atajos.

La regla: el escape hatch NO salta fases, las ACELERA según los
artefactos que ya existan en el RAG. Si el RAG ya tiene spec + design +
tasks, no tiene sentido forzar al MIM a responder preguntas que ya
están respondidas.

**Gradiente de certeza** — el SM evalúa autónomamente hasta dónde
puede avanzar sin consultar al MIM:

| Input | Certeza | FF del SM | Se detiene en... |
|-------|---------|-----------|------------------|
| "Hazme el uber de lanchas" | Baja | Crea idea + formula preguntas | Espera respuestas del MIM |
| "Agrega auth con JWT" | Media | Idea + spec parcial + preguntas de diseño | Decisiones del MIM |
| "Crea módulo OTEL" | Alta | Hasta handoff o ejecución | Casi sin preguntas |

El SM juzga: "¿la solución es determinista dado el contexto que tengo?"
Si sí → avanza. Si no → formula preguntas precisas y espera.

*~~Pendiente: documentar la mecanica de fast-forward y el gradiente de
certeza en behavior file~~ → RESUELTO en review-002 C2 (checklist de
4 factores con scoring 0/1/2, thresholds numericos, 3 ejemplos de
frontera)*

### ~~C4. El modelo de dos modos no cubre 8 etapas~~ → DIFERIDO

**Decisión**: foco actual es definir completamente el modo planificación.
Las etapas 6-8 (Verificar, Aceptar, Retrospectiva) se reconcilian
después, cuando el modo planificación esté sólido.

*Fuente: coherencia (#2, #8, #17) — se aborda en fase posterior*

### ~~C5. El costo de tokens es potencialmente prohibitivo~~ → RESUELTO (mitigación por diseño + exploración)

**Resolución**: el diseño del SM + RAG ya mitiga esto por construcción:

1. **Contexto acotado**: los agentes reciben slices del RAG, no archivos
   completos. Tareas simples con given/when/then + output esperado +
   criterios de cumplimiento.
2. **Delegación a modelos económicos**: explorar Ollama (modelos locales)
   y AI containers de Docker para tareas mecánicas (implementar,
   investigar, ejecutar tests).
3. **Evidencia empírica**: proyectos complejos entregados con disciplina
   de tokens. No es prohibitivo cuando se hace bien.

El costo no se reduce por magia — se reduce porque el SM da tareas bien
acotadas en vez de dumping de contexto.

*Estrategia de delegacion multi-modelo → RESUELTO. Documentada en
`operational-model.md`: tier Local (Docker/Ollama) para tareas
mecanicas (TPM, checks estructurales), tier Cloud (Claude/Codex) para
tareas de razonamiento (roles, coordinacion del SM, review adversarial).
Criterio de seleccion: determinismo vs juicio (commit 500b79f).*

---

## HIGH — Reevaluados post review-002/003

### ~~H1. El TPM no existe en `operational-model.md`~~ → RESUELTO

*Fuente: coherencia (#10) → resuelto en review-002 H9*

TPM agregado a `operational-model.md` como capa entre roles y
adaptadores, con referencia cruzada a `artifact-model.md`.

### ~~H2. Gate validation es semánticamente vacía~~ → RESUELTO

*Fuente: viabilidad técnica (#4) → resuelto en review-002 H4*

Gates separados en completitud estructural (TPM, mecanico) +
completitud semantica (rol validador, juicio). Ownership Matrix
ahora asigna validadores independientes a los 6 artefactos.

### H3. La escalación ejecución → planificación no tiene protocolo — DIFERIDO (depende de Modo de Ejecución)

*Fuente: viabilidad técnica (#5), producto (#8)*

`escalate(gap)` existe en la API del SM con 3 targets (Idea, Spec,
Design) segun tipo de gap. Protocolo completo depende del diseno
del Execution Mode (pendiente).

*Ver Pendientes Unificados (review-003) → A1.*

### ~~H4. Tablas de roles inconsistentes entre archivos~~ → RESUELTO

*Fuente: coherencia (#3, #4, #5, #6, #7, #12) → resuelto en review-002 C6, H8*

QA eliminado de Fase 3, DevSecOps/QA como condicionales en Fase 4,
tablas reconciliadas entre los 4 documentos.

### H5. Preguntas pre-definidas son rígidas — DIFERIDO (concern de producto)

*Fuente: producto (#5)*

Las 6 preguntas de Fase 1 (linea 323 de behavior.md) siguen siendo
una lista estatica. El fast-forward evalua certeza, pero las preguntas
del PO al MIM no se adaptan al contexto del proyecto.

*Ver Pendientes Unificados (review-003) → A2.*

### H6. Carga cognitiva insostenible para nuevos usuarios — DIFERIDO (concern de producto)

*Fuente: producto (#4)*

El framework tiene 4 documentos interrelacionados con ~4000 lineas.
No hay onboarding guide ni quick-start. Concern de producto, no
auto-fixable.

*Ver Pendientes Unificados (review-003) → A3.*

### ~~H7. Terminología SDD residual en modelo operativo~~ → RESUELTO

*Fuente: coherencia (#1, #9) → resuelto en review-002 C7*

"Explorar/Proponer" reemplazado con nombres canonicos (Idea, Spec,
Design, Tasks, Handoff). Nota: quedan 3 menciones legitimas de "SDD"
en operational-model.md referentes al artifact store, no a fases.

---

## MEDIUM — Reevaluados post review-002/003

### ~~M1. RAG fuera del repo rompe colaboración de equipo~~ → RESUELTO por diseno

La capa de adaptadores de persistencia (artifact-model.md sec
"Adaptadores de Persistencia") resuelve esto por arquitectura:
interfaz universal de 8 operaciones con 6 implementaciones (Local,
Engram, Jira/Asana/Linear, DBMS, Git Repo, MS Project/Basecamp).
Colaboracion multi-dev = elegir adaptador con soporte de concurrencia
(Git Repo, Jira, DBMS). El adaptador local es el MVP, no el techo.

### ~~M2. La regla "SM nunca toca archivos" es unenforceable~~ → RESUELTO por diseno

La delegacion via TPM (escritura) y sub-agentes (lectura) ES el
mecanismo de enforcement. No hay enforcement a nivel de runtime, pero
el protocolo de delegacion lo hace innecesario si se sigue.

### ~~M3. El handoff asume planificación perfecta~~ → MITIGADO

Smoke test (review-003 C5), verifyConsistency cascading (review-002
H5), y mid-planning edit protocol (review-003 H4) cubren los paths
de fallo. No es "planificacion perfecta" — es planificacion validada.

### ~~M4. No hay feedback del usuario sobre el proceso~~ → RESUELTO

Fase 8 (Retrospectiva) ahora tiene contrato completo en behavior.md:
stop/start/continue doing + agreements. El MIM tiene un punto formal
para dar feedback sobre el proceso. Los agreements son meta-config
que afinan el comportamiento del SM en el siguiente ciclo.

### ~~M5. Path del RAG: `docs/` vs. `artifacts/`~~ → RESUELTO

*Resuelto en review-002 M7*. Canonicalizado a `docs/`.

### ~~M6. Interfaz del adaptador insuficiente~~ → RESUELTO

*Resuelto en review-002 M8*. operational-model.md referencia la
interfaz completa de 7 ops definida en artifact-model.md.

---

## Patrones observados

Los hallazgos convergen en tres temas transversales:

1. ~~**Diseno para estado que no existe**~~ → **RESUELTO**. State machine
   derivada del RAG (C1), supervision post-hoc en vez de heartbeat (C2).
   El diseno ahora asume agentes efimeros y fire-and-forget.

2. **Ceremonia desproporcionada** → **RESUELTO**. Fast-forward contextual
   (C3) con scoring numerico (review-002 C2), tiers de activacion
   (condensed contracts), Pattern B reduce dispatches, y estrategia de
   delegacion multi-modelo documentada (C5, commit 500b79f).

3. ~~**Dos documentos divergentes**~~ → **RESUELTO**. operational-model.md
   reconciliado en review-002 (C7, C9, H9, M7, M8) y review-003 (H10,
   M1). Regla de propagacion adoptada (review-003 Patron 1) previene
   recurrencia.
