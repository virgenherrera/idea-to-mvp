# Validation: Virgil v2 vs swarm-forge

Análisis comparativo adversarial entre la idea de Virgil v2
(`idea-virgil-v2.md`) y SwarmForge de Robert C. Martin (branch `main`,
documentaria + scripts operativos). El objetivo no es copiar swarm-forge:
es detectar gaps reales en Virgil y confirmar sus ventajas legítimas.

## 1. Architectural Comparison

| Dimensión | swarm-forge | Virgil v2 |
|-----------|-------------|-----------|
| Identidad / propósito | Plataforma de orquestación de swarms de agentes sobre tmux. "Disciplined agents build better software". Orquesta el loop de delivery. | Herramienta de metodología + ownership. Puente entre la metodología idea-to-mvp y cualquier agente. No es un agente ni un orquestador. |
| Modelo de agentes | Multi-agente explícito: 2, 4 o 6 roles fijos por branch (specifier, coder, cleaner, architect, hardener, QA), cada uno en su propio git worktree y sesión tmux. Backends: claude, codex, copilot, grok. | Agent-agnostic: un solo agente (o un swarm externo) asume el rol de SM y consume skills, hooks y el binding layer. Virgil no lanza ni posee procesos de agentes. |
| Mecanismo de handoff | Archivos `.handoff` durables entregados por un daemon (Babashka). Dos tipos: `git_handoff` (puntero a commit de 10 hex + task name) y `note` (una línea, máx 80 chars). Colas `outbox/sent/failed/inbox` con estados por ubicación de archivo y timestamps de auditoría. | Un solo `handoff.md` autocontenido: US completa, justificación de negocio, ACs verificables, decisiones arquitectónicas, tasks con dependencias resueltas. Habilita lanes paralelos en lugar de cadena serial. |
| Enforcement de calidad | Baterías de métricas operativas: TDD, coverage, CRAP (crap4go/clj/java), DRY analysis, mutation testing de lenguaje, Gherkin acceptance mutation, property tests. Roles dedicados a revisión (cleaner, architect, hardener, QA). | Tres niveles de enforcement: advisory (AGENTS.md), contextual (skills/MCP), determinístico (hooks). Verificación vía `virgil verify` / `virgil coverage` sobre el binding layer. Las métricas de calidad se mencionan pero no están diseñadas. |
| Profundidad metodológica | Empieza en Gherkin: el specifier convierte intención del usuario en specs de aceptación. No hay capa de negocio, ni idea, ni design doc, ni justificación del "por qué". | Cadena e2e: idea → spec → design → tasks → handoff → implementación → operación (opcional). Cada fase produce artefactos verificables y groomeados antes de escribir código. |
| Trazabilidad | Implícita: task name que viaja en los handoffs + Gherkin acceptance tests. No hay grafo requirement → code → test, ni detección de staleness. | Explícita y central: grafo REQUIREMENT/CODE_ARTIFACT/TEST_ARTIFACT/BINDING con confidence (declared/inferred/verified) y lifecycle (ACTIVE/STALE/BROKEN/ARCHIVED). |
| Tooling | zsh + tmux + Babashka + git worktrees. Distribución por `curl | tar` de branches de GitHub. Dependencias de runtime: bb, tmux, terminal adapters, AppleScript/wt.exe. macOS/Linux-céntrico. | Go 100% pure (CGO_ENABLED=0), GoTreeSitter, modernc sqlite, go:embed, GoReleaser → Homebrew. Binario estático multiplataforma, sin runtime externo. |
| Operación / deployment | Ninguno. El flujo termina en QA/completion notification. Soporte operativo limitado al host del swarm (sleep inhibitors, window watchdog, cleanup). | Fase de operación opcional como facade/plugin. Además: brownfield takeover, mantenimiento incremental del grafo vía git hooks, ownership post-delivery. |

## 2. What swarm-forge does that Virgil doesn't (GAPS)

### 2.1 Motor de métricas de calidad operativo (GAP CRÍTICO)

- **Qué es**: swarm-forge operacionaliza exactamente la cita de Uncle Bob que
  Virgil declara como dogma: mutation testing (`mutate4go`), CRAP score
  (`crap4go`), DRY analysis (`dry4go`), Gherkin acceptance mutation
  (`gherkin-mutator`), property tests y coverage, con reglas de cuándo corren
  y qué rol las posee (`engineering.prompt`).
