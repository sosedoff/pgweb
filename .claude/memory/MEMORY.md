# Memory Index

- [Aggregate Panel Implementation State](project_aggregate_panel.md) — feature/aggregate-panel branch: 10 commits done, 5 final fixes applied but uncommitted; next step is build+commit then finishing-a-development-branch
- Always use following build command:

```bash
GOROOT=/opt/homebrew/Cellar/go/1.26.2/libexec \
GOPROXY=https://proxy.golang.org,direct \
GONOSUMDB='*' \
GOOS=linux GOARCH=amd64 \
go build -o ./bin/pgweb_linux_amd64
```

>Output: `bin/pgweb_linux_amd64` (~28MB, ELF 64-bit, statically linked)

- local branch stays local and not to be pushed to remote. Directive is set via .git/hooks/pre-push

- [Query Tab Button Bug](bug_query_tab_buttons.md) — Run Query button blocked by #output overlay after switching from Rows tab; fixed in showQueryPanel()

