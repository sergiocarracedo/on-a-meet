# Quick Task 017: YAML on/off command quoting

## Task

**DESCRIPTION:** The onboard-generated config.yml is missing on-command and off-command fields. Fix the YAML output to always include both fields with proper double-quoting and escaped inner double quotes.

### Must haves (truths)
- `on-command` and `off-command` must always appear in the generated YAML (even if empty)
- Values must be double-quoted with `\"` escaping for embedded double quotes
- The config.yaml.example already shows the expected format: `on-command: "..."`

### Files
- `cmd/onboard.go`

### Action

1. Add a `yamlQuotedString` type implementing `yaml.Marshaler` that forces `yaml.DoubleQuotedStyle`
2. Change `writeConfig.OnCmd` and `writeConfig.OffCmd` from `string` to `yamlQuotedString`
3. Remove `omitempty` from their YAML tags
4. Convert `cfg.OnCmd`/`cfg.OffCmd` with `yamlQuotedString()` in the apply path

### Verify

```bash
go build ./... && go vet ./... && go test ./...
```

### Done

When `yaml.Marshal(writeConfig{...})` produces:
```yaml
on-command: "value"
off-command: "value"
```
with double quotes and escaped inner quotes, even for empty values.
