---
id: execution/index
title: "Modelo de Ejecución — Modo 2"
mode: execution
type: index
tags: [ejecución, fases, roles, delegación, ciclo-iterativo]
---

# Modelo de Ejecución — Modo 2

← [Índice principal](../README.md)

> Del handoff al codigo funcional. Este modo consume el contrato producido
> por el Modo 1 (planificacion) y produce codigo implementado, probado y
> refactorizado en el working tree del repositorio destino.

---

## Vision General

El Modo 2 transforma un `handoff.md` aprobado en codigo funcional
mediante cinco fases estructurales. La secuencia **Contrato - Red -
Green - Refactor - Accept** no es un ciclo micro por funcion --- es la
columna vertebral macro de toda la ejecucion.

> **TDD a nivel batch, no micro**: Este framework reinterpreta
> Red-Green-Refactor como fases macro — primero toda la suite de tests,
> luego toda la implementación, luego todo el refactoring. Esto es TDD
> por lotes (batch), no TDD función-por-función (micro). El trade-off:
> se pierde feedback rápido entre test individual y su implementación,
> pero se gana una especificación ejecutable completa y la capacidad de
> paralelizar el trabajo del Test Engineer y el Contract Architect. Para
> tareas de alta complejidad algorítmica, el Implementor puede usar TDD
> micro dentro de la Fase Green como herramienta complementaria.

```mermaid
flowchart LR
    HANDOFF["handoff.md\n(entrada)"]

    subgraph EXEC["MODO EJECUCION"]
        direction LR
        C["Pre-Fase\nContratos"]
        R["Fase Red\nTests"]
        G["Fase Green\nImplementacion"]
        RF["Fase Refactor\nCalidad"]
        A["Fase Accept\nCertificacion QA"]
        C --> R --> G --> RF --> A
    end

    CODE["Codigo funcional\n(salida)"]

    HANDOFF -->|"consume"| C
    A -->|"produce"| CODE

    A -.->|"rechazado:\nre-trabajo"| RF
    A -.->|"gap de tests"| R
    A -.->|"gap de contrato"| C
    RF -.->|"regresión → revertir"| RF
    G -.->|"test incorrecto"| R
    R -.->|"contrato ambiguo"| C
    C -.->|"gap critico"| HANDOFF
```

### Tabla de fases

| Fase | Entrada | Salida | Actores |
|------|---------|--------|---------|
| Pre-Fase: Contratos | `handoff.md` (spec, design, tasks) | Contratos formales (API, DB, interfaces) | Orquestador + Contract Architect |
| Red | Contratos + ACs de `spec.md` | Test Plan + Test Contract + Test Implementation (todos fallan) + cobertura configurada | Test Engineer |
| Green | Tests rojos + contratos | Codigo que pasa todos los tests | Implementor |
| Refactor | Codigo verde + tests | Codigo limpio, revisado, alineado a `design.md` | Reviewers (multiples) |
| Accept: Certificación QA | Código refactorizado + test reports + reports de Reviewers + handoff.md + documentación operativa (si requerida por handoff) | Certificación formal de cumplimiento del handoff | QA (Modo 2) |

---

## Ciclo Iterativo

### Red-Green-Refactor dentro de una iteracion

El ciclo macro (Contrato → Red → Green → Refactor) puede ejecutarse
multiples veces dentro de un proyecto:

```mermaid
flowchart TD
    subgraph ITER_1["Iteracion 1 (tasks T-01..T-05)"]
        direction LR
        C1["Contratos"] --> R1["Red"] --> G1["Green"] --> RF1["Refactor"] --> A1["Accept"]
    end

    subgraph ITER_2["Iteracion 2 (tasks T-06..T-10)"]
        direction LR
        C2["Contratos\n(incrementales)"] --> R2["Red\n(nuevos tests)"] --> G2["Green"] --> RF2["Refactor"] --> A2["Accept"]
    end

    ITER_1 --> ITER_2
    ITER_2 -.->|"si hay mas\ntasks"| ITER_N["Iteracion N..."]
```

### Cuando re-entrar a Red

| Situacion | Accion |
|-----------|--------|
| Nuevos ACs descubiertos durante Green | Volver a Red: escribir tests para los nuevos ACs |
| Nuevo contrato necesario (integracion no prevista) | Volver a Pre-Fase: definir contrato, luego Red |
| Bug descubierto durante Refactor | Escribir test que reproduzca el bug (Red), corregir (Green) |
| Requisito cambiado por el MIM | Escalar a planificacion si es structural. Si es menor: actualizar contrato, Red, Green |

