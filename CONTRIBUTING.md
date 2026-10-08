# Contributing to alfred-timein

Thank you for your interest in contributing to alfred-timein!

## Commit Message Format

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automated changelog generation.

### Format

```
<type>: <description>
```

### Types

- **feat**: A new feature
- **fix**: A bug fix
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **chore**: Changes to build process, auxiliary tools, or maintenance
- **ci**: Changes to CI/CD configuration
- **docs**: Documentation changes
- **test**: Adding or updating tests

### Examples

```
feat: add timezone lookup caching for improved performance
fix: resolve crash when geocoding returns no results
refactor: extract timezone finder interface for better testability
chore: update dependencies to latest versions
ci: add automated testing for pull requests
docs: update README with installation instructions
test: add unit tests for cache functionality
```

## Development Workflow

1. **Create feature branch** from `dev`
2. **Make changes** with conventional commit messages
3. **Test locally** with `make lint test live`
4. **Submit PR** to `dev` branch
5. **Automated testing** runs on all PRs
6. **Beta releases** are created automatically from `dev`
7. **Stable releases** are created when `dev` is merged to `main`

## Testing

```bash
# ES2020 syntax ceiling (older macOS JavaScriptCore); fetches acorn via npx
make lint

# Offline: unit tests (test/unit.js) + golden tests (test/run.sh), no network
make test

# Golden tests against real OpenStreetMap/Apple services
make live

# Rewrite goldens after an intended output change
UPDATE=1 make test live

# Package the Alfred workflow
make alfredworkflow
```

CI runs `make test` on the oldest supported GitHub-hosted macOS too, before testing and releasing on the latest.

Everything lives in `timein.js`; run it directly with `./timein.js <place>`.

## Release Process

- **Beta releases**: Automatic on `dev` branch pushes
- **Stable releases**: Manual tag creation on `main` branch
- **Changelogs**: Generated from commit messages; betas list everything since the last stable release