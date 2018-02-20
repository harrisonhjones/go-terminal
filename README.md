# `go-terminal`: A terminal helper for Go

[![GoDoc][1]][2]
[![GoCard][3]][4]

[1]: https://godoc.org/harrisonhjones.com/go-terminal?status.svg
[2]: https://godoc.org/harrisonhjones.com/go-terminal
[3]: https://goreportcard.com/badge/harrisonhjones.com/go-terminal?status.svg
[4]: https://goreportcard.com/report/harrisonhjones.com/go-terminal?status.svg

A simple terminal helper for Go. Check out the [GoDoc](https://godoc.org/harrisonhjones.com/go-terminal).

## Examples

### Pausing

You can use `terminal.Pause()` to pause operation and wait for the operator to press the "enter" key.

```go
terminal.Pause() // Outputs: "Press enter to continue" and waits.
terminal.Pause("Press enter") // Outputs: "Press enter" and waits.
```

### Collecting Input

You can use `terminal.Input()` to collect user input.

```go
val, err := terminal.Input("foo") // Outputs: "foo (required): " & waits. Returns with val set to the input.
val, err := terminal.Input("foo", terminal.Optional) // Outputs: "foo (optional): " & waits. Returns with val set to the input.
val, err := terminal.Input("foo", "bar") // Outputs: "foo (optional, default: bar): " & waits. Returns with val set to the input or "bar" if the input was blank.
```

You can also use `terminal.MustInput()` in situations where you don't want to check the `error` returned from `terminal.Input` and would rather just panic if there is an error.

```go
val := terminal.MustInput("foo") // Outputs: "foo (required): " & waits. Returns with val set to the input.
val := terminal.MustInput("foo", terminal.Optional) // Outputs: "foo (optional): " & waits. Returns with val set to the input.
val := terminal.MustInput("foo", "bar") // Outputs: "foo (optional, default: bar): " & waits. Returns with val set to the input or "bar" if the input was blank.
```

