---
id: artifact-system
title: "Sistema de Artifacts"
mode: framework
type: reference
tags: [artifacts, build, coverage, reportes, outputs, gitignore]
---

# Sistema de Artifacts

← [Índice](README.md)

> El sistema de artifacts es una convención que declara DÓNDE aparecen
> los outputs generados por el proyecto (builds, reportes de cobertura,
> documentación API, etc.) y QUÉ scripts los producen. Es el complemento
> del [sistema de ecos](echo-system.md): el eco define los pasos, el
> sistema de artifacts define dónde aterrizan los resultados.

---

## Distinción Semántica

En este framework, "artifact" tiene dos significados que no deben
confundirse:

| Tipo | Qué es | Dónde vive | Quién lo gestiona |
|------|--------|------------|-------------------|
| **Artifact de planificación** | Documentos del Modo 1 (idea.md, spec.md, design.md, tasks.md, handoff.md, ops-runbook.md) | [Artifact store](planning/artifacts/README.md) (fuera del repo) | [TPM](overview.md) |
| **Artifact de build** | Outputs generados por el eco (compilados, reportes, documentación API) | Carpeta gitignoreada dentro del repo | Sistema de artifacts (este documento) |

Este documento define el segundo tipo. Cuando el resto de la
documentación dice "artifact" sin calificador, se refiere al primer
tipo (planificación). Cuando se refiere a outputs de build, se usa
"artifact de build" o se referencia este documento.

---

## La Convención

El sistema de artifacts establece tres elementos:

1. **Una carpeta gitignoreada** como destino único y predecible para
   todos los outputs generados.
2. **Un catálogo de tipos** que enumera qué artifacts produce el
   proyecto.
3. **Un mapa de generación** que vincula cada tipo de artifact con el
   script o comando que lo produce.

```mermaid
flowchart TD
    subgraph ECHO["Sistema de Ecos (5 pasos)"]
        direction LR
        S["Setup"] --> B["Build"] --> ST["Static"] --> DT["Dynamic"] --> E["E2E"]
    end

    subgraph ARTIFACTS["Carpeta de Artifacts (gitignoreada)"]
        direction TB
        A_BUILD["Apps buildeadas\n(compilados, bundles)"]
        A_COV["Reportes de cobertura\n(HTML, JSON, LCOV)"]
        A_TEST["Reportes de tests\n(JUnit XML, JSON)"]
        A_API["Documentación API\n(OpenAPI, AsyncAPI)"]
        A_BUNDLE["Bundle analysis\n(tamaño, tree shaking)"]
    A_OTHER["Otros artifacts\n(PDFs, assets, etc.)"]
    end

    B -->|"produce"| A_BUILD
    B -->|"produce"| A_API
    B -->|"produce"| A_OTHER
    DT -->|"produce"| A_COV
    DT -->|"produce"| A_TEST
    E -->|"produce"| A_TEST
```

### Ubicación de la carpeta

