# Virgil

Metodología e2e para desarrollo de software asistido por IA.

Alternativa a un AGENTS.md monolítico: en lugar de un archivo de 1000+
líneas que prescribe cómo opera un agente, Virgil define qué hacer —
conceptos, definiciones y metodología — en documentos modulares bajo `docs/`.

## Scope actual

Solo definiciones, conceptos y metodología. No hay implementación.
El CLI que consume esta metodología se construirá cuando el dogma esté
estabilizado.

## Documentación

`docs/` contiene la metodología completa:

- `docs/overview.md` — vista general, dogma rector, adopción progresiva
- `docs/execution/` — ciclo Red/Green/Refactor/Accept
- `docs/planning/` — fases, artefactos, comportamiento del SM
- `docs/operation/` — facade opcional post-ejecución
- `docs/glossary.md` — terminología
- `docs/agile-adaptations.md` — trade-offs explícitos vs Agile clásico
- `docs/echo-system.md` — sistema de hooks y verificación
