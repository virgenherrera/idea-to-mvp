---
id: echo-system
title: "Sistema de Ecos"
mode: framework
type: reference
tags: [ecos, pipeline, homogeneidad, CI/CD, hooks, ambientes, bumpDependencies]
---

# Sistema de Ecos

← [Índice](README.md)

> Un eco es una secuencia determinista de 5 pasos que se ejecuta en
> todo ambiente — dev, QA, CI, CD. La garantía es estructural: los
> mismos pasos corren en el mismo orden en cada ambiente. Lo que varía
> es el **scope** (dev prioriza feedback rápido, CI prioriza
> completitud) — pero ningún paso se omite ni se reordena.

---

## Por Qué Existe

El framework define qué construir (Modo 1), cómo construirlo (Modo 2) y
cómo operarlo (Modo 3). Pero ninguno de esos modos define **cómo
verificar que el entorno de trabajo es confiable** en cada ambiente donde
el código se ejecuta.

Sin un pipeline determinista compartido entre ambientes:

- Un test puede pasar en dev porque las dependencias están cacheadas,
  y fallar en CI porque no se ejecutó el setup.
- Un linter puede correr en dev pero no en CI, permitiendo que código
  con violaciones llegue a producción.
- Un build puede funcionar localmente con una versión flotante de una
  dependencia, y romperse cuando CI instala una versión diferente.

El eco elimina estas discrepancias. No es un "nice to have" de DevOps
— es infraestructura fundacional que habilita la confiabilidad de todo
lo que el framework promete.

---

## Los 5 Pasos

El eco siempre tiene 5 pasos, siempre en este orden. Cada paso tiene un
propósito, una entrada, una salida y un criterio de fallo binario: pasa
o no pasa.

```mermaid
flowchart LR
    S["1. Setup"]
    B["2. Build"]
    ST["3. Static\nTest"]
    DT["4. Dynamic\nTest"]
    E2E["5. E2E\nTest"]

    S --> B --> ST --> DT --> E2E

    S -.-|"deps instaladas\naudit limpio"| B
    B -.-|"build\nexitoso"| ST
    ST -.-|"código\nlimpio"| DT
    DT -.-|"tests pasan\ncoverage cumple"| E2E
```

### Paso 1 — Setup

| Atributo | Valor |
|----------|-------|
| Propósito | Garantizar que las dependencias están instaladas y libres de vulnerabilidades conocidas |
| Entrada | Manifiesto de dependencias + lockfile |
| Salida | Dependencias instaladas, sin vulnerabilidades conocidas de severidad crítica |
| Fallo | Dependencia faltante, lockfile desactualizado, vulnerabilidad crítica sin fix disponible |

El setup incluye la instalación de dependencias y, cuando el ecosistema
lo soporte, la auditoría de seguridad (equivalente a `audit fix`). En
ecosistemas sin herramienta de auditoría para dependencias, el paso se
limita a instalación y verificación de lockfile — la ausencia de
auditoría se documenta en `design.md` como limitación de stack. Este
paso se ejecuta siempre. En proyectos sin dependencias externas, el paso
valida que el manifiesto refleja esa decisión (o ejecuta un no-op
documentado).

### Paso 2 — Build

| Atributo | Valor |
|----------|-------|
| Propósito | Transformar el código fuente en artifacts ejecutables o distribuibles |
| Entrada | Código fuente + dependencias instaladas |
| Salida | Artifacts de build (compilados, transpilados, bundleados) |
| Fallo | Error de compilación, error de tipos, error de bundling |
| Condicional | Proyectos puramente interpretados sin paso de build pueden marcar este paso como no-op documentado |

El build produce los artifacts que el
[sistema de artifacts](artifact-system.md) define dónde aterrizan.

### Paso 3 — Static Test

| Atributo | Valor |
|----------|-------|
| Propósito | Verificar que el código cumple las reglas de estilo, formato y análisis estático del proyecto |
| Entrada | Código fuente |
| Salida | Código sin violaciones de linting ni formato |
| Fallo | Violación de regla de linting, error de formato, warning tratado como error |

