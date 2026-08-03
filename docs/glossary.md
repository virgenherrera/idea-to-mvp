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
> Ordenados alfabéticamente. Los nombres canónicos usan la convención
> camelCase en inglés.

---

| Término | Definición |
|---------|------------|
| **AAA** (Arrange-Act-Assert) | Patrón estructural obligatorio para todo test. Tres bloques separados: preparar estado, ejecutar operación, verificar resultado. Si un test necesita más de un Act, son dos tests. |
| **abuseCases** | Contrapartida adversarial de los ACs positivos. Para cada AC con entrada de datos, el testPlan debe incluir casos de payload vacío, corrupto, inválido, inyecciones (SQL, NoSQL, XSS, prompt), campos extra no declarados, y abuso de autenticación/autorización. complianceByDesign. |
| **AC** (Acceptance Criterion) | Criterio de aceptación en formato given/when/then que define cuándo un requisito se considera cumplido. |
| **accept** | Gate final de cada iteración de execution donde QA certifica que el producto cumple lo estipulado en el handoff. No valida código — valida producto contra contrato. |
| **adapter** | Implementación pluggable de la universalInterface del artifactStore (local, engram, híbrido, DBMS, etc.). |
| **ADR** (Architecture Decision Record) | Registro de una decisión arquitectónica con contexto, alternativas evaluadas y justificación. Parte de `design.md`. |
| **appTest** | Test que ejerce el stack real de la aplicación sin mocks. Boundary = la app. Tier primario del framework. Detecta droppableCode vía coverage alto obligatorio. |
| **artifactStore** | Capa de persistencia donde viven los artefactos de planificación. Fuera del repo destino. Accesible vía universalInterface de 9 operaciones. |
| **boundaryModel** | Criterio que determina el tipo de test según dónde se ubica la frontera del mock: File (unit, prohibido), Module (integración, derivado), App (servicio/componente, desarrollo explícito), Solution (E2E, desarrollo explícito). Eje organizador de la Fase red. |
| **buildArtifact** | Output generado por el echo: compilados, reportes de cobertura, documentación API, etc. Vive en carpeta gitignoreada dentro del repo (o en registro remoto para artifacts cloud). No confundir con artifacts de planificación (idea.md, spec.md, etc.) que viven en el artifactStore. Ver [artifact system](artifact-system.md). |
| **builderPattern** | Patrón para construir datos de test mediante factories reutilizables. Centraliza la creación de datos y permite variar solo lo relevante al caso. Evita datos hardcodeados en el cuerpo del test. |
| **bumpDependencies** | Automatización habilitada por el echo determinista. Patrón: bump dependencias → ejecutar echo completo → si verde → commit automático. Aborda la tensión entre versiones pinneadas y dependencias modernas. La mecánica concreta se porta a la plataforma del proyecto. Ver [echo system](echo-system.md). |
| **circuitBreaker** | Mecanismo de protección: si 3 delegaciones consecutivas al mismo rol fallan, el SM detiene la cadena y escala al MIM. |
| **compactRules** | Reglas de código y convenciones del proyecto, extraídas del skill registry, que el orquestador inyecta pre-digeridas en cada subAgent. |
| **complianceByDesign** | Principio del framework: si los tests asiertan DTOs estrictamente, incluyen abuseCases, y verifican structuralCompliance (persistencia, frontend, IaC), se obtiene verificación de la CAPA DE DATOS del compliance regulatorio (HIPAA, PCI DSS, GDPR, SOC 2, WCAG, ADA) como efecto secundario — sin suites adicionales ni rewrites. No cubre controles organizacionales, físicos, legales ni procedurales de cada regulación. |
| **compositeAgent** | subAgent que asume múltiples personalidades secuencialmente dentro de un worktree (testEngineer → implementor → reviewer). Se usa en ejecución paralela para evitar conflictos de filesystem. |
| **contractArchitect** | Rol de execution que define contratos formales (APIs, schemas, interfaces) a partir del handoff. Activo en prePhase. |
| **DAG** (Directed Acyclic Graph) | Grafo de dependencias entre tareas en `tasks.md`. Define el orden de ejecución y los lanes paralelos. |
| **delegationContract** | Estructura con campos obligatorios (rol, personalidad, contexto, input, output, restricciones, status report) que el SM usa para lanzar un subAgent. |
| **droppableCode** | Código con 0% de cobertura en appTests. Si ningún test lo toca a través de interacciones reales de producto, no tiene justificación de existir. Candidato a eliminación. |
| **E2E** (End-to-End) | Test que ejerce la solución completa desplegada, multi-servicio, con cero mocks. Boundary = la solución. Segundo tier explícito del framework. Se ejecuta en deploys, tags y merges a main/develop. |
| **echo** | Secuencia determinista de 5 pasos (Setup → Build → Static Test → Dynamic Test → E2E Test) que se ejecuta en todo ambiente (dev, QA, CI, CD) con el mismo orden pero scope variable. Garantiza homogeneidad estructural de ambientes. El echo es obligatorio (TINA). Ver [echo system](echo-system.md). |
| **execution** | Modo 2 del framework. Transforma el handoff en código funcional mediante el ciclo red/green/refactor + accept. Opera sobre el working tree del repositorio. |
| **executionOrchestrator** | Coordinador de execution. Análogo al SM de planning pero opera sobre el working tree del repositorio, no sobre el artifactStore. Delega, no ejecuta. |
| **fastForward** | Mecanismo que permite al SM avanzar múltiples fases cuando el gradiente de certeza (F1-F4) es alto. Aplica al inicio y mid-cycle. |
| **gate** | En planning: punto de validación donde un artefacto debe alcanzar `approved` para desbloquear la siguiente fase. En execution: checkpoint operacional (tests pasan, cobertura cumple umbral, QA certifica). |
| **handoff** | Artefacto (`handoff.md`) que actúa como contrato autocontenido entre planning y execution. Portable y acotado. |
| **highValueTesting** | Filosofía del framework: solo tests que ejercen interacciones REALES de producto aportan valor. Tests con mocks extensivos dan falsa confianza. |
| **implementor** | Rol de execution que escribe código para pasar los tests. Activo en Fase green. Pragmático, sin perfeccionismo prematuro. |
| **lane** | Rama paralela en el DAG de tareas que puede ejecutarse concurrentemente con otras. |
| **MIM** (Mind behind the Idea and the Money) | El humano que dirige el proyecto. Decide, aprueba y desbloquea. Es el nodo de decisión final. |
| **miniPDC** | Versión abreviada del PDC. Se aplica cuando la delegación es de bajo riesgo y el resultado es verificable mecánicamente. |
| **operation** | Modo 3 del framework. Modo opcional y reactivo donde el MIM se convierte en usuario del producto y el agente en operationalAssistant. Sin fases ni ceremonia. Se activa solo si el proyecto tiene superficie operativa. |
| **operationalAssistant** | Rol del agente en operation. Ejecuta lo que el usuario pide dentro del contexto del proyecto construido. No planifica ni construye — opera. |
| **opsRunbook** | Artefacto (`ops-runbook.md`) que documenta operación, monitoreo, troubleshooting y escalación para proyectos con servicios vivos. Respaldado por ISO 20000 + ITIL 4. |
| **patternA** | Estrategia de retrieval donde el SM busca, cura y re-inyecta contexto en el prompt del subAgent. Útil para búsquedas fuzzy o fan-out alto (8+). |
| **patternB** | Estrategia de retrieval donde el SM pasa solo topic_keys y el subAgent consulta el artifactStore directamente. 6x más barato que patternA. |
| **PDC** (Post-Delegation Checkpoint) | Protocolo de 4 pasos que el SM ejecuta después de cada retorno de subAgent: ECHO (coherencia), VERIFY (cobertura), MARK (persistir), DECIDE (siguiente acción). |
| **planning** | Modo 1 del framework. Transforma la idea en un handoff validado mediante fases de planificación, gates de aprobación y el artifactStore. |
| **POM** (Page Object Model) | Patrón para tests con interfaz (UI, CLI). Abstrae interacciones mecánicas en objetos reutilizables. El test describe intención, el POM ejecuta mecánica. |
| **prePhase** | Primera etapa de execution. Define contratos formales antes de escribir tests o código. Habilita desarrollo paralelo. |
| **QA** | Rol de execution activo en accept. Verifica producto contra handoff. No escribe tests ni corrige código — certifica. Mecanismo de certificación definido por el consumidor del framework. |
| **RAG** (Retrieval-Augmented Generation) | El artifactStore usado como fuente de contexto acotado para los agentes. Los agentes consultan slices específicos, no el codebase completo. |
| **red** / **green** / **refactor** | Las 3 fases macro de execution. red = escribir tests (fallan). green = escribir código (tests pasan). refactor = review de calidad (tests siguen pasando). |
| **reviewer** | Rol de execution activo en refactor. Tres variantes: Arquitectura (SOLID, DRY, KISS), Seguridad (OWASP), Performance (memory leaks, N+1). |
| **schemaStrictAssertions** | Regla de disciplina de test: toda aserción sobre objetos de respuesta verifica la forma COMPLETA del DTO (campos presentes, campos ausentes, tipos). Detecta campos extra que podrían violar compliance regulatorio. |
| **semanticDrift** | Desalineación entre un artefacto downstream y su upstream después de que el upstream fue modificado. Detectado por `verifyConsistency`. |
| **SM** (Session Manager / Orquestador) | El agente principal que actúa como facade del proyecto. Orquesta fases, convoca roles, valida gates y controla transiciones. No produce contenido. No es un Scrum Master en el sentido del Scrum Guide. |
| **smProcess** | subAgent especializado que el SM instancia para extraer reglas de proceso de archivos del proyecto (e.g., reglas de un challenge). Instanciado bajo demanda, no es un rol permanente. |
| **structuralCompliance** | Tests que verifican la ESTRUCTURA de cada capa arquitectónica, no su comportamiento. Tres dimensiones condicionales (se activan solo si el proyecto las tiene): persistencia (schema, hashing, cifrado), frontend (A11y, i18n, responsive), infraestructura (versiones exactas, env vars validadas, fail-fast). Tageados como `structural`, ejecutados en CI. |
| **subAgent** | Agente instanciado por el SM con un delegationContract acotado. Recibe personalidad, contexto y restricciones específicas para su tarea. Los roles del equipo son subAgents. |
| **testContract** | Manifiesto enumerable por sujeto bajo prueba. Cada entrada vincula un caso de test con un nombre inmutable y trazable a un AC. Previene código de test spaghetti. |
| **testEngineer** | Rol de execution que escribe la suite de tests mapeada a ACs y contratos. Activo en Fase red. Escéptico, prioriza appTests (stack real) sobre cualquier forma de mocking; unit prohibido, integración derivada por filtrado. |
| **testImplementation** | Tests ejecutables que referencian el testContract. Incluyen appTests (stack real, sin mocks) y tests E2E (solución desplegada). Todos deben fallar al finalizar red (no hay implementación aún). |
| **testPlan** | Meta-documento que mapea ACs a casos de prueba, asigna boundaries (App/E2E), y etiqueta tests para filtrado (smoke, critical, regression). |
| **tier** | Nivel de ceremonia del framework (Ligero, Estándar, Completo). El SM determina el tier usando el score F1-F4. Los tiers escalan ceremonia; en tier Ligero, los artefactos pueden comprimirse en un documento único (`plan.md`) manteniendo las secciones de contenido ISO. |
| **TPM** (Technical Program Manager) | Agente de infraestructura que actúa como DBMS del artifactStore. Único actor que escribe en el store. Valida integridad, formato y completitud con criterio editorial. |
| **transition()** | Operación del adapter que cambia el estado de un artefacto en la state machine (draft, review, approved, rejected, cancelled). |
| **universalInterface** | Las 9 operaciones que todo adapter del artifactStore debe implementar: `ingest`, `save`, `read`, `search`, `list`, `delete`, `verifyConsistency`, `history`, `transition`. |
| **workItem** (L0-L4) | Jerarquía de descomposición de trabajo: L0 (Initiative), L1 (Feature), L2 (Requirement), L3 (Activity), L4 (Sub-activity). Definidos en `tasks.md`. |
| **worktree** | Mecanismo de git que crea un directorio de trabajo aislado con su propia rama. El executionOrchestrator los usa para ejecutar lanes paralelos sin conflictos de archivos. |
