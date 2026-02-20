# goreleaser + GitHub Actions Release Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Set up goreleaser and GitHub Actions to automatically publish two platform-specific `.alfredworkflow` release artifacts on every `v*` git tag push.

**Architecture:** goreleaser builds the Go binary for `darwin/amd64` and `darwin/arm64` separately, packages each into a zip named `metacpan-{version}-darwin-{arch}.alfredworkflow` with `icon.png` and `info.plist`, then creates a GitHub release. A GitHub Actions workflow triggers this on tag push. The version is injected at build time via ldflags so `main.go` doesn't need to be updated on each release.

**Tech Stack:** Go, goreleaser, GitHub Actions

---

### Task 1: Wire version via ldflags in main.go

**Files:**
- Modify: `cmd/workflow/main.go`

**Step 1: Replace hardcoded version with a variable**

In `cmd/workflow/main.go`, replace:
```go
app.Version = "0.9.0"
```
with:
```go
app.Version = version
```

And add a package-level variable before `main()`:
```go
var version = "dev"
```

**Step 2: Build and verify it compiles**

```bash
go build ./...
```
Expected: no output, exits 0.

**Step 3: Verify version flag works**

```bash
go run cmd/workflow/main.go --version
```
Expected: output contains `dev` (the default).

**Step 4: Commit**

```bash
git add cmd/workflow/main.go
git commit -m "Use ldflags-injectable version variable in main.go"
```

---

### Task 2: Create .goreleaser.yaml

**Files:**
- Create: `.goreleaser.yaml`

**Step 1: Create the file**

```yaml
version: 2

project_name: metacpan

builds:
  - id: workflow
    main: ./cmd/workflow
    binary: workflow
    ldflags:
      - -s -w -X main.version={{.Version}}
    goos:
      - darwin
    goarch:
      - amd64
      - arm64

archives:
  - id: alfredworkflow
    format: zip
    name_template: "metacpan-{{ .Version }}-{{ .Os }}-{{ .Arch }}"
    files:
      - icon.png
      - info.plist

release:
  github:
    owner: oalders
    name: alfred-metacpan-workflow

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^chore:"
```

**Step 2: Validate the config**

```bash
goreleaser check
```
Expected: `config is valid` (install goreleaser first if needed: `brew install goreleaser`).

**Step 3: Do a dry-run snapshot build**

```bash
goreleaser release --snapshot --clean
```
Expected: `dist/` directory created with two `.zip` files named `metacpan-*-darwin-amd64.zip` and `metacpan-*-darwin-arm64.zip`.

**Step 4: Verify archive contents**

```bash
unzip -l dist/metacpan-*darwin-amd64*.zip
```
Expected: listing shows `workflow`, `icon.png`, `info.plist`.

**Step 5: Commit**

```bash
git add .goreleaser.yaml
git commit -m "Add goreleaser config for darwin amd64/arm64 alfredworkflow archives"
```

---

### Task 3: Rename archive extension to .alfredworkflow

**Files:**
- Modify: `.goreleaser.yaml`

The goreleaser `format` field supports `zip` but not custom extensions. We need the file to end in `.alfredworkflow`. Use a goreleaser `after` hook to rename the zips.

**Step 1: Add after hooks to .goreleaser.yaml**

Replace the `archives` section with:

```yaml
archives:
  - id: alfredworkflow
    format: zip
    name_template: "metacpan-{{ .Version }}-{{ .Os }}-{{ .Arch }}"
    files:
      - icon.png
      - info.plist

after:
  hooks:
    - cmd: bash -c 'for f in dist/metacpan-*.zip; do mv "$f" "${f%.zip}.alfredworkflow"; done'
```

**Step 2: Run snapshot build again**

```bash
goreleaser release --snapshot --clean
```
Expected: `dist/` contains `metacpan-*-darwin-amd64.alfredworkflow` and `metacpan-*-darwin-arm64.alfredworkflow`.

**Step 3: Verify archive contents still correct**

```bash
unzip -l dist/metacpan-*darwin-amd64*.alfredworkflow
```
Expected: `workflow`, `icon.png`, `info.plist`.

**Step 4: Commit**

```bash
git add .goreleaser.yaml
git commit -m "Rename release archives to .alfredworkflow extension"
```

---

### Task 4: Create GitHub Actions release workflow

**Files:**
- Create: `.github/workflows/release.yml`

**Step 1: Create the directory**

```bash
mkdir -p .github/workflows
```

**Step 2: Create the workflow file**

```yaml
name: Release

on:
  push:
    tags:
      - "v*"

permissions:
  contents: write

jobs:
  release:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - uses: goreleaser/goreleaser-action@v6
        with:
          version: "~> v2"
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

Note: `runs-on: macos-latest` is required since we're building darwin binaries. `fetch-depth: 0` is required by goreleaser to generate the changelog.

**Step 3: Commit**

```bash
git add .github/workflows/release.yml
git commit -m "Add GitHub Actions workflow to publish releases on tag push"
```

---

### Task 5: Test the full release flow

**Step 1: Push the branch and verify Actions config is valid**

```bash
git push origin modernize
```

Check GitHub Actions tab — no workflow should run yet (no tag pushed).

**Step 2: Tag and push**

```bash
git tag v0.9.0
git push origin v0.9.0
```

**Step 3: Verify the release**

Go to `https://github.com/oalders/alfred-metacpan-workflow/releases` and confirm:
- A release named `v0.9.0` exists
- Two assets are attached: `metacpan-0.9.0-darwin-amd64.alfredworkflow` and `metacpan-0.9.0-darwin-arm64.alfredworkflow`
- Each archive unzips to contain `workflow`, `icon.png`, `info.plist`
- Changelog is populated from commits
