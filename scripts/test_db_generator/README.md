This is a quick and dirty go script for generating a contrived database for testing purposes.

Edit the `config.yml` file to your liking. The numbers indicate the number of objects to generate, the `naming` section indicates the files from which to generate names.

May cause unexpected behaviour if run against an existing database file.

To run - from the `test_db_generator`:
`go run .`

The database file will be generated in the current directory.