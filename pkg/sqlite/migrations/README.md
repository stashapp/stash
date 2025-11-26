# Creating a migration

1. Create new migration file in the migrations directory with the format `NN_description.up.sql`, where `NN` is the next sequential number.

2. Update `pkg/sqlite/database.go` to update the `appSchemaVersion` value to the new migration number.

For migrations requiring complex logic or config file changes, see existing custom migrations for examples.