La ubicación la define el proyecto en
[`design.md`](planning/artifacts/schemas.md) (sección "Restricciones de
infraestructura"). El nombre y la ruta son decisión del proyecto — el
framework solo exige que:

- Sea una carpeta dedicada y gitignoreada.
- Esté documentada en `design.md`.
- Todo output generado apunte a ella (o a un subdirectorio dentro de
  ella).

---

## Catálogo de Tipos

Los siguientes tipos de artifacts son comunes en la mayoría de
proyectos. El catálogo concreto de cada proyecto se define en
`design.md`.

| Tipo | Descripción | Producido por (paso del eco) | Consumido por |
|------|-------------|------------------------------|---------------|
| **Apps buildeadas** | Código compilado, transpilado, bundleado, listo para deploy o distribución | Paso 2 (Build) | Pipeline de deployment, Modo 3 |
| **Reportes de cobertura** | Cobertura de tests por archivo y agregada. Formatos: HTML para lectura, JSON/LCOV para procesamiento | Paso 4 (Dynamic Test) | [Fase Accept](execution/accept.md) (QA verifica umbral), SOC 2 (evidencia de controles) |
| **Reportes de tests** | Resultados de ejecución de tests. Formato estructurado para parsing automatizado | Paso 4 (Dynamic Test) + Paso 5 (E2E) | [Fase Accept](execution/accept.md) (QA verifica que todos pasan) |
| **Documentación API** | Especificaciones de API generadas desde contratos o código. OpenAPI, AsyncAPI, GraphQL schema | Paso 2 (Build) | Consumidores de la API, [Pre-Fase Contratos](execution/contracts.md) |
| **Bundle analysis** | Análisis de tamaño de bundles, tree shaking, dependencias incluidas | Paso 2 (Build) | [Fase Refactor](execution/refactor.md) (revisión de performance) |
| **Otros outputs** | Cualquier output generado que el proyecto necesite rastrear (PDFs, assets procesados, etc.) | Variable | Variable |

---

## Mapa de Generación

El mapa de generación es una tabla que el proyecto declara en
`design.md`, vinculando cada tipo de artifact con:

- El script o comando que lo genera.
- La ruta de destino dentro de la carpeta de artifacts.
- El paso del eco donde se produce.

```markdown
| Artifact | Script | Destino | Paso del eco |
|----------|--------|---------|--------------|
| App backend | {comando de build} | {carpeta}/backend | 2. Build |
| App frontend | {comando de build} | {carpeta}/frontend | 2. Build |
| Coverage | {comando de test} | {carpeta}/coverage | 4. Dynamic |
| Test report | {comando de test} | {carpeta}/reports | 4. Dynamic |
| OpenAPI spec | {comando de generación} | {carpeta}/api | 2. Build |
```

El mapa no es prescriptivo en su contenido (cada proyecto tiene
diferentes artifacts) — es prescriptivo en su existencia: todo proyecto
que adopte el sistema de artifacts DEBE tener este mapa documentado.

---

## Conexión con el Framework

### Lo que el framework ya exige y este sistema formaliza

El framework menciona reportes y outputs en múltiples lugares sin
definir dónde viven ni en qué formato se producen:

| Referencia existente | Documento | Qué formaliza el sistema de artifacts |
|---------------------|-----------|---------------------------------------|
| "Coverage report" como input de QA | [accept.md](execution/accept.md) | Ubicación y formato del reporte de cobertura |
| "Test reports" como evidencia | [accept.md](execution/accept.md) | Ubicación y formato de los reportes de tests |
| "Reporte de cobertura" como evidencia SOC 2 | [red.md](execution/red.md) | Dónde se genera y persiste el reporte |
| Herramienta de coverage "debe reportar por archivo y agregada" | [red.md](execution/red.md) | Dónde aterriza ese reporte |
| "Código droppable" identificado via coverage report | [accept.md](execution/accept.md) | El reporte vive en ubicación conocida |

### Ciclo de vida de un artifact de build

```mermaid
sequenceDiagram
    participant ECO as Eco (paso N)
    participant FS as Carpeta de Artifacts
    participant QA as Fase Accept
    participant CI as Pipeline CI/CD

    ECO->>FS: Genera artifact en ruta predecible
    Note over FS: Gitignoreado, regenerable

    QA->>FS: Lee reportes de cobertura y tests
    QA->>QA: Verifica umbrales y resultados

    CI->>FS: Lee artifacts de build
    CI->>CI: Despliega si eco es verde
```

Los artifacts de build son **efímeros y regenerables**. No se versionan
en git — se regeneran en cada ejecución del eco. La fuente de verdad es
el código fuente + el eco, no los artifacts generados.

---

## Obligatorio en Principio, Flexible en Implementación

El sistema de artifacts es **obligatorio** en su requisito central: todo
proyecto DEBE saber dónde están sus artifacts y qué los genera. El mapa
de generación y la documentación de ubicaciones son prescriptivos.

Lo que es **flexible** es la implementación de la carpeta centralizada
— la convención de una única carpeta gitignoreada como destino. Esta
convención es la recomendación por defecto, pero no siempre puede
aplicarse tal cual. Existen limitaciones conocidas:

### Limitaciones de stack

Algunos stacks hardcodean la ubicación de sus outputs y no permiten
redireccionarlos sin efectos secundarios:

- Frameworks con directorio de build fijo que rompen al moverlo.
- Herramientas que generan outputs en la raíz del proyecto sin opción
  de redireccionamiento.

### Limitaciones de proyecto

Algunos proyectos producen outputs que no encajan en el modelo de
carpeta gitignoreada:

- Repositorios que auto-generan archivos trackeados (READMEs dinámicos,
  assets procesados que deben committearse).
- Proyectos donde el output ES el repositorio (sites estáticos,
  generadores de documentación).

### Artifacts remotos

Algunos proyectos producen artifacts cuyo destino no es el filesystem
local sino un registro o almacén externo:

- Imágenes Docker/OCI pusheadas a un registry.
- Packages publicados a un registry de paquetes (npm, PyPI, crates.io).
- Binarios de apps móviles subidos a app stores.
- Helm charts pusheados a chart repositories.

Estos artifacts se documentan en el mapa de generación con su destino
remoto (ej: `{registry}:{tag}`) en vez de una ruta local. El principio
se mantiene: el proyecto SABE dónde está cada artifact y qué lo genera.

### Multi-target

Para proyectos con múltiples targets de compilación (cross-platform,
multi-arch, library + CDN bundle), cada target se documenta como una
fila separada en el mapa de generación. El paso del eco es el mismo;
el script y destino varían por target.

### Cómo manejar las limitaciones

Cuando el sistema de artifacts no puede aplicarse limpiamente:

1. **Documentar la excepción** en `design.md` — qué artifact no puede
   redireccionarse y por qué.
2. **Documentar la ubicación real** — dónde aparece el artifact
   efectivamente.
3. **Mantener el mapa de generación** — el mapa sigue existiendo, solo
   cambia la columna de destino para reflejar la ubicación real en vez
   de la carpeta centralizada.

La excepción no elimina el sistema — lo adapta. Lo que no puede variar
es que el proyecto SEPA dónde están sus artifacts y qué los genera.

---

## Dónde Se Configura

| Qué | Dónde | Cuándo |
|-----|-------|--------|
| Nombre y ruta de la carpeta de artifacts | `design.md` — Restricciones de infraestructura | Fase 3 (Diseñar) |
| Catálogo de tipos del proyecto | `design.md` — Restricciones de infraestructura | Fase 3 (Diseñar) |
| Mapa de generación (script → destino) | `design.md` — Restricciones de infraestructura | Fase 3 (Diseñar) |
| Excepciones documentadas | `design.md` — Restricciones de infraestructura | Fase 3 (Diseñar) |
| Entrada en `.gitignore` | Working tree del repo | Modo 2 (Pre-Fase o Green) |

---

## Índice de Documentos Relacionados

| Documento | Relación con este |
|-----------|-------------------|
| [Sistema de ecos](echo-system.md) | El eco produce los artifacts; este sistema define dónde aterrizan |
| [Fase Red](execution/red.md) | Define requisitos de cobertura y reportes que se convierten en artifacts |
| [Fase Accept](execution/accept.md) | QA consume los reportes de cobertura y tests como inputs de certificación |
| [Fase Refactor](execution/refactor.md) | Revisión de performance consume bundle analysis; revisión de seguridad consume audit reports |
| [Contratos](execution/contracts.md) | Contract-first puede generar documentación API como artifact |
| [Schemas](planning/artifacts/schemas.md) | `design.md` es donde se configura el sistema de artifacts |
| [Artefactos de planificación](planning/artifacts/README.md) | Distinción semántica: artifacts de planificación ≠ artifacts de build |
