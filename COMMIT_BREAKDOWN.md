# Commit Breakdown (Refactor / Tests / Docs)

This splits the current work into clean, reviewable commits.

> Note: `README.md` is currently modified in the working tree but is **not** part of this change set.

---

## 1) refactor(app): dependency injection + phase gating

### Files
- `go.mod`
- `internal/app/app.go`
- `internal/app/commands.go`
- `internal/app/handlers.go`
- `internal/app/input.go`
- `internal/app/deps.go`
- `internal/app/state.go`
- `internal/pick/pick_darwin.go`
- `internal/pick/script.go`

### Suggested commands
```bash
git add go.mod \
  internal/app/app.go \
  internal/app/commands.go \
  internal/app/handlers.go \
  internal/app/input.go \
  internal/app/deps.go \
  internal/app/state.go \
  internal/pick/pick_darwin.go \
  internal/pick/script.go

git commit -m "refactor(app): add dependency boundaries, phase gating, and copy readiness checks"
```

---

## 2) test(app,pick): raise app coverage and add pick sanitization tests

### Files
- `internal/app/commands_test.go`
- `internal/app/input_test.go`
- `internal/app/state_test.go`
- `internal/app/display_test.go`
- `internal/app/handlers_test.go`
- `internal/pick/script_test.go`

### Suggested commands
```bash
git add internal/app/commands_test.go \
  internal/app/input_test.go \
  internal/app/state_test.go \
  internal/app/display_test.go \
  internal/app/handlers_test.go \
  internal/pick/script_test.go

git commit -m "test(app,pick): add phase/copy gating tests and improve app coverage"
```

---

## 3) docs: review + commit planning

### Files
- `CODE_REVIEW.md`
- `COMMIT_BREAKDOWN.md`

### Suggested commands
```bash
git add CODE_REVIEW.md COMMIT_BREAKDOWN.md

git commit -m "docs: update code review status and add commit breakdown"
```

### Note on `/agent`

`agent/` is currently ignored by `.gitignore`, so the updated
`agent/STATE-MACHINE.md` is a local/spec artifact unless you intentionally force-add it.

Optional force-add:
```bash
git add -f agent/STATE-MACHINE.md
git commit -m "docs(agent): update runtime phase state machine"
```

---

## Optional: keep README separate

If you want your current README edits committed independently:

```bash
git add README.md
git commit -m "docs(readme): refresh project intro and disclaimer"
```
