# `flack.nix`

## Project

Go CLI (`github.com/YPares/flack`) that wraps `nix profile` / `nix flake` for stack-based profile management. Uses cobra for CLI, charmbracelet (huh/bubbletea) for interactive TUI.

## Commands

**Build and run ONLY via Nix:**
- `nix build .#flack` — builds the binary (vendorHash pinned in `flake.nix`; update it if `go.mod` changes)
- `nix run .#flack` — runs without installing

## Structure

- `cmd/flack/main.go` — single entrypoint, delegates to `cli.Execute()`
- `internal/cli/cli.go` — all cobra commands: show (default), push, pop, move, upgrade, update, inputs
- `internal/nix/nix.go` — spawns `nix` CLI; JSON parsing for profile list and flake metadata
- `internal/profile/profile.go` — profile stack logic (List, Push, Pop, Move, Upgrade)
- `internal/flake/flake.go` — flake input inspection and update
- `internal/tui/tui.go` — TUI select widgets (huh)
- `internal/display/table.go` — terminal table rendering (tablewriter)

## Gotchas

- No tests, no CI, no linter config — verify manually with `nix build` and `nix run`
- `moveCmd.Flags().SetInterspersed(false)` — flag parsing is intentionally disabled for the move subcommand so negative numbers as positional args work correctly
