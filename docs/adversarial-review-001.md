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

### C3. No hay escape hatch

El SM bloquea avances sin alternativa. Un desarrollador senior que sabe
exactamente qué quiere construir debe responder 6 preguntas de negocio +
5 fases de planificación antes de escribir una línea de código. No existe
`/skip-to-execution` ni `/fast-track` con acknowledgment de riesgos.

*Fuente: producto (#1)*

### C4. El modelo de dos modos no cubre 8 etapas

Las etapas 6-8 (Verificar, Aceptar, Retrospectiva) son huérfanas: no
caben en Planificación (que termina en Handoff) ni en Ejecución (que según
el modelo solo involucra orquestador + minions). Pero el behavior file
muestra al scrum team completo activo en esas etapas — contradiciendo
"el scrum team está EN SILENCIO en ejecución."

Se necesita: un tercer modo, o redefinir dónde terminan los modos, o
asignar explícitamente estas etapas.

*Fuente: coherencia (#2, #8, #17)*

### C5. El costo de tokens es potencialmente prohibitivo

Cada fase invoca SM + 1-3 roles + TPM. Ejemplo conservador: 8 fases × ~5
agentes = ~40 invocaciones. A ~10K tokens cada una = 400K tokens solo en
overhead de orquestación, antes de escribir código. Para un "CRUD con auth"
de 20 líneas, el framework podría costar 10-20x más que un agente único.

*Fuente: viabilidad técnica (#3), producto (#3, #6)*

---

## HIGH — Requieren especificación significativa

### H1. El TPM no existe en `operational-model.md`

El behavior file hace del TPM un componente indispensable (el SM delega
TODO el acceso al RAG al TPM). El modelo operativo no lo menciona. Un
implementador que siga solo el modelo operativo no sabe quién gestiona
el RAG.

*Fuente: coherencia (#10)*

### H2. Gate validation es semánticamente vacía

El SM valida gates vía TPM, pero el TPM reporta completitud estructural
("6/6 secciones"), no calidad semántica. Nadie tiene AMBOS: autoridad para
bloquear Y visibilidad del contenido. El SM tiene autoridad pero no lee
contenido. El TPM lee contenido pero no tiene autoridad de proceso.

*Fuente: viabilidad técnica (#4)*

### H3. La escalación ejecución → planificación no tiene protocolo concreto

La state machine muestra flechas `Execution → Spec`, `Execution → Design`,
pero: ¿cómo comunica el orquestador de ejecución un gap al SM? Pueden ser
sesiones diferentes, contextos diferentes. No hay formato de gap report,
no hay trigger definido, no hay mecanismo de recovery.

*Fuente: viabilidad técnica (#5), producto (#8)*

### H4. Tablas de roles inconsistentes entre archivos

Las fases de activación de QA, UX, Dev Lead, DevSecOps y SM difieren entre
`operational-model.md` y `behavior-scrum-master-routing.md`. Además, el
mapa de convocatoria (diagrama Mermaid) difiere de la matriz detallada
dentro del mismo archivo. Tres fuentes, tres respuestas diferentes a "¿a
quién convoco en fase X?"

*Fuente: coherencia (#3, #4, #5, #6, #7, #12)*

### H5. Preguntas pre-definidas son rígidas

6 preguntas fijas en Fase 1, 5 en Fase 2, etc. No hay mecanismo para
detectar que el input ya respondió algunas. El 60% de las preguntas
podrían estar respondidas en el input original y se preguntan de todas
formas.

*Fuente: producto (#5)*

### H6. Carga cognitiva insostenible para nuevos usuarios

7 roles + RAG + gates + heartbeats + contratos + state machine + handoffs +
tiers + 2 modos = 15+ conceptos interconectados antes de poder usar la
herramienta. Los roles deberían ser transparentes (implementation detail),
no parte de la interfaz visible.

*Fuente: producto (#4)*

### H7. Terminología SDD residual en modelo operativo

El diagrama de límites usa "explorar, proponer" que no son fases del modelo.
La tabla del PO referencia "Propuesta" que no existe. Son remanentes de
SDD que contradicen la propia separación gobernanza/operacional.

*Fuente: coherencia (#1, #9)*

---

## MEDIUM — Fricción de diseño, no bloquean

### M1. RAG fuera del repo rompe colaboración de equipo

`~/.idea-to-mvp/` es local al usuario. Sin convención de compartir, versionar
o resolver conflictos en artefactos de planificación entre personas.

### M2. La regla "SM nunca toca archivos" es unenforceable

El SM tiene acceso a Read/Write/Bash. La única restricción es el system
prompt. Después de compaction puede olvidarla.

### M3. El handoff asume planificación perfecta

"Autocontenido" asume 100% de decisiones capturadas. En la práctica, la
ejecución siempre descubre info nueva. Falta un "decision budget" para
que ejecución tome decisiones tácticas sin escalar.

### M4. No hay feedback del usuario sobre el proceso

La Retrospectiva evalúa el producto, no la experiencia con el framework.
Sin mecanismo para que el usuario diga "DevSecOps fue inútil para este
proyecto" y que eso persista.

### M5. Path del RAG: `docs/` vs. `artifacts/`

Dos nombres diferentes para el mismo concepto en el modelo operativo.

### M6. Interfaz del adaptador insuficiente

4 operaciones (guardar, leer, buscar, listar) vs. 6 del TPM (incluye
marcar completo y verificar consistencia). Faltan metadatos.

---

## Patrones observados

Los hallazgos convergen en tres temas transversales:

1. **Diseño para estado que no existe**: la state machine, el heartbeat, y
   la escalación asumen un proceso persistente y bidireccional. El SM es
   efímero y los sub-agentes son fire-and-forget. Hay que diseñar para la
   realidad, no para la API deseada.

2. **Ceremonia desproporcionada**: 8 fases, 7 roles, gates obligatorios,
   contratos de 6 campos. Funciona para un proyecto complejo con equipo.
   Destruye la adopción para un solo desarrollador con un fin de semana.
   Los tiers de activación se mencionan pero nunca se definen
   concretamente.

3. **Dos documentos divergentes**: el modelo operativo y el behavior file
   evolucionaron en paralelo sin reconciliación. Tablas de roles, fases,
   modos, y terminología difieren. Se necesita una fuente de verdad única
   o una reconciliación explícita.
