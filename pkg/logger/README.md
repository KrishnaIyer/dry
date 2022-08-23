# Logger

This is a highly opinionated wrapper around the [`zap`](https://github.com/uber-go/zap) logger.

## Usage

Create a new logger near the entrypoint of your code.

```go
logger, err := logger.New(ctx, false) // Without debug
```

Make sure to call the `Clean` function when you're done with the usage.

```go
defer logger.Clean()
```

You can now add this logger to a context.
```go
ctx = NewContextWithLogger(ctx, logger)
```

If you pass this context correctly to the call site, you can fetch the logger from the context and use its functions.

```go
logger = NewLoggerFromContext(ctx)
logger.Info("This is an info message")
```

One of the primary reasons for creating this package is to easily append fields to a logger for additional debugging.

You can set global fields to the logger.

```go
logger = logger.WithField("name", name) // All subsequent calls to this logger will log the name field.
```

You can also set local fields only to the log message.

```go
logger.WithField("name", name).Info("Login user")
```

The logger also supports multiple fields.
```go
	logger.WithFields(Fields(
		"test-key", "test-value",
		"test-other-key", 1,
	)).Info("This is an info message with multiple field")
```
