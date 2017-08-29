# Sentinel Import Test Framework

This folder contains a library for testing imports that are written in Go.

This works by building the import binary dynamically during `go test` and
executing your test policy. The policy must pass. If the policy fails, the
failure trace is logged and shown. Execution is done via the publicly available
`sentinel` binary.

## Example

You can see an example in the `import_test.go` file in this folder. This
test actually runs as part of the unit tests to verify the behavior.
