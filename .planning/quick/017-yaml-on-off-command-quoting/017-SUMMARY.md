# Quick Task 017 Summary

**Task:** YAML on/off command quoting in onboard --apply
**Completed:** 2026-05-29

## What was done

The interactive onboard wizard's `--apply` path now always includes `on-command` and `off-command` in the generated `/etc/on-a-meet/config.yaml`, with values wrapped in double quotes and inner double quotes properly escaped as `\"`.

## Files changed

- `cmd/onboard.go`: Added `yamlQuotedString` type with `MarshalYAML` forcing `yaml.DoubleQuotedStyle`; changed `writeConfig.OnCmd`/`OffCmd` types from `string` to `yamlQuotedString`; removed `omitempty` from their yaml tags

## Commit

`19eb82d`
