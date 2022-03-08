# Contributing

Thanks for your interest in contributing!

## Contributing issues

Please feel free to open an issue if you think something is working incorrectly,
if you have a feature request, or if you just have any questions.  No templates
are in place, all I ask is that you provide relevant information if you believe
something is working incorrectly so we can sort it out quickly.

## Contributing code

All contributions should have an associated issue.  If you are at all unsure
about how to solve the issue, please ask!  I'd rather chat about possible
solutions than have someone spend hours on a PR that requires a lot of major
changes.

Test coverage is important.  If you end up with a small (<1%) drop, I'm happy to
help cover the gap, but generally any new features or changes should have some
tests as well.  If you're not sure how to test something, feel free to ask!

Linting can be done with `make lint`.  Running `fmt` can be done with `make fmt`.
Tests can be run with `make test`.

Doing all these before submitting the PR can help you pass the existing
gatekeeping tests without having to wait for them to run every time.

The name of the PR, commit messages, and branch names aren't very important,
there are no special triggers or filters or anything in place that depend on names.
The name of the PR should at least reasonably describe the change, because this
is what people will see in the commit history and what goes into the change logs.

Exported functions should generally follow the pattern of returning a `Model`
and not use a pointer receiver.  This matches more closely with the flow of Bubble
Tea (and Elm), and discourages users from making mutable changes with unintended
side effects.  Unexported functions are free to use pointer receivers for both
simplicity and performance reasons.
