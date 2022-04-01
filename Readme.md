#SQL Exec

SQL Exec is a small tool to run SQL queries against a database.

## Build
`env GOOS=windows GOARCH=amd64 go build -o sqlexec.exe main.go `

## Usage
`sqlexec [query | exec] "SELECT * FROM users" <configfilename>`

exec just execeutes the query and prints the affected rows.

query will execute the query and print the results.

