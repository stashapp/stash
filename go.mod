module github.com/stashapp/stash

require (
	github.com/99designs/gqlgen v0.10.1
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/bmatcuk/doublestar v1.1.5
	github.com/disintegration/imaging v1.6.1
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/gobuffalo/packr/v2 v2.0.2
	github.com/golang-migrate/migrate/v4 v4.7.0
	github.com/gorilla/websocket v1.4.1
	github.com/h2non/filetype v1.0.10
	github.com/jmoiron/sqlx v1.2.0
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.5.0
	github.com/vektah/gqlparser v1.2.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/image v0.0.0-20190118043309-183bebdce1b2 // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

go 1.13
