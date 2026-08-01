# Fast-Forward y Tiers de Activación

← [Índice principal](../../README.md) | [Planificación](../README.md) | [Comportamiento SM](README.md)

## Fast-Forward Contextual — Gradiente de Certeza

El SM no avanza siempre una fase a la vez. Al recibir un input, evalúa
**qué tan determinista es la solución dado el contexto existente** y
avanza proporcionalmente:

```mermaid
flowchart TD
    INPUT["Input del MIM"] --> EVAL["SM evalúa gradiente\nde certeza"]
    EVAL -->|"Baja\n(dominio desconocido)"| LOW["Crea idea.md\n+ formula preguntas"]
    EVAL -->|"Media\n(estándar con decisiones)"| MED["Idea + spec parcial\n+ preguntas de diseño"]
    EVAL -->|"Alta\n(estándar determinista)"| HIGH["Hasta handoff\no ejecución"]
    LOW --> WAIT["⏳ Espera respuestas\ndel MIM"]
    MED --> WAIT
    HIGH --> EXEC["▶ Avanza sin preguntar"]
```

### Reglas del gradiente

| Certeza | Criterio del SM | Hasta dónde avanza | Ejemplo |
|---------|-----------------|--------------------|---------|
| **Baja** | Dominio desconocido, requisitos ambiguos, no hay app existente | Idea + preguntas | "Hazme el uber de lanchas" |
| **Media** | Estándar conocido pero con decisiones pendientes | Idea + spec parcial + preguntas específicas | "Agrega auth con JWT" |
| **Alta** | Estándar abierto, app existente en el RAG, patrones bien definidos | Hasta handoff o ejecución directa | "Crea módulo OTEL" |

### Quién decide

**El SM decide autónomamente** usando un checklist de 4 factores.
No es el MIM quien dice "ve en fast-forward" — el SM evalúa y decide.

### Checklist de certeza (obligatorio, auditable)

El SM evalúa 4 factores y asigna 0, 1, o 2 puntos a cada uno:

| Factor | 0 puntos | 1 punto | 2 puntos |
|--------|----------|---------|----------|
| **F1. Artefactos existentes** | RAG vacío | 1-2 artefactos upstream | spec + design + tasks aprobados |
| **F2. Estandarización** | Dominio custom sin estándar | Estándar con variantes (auth, API) | Estándar abierto puro (OTEL, i18n, linting) |
| **F3. Ambigüedad de dominio** | Infinitas interpretaciones ("uber de X") | Dominio acotado con decisiones pendientes | Dominio determinista (agregar módulo X a app existente) |
| **F4. Referencia existente** | Sin codebase ni precedentes | Codebase existe pero no cubre este dominio | Codebase con patrones/stack que aplican directamente |

> **Nota F1**: un artefacto que existe pero no está aprobado
> cuenta como 0.5 puntos. "No aprobado" = el TPM reporta que faltan
> secciones requeridas o que el artefacto está en borrador/en revisión.
> **Cap**: la suma de puntos de artefactos no aprobados tiene un techo
> de **1 punto** para F1, independientemente de cuántos existan. Esto
> previene que N borradores incompletos alcancen el mismo score (F1=2)
> que artefactos validados y aprobados. Para alcanzar F1=2, los
> artefactos upstream deben estar aprobados.

**Thresholds**:

| Score total | Certeza | Hasta dónde avanza |
|-------------|---------|-------------------|
| 0–2 | **Baja** | Idea + preguntas al MIM |
| 3–5 | **Media** | Idea + spec parcial + preguntas específicas |
| 6–8 | **Alta** | Hasta handoff o ejecución directa |

**El SM DEBE registrar el score en su reasoning** (no solo la
conclusión) para que la decisión sea auditable:

> *"F1=0 (RAG vacío), F2=1 (JWT es estándar con variantes), F3=1
> (auth es acotado pero hay decisiones), F4=2 (codebase existente
> con Express). Total: 4 → Media. Avanzo a idea + spec parcial."*

El SM instruye al TPM para persistir el score F1-F4 y el reasoning en
`idea.md` sección "Decisiones tomadas" como entrada con formato:
`[FAST-FORWARD] F1={n}, F2={n}, F3={n}, F4={n}. Total={n} → {certeza}.
Razón: {resumen}.` Esto garantiza auditabilidad cross-session.

### Ejemplos resueltos de frontera

| Input | F1 | F2 | F3 | F4 | Total | Certeza | Acción |
|-------|----|----|----|----|-------|---------|--------|
| "Hazme el uber de lanchas" | 0 | 0 | 0 | 0 | 0 | Baja | Idea + preguntas |
| "Agrega auth con JWT" (sin codebase) | 0 | 1 | 1 | 0 | 2 | Baja | Idea + preguntas |
| "Agrega auth con JWT" (codebase Express existente) | 0 | 1 | 1 | 2 | 4 | Media | Idea + spec parcial |
| "Crea módulo OTEL" (codebase con NestJS) | 0 | 2 | 2 | 2 | 6 | Alta | Hasta handoff |
| "Epic X ya groomeado" (spec+design+tasks en RAG) | 2 | 2 | 2 | 2 | 8 | Alta | Fast-forward a ejecución |
| "Implementa pagos con Stripe" (sin codebase) | 0 | 1 | 1 | 0 | 2 | Baja | F2=1: Stripe es estándar PERO tiene variantes (checkout, elements, custom). F3=1: pagos es acotado pero requiere decisiones (moneda, suscripciones, webhooks). |
| "Agrega logging con Winston" (codebase Node existente) | 0 | 2 | 2 | 2 | 6 | Alta | F2=2: Winston es estándar abierto sin variantes significativas. F3=2: logging es determinista — configuración, transports, formato. |
| "Migra de REST a GraphQL" (API existente) | 1 | 1 | 0 | 2 | 4 | Media | F2=1: GraphQL es estándar PERO cada migración es diferente. F3=0: infinitas interpretaciones — qué endpoints migrar, schema design, N+1. |

