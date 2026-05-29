# Quick Task 004 Summary

**Task:** Fix --detect flag Viper bind clash between detect.go and list.go
**Completed:** 2026-05-29

## What was done
Fixed `--detect` flag binding cross-contamination. Both detect.go and list.go called `viper.BindPFlag("detect-method", ...)` in init(), but Go runs init() in filename order — list.go overwrote detect.go's binding. Viper then returned listCmd's flag default ("v4l2") instead of detectCmd's flag value ("lsof").

Fix: bypass Viper's BindPFlag for detect-method. Read flag value directly from cmd.Flags().Changed + cmd.Flags().Lookup("detect").Value.String() when the flag was explicitly set. Fall back to viper.GetString for config file support. Set default in root.go initConfig().

## Files changed
- `cmd/detect.go`: Changed detectConfig struct to accept onCmd/offCmd/timeout/camera. configFromViper() reads detectMethod from flag directly when changed, falls back to viper for config file.
- `cmd/list.go`: Same pattern — checks cmd.Flags().Lookup("detect").Changed.
- `cmd/root.go`: Added viper.SetDefault("detect-method", "v4l2") in initConfig().

## Commit
5e9e50d
