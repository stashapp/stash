package countries

// Value - for database/sql/driver.Valuer compatibility
// Value is a value that drivers must be able to handle.
// It is either nil, a type handled by a database driver's NamedValueChecker
// interface, or an instance of one of these types:
//
//   int64
//   float64
//   bool
//   []byte
//   string
//   time.Time
//
// If the driver supports cursors, a returned Value may also implement the Rows interface
// in this package. This is used, for example, when a user selects a cursor
// such as "select cursor(select * from my_table) from dual". If the Rows
// from the select is closed, the cursor Rows will also be closed.
type Value interface{}