- **¿Relevante para Virgil?**: máximamente. El dogma de Virgil dice
  "measure test coverage, dependency structure, cyclomatic complexity,
  module sizes, mutation testing". El documento de idea diseña en detalle la
  trazabilidad, pero las métricas aparecen en una sola línea sin diseño. Sin
  mutation testing, el binding layer puede reportar "AC-3 → test passing" con
  tests vacuos: trazabilidad sin verificación de fuerza de los tests es una
  verificación débil. Virgil hoy cubre la mitad del dogma que reclama.
- **Recomendación**: **incorporar adaptado**. Virgil no debe construir sus
  propias herramientas de mutación/CRAP (eso es el ecosistema de Uncle Bob);
  debe definir un contrato de métricas: `virgil metrics` que orqueste
  herramientas existentes por lenguaje y persista resultados junto al
  binding, con thresholds configurables por tier.

### 2.2 Semántica de coordinación concurrente

- **Qué es**: swarm-forge define con precisión quirúrgica cómo múltiples
  agentes trabajan sin pisarse: un worktree por rol, colas con moves atómicos,
  lock de secuencia, estados por ubicación de archivo, regla "un solo item
  in_process", batch mode, y recuperación tras restart.
- **¿Relevante para Virgil?**: sí. Virgil promete "un handoff, ejecución
  paralela por lanes independientes", pero el documento no define: cómo un
  lane reclama una task, qué pasa si dos lanes tocan el mismo módulo, cómo se
  serializan las escrituras concurrentes al grafo SQLite (lanes + hooks
  post-commit simultáneos), ni la disciplina de merge entre lanes. La promesa
  de paralelismo es hoy una afirmación, no un diseño.
- **Recomendación**: **adaptar**. No se necesita un daemon ni tmux, pero sí
  un contrato mínimo de claiming (estado por task dentro del handoff o del
  grafo: `pending | claimed | done`) y una política de concurrencia de
  escritura para el binding store.

### 2.3 Semántica de fallo y recuperación

- **Qué es**: en swarm-forge todo estado sobrevive a un crash: inbox durable,
  `in_process/` que obliga a retomar antes de aceptar trabajo nuevo, `failed/`
  con diagnósticos, timestamps de lifecycle (`dequeued_at`, `completed_at`),
  y regla explícita de "on restart, run `ready_for_next.sh`".
- **¿Relevante para Virgil?**: sí. ¿Qué pasa cuando un agente muere a mitad
  de la ejecución del handoff? ¿Cómo distingue el siguiente agente (o el
  mismo tras compaction) qué tasks están completas, a medias, o abandonadas?
  El binding layer registra lo implementado, pero no hay estado de ejecución
  del handoff en sí.
- **Recomendación**: **incorporar adaptado**: el handoff necesita execution
  state versionado (checklist con timestamps) para que la reanudación sea
  determinística y auditable.

### 2.4 Validación estricta de artefactos como gate de máquina

- **Qué es**: `swarm_handoff.sh` es un validation boundary real: rechaza
  headers reservados, valida prioridad, resuelve el commit contra git,
  canonicaliza, y emite guía de reparación línea por línea. Un handoff
  malformado no entra al sistema, punto.
- **¿Relevante para Virgil?**: sí. `handoff.md` es "un contrato de ejecución
  autocontenido", pero nada en la idea valida que ese contrato esté completo
  antes de habilitar código: ¿todos los ACs tienen ID? ¿todas las tasks tienen
  dependencias resueltas? ¿cada AC referencia spec? Hoy el gate "ZERO código
  hasta handoff aprobado" es advisory (aprobación humana) sin verificación
  mecánica de completitud.
- **Recomendación**: **incorporar**: `virgil handoff lint` (o hook pre-fase)
  que valide el schema del handoff con errores reparables, al estilo del
  reporte de `swarm_handoff.sh`. Encaja naturalmente en la capa
  "determinística" que Virgil ya postula.

### 2.5 Revisión independiente por contexto separado

