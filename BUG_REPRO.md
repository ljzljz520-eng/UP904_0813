# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	slowpreview/cmd/slowpreview	[no test files]
?   	slowpreview/ingest	[no test files]
?   	slowpreview/report	[no test files]
--- FAIL: TestPreviewUsesSelectedSpeed (0.17s)
    preview_regression_test.go:31: selected half speed missing from command: ffmpeg -hide_banner -y -i training.mp4 -ss 1.000 -t 4.000 -vf crop=start=1000:end=5000,scale=1920:1080 -an previews/preview-regression.mp4
FAIL
FAIL	slowpreview	0.171s
ok  	slowpreview/analysis	0.004s
ok  	slowpreview/catalog	0.001s
ok  	slowpreview/domain	0.006s
ok  	slowpreview/engine	0.005s
ok  	slowpreview/playback	0.007s
ok  	slowpreview/service	0.172s
ok  	slowpreview/store	0.152s
ok  	slowpreview/timeline	0.001s
ok  	slowpreview/validate	0.001s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/slowpreview): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/slowpreview): exit `0`