### Escalacion a Modo 1 (planificacion)

```mermaid
flowchart TD
    EXEC["Ejecucion en curso"]
    GAP{{"Gap detectado"}}

    GAP -->|"AC ambiguo"| SPEC["Escalar a Fase 2\n(Especificar)"]
    GAP -->|"Arquitectura insuficiente"| DESIGN["Escalar a Fase 3\n(Disenar)"]
    GAP -->|"Tarea faltante"| TASKS["Escalar a Fase 4\n(Desglosar)"]
    GAP -->|"Contradicción entre\nartefactos fuente"| SOURCE["Escalar a la fase\nque posee el artefacto\ncontradictorio"]

    SPEC --> RESOLVE["Modo 1 resuelve el gap"]
    DESIGN --> RESOLVE
    TASKS --> RESOLVE
    SOURCE --> RESOLVE

    RESOLVE --> EXEC
```

> **Nota**: Si la contradicción es entre ACs de spec → Fase 2. Si es
> entre decisiones de design → Fase 3. Si es entre tareas → Fase 4. Fase
> 5 se re-ejecuta DESPUÉS de resolver la contradicción upstream para
> regenerar el handoff con los artefactos corregidos.

El Orquestador de Ejecucion NO resuelve gaps de planificacion --- los
escala. El Modo 1 tiene los roles y la ceremonia para resolverlos. El
Modo 2 opera con lo que recibe; si lo que recibe es insuficiente, lo
devuelve.

---

## Conexion con Modo 1

### Como el handoff alimenta la Pre-Fase

```mermaid
flowchart LR
    subgraph HANDOFF_CONTENT["handoff.md (contiene)"]
        direction TB
        H_SPEC["ACs de spec.md"]
        H_DESIGN["Stack y arquitectura\nde design.md"]
        H_TASKS["Tareas ordenadas\nde tasks.md"]
    end

    subgraph PREFASE["Pre-Fase: Contratos"]
        direction TB
        P_API["API Contracts\n(de ACs + arquitectura)"]
        P_DB["DB Schema\n(de design.md)"]
        P_INT["Interfaces\n(de arquitectura)"]
        P_EVENTS["Event Schemas\n(si aplica)"]
    end

    H_SPEC -->|"ACs → endpoints\ny validaciones"| P_API
    H_DESIGN -->|"modelos y\nrelaciones"| P_DB
    H_DESIGN -->|"capas y\ndependencias"| P_INT
    H_DESIGN -->|"eventos\ny mensajes"| P_EVENTS
    H_TASKS -->|"prioriza orden\nde definición"| P_API
```

> **Nota**: `Estrategia de pruebas` alimenta la Fase Red (no la
> Pre-Fase). `DAG + lanes` alimenta al Orquestador para decisiones de
> paralelismo (ver [git-strategy.md](git-strategy.md)). No se incluyen
> en el diagrama de la Pre-Fase porque no alimentan directamente la
> definición de contratos.

### Artefactos del Modo 1 que informan al Modo 2

| Artefacto Modo 1 | Como lo usa el Modo 2 |
|-------------------|----------------------|
| `spec.md` (ACs) | Cada AC se convierte en uno o mas tests. Trazabilidad directa. |
| `design.md` (arquitectura) | Define la estructura del codigo. El Refactor verifica alineacion. |
| `design.md` (ADRs) | Decisiones tecnicas que restringen la implementacion. |
| `tasks.md` (DAG) | Orden de ejecucion. Lanes paralelos. Ruta critica. |
| `tasks.md` (work items) | Cada L3/L4 es una unidad de trabajo en Green. |
| `handoff.md` | Contrato autocontenido. Punto de entrada del Modo 2. |

### Feedback loop: ejecucion → planificacion

| Evento en ejecucion | Feedback a planificacion |
|---------------------|---------------------------|
| AC no implementable como esta escrito | `spec.md` necesita reformulacion (Fase 2) |
| Arquitectura insuficiente para un AC | `design.md` necesita ADR adicional (Fase 3) |
| Tarea faltante descubierta | `tasks.md` necesita actualizacion (Fase 4) |
| Contradiccion entre ACs | `spec.md` tiene conflicto interno (Fase 2) |
| Dependencia externa no documentada | `design.md` necesita componente (Fase 3) |

