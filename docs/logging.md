# Logging

The 3scale SaaS operator uses the [go.uber.org/zap](https://pkg.go.dev/go.uber.org/zap) library for logs.

## Configuration

The operator logger supports the following environment variables:

| Variable      | Format | Default    | Values                                                 |
| ------------- | ------ | ---------- | ------------------------------------------------------ |
| LOG_MODE      | string | production | `production` / `development`                           |
| LOG_ENCODING  | string | -          | `json` / `console`                                     |
| LOG_LEVEL     | string | -          | `debug`,`info`,`warn`,`error`,`dpanic`,`panic`,`fatal` |
| LOG_VERBOSITY | int8   | 0          | `0-10`                                                 |

### LOG_MODE

Log mode defaults to `production` and that configures the logger with:
- uses a JSON encoder
- writes to standard error
- enables sampling
- Stacktraces are automatically included on logs of ErrorLevel and above.

When set to `development`, it enables development mode:
- uses a console encoder
- writes to standard error
- disables sampling
- makes DPanicLevel logs panic
- Stacktraces are automatically included on logs of WarnLevel and above.

### LOG_ENCODING

If not set, will be configured by `LOG_MODE` profile: `json` for `production` and `console` for `development`.

Can be overrided by setting up the `LOG_ENCODING` variable.

### LOG_LEVEL

Defaults to `debug` in `development` mode and `info` for production.

Can be overrided by setting up the `LOG_LEVEL` variable.

More information in [zapcore#Level](https://pkg.go.dev/go.uber.org/zap@v1.21.0/zapcore#Level).

### LOG_VERBOSITY

If `LOG_LEVEL` is set to `debug`, via `LOG_LEVEL` or with `LOG_MODDE` set to `development`, allows increasing the log verbosity.

Allows values from 1 to 10.