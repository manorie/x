# Project X

Two different protocols.

- http / https

Methods

- GET / POST / PUT / DELETE

Intervals

- 30s to 24h customizable completely

HTTP check result

- status code
- namelookup time
- connect time
- content size
- content transfer start
- content transfer end

Persistance

- InfluxDB ? Time series data

Common Objects

- Status Checker
- Request Result

Status Checker

- UUID
- Timeout
- End point
- Method
- Allow redirects
- Interval (Seconds)
- Active?

BE

- check.go
  - new
- checker.go
  - new
  - delete
  - start
  - stop
  - checkers (populate on init)
  - sych checkers with DB?
- client.go
  - run(for checker)
- store.go
  - sync checkers?
  - write(for checker)
  - find_since(for checker, time)
  - remove(checker uuid)
- api.go
  - add_checker
  - stop_checker
  - run_checker
  - delete_checker
  - update_checker
  - checks_since (for checker, time)
- logger.go
- errors.go
- configuration.toml
- notifier.go
  - ?
- main.go
  - start configuration
  - start checkers on run
  - handle api requests