---

## Roles del Modo de Ejecucion

### Tabla de roles

| Rol | Personalidad | Fase activa | Responsabilidad | Input | Output |
|-----|-------------|-------------|-----------------|-------|--------|
| **Orquestador de Ejecucion** | Metodico, orientado a flujo. Delega, no ejecuta. Analogo al SM en planificacion. | Todas | Lee handoff, coordina las 5 fases, delega a sub-agentes, valida resultados, gestiona commits. | `handoff.md` + AGENTS.md del repo | Codigo implementado en el working tree |
| **Contract Architect** | Preciso, orientado a interfaces. Piensa en consumidores del contrato. | Pre-Fase | Define contratos formales basados en la arquitectura y ACs. | `design.md` + `spec.md` (via handoff) | Contratos tipados (OpenAPI, schemas, interfaces) |
| **Test Engineer** | Esceptico, orientado a cobertura real. Prioriza App tests (stack real) sobre cualquier forma de mocking; unit prohibido, integración derivada por filtrado. | Red | Escribe la suite completa de tests mapeada a ACs y contratos. | Contratos + ACs | Test Plan + Test Contract + Test Implementation (todos fallan) + coverage config |
| **Implementor** | Pragmatico, orientado a "que funcione". Sin perfeccionismo prematuro. | Green | Escribe codigo que pase los tests. Commits frecuentes. | Tests rojos + contratos | Codigo que pasa los tests |
| **Reviewer (Arquitectura)** | Critico, orientado a mantenibilidad. Compara contra design.md. | Refactor | Revisa alineacion arquitectonica, SOLID, DRY, KISS, patrones. | Codigo verde + design.md | Reporte de revision + sugerencias de refactor |
| **Reviewer (Seguridad)** | Paranoico constructivo. Busca vulnerabilidades. | Refactor | Revisa OWASP Top 10, secrets, dependencias, surface area. | Codigo verde + spec.md (no-funcionales) | Reporte de seguridad |
| **Reviewer (Performance)** | Analitico, orientado a metricas. Busca bottlenecks. | Refactor | Revisa memory leaks, N+1, operaciones bloqueantes. | Codigo verde | Reporte de performance |
| **QA (Modo 2)** | Verificador exhaustivo. Contrasta producto contra handoff. No asume que "tests pasan" es suficiente. | Accept | Verifica que CADA AC del handoff se cumple en el producto, que la cobertura no bajó, que el comportamiento de producto es el esperado, y que la documentación operativa declarada en el handoff existe. Certifica formalmente. | Código refactorizado + test reports + handoff.md + documentación operativa (si requerida por handoff) | Certificación formal (mecanismo definido por el consumidor del framework) |
| **MIM** | Humano. Decide, aprueba, desbloquea. | Todas (on-demand) | Aprueba contratos, resuelve ambiguedades, acepta resultado final. | Reportes del Orquestador | Decisiones y aprobaciones |

### Mapeo a roles de planificacion

```mermaid
flowchart LR
    subgraph PLAN["Modo 1 (Planificacion)"]
        direction TB
        SM_P["SM\n(orquesta)"]
        DEV_P["Dev Lead\n(arquitectura)"]
        QA_P["QA\n(testeabilidad)"]
        SEC_P["DevSecOps\n(seguridad)"]
    end

    subgraph EXEC_R["Modo 2 (Ejecucion)"]
        direction TB
        OE["Orquestador\nde Ejecucion"]
        CA["Contract\nArchitect"]
        TE["Test\nEngineer"]
        IMP["Implementor"]
        REV["Reviewers"]
        QA_E["QA\n(Certificación)"]
    end

    SM_P -.->|"analogo"| OE
    DEV_P -.->|"informa"| CA
    QA_P -.->|"informa"| TE
    DEV_P -.->|"informa"| IMP
    SEC_P -.->|"informa"| REV
    QA_P -.->|"informa"| QA_E
```

Los roles del Modo 2 NO son los mismos que los del Modo 1. En
planificacion, los roles son **lentes de revision** que evaluan
artefactos. En ejecucion, los roles son **ejecutores** que producen
codigo. La relacion es de **influencia** (las decisiones de planificacion
guian la ejecucion), no de identidad.

---

## Modelo de Delegacion del Orquestador

