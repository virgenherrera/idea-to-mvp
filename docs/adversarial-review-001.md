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

*Pendiente: documentar la mecánica de fast-forward y el gradiente de
certeza en behavior file*

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

*Pendiente: documentar estrategia de delegación multi-modelo*

---

## HIGH — TBD (pendiente: definir capa de metodología primero)

> **Nota**: todos los hallazgos HIGH se reevalúan después de definir la
> capa de abstracción de metodología (Scrum/Kanban/Shape Up/PI Planning).
> La metodología elegida afecta roles, gates y ceremonia — pero NO el
> modelo de artefactos del RAG.

### H1. El TPM no existe en `operational-model.md` — TBD

*Fuente: coherencia (#10)*

### H2. Gate validation es semánticamente vacía — TBD

*Fuente: viabilidad técnica (#4)*

### H3. La escalación ejecución → planificación no tiene protocolo — TBD

*Fuente: viabilidad técnica (#5), producto (#8)*

### H4. Tablas de roles inconsistentes entre archivos — TBD

*Fuente: coherencia (#3, #4, #5, #6, #7, #12)*

### H5. Preguntas pre-definidas son rígidas — TBD

*Fuente: producto (#5)*

### H6. Carga cognitiva insostenible para nuevos usuarios — TBD

*Fuente: producto (#4)*

### H7. Terminología SDD residual en modelo operativo — TBD

*Fuente: coherencia (#1, #9)*

---

## MEDIUM — TBD

### M1. RAG fuera del repo rompe colaboración de equipo — TBD
### M2. La regla "SM nunca toca archivos" es unenforceable — TBD
### M3. El handoff asume planificación perfecta — TBD
### M4. No hay feedback del usuario sobre el proceso — TBD
### M5. Path del RAG: `docs/` vs. `artifacts/` — TBD
### M6. Interfaz del adaptador insuficiente — TBD

---

## Patrones observados

Los hallazgos convergen en tres temas transversales:

1. ~~**Diseño para estado que no existe**~~ → **RESUELTO**. State machine
   derivada del RAG (C1), supervisión post-hoc en vez de heartbeat (C2).
   El diseño ahora asume agentes efímeros y fire-and-forget.

2. **Ceremonia desproporcionada** → **MITIGADO parcialmente**. Fast-forward
   contextual (C3) evita forzar fases cuando los artefactos ya existen.
   Delegación a modelos económicos (C5) reduce costo. Pendiente: definir
   concretamente los tiers de activación.

3. **Dos documentos divergentes** → **TBD**. Se reevalúa después de
   definir la capa de metodología. La metodología determina roles y
   ceremonia; el modelo de artefactos del RAG es independiente.
