# AGENTS.md

Guidance for AI agents (and new contributors) working on govips: Go bindings for [libvips](https://www.libvips.org/) built on cgo. Most bugs here are C-side bugs: reference counting, vector lengths, version drift between libvips releases. Read the conventions below before touching `vips/*.c`.

## Setup

Requires libvips headers installed:

```bash
# macOS
brew install vips
# Debian/Ubuntu (libheif-plugin-libde265 is needed for HEIC test fixtures)
sudo apt-get install libvips-dev libheif-plugin-libde265
```

## Commands

```bash
go build ./...        # build everything
make test             # go test -v -coverprofile=profile.cov ./...
make clean-cache      # purge go build/test caches; try this when cgo builds act weird
go generate ./vips    # regenerate generated.{go,c,h} via cmd/vipsgen
```

CI (`.github/workflows/build.yml`) runs `go build` and `go test -v ./...` on ubuntu-latest for every push and PR. CI is the authoritative gate; local macOS runs have known environmental failures (see below).

## Layout

```
vips/                 the entire library (single package)
  generated.{go,c,h}  machine-generated bindings from cmd/vipsgen. NEVER hand-edit;
                      change cmd/vipsgen and run go generate instead.
  operations.{go,c,h} hand-written bridge for ops needing custom logic
  foreign.{go,c,h}    load/save: LoadParams/SaveParams structs + per-format option setters
  stream.{go,c,h}     streaming I/O: VipsSourceCustom/VipsTargetCustom callback bridge
  image.go            ImageRef type, constructors, lifecycle
  image_*.go          Go API grouped by concern (transform, color, pixel, export,
                      metadata, composite, icc)
  image_export.go     ExportParams mapping; *ParamsFromExport helpers shared with streaming
cmd/vipsgen/          the code generator
resources/            test fixtures + golden files
examples/             runnable examples (separate go modules)
```

`ImageRef` is a mutable handle: operations replace the underlying `VipsImage` via `setImage` (which unrefs the old one). `Close` and a GC finalizer release the ref. `r.lock` guards only lifecycle/ownership transitions (`Close`, `SetKill`, `setImage`, streaming materialize/save); ordinary transform/export/metadata methods take no lock, just `runtime.KeepAlive(r)`, so an `ImageRef` is NOT safe for concurrent use. Never call methods (including `Close`) on the same `ImageRef` from multiple goroutines.

## cgo and memory rules

These are the rules that have caused real crashes when broken:

- **libvips operations take their own references on their inputs.** Never `g_object_unref` a `VipsImage` you received from a caller's `ImageRef`; the ref belongs to the ImageRef and is released in `Close`/finalize. A single wrong unref produces intermittent use-after-free SIGSEGVs at unrelated call sites, often much later (see #531; the crash wandered across `find_trim`, `heifsave`, `stats`).
- **Background/ink vectors must match the image's band count.** libvips rejects (or corrupts) longer vectors on greyscale images. Use `background_for_bands` in `operations.c` (1 band -> `{r}`, 2 -> `{r, a}`, 3 -> `{r, g, b}`, 4+ -> RGBA); don't hardcode 3/4-element arrays (see #534).
- **Version-guard new libvips properties.** Setting an unknown property fails at the C layer with an opaque error. Check `MajorVersion`/`MinorVersion`/`MicroVersion` in Go and return a clear error naming the required libvips version (pattern: `WebpExportParams.TargetSize`, requires 8.17.4+). Guards belong in shared helpers so every path (buffer and streaming) gets them.
- **Report the format from the file signature, not the loader.** Several loaders serve multiple formats (heifload: HEIF and AVIF; magickload: BMP, PSD, ICO). Use `DetermineImageType` on header bytes when the loader alone is ambiguous (see #540).
- **Errors**: C functions return non-zero and leave the message in the vips error buffer; Go reads it via `handleVipsError`/`handleImageError`. Free C strings you allocate (`C.CString` + `defer C.free`), and return cleanup funcs from param builders that allocate.
- The streaming callback registry (`stream.go`) uses integer handles across the cgo boundary, never Go pointers. Per-entry mutexes serialize callbacks; the global registry lock is never held during I/O.

## Testing

- **Golden files**: `goldenTest` (see `vips/image_golden_helpers_test.go`) compares output against `resources/<fixture>.<TestName minus "TestImage_">-<os>_<arch>_libvips-<version>.golden.<ext>` (e.g. `jpg-24bit.Crop-macos-15_arm64_libvips-8.18.0.golden.jpeg`). On a platform/libvips combination with no golden yet, the first run writes one; commit it. A mismatch writes a `.failed` file next to the golden for diffing. Assertion-only tests (via `GetPoint`, dimensions) avoid new binaries; prefer them when a golden adds no value.
- **Leak detection**: `Startup` installs a live-ImageRef counter. Use `OpenImageRefs()` deltas or `AssertNoLeaks(t)`; close every ImageRef a test creates.
- **Synthesized fixtures** beat committed binaries: e.g. a 2-band grey+alpha image is `ToColorSpace(InterpretationBW)` + `BandJoinConst([]float64{255})`.
- **Byte identity**: streaming savers (`SaveToWriter*`) must produce byte-identical output to their `Export*` counterparts; parity tests enforce this. If you touch save-param mapping, keep the shared `*ParamsFromExport` / `newSaveParams*` helpers as the single source of truth for both paths.
- For bug fixes, demonstrate the failure first: run the new test against the pre-fix code, then with the fix.
- Aim tests at behavior (outputs, pixels, errors), not call counts.

## Known local quirks (macOS)

Verified against clean master; never attribute these to a change under review:

- HEIC-input tests fail ("unsupported image format") when homebrew's vips-heif module can't dlopen (missing `libx265` dylib). `brew reinstall x265 libheif` usually fixes it.
- A full local `go test ./vips` can abort with SIGSEGV in `vips_object_print_all` (`TestImageRef_Close` -> `PrintObjectReport`), which silently skips all tests that would run after it. Use `-run` filters for the area you're touching, and trust CI for the full suite.

## Workflow

- All changes on branches/worktrees, never directly on master. Squash-merge PRs; reference issues (`Fixes #N`).
- Issues and PRs live on GitHub: https://github.com/davidbyttow/govips
- Public API is v2 (`github.com/davidbyttow/govips/v2/vips`); keep additions backward compatible. New export options go on the format-specific `*ExportParams` structs, not the deprecated generic `ExportParams`.
