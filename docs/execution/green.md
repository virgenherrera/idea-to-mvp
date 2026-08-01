---
id: execution/green
title: "Fase Green — Implementación"
mode: execution
type: process
tags: [implementación, commits, tdd-micro, código, escalación]
---

# Fase Green — Implementación

← [Índice principal](../README.md) | [Ejecución](README.md)

```mermaid
sequenceDiagram
    participant OE as Orquestador de Ejecución
    participant IMP as Implementor
    participant WT as Working Tree
    participant CI as Tests

    OE->>IMP: Contrato: implementar código<br/>que pase los tests rojos
    activate IMP
    IMP->>WT: Escribe código
    IMP->>CI: Ejecuta tests
    CI-->>IMP: Resultados

    alt Test incorrecto detectado
        IMP-->>OE: "Test X verifica comportamiento<br/>incorrecto según AC-Y"
        OE->>OE: Decide: re-delegar a<br/>Test Engineer o autorizar fix
    end

    IMP-->>OE: Status Report + commits
    deactivate IMP
```

> **Input de Red**: El Implementor recibe la Capa 3 (Test Implementation)
> como entrada directa — los tests ejecutables que debe hacer pasar. Las
> Capas 1 (Test Plan) y 2 (Test Contract) proporcionan trazabilidad pero no
> son input operativo de Green.

## Reglas de Green

La unica meta es hacer que los tests pasen. Nada mas.

```mermaid
flowchart TD
    START["Tests rojos"]
    WRITE["Escribir codigo\nque pase el test"]
    RUN["Ejecutar tests"]
    CHECK{{"¿Pasan?"}}
    COMMIT["Commit\n(incremento verde)"]
    NEXT{{"¿Quedan tests\nrojos?"}}
    DONE["Todos los tests pasan\n✅ GREEN"]

    START --> WRITE
    WRITE --> RUN
    RUN --> CHECK
    CHECK -->|"No"| FIX{{"¿Test incorrecto?"}}
    FIX -->|"Si"| FIX_TEST["Escalar al Orquestador\n(volver a Red)"]
    FIX -->|"No"| WRITE
    FIX_TEST --> WRITE
    CHECK -->|"Si"| COMMIT
    COMMIT --> NEXT
    NEXT -->|"Si"| WRITE
    NEXT -->|"No"| DONE
```

| Regla | Descripcion |
|-------|-------------|
| **Lo primero que funcione** | Codigo feo, duplicado, con magic numbers --- todo vale si los tests pasan |
| **Sin optimizacion prematura** | No abstraer, no generalizar, no "mejorar". Eso es la siguiente fase |
| **Cumplir contratos** | El codigo DEBE respetar los contratos definidos en la Pre-Fase |
| **Commits frecuentes** | Cada test que pasa = un posible commit. Incrementos verdes pequenos |
| **Test incorrecto → corregir test** | Si un test verifica algo equivocado, arreglarlo ANTES de implementar |

> **Excepción: TDD micro para complejidad algorítmica**: Para tareas de
> alta complejidad algorítmica (algoritmos, parsers, cálculos financieros),
> el Implementor puede usar TDD micro (test-implement-test por función)
> como herramienta complementaria dentro de Green. Esta excepción no
> aplica a código de aplicación estándar (CRUD, endpoints, flujos de UI).

---

## Estrategia de commits

```plaintext
feat: implement login endpoint (passes auth-login-success test)
feat: implement login validation (passes auth-login-invalid-credentials test)
feat: implement token refresh (passes auth-token-refresh test)
```

Cada commit referencia que test(s) pasa. Esto crea trazabilidad entre
implementacion y especificacion ejecutable.

---

## Cuando corregir tests vs corregir codigo

```mermaid
flowchart TD
    FAIL["Test falla"]
    Q1{{"¿El test verifica\nel comportamiento correcto\nsegun el AC?"}}
    Q1 -->|"Si"| FIX_CODE["Corregir el CODIGO\n(el test esta bien)"]
    Q1 -->|"No"| Q2{{"¿El AC esta mal\no el test lo\ninterpreta mal?"}}
    Q2 -->|"Test mal escrito"| FIX_TEST["Escalar al Orquestador\n(re-delegar a Test Engineer\no autorizar correccion)"]
    Q2 -->|"AC ambiguo"| ESCALATE["Escalar al Orquestador\n→ re-evaluar contrato"]
```

> **Separación de responsabilidades**: El Implementor NO corrige tests directamente. Si sospecha
> que un test es incorrecto, escala al Orquestador con evidencia (qué test, qué AC contradice,
> por qué). El Orquestador decide si re-delega al Test Engineer o autoriza la corrección in situ.
