CHANGELOG
=========

- Defines 6 assertions: Equal, Nil, NotNil, True, False, Error
- Adds go docs [2026-04-22 Wed]

IDEAS 💡 & TODOS ☑️
===================

- [ ] Introduce versioning
- [ ] Investigate using https://pkg.go.dev/github.com/google/go-cmp/cmp 
- [ ] Learn about Go module versioning
  https://go.dev/doc/modules/version-numbers
  https://go.dev/doc/modules/release-workflow

FIELDNOTES 📝
=============

## Wed, Apr.29th, 2026

- [ ] Make attempt at v2 offering an API without testing.T as first parameter

### Guide to publishing v2 of a Go module

Here's a complete guide to organizing your Go code so it can be published as `github.com/al3xandru/assert/v2`.

#### 1. Repository & Directory Structure

There are two strategies for a v2 module. The **Major Subdirectory** strategy is recommended for maximum compatibility:

```
assert/                        ← GitHub repo root
├── README.md
├── LICENSE
│
├── go.mod                     ← v1 module (github.com/al3xandru/assert)
├── assert.go
├── assert_test.go
│
└── v2/                        ← v2 lives in a subdirectory
    ├── go.mod                 ← v2 module (github.com/al3xandru/assert/v2)
    ├── assert.go
    └── assert_test.go
```

Alternatively, you can use the **Major Branch** strategy (simpler, no subdirectory), 
where v2 lives on a separate git branch 
and the repo root holds the v2 code with a `/v2`-suffixed module path. 
Both are valid — subdirectory is safer for older Go compatibility.

#### 2. The `go.mod` File

Inside `v2/go.mod` (or at repo root if using the branch strategy):

```
module github.com/al3xandru/assert/v2

go 1.21
```

A module released at major version 2 or higher is required to have a major version suffix on its module path. 
For module versions v2 and later, the module path value must end with the major version number, such as `/v2`.


#### 3. Update All Internal Import Paths

If your v2+ module has multiple packages, you should update them to use the new `/v2` import path; 
otherwise your v2+ module will depend on your v0/v1 module.

Every `.go` file inside `v2/` must use the new path:

```go
// v2/assert.go
package assert

import (
    "github.com/al3xandru/assert/v2/helpers" // internal subpackage — note /v2
)
```

#### 4. Package Declaration

The package name itself does **not** change — only the module path does:

```go
// v2/assert.go
package assert   // ← still "assert", not "assertv2"
```

#### 5. Tagging a Release

If the module is defined at the root of the repository, the prefix is empty and the tag is just the version. However, in multi-module repositories, the prefix distinguishes versions for different modules — the prefix is the directory within the repository where the module is defined.

So for the **subdirectory strategy**, tags must be prefixed with the subdirectory name:

```bash
# For v2/ subdirectory strategy:
git tag v2/v2.0.0
git push origin v2/v2.0.0

# For branch strategy (v2 code is at repo root on a v2 branch):
git tag v2.0.0
git push origin v2.0.0
```

---

## 6. Publishing to pkg.go.dev

After pushing the tag, trigger indexing:

```bash
GOPROXY=proxy.golang.org go list -m github.com/al3xandru/assert/v2@v2.0.0
```

Once you've published the module, developers importing its packages will be able to resolve a dependency by running commands such as `go get`. Note: don't change a tagged version after publishing — Go tools authenticate a downloaded module against the first downloaded copy. Instead of changing previously published code, publish a new version.

---

## 7. How Users Import It

```go
import "github.com/al3xandru/assert/v2"
```

```bash
go get github.com/al3xandru/assert/v2@latest
```

---

## Best Practices Checklist

Before publishing, follow these best practices: include a **LICENSE** file with minimal restrictions for easy use and redistribution; add **package-level doc comments** (Go uses these as documentation on pkg.go.dev); use **semver tags** — stable releases should start at `v2.0.0` or above, which gives developers confidence.

---

## References

- **Go Modules Reference (v2+ paths)**: https://go.dev/ref/mod
- **Go Blog: v2 and Beyond**: https://go.dev/blog/v2-go-modules
- **Organizing a Go module**: https://go.dev/doc/modules/layout
- **Publishing a module**: https://go.dev/doc/modules/publishing
- **go.mod file reference**: https://go.dev/doc/modules/gomod-ref

## Wed, Apr.22nd, 2026

Want to move the utility code to capture console output to this project.

- [ ] Learn about Go module versioning
  https://go.dev/doc/modules/version-numbers
  https://go.dev/doc/modules/release-workflow

## Mon, Apr.20th, 2026

Introducing the library.

