# alfred-timein Project Memory

## Project Overview
A Go-based Alfred workflow for timezone lookups and time conversion. Uses Clean Architecture principles with comprehensive testing.

## Key Guidelines

### Code Standards
- **Architecture**: Follow Clean Architecture (Domain, Use Cases, Adapters, Frameworks)
- **Testing**: All changes must include BDD scenarios and unit tests
- **Performance**: Maintain cache optimization (6ms cache hits, 400ms+ cache misses acceptable)
- **Dependencies**: Prefer reliability over binary size (tzf library is acceptable at 15MB for offline accuracy)

### Commit Messages
- **Format**: Use Conventional Commits (`feat:`, `fix:`, `chore:`, `ci:`, `refactor:`, `docs:`)
- **Automation**: Commit messages auto-generate changelogs for releases
- **Examples**: `feat: add timezone caching`, `fix: resolve geocoding error`

### Release Process
- **Development**: Work on `dev` branch, auto-creates beta releases
- **Stable**: Tag releases on `main` branch for stable versions
- **Naming**: Project is "alfred-timein" (never "Alfred TimeIn")
- **Versioning**: Semantic versioning for stable, `v0.0.0-beta.X+sha` for dev

### Build and CI
- **Commands**: `make test-all` (runs unit + BDD), `make build`, `make alfredworkflow`
- **BDD Tags**: Use `go test -tags=bdd` to exclude test deps from production builds
- **Artifacts**: CI uploads binaries and workflows for manual testing

### Code Quality
- **No Emojis**: Avoid emojis and LLM-style language in release notes or documentation
- **Professional Tone**: Follow standard open source project conventions
- **Binary Size**: Optimized images (249KB icon), but prioritize functionality over size

### Directory Structure
```
cmd/           # CLI entry points (geotz, timein)
internal/      # Clean Architecture layers
  domain/      # Business entities and rules
  usecases/    # Application logic
  adapters/    # External interfaces (cache, geocoder, timezone, presenter)
features/      # BDD test scenarios
workflow/      # Alfred workflow assets
```

### Testing Philosophy
- **BDD**: Living documentation of business requirements
- **Coverage**: Unit tests for logic, BDD for end-to-end workflows
- **Clean Output**: No undefined scenarios in BDD to avoid pollution

### Performance Requirements
- **Cache Hits**: Under 10ms response time
- **Cache Misses**: Network-dependent, optimize where possible
- **Offline First**: Timezone data should work without internet (tzf library)

## Current State
- Clean Architecture implemented with performance optimizations
- Comprehensive BDD test suite (12 scenarios, 64 steps)
- Dev branch workflow with automated beta releases
- Optimized assets and build processes