Vista alternativa: el quadrant chart ubica cada ejemplo según cuánto
contexto existe (eje X) y cuán determinista es el dominio (eje Y) — los
casos en el cuadrante superior derecho son los candidatos naturales a
fast-forward completo.

```mermaid
quadrantChart
    title Gradiente de Certeza
    x-axis Bajo Contexto --> Alto Contexto
    y-axis Dominio Ambiguo --> Dominio Determinista
    quadrant-1 Fast-forward completo
    quadrant-2 Spec + diseño parcial
    quadrant-3 Idea + preguntas al MIM
    quadrant-4 Spec + diseño parcial
    Uber de lanchas: [0.1, 0.1]
    Auth JWT sin codebase: [0.15, 0.85]
    Auth JWT con Express: [0.65, 0.85]
    Módulo OTEL NestJS: [0.8, 0.9]
    Epic groomeado: [0.95, 0.95]
```

### Fast-forward también aplica MID-CYCLE

No solo al inicio. Ejemplos:

- **Bug en producción** → MIM dice "esto tronó" → SM orquesta:
  reproduce → diagnostica → fix → promueve al ambiente apropiado.
  No pasa por Idea → Spec → Design.
- **Epic ya groomeado** → todo en el RAG → SM detecta artefactos
  aprobados → fast-forward directo a ejecución.

---

## Tiers de Activación

El SM determina el **tier de ceremonia** al inicio de cada ciclo usando el
score de fast-forward (F1-F4). El tier define cuánta ceremonia se aplica,
no qué artefactos se producen — los artefactos son universales.

```mermaid
flowchart TD
    SCORE["Score F1-F4\n(0-8 puntos)"] --> CHECK{{"Evaluar\nrango"}}
    CHECK -->|"0-2"| COMPLETO["Tier Completo\nCeremonia total"]
    CHECK -->|"3-5"| ESTANDAR["Tier Estándar\nCeremonia normal"]
    CHECK -->|"6-8"| LIGERO["Tier Ligero\nCeremonia mínima"]

    COMPLETO --> C_OUT["Todos los roles\nTodos los gates\nDispatch normal"]
    ESTANDAR --> E_OUT["3-4 roles por fase\nGates estándar\nFast-forward parcial"]
    LIGERO --> L_OUT["1-2 roles esenciales\nGates comprimidos\nDispatch ultra-comprimido"]
```

### Tabla de tiers

| Tier | Score | Ceremonia | Roles | Dispatch | Ideal para |
|------|-------|-----------|-------|----------|------------|
| **Ligero** | 6-8 | Mínima. SM puede comprimir múltiples fases en una sola delegación. | 1-2 roles (los estrictamente necesarios para la fase) | Comprimido o ultra-comprimido | Bugs, epics ya groomeados, estándar abierto puro |
| **Estándar** | 3-5 | Normal. Fases secuenciales con fast-forward parcial posible. | 3-4 roles según fase | Normal | Features nuevos, dominio acotado con decisiones pendientes |
| **Completo** | 0-2 | Total. Toda fase ejecutada, todo rol convocado, todo gate enforced. | Todos los roles default + posibles ad-hoc | Normal (sin compresión) | Productos nuevos, alta ambigüedad, regulados, misión crítica |

### Qué cambia por tier

| Aspecto | Ligero (6-8) | Estándar (3-5) | Completo (0-2) |
|---------|--------------|----------------|----------------|
| **Roles por fase** | 1-2 esenciales | 3-4 según fase | Todos + ad-hoc |
| **Gates** | Comprimidos (SM valida inline) | Estándar (PDC completo) | Estrictos (PDC + validación cruzada) |
| **Dispatch** | Ultra-comprimido: múltiples fases en una delegación | Normal: una fase por delegación | Normal: una fase por delegación, sin omisiones |
| **Smoke test handoff** | Omisible si el contexto es determinista | Requerido | Requerido + revisión adversarial |

### Reglas de escalación

- El SM determina el tier al INICIO del ciclo, basado en el score F1-F4.
- El tier puede **escalar** mid-cycle (Ligero → Estándar, Estándar →
  Completo) si la complejidad descubierta lo justifica.
- El tier **NUNCA** de-escala mid-cycle. La complejidad descubierta no se
  puede des-descubrir.
- **Triggers de escalación**:
  1. Tasa de fallo PDC > 50% en el tier actual (más de la mitad de las
     delegaciones retornan FAILED o PARTIAL sin progreso).
  2. El MIM solicita explícitamente más ceremonia.

#### Ejemplo concreto de escalación

> El SM inicia un ciclo en **Tier Ligero** (score 7: módulo OTEL en app
> NestJS existente). Durante la fase de diseño, el Dev Lead descubre que
> la integración requiere un custom exporter con lógica de retry no
> trivial. Dos delegaciones consecutivas retornan PARTIAL. El SM evalúa:
> 2/3 delegaciones con problemas → tasa > 50%. Escala a **Tier Estándar**:
> convoca QA para validar testeabilidad y DevSecOps para revisar el
> surface del exporter. El ciclo continúa con ceremonia normal desde este
> punto.

### Nota sobre artefactos

Los tiers afectan la **ceremonia** (roles convocados, gates aplicados,
patrón de dispatch), NO los **artefactos**. Independientemente del tier,
el ciclo produce los mismos artefactos (`idea.md`, `spec.md`,
`design.md`, `tasks.md`, `handoff.md`). Lo que cambia es cuántos ojos
los revisan y cuántos checkpoints se aplican antes de aprobarlos.