- **Qué es**: cleaner/architect/hardener/QA son contextos frescos que revisan
  el trabajo del coder sin el sesgo del contexto que lo produjo.
- **¿Relevante para Virgil?**: parcialmente. Virgil rechaza el pipeline
  serial con razones válidas (mini-waterfall, churn de forwarding), pero la
  función de revisión adversarial por contexto independiente no queda
  reemplazada por nada explícito: `virgil verify` verifica bindings, no
  conformidad semántica con el design.
- **Recomendación**: **exclusión deliberada aceptable**, condicionada a
  cerrar el gap 2.1. Si las métricas (mutation, CRAP, complexity) son el
  reemplazo del review humano/multi-agente, deben existir de verdad. Sin
  ellas, Virgil elimina el review sin sustituto.

### 2.6 Neutralidad de agente operacionalizada

- **Qué es**: swarm-forge no solo declara multi-backend: `swarmforge.bb`
  construye el launch command concreto para claude, codex, copilot y grok,
  incluyendo el mapeo de permission modes por backend.
- **¿Relevante para Virgil?**: sí. Virgil declara independencia de
  ecosistema, pero sus mecanismos de delivery (skills, slash commands,
  AGENTS.md, hooks, MCP) son de facto Claude-céntricos. ¿Cómo llega la
  metodología a Cursor o Copilot, que no tienen skills ni hooks equivalentes?
- **Recomendación**: **incorporar adaptado**: matriz de adapters de delivery
  por ecosistema (qué mecanismo nativo mapea a advisory/contextual/
  determinístico en cada agente), aunque la v1 solo implemente el adapter de
  Claude.

### 2.7 Infraestructura de host (excluir)

- **Qué es**: sleep inhibitors, window watchdog, terminal adapters, cleanup
  de tmux.
- **¿Relevante?**: no. Virgil no posee procesos de agentes. **Exclusión
  deliberada correcta** — es la consecuencia natural de ser tool y no
  orquestador.

### 2.8 Tests del propio tooling (menor)

- swarm-forge trae `test/swarmforge/*.clj` y un task `bb test` para sus
  helpers. El documento de idea de Virgil no menciona cómo se testea Virgil
  mismo. Es una preocupación de implementación, no de arquitectura, pero un
  tool cuyo pitch es "verificación" debe nacer con su propia suite ejemplar.

## 3. What Virgil does that swarm-forge doesn't (ADVANTAGES)

### 3.1 Cadena de justificación de negocio

- **Qué es**: idea → spec → design → tasks → handoff, con el "por qué" de
  negocio y decisiones arquitectónicas documentadas antes de una línea de
  código. swarm-forge empieza en Gherkin: captura comportamiento, no
  intención de negocio ni tradeoffs de diseño.
- **Por qué importa**: cuando el requerimiento cambia, swarm-forge solo puede
  re-especificar comportamiento; Virgil puede razonar sobre el impacto desde
  la intención original.
- **¿Diferenciador sostenible?**: sí. Es ortogonal al loop de delivery y no
  hay señal de que swarm-forge quiera esa capa.

### 3.2 Brownfield takeover

- **Qué es**: `/virgil-takeover` deriva docs funcionales y baseline de
  trazabilidad desde un codebase existente. swarm-forge no tiene historia
  brownfield alguna: asume que el swarm construye desde cero o que el
  contexto ya está en git.
- **Por qué importa**: la mayoría del software real es brownfield. Un tool
  que solo funciona en greenfield tiene un mercado marginal.
- **¿Sostenible?**: sí, y es además la parte técnicamente más difícil de
  replicar (crawl + inferencia + grafo).

### 3.3 Binding layer con confidence y lifecycle

- **Qué es**: grafo explícito requirement ↔ code ↔ test con niveles de
  confianza y detección de staleness. En swarm-forge la trazabilidad es un
  task name de 80 caracteres viajando en headers de handoff.
- **Por qué importa**: es la única forma de responder "¿AC-3 sigue
  implementado?" seis meses después, sin re-leer el código. Es lo que hace
  posible "manage from a higher level" fuera del momento del delivery.
- **¿Sostenible?**: sí — es el core técnico del producto. Riesgo: la
  categoría existe (Sourcegraph, code intelligence); la ventaja de Virgil es
  atar el grafo a artefactos metodológicos, no solo a símbolos.

