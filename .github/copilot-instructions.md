Follow clean coding principles. Avoid creating code blocks longer than 30
lines. Avoid high cyclomatic complexity.

Write idiomatic code. Use features from the latest version of Go. Stick to the
standard library where possible, unless otherwise instructed.

Use test driven development: I will describe a feature, you will create a test
file if necessary, come up with a suitable public API for the feature and write
tests for it (using the github.com/smartystreets/goconvey library). Then you'll
ask me to check the tests and let me adjust them if needed. Then I'll say
"go ahead", and you'll make the tests pass by writing the implementation in
another file.

When you create a new file, include the comment block at the top of main.go.

Wrap external dependencies in our own mockable interface.

Use github.com/go-sql-driver/mysql for working with MySQL.
Use github.com/joho/godotenv to load config such as database connection details
and server ports to listen on from a .env file.

Document all public structs and funcs. Documentation comments must end with a
full stop.