El análisis estático no es cosmético — es la primera línea de defensa
contra patrones problemáticos, imports no usados, variables sin
declarar, y violaciones de convenciones del equipo. La configuración
de las herramientas se define en
[`design.md`](planning/artifacts/schemas.md) (sección "Restricciones de
infraestructura").

### Paso 4 — Dynamic Test

| Atributo | Valor |
|----------|-------|
| Propósito | Ejecutar la suite de tests del proyecto y verificar cobertura |
| Entrada | Artifacts de build (o código fuente) + suite de tests |
| Salida | Reporte de tests + reporte de cobertura |
| Fallo | Test que falla, cobertura por debajo del umbral del proyecto |

Este paso ejecuta los tests definidos en la
[Fase Red](execution/red.md) — App tests como tier primario,
integración como derivado. El umbral de cobertura es obligatorio para
stacks con herramientas de coverage maduras y nunca puede bajarse (ver
[Fase Red](execution/red.md), sección "Código Droppable"). Para
stacks donde coverage no es medible o semánticamente relevante (IaC,
data pipelines), `design.md` DEBE declarar la métrica de verificación
alternativa (ej: tasa de compliance de políticas, mutation testing
score, conformance de contratos). La métrica alternativa tiene la misma
regla: una vez establecida, no puede bajarse.

### Paso 5 — E2E Test

| Atributo | Valor |
|----------|-------|
| Propósito | Verificar la solución completa desplegada, multi-servicio, con cero mocks |
| Entrada | Solución desplegada en un ambiente accesible |
| Salida | Reporte de tests E2E |
| Fallo | Escenario E2E que falla |
| Condicional | Solo si el proyecto tiene superficie E2E. Si no aplica, se documenta como excepción |

E2E se ejecuta en deploys, tags y merges a ramas principales (ver
[tabla de pipeline placement](execution/red.md) para la distribución
detallada de qué corre cuándo).

---

## Homogeneidad de Ambientes

La propiedad fundamental del eco es que los mismos 5 pasos se ejecutan
en todo ambiente, en el mismo orden. Lo que varía entre ambientes es el
**scope** de cada paso y el **trigger** que lo invoca — no los pasos ni
su secuencia. Dev prioriza feedback rápido (scope selectivo), CI
prioriza completitud (scope amplio), CD prioriza confianza total (scope
completo + smoke post-deploy).

```mermaid
flowchart TD
    subgraph DEV["Ambiente: Dev"]
        direction LR
        D1["Setup"] --> D2["Build"] --> D3["Static"] --> D4["Dynamic\n(módulo tocado)"] --> D5["E2E\n(si aplica)"]
    end

    subgraph CI["Ambiente: CI (PR)"]
        direction LR
        C1["Setup"] --> C2["Build"] --> C3["Static"] --> C4["Dynamic\n(todos los módulos)"] --> C5["E2E\n(si aplica)"]
    end

    subgraph CD["Ambiente: CD (deploy)"]
        direction LR
        CD1["Setup"] --> CD2["Build"] --> CD3["Static"] --> CD4["Dynamic\n(completo)"] --> CD5["E2E\n(smoke +\ncompleto)"]
    end

    DEV -.->|"mismo eco\ndiferente scope"| CI
    CI -.->|"mismo eco\ndiferente scope"| CD
```

### Diferencias por ambiente

| Aspecto | Dev (hooks) | CI (PR) | CD (deploy) |
|---------|-------------|---------|-------------|
| Trigger | Pre-commit / pre-push | Push a PR, merge request | Tag, merge a main/develop |
| Scope del paso 4 | Módulo tocado | Todos los módulos afectados | Suite completa |
| Scope del paso 5 | Opcional (subset smoke) | Suite E2E si aplica | Suite E2E completa + smoke post-deploy |
| Velocidad vs confianza | Prioriza feedback rápido | Balance | Prioriza confianza total |

La tabla de distribución detallada está en la
[Fase Red](execution/red.md) (sección "Pipeline placement").

---

## Enforcement

El eco no es una recomendación — es obligatorio. El mecanismo de
enforcement depende del ambiente:

```mermaid
flowchart TD
    subgraph ENFORCEMENT["Mecanismo por ambiente"]
        direction TB

        HOOKS["Dev: Git hooks\n(pre-commit, pre-push)"]
        PIPELINE["CI: Pipeline stages\n(configuración del CI system)"]
        GATES["CD: Deployment gates\n(eco completo como condición)"]
    end

    HOOKS -->|"mismos pasos"| PIPELINE
    PIPELINE -->|"mismos pasos"| GATES
```

### En desarrollo (git hooks)

Los hooks de git son el enforcement local. La distribución entre
pre-commit y pre-push es una decisión de proyecto documentada en
[`design.md`](planning/artifacts/schemas.md) (sección "Restricciones de
infraestructura") y declarada en el
[`handoff.md`](planning/artifacts/schemas.md) (sección "Restricciones de
ejecución").

Distribución por defecto:

| Hook | Pasos que ejecuta | Justificación |
|------|-------------------|---------------|
| Pre-commit | 3 (static test) | Feedback inmediato sobre formato y lint |
| Pre-push | 1 → 2 → 3 → 4 (selectivo) | Verificación completa antes de compartir código |

Esta distribución es un default — la distribución exacta la decide el
proyecto y se documenta en `design.md`. Algunos proyectos pueden incluir
tests de App (paso 4, módulo tocado) en pre-commit para feedback más
rápido (ver [tabla de pipeline placement](execution/red.md)). El
principio invariante: **nunca pushear código que no pase el eco** (al
menos hasta el paso 4).

### Presupuesto de tiempo

Cuando el eco completo (pasos 1-4) excede un tiempo tolerable para el
workflow del desarrollador (ej: monorepos grandes, builds compilados),
el proyecto define un **presupuesto de tiempo** para el pre-push en
`design.md`. Los pasos que no caben en el presupuesto se difieren a CI,
documentando explícitamente el trade-off: el developer puede pushear
código que CI podría rechazar. El eco sigue corriendo completo en CI —
el presupuesto solo afecta la distribución local en hooks.

### En CI

El pipeline de CI ejecuta los 5 pasos. El scope de cada paso depende
del trigger: en PRs, el paso 5 puede limitarse al subset de seguridad;
en merges a ramas principales, la suite E2E completa (ver
[tabla de pipeline placement](execution/red.md)). Si algún paso falla,
el pipeline se detiene — no hay punto en correr tests dinámicos si el
build falló, ni E2E si los tests de App no pasan.

### En CD

El deployment gate exige eco verde completo como precondición. Post-
deploy, un subset de E2E (smoke) verifica que el despliegue fue exitoso.

---

## Conexión con el Framework

El eco es transversal — se define, implementa, verifica y explota a lo
largo de los tres modos.

```mermaid
flowchart TD
    subgraph MODE1["Modo 1 — Planificación"]
        direction TB
        M1_DESIGN["design.md define:\n- Herramientas de cada paso\n- Distribución de hooks\n- Umbral de cobertura"]
        M1_HANDOFF["handoff.md declara:\n- Compliance del eco como\n  restricción de ejecución\n- Hooks requeridos"]
    end

    subgraph MODE2["Modo 2 — Ejecución"]
        direction TB
        M2_RED["Fase Red: suite de tests\n(paso 4 y 5 del eco)"]
        M2_GREEN["Fase Green: implementación\nno debe romper pasos 1-4"]
        M2_REFACTOR["Fase Refactor: quality gates\nalineados con paso 3"]
        M2_ACCEPT["Fase Accept: QA verifica\nque el eco completo pasa"]
    end

    subgraph MODE3["Modo 3 — Operación"]
        direction TB
        M3_BUMP["bumpDependencies:\nautomatización habilitada\npor eco determinista"]
    end

    MODE1 --> MODE2 --> MODE3
```

### Dónde se configura cada aspecto

| Qué | Dónde | Cuándo |
|-----|-------|--------|
| Herramientas de cada paso | `design.md` — Restricciones de infraestructura | Fase 3 (Diseñar) |
| Distribución de hooks | `design.md` — Restricciones de infraestructura | Fase 3 (Diseñar) |
| Umbral de cobertura | `design.md` — Restricciones de infraestructura | Fase 3 (Diseñar) |
| Compliance del eco | `handoff.md` — Restricciones de ejecución | Fase 5 (Handoff) |
| Implementación de hooks | Working tree del repo | Modo 2 (Pre-Fase o Green) |
| Verificación del eco | Fase Accept | Modo 2 (Accept) |
| Explotación (bumpDeps) | Operación | Modo 3 |

---

## Automatización Habilitada: bumpDependencies

Cuando el eco es determinista y confiable, habilita una automatización
fundamental: la actualización automatizada de dependencias.

```mermaid
flowchart LR
    BUMP["Bump\ndependencias"]
    ECO["Ejecutar\neco completo\n(5 pasos)"]
    CHECK{{"¿Todo\nverde?"}}
    COMMIT["Commit\nautomático"]
    REPORT["Reporte de\nfallo"]

    BUMP --> ECO --> CHECK
    CHECK -->|"Sí"| COMMIT
    CHECK -->|"No"| REPORT
```

El patrón es simple:

1. Actualizar una o más dependencias en los manifiestos del proyecto
2. Ejecutar el eco completo (los 5 pasos)
3. Si todo pasa → commit automático
4. Si algo falla → reporte para intervención manual

Este patrón aborda la tensión inherente del framework:

- La [Fase Red](execution/red.md) exige versiones exactas (pinned) para
  builds reproducibles.
- La [Fase Red](execution/red.md) exige dependencias modernas (última
  versión estable).
- La [Fase Refactor](execution/refactor.md) verifica dependencias sin
  CVEs conocidos.

Sin un mecanismo de actualización, las versiones pinneadas se vuelven
obsoletas y vulnerables. El eco determinista es lo que hace viable la
actualización automatizada — sin él, bumping es una apuesta.

### Consideraciones del patrón

| Aspecto | Guía |
|---------|------|
| Patch / minor | Automatizables — el eco verde confirma compatibilidad |
| Major (breaking) | Requieren migración manual — tratarlos como trabajo planificado (Modo 1), no como bump automatizado |
| Peer dependencies | Deben bumpearse atómicamente como grupo (ej: react + react-dom + @types/react) |
| Proyectos polyglot | Cada package manager tiene su propio manifiesto; los bumps pueden necesitar coordinación entre managers |
| Frecuencia y agrupación | Decisión de proyecto documentada en `design.md` |

La mecánica concreta (herramienta de bump, estrategia de agrupación,
frecuencia) se porta a la plataforma del proyecto. El patrón es
universal; las decisiones de implementación no.

---

## Adaptabilidad

El eco es prescriptivo en su estructura (5 pasos, en orden) pero
adaptable en su contenido:

| Aspecto | Fijo | Adaptable |
|---------|------|-----------|
| Número de pasos | 5, siempre | — |
| Orden de pasos | Setup → Build → Static → Dynamic → E2E | — |
| Criterio de éxito | Binario (pasa / no pasa) | — |
| Herramientas de cada paso | — | Definidas en `design.md` por proyecto |
| Scope por ambiente | — | Dev (selectivo) vs CI (completo) |
| Distribución en hooks | — | Pre-commit vs pre-push vs ambos |
| Pasos condicionales | — | Cualquier paso puede ser no-op documentado cuando no aplica al stack |

### Unidad de ejecución

El eco opera a nivel de **unidad independientemente buildeable y
testeable**. En un proyecto simple, esa unidad es el proyecto completo.
En un monorepo con múltiples packages, cada package tiene su propia
instancia del eco. En un proyecto polyglot (ej: backend Go + frontend
TypeScript), cada stack tiene su propio eco con sus propias herramientas.

| Estructura | Unidad del eco | Orquestación |
|------------|----------------|--------------|
| Proyecto simple | El proyecto completo | Directa (1 eco) |
| Monorepo (workspaces) | Cada package independiente | El orchestrator del monorepo ejecuta ecos selectivamente por packages afectados |
| Polyglot | Cada stack | Cada stack define sus herramientas; el eco de proyecto los orquesta |

En dev (hooks), el eco corre solo para las unidades afectadas por el
cambio. En CI, corre para todas las unidades afectadas más sus
dependientes. En CD, corren todos los ecos.

### Stacks no convencionales

El modelo de 5 pasos está diseñado para proyectos de software con ciclo
build-test. Para stacks donde los pasos no mapean directamente (IaC,
data pipelines, generadores de sitios estáticos), el proyecto define en
`design.md` cómo cada paso del eco se traduce a su contexto:

| Paso del eco | IaC (ejemplo) | Data pipeline (ejemplo) |
|--------------|---------------|-------------------------|
| Setup | Instalar providers/plugins | Instalar dependencias de pipeline |
| Build | `plan` / `preview` (validación, no artifact distribuible) | Compilar DAGs / transformaciones |
| Static | Linting de HCL/YAML, policy-as-code (OPA, Sentinel) | Linting de scripts, validación de schemas |
| Dynamic | Policy compliance tests, validación de plan | Tests de transformación con datos de prueba |
| E2E | Deploy a ambiente efímero + verificación | Ejecución end-to-end con dataset de prueba |

El principio se mantiene: 5 pasos, en orden, binarios. Lo que cambia es
QUÉ ejecuta cada paso, no la estructura.

### Excepciones documentadas

Cuando un paso no aplica al proyecto (ej: una librería algorítmica pura
sin superficie E2E, o un proyecto IaC donde Build produce un plan en
vez de un artifact distribuible), se documenta como excepción en
`design.md` y se declara en `handoff.md`. La excepción sigue el formato
estándar del framework (ver
[excepciones documentadas](agile-adaptations.md)).

El paso se marca como no-op en el eco, no se elimina. Los 5 pasos
siempre existen conceptualmente — un paso que no aplica ejecuta un
no-op exitoso, no desaparece.

---

## Gaps que Este Sistema Resuelve

| Gap identificado | Dónde existía | Cómo lo aborda el eco |
|------------------|---------------|----------------------|
| CI/CD integration TBD | [execution/README.md](execution/README.md) | El eco ES la definición del pipeline que CI ejecuta |
| Sin mención de análisis estático | 29 docs, 0 referencias a linting/formatting | Paso 3 (Static Test) lo formaliza como obligatorio |
| Hooks mencionados pero no especificados | [schemas.md](planning/artifacts/schemas.md) | Enforcement local del eco via pre-commit/pre-push |
| Sin concepto de homogeneidad de ambientes | Framework completo | Propiedad fundamental del eco: mismos pasos, todo ambiente |
| Tensión versiones pinned ↔ deps modernas | [red.md](execution/red.md) | bumpDependencies como patrón habilitado por eco determinista |
| CI como participante sin definición | [green.md](execution/green.md), [git-strategy.md](execution/git-strategy.md) | CI ejecuta el eco — ahora está definido |

---

## Índice de Documentos Relacionados

| Documento | Relación con este |
|-----------|-------------------|
| [Sistema de artifacts](artifact-system.md) | Paso 2 (Build) produce artifacts; pasos 4 y 5 producen reportes de tests y cobertura. El sistema de artifacts define DÓNDE aterrizan |
| [Fase Red](execution/red.md) | Define la suite de tests (paso 4 y 5) y la tabla de pipeline placement |
| [Fase Green](execution/green.md) | La implementación no debe romper los pasos 1-4 del eco |
| [Fase Refactor](execution/refactor.md) | Quality gates alineados con paso 3; verificación de CVEs |
| [Fase Accept](execution/accept.md) | QA verifica que el eco completo pasa como parte de la certificación |
| [Schemas](planning/artifacts/schemas.md) | `design.md` define las herramientas; `handoff.md` declara compliance |
| [Estrategia Git](execution/git-strategy.md) | Hooks y CI como participantes del lifecycle de worktrees |
