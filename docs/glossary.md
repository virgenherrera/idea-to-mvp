---
id: glossary
title: "Glosario"
mode: framework
type: reference
tags: [definiciones, términos, vocabulario]
---

# Glosario

← [Índice](README.md)

> Definiciones de los términos especializados del framework idea-to-mvp.
> Ordenados alfabéticamente.

---

| Término | Definición |
|---------|------------|
| **AAA (Arrange-Act-Assert)** | Patrón estructural obligatorio para todo test. Tres bloques separados: preparar estado, ejecutar operación, verificar resultado. Si un test necesita más de un Act, son dos tests. |
| **Abuse Cases / Negative Testing** | Contrapartida adversarial de los ACs positivos. Para cada AC con entrada de datos, el Test Plan debe incluir casos de payload vacío, corrupto, inválido, inyecciones (SQL, NoSQL, XSS, prompt), campos extra no declarados, y abuso de autenticación/autorización. Compliance-by-design. |
| **AC (Acceptance Criterion)** | Criterio de aceptación en formato given/when/then que define cuándo un requisito se considera cumplido. |
| **Accept (Fase)** | Gate final de cada iteración del Modo 2 donde QA certifica que el producto cumple lo estipulado en el handoff. No valida código — valida producto contra contrato. |
| **Adaptador** | Implementación pluggable de la interfaz universal del artifact store (local, engram, híbrido, DBMS, etc.). |
| **ADR (Architecture Decision Record)** | Registro de una decisión arquitectónica con contexto, alternativas evaluadas y justificación. Parte de `design.md`. |
| **Agente compuesto** | Sub-agente que asume múltiples personalidades secuencialmente dentro de un worktree (Test Engineer → Implementor → Reviewer). Se usa en ejecución paralela para evitar conflictos de filesystem. |
| **App Test / Service Test** | Test que ejerce el stack real de la aplicación sin mocks. Boundary = la app. Tier primario del framework. Detecta código droppable vía coverage alto obligatorio. |
| **Artifact de build** | Output generado por el eco: compilados, reportes de cobertura, documentación API, etc. Vive en carpeta gitignoreada dentro del repo (o en registro remoto para artifacts cloud). No confundir con artifacts de planificación (idea.md, spec.md, etc.) que viven en el artifact store. Ver [sistema de artifacts](artifact-system.md). |
| **Artifact store** | Capa de persistencia donde viven los artefactos de planificación. Fuera del repo destino. Accesible vía interfaz universal de 9 operaciones. |
| **Asistente operativo** | Rol del agente en Modo 3. Ejecuta lo que el usuario pide dentro del contexto del proyecto construido. No planifica ni construye — opera. |
| **Boundary (modelo de)** | Criterio que determina el tipo de test según dónde se ubica la frontera del mock: File (unit, prohibido), Module (integración, derivado), App (servicio/componente, desarrollo explícito), Solution (E2E, desarrollo explícito). Eje organizador de la Fase Red. |
| **Builder Pattern (testing)** | Patrón para construir datos de test mediante factories reutilizables. Centraliza la creación de datos y permite variar solo lo relevante al caso. Evita datos hardcodeados en el cuerpo del test. |
| **bumpDependencies** | Automatización habilitada por el eco determinista. Patrón: bump dependencias → ejecutar eco completo → si verde → commit automático. Aborda la tensión entre versiones pinneadas y dependencias modernas. La mecánica concreta se porta a la plataforma del proyecto. Ver [sistema de ecos](echo-system.md). |
| **Circuit breaker** | Mecanismo de protección: si 3 delegaciones consecutivas al mismo rol fallan, el SM detiene la cadena y escala al MIM. |
| **Código droppable** | Código con 0% de cobertura en App tests. Si ningún test lo toca a través de interacciones reales de producto, no tiene justificación de existir. Candidato a eliminación. |
| **Compact rules** | Reglas de código y convenciones del proyecto, extraídas del skill registry, que el orquestador inyecta pre-digeridas en cada sub-agente. |
| **Compliance-by-Design** | Principio del framework: si los tests asiertan DTOs estrictamente, incluyen casos adversariales, y verifican compliance estructural (persistencia, frontend, IaC), se obtiene verificación de la CAPA DE DATOS del compliance regulatorio (HIPAA, PCI DSS, GDPR, SOC 2, WCAG, ADA) como efecto secundario — sin suites adicionales ni rewrites. No cubre controles organizacionales, físicos, legales ni procedurales de cada regulación. |
| **Compliance estructural** | Tests que verifican la ESTRUCTURA de cada capa arquitectónica, no su comportamiento. Tres dimensiones condicionales (se activan solo si el proyecto las tiene): persistencia (schema, hashing, cifrado), frontend (A11y, i18n, responsive), infraestructura (versiones exactas, env vars validadas, fail-fast). Tageados como `structural`, ejecutados en CI. |
| **Contract Architect** | Rol del Modo 2 que define contratos formales (APIs, schemas, interfaces) a partir del handoff. Activo en la Pre-Fase. |
| **Contrato de delegación** | Estructura con campos obligatorios (rol, personalidad, contexto, input, output, restricciones, status report) que el SM usa para lanzar un sub-agente. |
| **DAG (Directed Acyclic Graph)** | Grafo de dependencias entre tareas en `tasks.md`. Define el orden de ejecución y los lanes paralelos. |
| **Drift semántico** | Desalineación entre un artefacto downstream y su upstream después de que el upstream fue modificado. Detectado por `verifyConsistency`. |
| **E2E (End-to-End)** | Test que ejerce la solución completa desplegada, multi-servicio, con cero mocks. Boundary = la solución. Segundo tier explícito del framework. Se ejecuta en deploys, tags y merges a main/develop. |
| **Eco** | Secuencia determinista de 5 pasos (Setup → Build → Static Test → Dynamic Test → E2E Test) que se ejecuta en todo ambiente (dev, QA, CI, CD) con el mismo orden pero scope variable. Garantiza homogeneidad estructural de ambientes. El eco es obligatorio (TINA). Ver [sistema de ecos](echo-system.md). |
| **Fast-forward** | Mecanismo que permite al SM avanzar múltiples fases cuando el gradiente de certeza (F1-F4) es alto. Aplica al inicio y mid-cycle. |
| **Gate** | En planificación: punto de validación donde un artefacto debe alcanzar `approved` para desbloquear la siguiente fase. En ejecución: checkpoint operacional (tests pasan, cobertura cumple umbral, QA certifica). |
| **Handoff** | Artefacto (`handoff.md`) que actúa como contrato autocontenido entre Modo 1 (planificación) y Modo 2 (ejecución). Portable y acotado. |
| **Implementor** | Rol del Modo 2 que escribe código para pasar los tests. Activo en Fase Green. Pragmático, sin perfeccionismo prematuro. |
| **Interfaz universal** | Las 9 operaciones que todo adaptador del artifact store debe implementar: `ingest`, `save`, `read`, `search`, `list`, `delete`, `verifyConsistency`, `history`, `transition`. |
| **Lane** | Rama paralela en el DAG de tareas que puede ejecutarse concurrentemente con otras. |
| **MIM (Mind behind the Idea and the Money)** | El humano que dirige el proyecto. Decide, aprueba y desbloquea. Es el nodo de decisión final. |
| **Mini-PDC** | Versión abreviada del Post-Delegation Checkpoint. Se aplica cuando la delegación es de bajo riesgo y el resultado es verificable mecánicamente. |
| **Modo 3 (Operación)** | Modo opcional y reactivo donde el MIM se convierte en usuario del producto y el agente en asistente operativo. Sin fases ni ceremonia. Se activa solo si el proyecto tiene superficie operativa. |
| **Ops-runbook** | Artefacto (`ops-runbook.md`) que documenta operación, monitoreo, troubleshooting y escalación para proyectos con servicios vivos. Respaldado por ISO 20000 + ITIL 4. |
| **Orquestador de Ejecución** | Coordinador del Modo 2. Análogo al SM de planificación pero opera sobre el working tree del repositorio, no sobre el artifact store. Delega, no ejecuta. |
| **Pattern A** | Estrategia de retrieval donde el SM busca, cura y re-inyecta contexto en el prompt del sub-agente. Útil para búsquedas fuzzy o fan-out alto (8+). |
| **Pattern B** | Estrategia de retrieval donde el SM pasa solo topic_keys y el sub-agente consulta el artifact store directamente. 6x más barato que Pattern A. |
| **PDC (Post-Delegation Checkpoint)** | Protocolo de 4 pasos que el SM ejecuta después de cada retorno de sub-agente: ECHO (coherencia), VERIFY (cobertura), MARK (persistir), DECIDE (siguiente acción). |
| **POM (Page Object Model)** | Patrón para tests con interfaz (UI, CLI). Abstrae interacciones mecánicas en objetos reutilizables. El test describe intención, el POM ejecuta mecánica. |
| **Pre-Fase** | Primera etapa del Modo 2. Define contratos formales antes de escribir tests o código. Habilita desarrollo paralelo. |
| **QA (Modo 2)** | Rol del Modo 2 activo en Fase Accept. Verifica producto contra handoff. No escribe tests ni corrige código — certifica. Mecanismo de certificación definido por el consumidor del framework. |
| **RAG (Retrieval-Augmented Generation)** | El artifact store usado como fuente de contexto acotado para los agentes. Los agentes consultan slices específicos, no el codebase completo. |
| **Red/Green/Refactor** | Las 3 fases macro del Modo 2. Red = escribir tests (fallan). Green = escribir código (tests pasan). Refactor = review de calidad (tests siguen pasando). |
| **Reviewer** | Rol del Modo 2 activo en Fase Refactor. Tres variantes: Arquitectura (SOLID, DRY, KISS), Seguridad (OWASP), Performance (memory leaks, N+1). |
| **Schema-Strict Assertions** | Regla de disciplina de test: toda aserción sobre objetos de respuesta verifica la forma COMPLETA del DTO (campos presentes, campos ausentes, tipos). Detecta campos extra que podrían violar compliance regulatorio. |
| **SM (Session Manager / Orquestador)** | El agente principal que actúa como facade del proyecto. Orquesta fases, convoca roles, valida gates y controla transiciones. No produce contenido. No es un Scrum Master en el sentido del Scrum Guide. |
| **SM-Process** | Sub-agente especializado que el SM instancia para extraer reglas de proceso de archivos del proyecto (e.g., reglas de un challenge). Instanciado bajo demanda, no es un rol permanente. |
| **Sub-agente** | Agente instanciado por el SM con un contrato de delegación acotado. Recibe personalidad, contexto y restricciones específicas para su tarea. Los roles del equipo son sub-agentes. |
| **Test Contract** | Manifiesto enumerable por sujeto bajo prueba. Cada entrada vincula un caso de test con un nombre inmutable y trazable a un AC. Previene código de test spaghetti. |
| **Test Engineer** | Rol del Modo 2 que escribe la suite de tests mapeada a ACs y contratos. Activo en Fase Red. Escéptico, prioriza App tests (stack real) sobre cualquier forma de mocking; unit prohibido, integración derivada por filtrado. |
| **Test Implementation (Capa 3)** | Tests ejecutables que referencian el Test Contract. Incluyen tests de App (stack real, sin mocks) y tests E2E (solución desplegada). Todos deben fallar al finalizar la Fase Red (no hay implementación aún). |
| **Test Plan** | Meta-documento que mapea ACs a casos de prueba, asigna boundaries (App/E2E), y etiqueta tests para filtrado (smoke, critical, regression). |
| **Testing de Alto Valor** | Filosofía del framework: solo tests que ejercen interacciones REALES de producto aportan valor. Tests con mocks extensivos dan falsa confianza. |
| **Tier (de activación)** | Nivel de ceremonia del framework (Ligero, Estándar, Completo). El SM determina el tier usando el score F1-F4. Los tiers escalan ceremonia; en Tier Ligero, los artefactos pueden comprimirse en un documento único (`plan.md`) manteniendo las secciones de contenido ISO. |
| **TPM (Technical Program Manager)** | Agente de infraestructura que actúa como DBMS del artifact store. Único actor que escribe en el store. Valida integridad, formato y completitud con criterio editorial. |
| **`transition()`** | Operación del adaptador que cambia el estado de un artefacto en la state machine (draft, review, approved, rejected, cancelled). |
| **Work item (L0-L4)** | Jerarquía de descomposición de trabajo: L0 (Initiative), L1 (Feature), L2 (Requirement), L3 (Activity), L4 (Sub-activity). Definidos en `tasks.md`. |
| **Worktree** | Mecanismo de git que crea un directorio de trabajo aislado con su propia rama. El Orquestador de Ejecución los usa para ejecutar lanes paralelos sin conflictos de archivos. |