El Orquestador de Ejecucion sigue el mismo patron de delegacion
documentado en
[Comportamiento SM](../planning/behavior/README.md):
contratos de delegacion con campos obligatorios, Status Report, y PDC
(Post-Delegation Checkpoint).

```mermaid
sequenceDiagram
    participant OE as Orquestador de Ejecucion
    participant SUB as Sub-agente (rol)
    participant REPO as Working Tree

    OE->>SUB: Contrato de delegacion<br/>(rol, personalidad, contexto,<br/>input, output, restricciones)
    activate SUB
    SUB->>REPO: Lee/escribe en el working tree
    SUB-->>OE: Resultado + Status Report
    deactivate SUB

    Note over OE: PDC (4 pasos)

    OE->>OE: 1. ECHO: ¿coherente con contrato?
    OE->>OE: 2. VERIFY: ¿tests pasan? ¿coverage?
    OE->>OE: 3. MARK: commit si aplica
    OE->>OE: 4. DECIDE: ¿siguiente fase o re-delegar?
```

### Diferencias con el SM de planificacion

| Aspecto | SM (Modo 1) | Orquestador (Modo 2) |
|---------|-------------|------------------------|
| Donde escribe | Artifact store (fuera del repo) | Working tree del repo |
| Que produce | Artefactos de planificacion | Codigo, tests, commits |
| Roles que convoca | Equipo (lentes) | Ejecutores (code writers) |
| Validacion | Gates de artefactos | Tests pasan + coverage |
| Escalacion | Al MIM | Al MIM o de vuelta a Modo 1 |

---

## Que Sigue

Areas dentro del Modo 2 que requieren definicion adicional:

| Area | Estado | Descripcion |
|------|--------|-------------|
| Contratos de delegacion detallados | TBD | Plantillas completas para cada rol del Modo 2 (como `roles/` en Modo 1) |
| ~~Paralelismo en ejecucion~~ | DEFINIDO | Ver [git-strategy.md](git-strategy.md) — worktrees por lane, detección de conflictos, merge controlado |
| ~~Commit strategy~~ | DEFINIDO | Ver [git-strategy.md](git-strategy.md) — convención por fase, trazabilidad AC→test→commit, squash policy |
| ~~CI/CD integration~~ | DEFINIDO | Ver [sistema de ecos](../echo-system.md) — pipeline determinista de 5 pasos, enforcement vía hooks y CI, homogeneidad de ambientes |
| Metricas de ejecucion | TBD | Coverage thresholds, tiempos de ciclo, tasa de re-delegacion |
| Modo 3 (Operacion) | DEFINIDO | Opcional y reactivo. El MIM opera el producto con asistencia del agente. Ver [modelo de operación](../operation/README.md) |

---

## Contenido de esta sección

Este documento se divide en siete páginas:

| Página | Contenido |
|--------|-----------|
| **README.md** (este documento) | Vision general, ciclo iterativo, conexion con Modo 1, roles, modelo de delegacion |
| [Contratos](contracts.md) | Pre-Fase: contract-first, tipos de contrato, desarrollo paralelo, validacion |
| [Fase Red](red.md) | Modelo de boundaries, arquitectura de 3 capas, tests derivados y código droppable |
| [Fase Green](green.md) | Reglas de implementacion, estrategia de commits, test vs codigo |
| [Fase Refactor](refactor.md) | Gate de calidad, dimensiones de revision, checklist |
| [Fase Accept](accept.md) | Certificación QA: qué verifica, contra qué, cómo certifica |
| [Estrategia Git](git-strategy.md) | Gitflow, worktrees para paralelismo, commits, merge strategy |

---

## Indice de Documentos Relacionados

| Documento | Relacion con este |
|-----------|-------------------|
| [Vista general](../overview.md) | Mapa del framework completo |
| [Modelo operativo](../planning/operational-model.md) | Define los dos modos y sus limites |
| [Artefactos](../planning/artifacts/README.md) | Define `handoff.md` (input de este modo) y los 6 artefactos |
| [Comportamiento SM](../planning/behavior/README.md) | Patron de delegacion y PDC que el Orquestador adapta |
| [Roles](../planning/roles/README.md) | Roles de planificacion que informan los roles de ejecucion |
| [Sistema de ecos](../echo-system.md) | Pipeline determinista de 5 pasos que CI ejecuta; enforcement vía hooks |
| [Sistema de artifacts](../artifact-system.md) | Dónde aterrizan los outputs del eco (builds, reportes, cobertura) |
