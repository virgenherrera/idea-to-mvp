# Pre-Fase: Definición de Contratos

← [Ejecución](README.md)

## Por que Contract-First

El contrato es la fuente de verdad compartida entre todos los actores de
ejecucion. Antes de escribir una linea de test o de implementacion, el
equipo define la interfaz publica del sistema. Esto habilita **desarrollo
paralelo agil**:

```mermaid
flowchart TD
    CONTRACT["Contrato definido\n(API, DB schema, interfaces)"]

    CONTRACT --> QA["Test Engineer\nescribe tests contra\nel contrato"]
    CONTRACT --> FE["Frontend\nconstruye contra\nel API contract"]
    CONTRACT --> BE["Backend\nimplementa detras\ndel contrato"]
    CONTRACT --> INFRA["Infra\nprepara migraciones\ny schemas"]

    QA --> MERGE["Merge:\ntodo converge\nen el contrato"]
    FE --> MERGE
    BE --> MERGE
    INFRA --> MERGE
```

**Principio: Contract over Methodology.** El contrato importa mas que el
proceso. Como implementes detras del contrato es tu problema --- pero
DEBES cumplirlo. El contrato es el acuerdo verificable; la metodologia
es la ceremonia alrededor.

---

## Tipos de contrato

| Tipo | Formato | Cuando aplica | Ejemplo |
|------|---------|----------------|---------|
| API Contract | OpenAPI 3.x / AsyncAPI | Proyecto con endpoints HTTP o eventos | `POST /auth/login` con request/response schema |
| SDK / Library Interface | TypeScript interfaces, Rust traits | Librerias, modulos reutilizables | `interface AuthService { login(credentials): Token }` |
| Database Schema | SQL DDL + migraciones | Proyecto con persistencia | `CREATE TABLE users (...)` con constraints |
| Event Schema | JSON Schema / AsyncAPI | Sistemas event-driven | `UserCreatedEvent { id, email, timestamp }` |
| Component Interface | Props/inputs tipados | Frontend con componentes | `LoginFormProps { onSubmit, initialValues }` |
| Connector / Adapter Interface | Ports & adapters | Integraciones con terceros | `interface PaymentGateway { charge(amount): Receipt }` |

---

## Flujo de definicion de contratos

```mermaid
sequenceDiagram
    participant OE as Orquestador de Ejecucion
    participant CA as Contract Architect
    participant MIM as MIM (Humano)

    OE->>OE: Lee handoff.md (spec, design, tasks)
    OE->>CA: Contrato: definir interfaces publicas<br/>basadas en design.md y spec.md
    activate CA
    CA->>CA: Extrae endpoints, schemas,<br/>interfaces de design.md
    CA->>CA: Mapea ACs de spec.md<br/>a contratos verificables
    CA-->>OE: Contratos definidos + Status Report
    deactivate CA

    OE->>OE: PDC: valida coherencia<br/>contratos vs handoff

    alt Contrato ambiguo o incompleto
        OE->>MIM: "El contrato de X tiene un gap:<br/>¿Y o Z?"
        MIM-->>OE: Decision
    end

    OE->>OE: Contratos aprobados → Fase Red
```

---

## Criterios de validacion del contrato

Un contrato esta listo cuando:

1. Cada AC de `spec.md` puede mapearse a al menos un contrato
2. Cada contrato tiene tipos definidos (request, response, error)
3. Los contratos son consistentes entre si (sin contradicciones)
4. Las dependencias entre contratos estan explicitas
5. El MIM aprobo los contratos que requieren decisiones de negocio
