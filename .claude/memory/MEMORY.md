# Memory Index

- [Aggregate Panel Implementation State](project_aggregate_panel.md) — feature/aggregate-panel branch: 10 commits done, 5 final fixes applied but uncommitted; next step is build+commit then finishing-a-development-branch
- Always use following build command:

```bash
GOROOT=/opt/homebrew/Cellar/go/1.26.1/libexec \
GOPROXY=https://proxy.golang.org,direct \
GONOSUMDB='*' \
GOOS=linux GOARCH=amd64 \
go build -o ./bin/pgweb_linux_amd64
```

>Output: `bin/pgweb_linux_amd64` (~28MB, ELF 64-bit, statically linked)
