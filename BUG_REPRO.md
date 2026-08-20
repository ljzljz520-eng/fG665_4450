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
--- FAIL: TestBusiness06Regression (0.01s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
[signal SIGSEGV: segmentation violation code=0x1 addr=0x30 pc=0x164b58]

goroutine 33 [running]:
testing.tRunner.func1.2({0x18efe0, 0x328d40})
	/usr/local/go/src/testing/testing.go:1872 +0x19c
testing.tRunner.func1()
	/usr/local/go/src/testing/testing.go:1875 +0x31c
panic({0x18efe0?, 0x328d40?})
	/usr/local/go/src/runtime/panic.go:783 +0x120
instrumentarchive/service.UnsafeAssembleDetail(...)
	/app/service/detail_bug.go:17
instrumentarchive/service.(*Service).Detail(0x4000091e40, {0x1bf816, 0x2}, {0x1bfbdf, 0x5})
	/app/service/service.go:141 +0x1d8
instrumentarchive.TestBusiness06Regression(0x4000150000)
	/app/integration_test.go:34 +0x230
testing.tRunner(0x4000150000, 0x1cf4c0)
	/usr/local/go/src/testing/testing.go:1934 +0xc8
created by testing.(*T).Run in goroutine 1
	/usr/local/go/src/testing/testing.go:1997 +0x364
FAIL	instrumentarchive	0.029s
ok  	instrumentarchive/api	0.008s
ok  	instrumentarchive/auth	0.001s
ok  	instrumentarchive/cmd/server	0.002s
ok  	instrumentarchive/model	0.002s
ok  	instrumentarchive/report	0.001s
--- FAIL: TestDetailMissingCalibration (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
[signal SIGSEGV: segmentation violation code=0x1 addr=0x30 pc=0x162ba8]

goroutine 8 [running]:
testing.tRunner.func1.2({0x18ede0, 0x328d40})
	/usr/local/go/src/testing/testing.go:1872 +0x19c
testing.tRunner.func1()
	/usr/local/go/src/testing/testing.go:1875 +0x31c
panic({0x18ede0?, 0x328d40?})
	/usr/local/go/src/runtime/panic.go:783 +0x120
instrumentarchive/service.UnsafeAssembleDetail(...)
	/app/service/detail_bug.go:17
instrumentarchive/service.(*Service).Detail(0x40000cfe20, {0x1bee57, 0x2}, {0x1bf236, 0x5})
	/app/service/service.go:141 +0x1d8
instrumentarchive/service.TestDetailMissingCalibration(0x4000003c00)
	/app/service/service_test.go:51 +0x230
testing.tRunner(0x4000003c00, 0x1ce948)
	/usr/local/go/src/testing/testing.go:1934 +0xc8
created by testing.(*T).Run in goroutine 1
	/usr/local/go/src/testing/testing.go:1997 +0x364
FAIL	instrumentarchive/service	0.018s
ok  	instrumentarchive/store	0.010s
ok  	instrumentarchive/workflow	0.015s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/server): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/server): exit `0`