### 3.4 Ownership post-delivery

- **Qué es**: hooks post-merge/post-commit/post-rewrite/post-checkout
  mantienen el grafo incremental; staleness marca desincronización doc-código.
  swarm-forge termina en la completion notification de QA — no tiene noción
  de vida del producto después del loop.
- **Por qué importa**: el costo real del software está en mantenimiento, no
  en la primera entrega.
- **¿Sostenible?**: sí.

### 3.5 Distribución y footprint

- **Qué es**: binario Go estático vía Homebrew vs zsh + tmux + Babashka +
  git worktrees + terminal adapters + `curl | tar` de un branch de GitHub.
- **Por qué importa**: fricción de adopción. swarm-forge exige una topología
  local frágil y macOS/Linux; Virgil se instala como cualquier dev tool.
- **¿Sostenible?**: ventaja real pero copiable; no es un moat, es higiene.

### 3.6 Enforcement por capas

- **Qué es**: la distinción advisory/contextual/determinístico. En
  swarm-forge, la constitución completa es prompt text (advisory: el agente
  "debe" releerla en cada handoff); lo único determinístico son los helper
  scripts del transporte.
- **Por qué importa**: los guardrails de metodología (no codear sin handoff,
  commit por refactor) como hooks son inescapables, no depende de que el
  agente obedezca un prompt.
- **¿Sostenible?**: sí, mientras Virgil realmente implemente los hooks (hoy
  es diseño, no código).

### 3.7 Paralelismo sin churn de forwarding

