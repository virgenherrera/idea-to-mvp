# Adaptaciones al Manifiesto Ágil para Agentes IA

← [Índice](README.md)

> Este framework toma vocabulario de Scrum y Agile pero opera con un modelo
> de delegación prescriptivo. La razón: los agentes IA carecen de
> persistencia entre sesiones, capacidad de auto-organización y confianza
> interpersonal. Estas diferencias hacen que ciertos principios ágiles
> requieran adaptación consciente, no abandono.

---

## Tabla de Cumplimiento de los 12 Principios

| # | Principio | Cumplimiento | Observación |
|---|-----------|-------------|-------------|
| 1 | Satisfacer al cliente con entrega continua de software valioso | Parcial | Pipeline largo antes de la primera entrega, pero el ciclo es iterativo (Retro → Idea). |
| 2 | Bienvenidos los cambios tardíos en los requisitos | Parcial | Mecanismos existen (`transition` a draft, re-convocación, `verifyConsistency`) pero son costosos operativamente. |
| 3 | Entregar software frecuentemente | Bien servido | Iteraciones dentro del Modo 2, commits frecuentes, ciclo Retro → Idea. |
| 4 | Negocio y desarrollo trabajan juntos diariamente | Bien servido | El MIM interactúa en todas las fases vía el SM. No hay "muro" entre negocio y desarrollo. |
| 5 | Construir proyectos alrededor de individuos motivados y darles confianza | Adaptado | La confianza se reemplaza por verificación sistémica (PDC). Ver justificación abajo. |
| 6 | Comunicación cara a cara como método más eficiente | No aplica | Los agentes IA no tienen "cara". El SM como intermediario estructurado es necesario. Ver justificación abajo. |
| 7 | Software funcionando como medida principal de progreso | Bien servido | Modelo de boundaries (App + E2E) con cadena de trazabilidad AC → Test Plan → Contrato de Test → Implementación → Cobertura. |
| 8 | Ritmo sostenible de desarrollo | No abordado | Sin mención explícita de límites de carga o throttling de agentes. |
| 9 | Excelencia técnica continua y buen diseño | Excelente | Refactor con 7 dimensiones de revisión, ADRs, respaldo ISO, gates de calidad. |
| 10 | Simplicidad: maximizar el trabajo no hecho | Bien servido | Fast-forward evita fases innecesarias, tiers de activación escalan ceremonia, roles se condensan. |
| 11 | Equipos auto-organizados producen las mejores arquitecturas | Adaptado | Prescripción vía contrato es necesaria porque los agentes IA no comparten contexto. Ver justificación abajo. |
| 12 | Reflexión y adaptación regular | Excelente | Fase 8 (Retrospectiva) completa con stop/start/continue/agreements. Alimenta el siguiente ciclo. |

---

## Adaptaciones Clave y Justificación

### Principio 5 — Confianza reemplazada por verificación

El principio original asume individuos con identidad persistente,
reputación acumulada y motivación intrínseca. Los agentes IA no tienen
ninguna de estas propiedades:

- No recuerdan interacciones previas (sin persistencia cross-session).
- No tienen reputación — cada instancia empieza desde cero.
- No tienen motivación — cumplen contratos, no objetivos personales.

El framework reemplaza "confianza" por el **PDC (Post-Delegation
Checkpoint)**: después de cada delegación, el SM verifica coherencia
(ECHO), cobertura (VERIFY), persiste el resultado (MARK) y decide el
siguiente paso (DECIDE). Esto no es desconfianza — es el equivalente
funcional de la confianza en un contexto donde la identidad no persiste.

### Principio 6 — Comunicación estructurada en vez de cara a cara

El principio original privilegia la comunicación de alto ancho de banda
entre humanos. Los agentes IA no se comunican entre sí — reciben contratos
y retornan resultados. El SM actúa como intermediario que:

- Traduce la intención del MIM en contratos de delegación con campos
  obligatorios.
- Recibe status reports estructurados de los sub-agentes.
- Usa el PDC como protocolo de verificación post-retorno.

La "conversación" entre agentes es un intercambio de contratos y
resultados, no un diálogo. Esto es una limitación inherente, no una
elección de diseño.

### Principio 11 — Prescripción necesaria por falta de contexto compartido

Los equipos humanos auto-organizados funcionan porque comparten contexto
implícito: cultura del equipo, decisiones previas, preferencias,
relaciones de confianza. Los agentes IA no comparten nada:

- Cada sub-agente recibe un contexto acotado por contrato.
- No saben qué están haciendo los otros sub-agentes.
- No pueden negociar entre sí ni ajustar su enfoque en tiempo real.

El modelo command-and-control vía contratos rígidos es una compensación
necesaria. El SM centraliza la coordinación que en un equipo humano sería
distribuida. Esto no es una elección ideológica — es la única forma de
producir resultados coherentes cuando los participantes no comparten
estado.

---

## Clarificación de Nomenclatura

### SM no es un Scrum Master

En este framework, **SM significa Session Manager (Orquestador)**, no
Scrum Master según el Scrum Guide. Las diferencias son sustanciales:

| Aspecto | Scrum Master (Scrum Guide) | SM (este framework) |
|---------|---------------------------|---------------------|
| Función | Servant leader, facilita al equipo | Facade, controla gates y delegación |
| Autoridad | No tiene autoridad sobre el producto | Decide convocación, valida outputs, bloquea avances |
| Producción | No produce artefactos | No produce contenido, pero controla transiciones |
| Equipo | Sirve al equipo auto-organizado | Comanda sub-agentes sin autonomía |

El SM del framework tiene funciones que un Scrum Master no tiene:
controlar gates, decidir convocación de roles, validar outputs, y bloquear
avances prematuros. Es más cercano a un **controller** que a un
facilitador.

### PO no es el Product Owner clásico

El MIM es el verdadero decisor de producto: sabe qué quiere, tiene la
visión y el presupuesto. El rol PO en este framework funciona como un
**Business Analyst proxy** que:

- Formaliza las ideas del MIM en artefactos estructurados.
- Desafía las ideas con preguntas de negocio.
- Prioriza requisitos y define ACs.
- No tiene autoridad final — el MIM decide.

---

## Scope del Framework

Este framework está optimizado para el caso **"1 humano (MIM) + N agentes
IA."** Las fases, artefactos y gates son reutilizables para equipos
humanos, pero el modelo de delegación (contratos rígidos, SM como único
punto de interacción, PDC como mecanismo de verificación) debe adaptarse
para contextos donde los participantes tienen persistencia, autonomía y
capacidad de comunicación directa.
