# Fase Green — Implementación

← [Ejecución](README.md)

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
    FIX -->|"Si"| FIX_TEST["Corregir test primero\n(volver a Red)"]
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
    Q2 -->|"Test mal escrito"| FIX_TEST["Corregir el TEST\n(reescribir segun el AC)"]
    Q2 -->|"AC ambiguo"| ESCALATE["Escalar al Orquestador\n→ re-evaluar contrato"]
```
