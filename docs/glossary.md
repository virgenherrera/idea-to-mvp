# Glosario

← [Índice](README.md)

> Definiciones de los términos especializados del framework idea-to-mvp.
> Ordenados alfabéticamente.

---

| Término | Definición |
|---------|------------|
| **AC (Acceptance Criterion)** | Criterio de aceptación en formato given/when/then que define cuándo un requisito se considera cumplido. |
| **Accept (Fase)** | Gate final del Modo 2 donde QA certifica que el producto cumple lo estipulado en el handoff. No valida código — valida producto contra contrato. |
| **Adaptador** | Implementación pluggable de la interfaz universal del artifact store (local, engram, híbrido, DBMS, etc.). |
| **ADR (Architecture Decision Record)** | Registro de una decisión arquitectónica con contexto, alternativas evaluadas y justificación. Parte de `design.md`. |
| **Agente compuesto** | Sub-agente que asume múltiples personalidades secuencialmente dentro de un worktree (Test Engineer → Implementor → Reviewer). Se usa en ejecución paralela para evitar conflictos de filesystem. |
| **App Test / Service Test** | Test que ejerce el stack real de la aplicación sin mocks. Boundary = la app. Tier primario del framework. Detecta código droppable vía coverage alto obligatorio. |
| **Artifact store** | Capa de persistencia donde viven los artefactos de planificación. Fuera del repo destino. Accesible vía interfaz universal de 9 operaciones. |
| **Circuit breaker** | Mecanismo de protección: si 3 delegaciones consecutivas al mismo rol fallan, el SM detiene la cadena y escala al MIM. |
| **Código droppable** | Código con 0% de cobertura en App tests. Si ningún test lo toca a través de interacciones reales de producto, no tiene justificación de existir. Candidato a eliminación. |
| **Compact rules** | Reglas de código y convenciones del proyecto, extraídas del skill registry, que el orquestador inyecta pre-digeridas en cada sub-agente. |
| **Contract Architect** | Rol del Modo 2 que define contratos formales (APIs, schemas, interfaces) a partir del handoff. Activo en la Pre-Fase. |
| **Contrato de delegación** | Estructura con campos obligatorios (rol, personalidad, contexto, input, output, restricciones, status report) que el SM usa para lanzar un sub-agente. |
| **DAG (Directed Acyclic Graph)** | Grafo de dependencias entre tareas en `tasks.md`. Define el orden de ejecución y los lanes paralelos. |
| **Drift semántico** | Desalineación entre un artefacto downstream y su upstream después de que el upstream fue modificado. Detectado por `verifyConsistency`. |
| **Fast-forward** | Mecanismo que permite al SM avanzar múltiples fases cuando el gradiente de certeza (F1-F4) es alto. Aplica al inicio y mid-cycle. |
| **Gate** | Punto de validación entre fases. Un artefacto debe alcanzar el estado `approved` (check estructural + check semántico) para que el SM desbloquee la siguiente fase. |
| **Handoff** | Artefacto (`handoff.md`) que actúa como contrato autocontenido entre Modo 1 (planificación) y Modo 2 (ejecución). Portable y acotado. |
| **Implementor** | Rol del Modo 2 que escribe código para pasar los tests. Activo en Fase Green. Pragmático, sin perfeccionismo prematuro. |
| **Interfaz universal** | Las 9 operaciones que todo adaptador del artifact store debe implementar: `ingest`, `save`, `read`, `search`, `list`, `delete`, `verifyConsistency`, `history`, `transition`. |
| **Lane** | Rama paralela en el DAG de tareas que puede ejecutarse concurrentemente con otras. |
| **MIM (Mind behind the Idea and the Money)** | El humano que dirige el proyecto. Decide, aprueba y desbloquea. Es el nodo de decisión final. |
| **Ops-runbook** | Artefacto (`ops-runbook.md`) que documenta operación, monitoreo, troubleshooting y escalación para proyectos con servicios vivos. Respaldado por ISO 20000 + ITIL 4. |
| **Orquestador de Ejecución** | Coordinador del Modo 2. Análogo al SM de planificación pero opera sobre el working tree del repositorio, no sobre el artifact store. Delega, no ejecuta. |
| **Pattern A** | Estrategia de retrieval donde el SM busca, cura y re-inyecta contexto en el prompt del sub-agente. Útil para búsquedas fuzzy o fan-out alto (8+). |
| **Pattern B** | Estrategia de retrieval donde el SM pasa solo topic_keys y el sub-agente consulta el artifact store directamente. 6x más barato que Pattern A. |
| **PDC (Post-Delegation Checkpoint)** | Protocolo de 4 pasos que el SM ejecuta después de cada retorno de sub-agente: ECHO (coherencia), VERIFY (cobertura), MARK (persistir), DECIDE (siguiente acción). |
| **Pre-Fase** | Primera etapa del Modo 2. Define contratos formales antes de escribir tests o código. Habilita desarrollo paralelo. |
| **QA (Modo 2)** | Rol del Modo 2 activo en Fase Accept. Verifica producto contra handoff. No escribe tests ni corrige código — certifica. Mecanismo de certificación definido por el consumidor del framework. |
| **RAG (Retrieval-Augmented Generation)** | El artifact store usado como fuente de contexto acotado para los agentes. Los agentes consultan slices específicos, no el codebase completo. |
| **Red/Green/Refactor** | Las 3 fases macro del Modo 2. Red = escribir tests (fallan). Green = escribir código (tests pasan). Refactor = review de calidad (tests siguen pasando). |
| **Reviewer** | Rol del Modo 2 activo en Fase Refactor. Tres variantes: Arquitectura (SOLID, DRY, KISS), Seguridad (OWASP), Performance (memory leaks, N+1). |
| **SM (Session Manager / Orquestador)** | El agente principal que actúa como facade del proyecto. Orquesta fases, convoca roles, valida gates y controla transiciones. No produce contenido. No es un Scrum Master en el sentido del Scrum Guide. |
| **Sub-agente** | Agente instanciado por el SM con un contrato de delegación acotado. Recibe personalidad, contexto y restricciones específicas para su tarea. Los roles del scrum team son sub-agentes. |
| **Test Contract** | Manifiesto enumerable por sujeto bajo prueba. Cada entrada vincula un caso de test con un nombre inmutable y trazable a un AC. Previene código de test spaghetti. |
| **Test Engineer** | Rol del Modo 2 que escribe la suite de tests mapeada a ACs y contratos. Activo en Fase Red. Escéptico, prioriza integración sobre unit. |
| **Test Plan** | Meta-documento que mapea ACs a casos de prueba, asigna boundaries (App/E2E), y etiqueta tests para filtrado (smoke, critical, regression). |
| **Testing de Alto Valor** | Filosofía del framework: solo tests que ejercen interacciones REALES de producto aportan valor. Tests con mocks extensivos dan falsa confianza. |
| **Tier (de activación)** | Nivel de ceremonia del framework (Ligero, Estándar, Completo). El SM determina el tier usando el score F1-F4. Los tiers escalan ceremonia, no artefactos. |
| **TPM (Technical Program Manager)** | Agente de infraestructura que actúa como DBMS del artifact store. Único actor que escribe en el store. Valida integridad, formato y completitud con criterio editorial. |
| **`transition()`** | Operación del adaptador que cambia el estado de un artefacto en la state machine (draft, review, approved, rejected, cancelled). |
| **Work item (L0-L4)** | Jerarquía de descomposición de trabajo: L0 (epic/proyecto), L1 (feature), L2 (story/tarea), L3 (sub-tarea), L4 (paso atómico). Definidos en `tasks.md`. |
| **Worktree** | Mecanismo de git que crea un directorio de trabajo aislado con su propia rama. El Orquestador de Ejecución los usa para ejecutar lanes paralelos sin conflictos de archivos. |
