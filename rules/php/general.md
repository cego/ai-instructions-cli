# General Engineering Guidelines

These principles apply to all projects, regardless of framework or programming language.  
They set the foundation for maintainable, predictable, and secure software.

---

## Code Quality

- Keep functions short and focused on a single responsibility.
- Name variables, classes, and methods with intention and clarity.
- Avoid premature optimization; prefer readable and predictable code.
- Remove dead code and unused dependencies when possible.

---

## Security

- Never trust input. Validate and sanitize all external data.
- Store secrets in environment variables or secret managers, not in source control.
- Use HTTPS by default and enforce secure cookie settings in production.
- Keep dependencies updated to avoid vulnerabilities.

---

## Architecture

- Separate domain logic from infrastructure (e.g., framework controllers vs. application services).
- Prefer composition over inheritance.
- Avoid global state unless strictly necessary; inject what you depend on.
- Keep configuration declarative and environment-specific.

---

## Testing

- Write tests for logic that can break or is business-critical.
- Use unit tests for pure logic and integration tests for framework/IO boundaries.
- Mock external services (HTTP, databases, queues) where appropriate.
- Avoid testing implementation details; focus on behavior.

---

## Tooling & Workflow

- Use code formatters and linters to enforce consistent style.
- Automate repetitive or error-prone tasks (CI, code generation, deployments).
- Use version control with meaningful commit messages.
- Keep documentation close to the code and ensure it reflects reality.

---

## Logging & Observability

- Log events, not noise: errors, warnings, state changes, and important actions.
- Use structured logging where possible for easier analysis.
- Emit metrics that reflect system health and performance.
- Add tracing when working with distributed systems or microservices.

---

Following these principles ensures that your project stays maintainable as it grows,
no matter which languages, frameworks, or tools are in play.
