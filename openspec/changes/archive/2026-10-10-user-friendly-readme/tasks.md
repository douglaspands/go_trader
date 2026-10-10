# Tasks

## 1. Rewrite README.md

- [x] 1.1 Draft the user-friendly header, motivation, target audience, and key features in `README.md` and verify all sections are present in Brazilian Portuguese.
- [x] 1.2 Write the quickstart installation and execution guide for Windows, Linux, and macOS in `README.md` and verify download links and terminal commands are accurate.
- [x] 1.3 Add practical CLI usage examples (quote retrieval, list, stock/REIT balance, combined security balance, CSV export) in `README.md` and verify command syntax and sample tables match CLI specifications.
- [x] 1.4 Add the FAQ section answering questions on data sources, execution safety, and internet connectivity in `README.md` and verify answers are clear and accurate.
- [x] 1.5 Write the dedicated "Para Desenvolvedores" section in `README.md` covering architecture, build commands, tests, exit codes, CSV formula sanitization, and AI agent guardrails, and verify all developer commands from the original README are preserved.

## 2. Review and Verification

- [x] 2.1 Verify markdown formatting and link integrity across `README.md` and confirm no broken anchors or missing sections remain.
- [x] 2.2 Verify repository integrity by running `gofmt -l .`, `go vet ./...`, and `go test -race ./...` to confirm documentation changes did not impact the build or offline test suite.
