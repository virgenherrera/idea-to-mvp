---
id: planning/behavior/recovery
title: "Recovery y Manejo de Fallos"
mode: planning
type: spec
tags: [recovery, circuit-breaker, fallos, escalación, rollback]
---

# Protocolo de Recuperación

← [Índice principal](../../README.md) | [Planificación](../README.md) | [Comportamiento SM](README.md)

## Recovery protocol (inicio de sesión)

```mermaid
sequenceDiagram
    participant SM as SM (nueva sesión)
    participant TPM as TPM

    SM->>TPM: "¿Qué artefactos existen y cuál es su estado?"
    TPM->>TPM: Escanea RAG
    TPM->>SM: "idea.md: aprobado, spec.md: aprobado, design.md: borrador (3/5 secciones)"
    SM->>SM: Deriva: estamos en Fase 3 (Diseñar), design.md en draft
    SM->>TPM: "¿Hay fallos registrados en el ciclo actual?"
    TPM->>TPM: Consulta history() filtrando action: failure
    TPM->>SM: "2 rechazos PDC en design.md (VERIFY), rol Dev Lead"
    SM->>SM: Ajusta estrategia: contrato más explícito o personalidad diferente
    SM->>SM: Siguiente acción: convocar Dev Lead con contrato ajustado
```

---

## Historial de Fallos

El circuit breaker protege intra-sesion, pero los fallos tambien se
registran cross-session en `history()` del artefacto afectado (ver
[TPM e Interfaz del Adaptador](../artifacts/tpm-adapter.md#historyartifact)).
Esto permite al SM aprender
de fallos anteriores al recuperar estado.

**Que se registra**: cada fallo se almacena como entrada en `history()`
con `action: "failure"` y metadata especifica del tipo:

| Tipo | Campos adicionales | Ejemplo |
|------|-------------------|---------|
| `pdc_rejection` | `step` (ECHO/VERIFY/MARK/DECIDE), `role`, `reason` | Rechazo en VERIFY: output no cubre ACs |
| `circuit_breaker` | `role`, `consecutive` | 3 fallos consecutivos del rol QA |
| `escalation` | `role`, `description`, `resolution` | Gap en diseno de auth, MIM proveyo ADR |
| `redelegation` | `role`, `reason`, `contract_delta` | Scope demasiado amplio, acotado a ACs 1-3 |

Formato de registro (todos comparten campos base `action: "failure"`,
`phase`, `timestamp`):

```yaml
# Ejemplo: rechazo PDC
{ action: "failure", type: "pdc_rejection", step: "VERIFY",
  role: "Dev Lead", reason: "output no cubre 2 de 5 ACs", phase: 3 }

# Ejemplo: circuit breaker
{ action: "failure", type: "circuit_breaker",
  role: "QA", consecutive: 3, phase: 6 }
```

**Como el SM usa el historial en recovery**:

1. Despues de derivar la fase actual, el SM pregunta al TPM:
   "Hay fallos registrados en el ciclo actual?"
2. El TPM consulta `history()` de los artefactos en progreso filtrando
   `action: "failure"`.
3. Si existen fallos, el SM ajusta la estrategia antes de re-delegar:
   - **Rechazo PDC recurrente** — contrato mas explicito, personalidad
     del rol ajustada, scope mas acotado.
   - **Circuit breaker previo** — cambiar enfoque del rol o escalar
     tier desde el inicio.
   - **Escalacion resuelta** — inyectar la resolucion del MIM como
     contexto explicito en el nuevo contrato.
   - **Re-delegacion previa** — aplicar el `contract_delta` que
     funciono como baseline del nuevo contrato.
