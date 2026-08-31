# contributing to torshell

## branch strategy

```
main (sacred)
 └── staging
      ├── feat/xxx
      ├── fix/xxx
      └── refactor/xxx
```

### `main` sacred, release only
- never push directly
- only receives merges from `staging`
- merges are done by **@ijuttt** when staging has accumulated enough updates

### `staging` integration branch
- direct push is allowed but pull requests are preferred
- all feature/fix branches should PR into `staging`
- PRs are reviewed and approved by **@ijuttt**
- version bumps can be pushed directly here

### feature/fix branches — your workspace
- create your own branch from `staging`:
  ```sh
  git checkout staging
  git pull origin staging
  git checkout -b feat/your-feature
  ```
- naming convention:
  - `feat/xxx` — new features
  - `fix/xxx` — bug fixes
  - `refactor/xxx` — code restructuring
  - `docs/xxx` — documentation
  - `test/xxx` — tests
- when done, open a pull request to `staging`
- **@ijuttt** reviews and approves

## commit messages

use conventional format:

```
feat: add something new
fix: resolve some bug
refactor: restructure something
docs: update documentation
test: add tests
chore: maintenance tasks
```

## atomic commits

one commit = one logical change. no monster commits.

### rules
- don't bundle unrelated changes into a single commit
- don't commit a whole feature in one giant blob
- do split work into small, self-contained commits that each make sense on their own
- do make each commit compile and pass tests (no broken intermediate states)

### Guideline

| Scope | One Commit Should Be |
|:---|:---|
| New type/struct | Define the type + its constructor |
| New function | The function + its direct helper (if small) |
| Bug fix | The fix + the test that catches it |
| Refactor | One structural change (rename, extract, move) |
| Docs | One file or one section update |

## code guidelines

- run `go vet` and `go fmt` before committing.
- keep packages focused — one responsibility per package.
- no global state — pass dependencies explicitly.
