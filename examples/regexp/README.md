# Regexp example

Shows how to use `WithFilterFunc` and `WithSearchFunc` to match rows against a
regular expression instead of the default substring matching, caching the
compiled pattern so it isn't rebuilt on every row on every keystroke.

Filtering a log table by `^(ERROR|WARN)$` isolates those levels exactly,
and searching for `\b[45]\d\d\b` jumps between rows mentioning a 4xx/5xx
status code — both patterns that plain substring or fuzzy matching can't
express.
