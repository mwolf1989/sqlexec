#SQL Exec

SQL Exec is a small tool to run SQL queries against a database.

## Build
`env GOOS=windows GOARCH=amd64 go build -o sqlexec.exe main.go `

## Usage
`sqlexec "SELECT * FROM users"`