- **Qué es**: swarm-forge obliga a cada rol intermedio a reenviar
  `git_handoff` aunque no haya cambiado nada funcional ("formatting-only...
  churn still require a forward down the chain"). Cada task paga la latencia
  y el costo de toda la cadena, siempre. Virgil elimina la cadena.
- **Por qué importa**: costo por task y throughput.
- **Advertencia**: esta ventaja solo es real si Virgil resuelve el gap 2.2;
  paralelismo sin semántica de coordinación es una promesa, no una ventaja.

## 4. Handoff Protocol Deep Comparison

Los dos "handoffs" resuelven problemas distintos, y esa es la clave del
análisis.

### El handoff de swarm-forge es transporte, no conocimiento

El payload entregado es literalmente:

```text
Re-read your role and constitution.

merge_and_process coder a1b2c3d9
```

Todo el contexto vive fuera del mensaje: en el commit de git, en el role
prompt y en la constitución. El protocolo es excelente en lo operativo:

- Validación estricta en el borde (`swarm_handoff.sh`): headers permitidos,
  prioridad 00-99, commit de 10 hex resuelto y canonicalizado contra git.
- Durabilidad y auditoría: estados por ubicación de archivo, timestamps de
  lifecycle completos, secuencia serializada con lock atómico.
- Recuperación determinística tras restart y wake-ups intencionalmente lossy.

Pero es deliberadamente pobre en intención: el receptor reconstruye el "qué
y por qué" desde el diff. Y el canal lateral (`note`) está capado a 80
caracteres y casi prohibido: ante ambigüedad el agente debe detenerse y
preguntar al humano. Es un protocolo de mensajería entre workers de una
cadena serial, con el humano como resolutor de conflictos.

### El handoff de Virgil es contrato, no mensaje

`handoff.md` transfiere intención: US definida, justificación de negocio,
ACs verificables, decisiones arquitectónicas, tasks con dependencias
resueltas, producto de cinco fases de grooming. Front-loads el contexto para
que N lanes ejecuten sin cascadas inversas.

Sus debilidades son exactamente las fortalezas del protocolo de swarm-forge:

- No hay validación mecánica del contrato (gap 2.4).
- No hay estado de ejecución ni audit trail dentro del handoff (gap 2.3).
- No hay semántica de claiming para los lanes paralelos (gap 2.2).
- No está definido qué pasa si el handoff cambia a mitad de ejecución
  (¿versionado? ¿invalidación de lanes en curso?).

### Veredicto de la comparación

Para la misión de Virgil, su enfoque es el correcto y el más fuerte: un
puntero a commit no puede sostener ejecución paralela ni ownership, porque
no transfiere intención. El modelo de swarm-forge acopla el conocimiento al
historial de git y a la relectura de prompts, lo cual solo funciona en una
cadena serial de roles fijos. Pero swarm-forge gana sin discusión en rigor
operativo del protocolo: validación en el borde, atomicidad, audit trail y
recuperación. La conclusión no es copiar el transporte; es dotar al contrato
de Virgil del mismo rigor de máquina que swarm-forge le da a su mensaje de
una línea.

## 5. Validation Verdict

La arquitectura de Virgil v2 es **sólida en su núcleo y incompleta en sus
bordes**. Las tres capas + binding layer, los dos flujos (greenfield y
takeover), la estrategia de crawl exhaustivo-una-vez/incremental-siempre, y
el delivery por mecanismos nativos con enforcement en capas son un diseño
coherente que swarm-forge no contradice en nada: ambos sistemas son
complementarios, no competidores (de hecho, un swarm tipo swarm-forge podría
ser un consumidor de la metodología de Virgil).

Gaps críticos que deben resolverse antes de implementar:

1. **Métricas de calidad sin diseño** (2.1). Es el gap más serio porque
   contradice el dogma fundacional declarado: Virgil promete la visión de
   Uncle Bob y hoy solo diseña la mitad de trazabilidad. Bloqueante
   conceptual.
2. **Semántica de coordinación paralela** (2.2). La ventaja competitiva
   número uno frente a pipelines seriales está afirmada, no diseñada.
   Bloqueante para la promesa de paralelismo.
3. **Validación y estado de ejecución del handoff** (2.3, 2.4). No
   bloqueante conceptual, pero debe entrar al diseño antes de implementar el
   comando de handoff.

**Confianza: 7/10.** El fundamento es correcto y las ventajas (negocio,
brownfield, binding, post-delivery) son reales y sostenibles. Se descuentan
tres puntos: uno por la mitad no diseñada del dogma (métricas), uno por el
paralelismo sin semántica de coordinación, y uno por la ausencia de rigor de
máquina en el contrato de handoff que la contraparte de referencia sí tiene
resuelto.

## 6. Recommendations

1. **Diseñar el contrato de métricas** (prioridad 1). Definir `virgil
   metrics` como orquestador de herramientas existentes por lenguaje
   (mutation, complexity, CRAP-like, dependency structure), con resultados
   persistidos junto al binding y thresholds por tier. Sin esto, `virgil
   verify` certifica presencia de tests, no fuerza de tests.
2. **Especificar la semántica de lanes** (prioridad 1). Estado de claiming
   por task (`pending | claimed | done`, con owner y timestamp), política de
   escritura concurrente al store SQLite (WAL + serialización de writes), y
   disciplina de merge entre lanes. Un apéndice de una página en el design
   doc basta para desbloquear la implementación.
3. **Agregar `virgil handoff lint`**. Schema del handoff (ACs con ID, tasks
   con dependencias resueltas, referencias a spec/design válidas) validado
   mecánicamente, con errores reparables al estilo del reporte de
   `swarm_handoff.sh`. Es la materialización natural del gate "ZERO código
   hasta handoff aprobado" en la capa determinística.
4. **Dar execution state y versionado al handoff**. Checklist con timestamps
   dentro del artefacto (o del grafo) para reanudación determinística tras
   crash/compaction, y regla explícita de qué ocurre con lanes en curso si
   el handoff cambia.
5. **Definir la matriz de delivery por ecosistema**. Documentar qué
   mecanismo de cada agente (Claude skills/hooks, Cursor rules, Copilot
   instructions) mapea a advisory/contextual/determinístico, aunque v1 solo
   implemente Claude. Convierte la "independencia" declarada en un contrato
   verificable.
6. **Nacer con suite de tests propia**. Un tool cuyo pitch es verificación
   pierde credibilidad si no se auto-verifica; incluirlo como requisito del
   design, no como afterthought.
7. **No incorporar**: daemon de mensajería, tmux/worktrees por rol, roles
   fijos de pipeline, ni infraestructura de host. Son consecuencias del
   modelo de orquestador de procesos que Virgil correcta y deliberadamente
   no es.
