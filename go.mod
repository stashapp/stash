module github.com/stashapp/stash

require (
	github.com/99designs/gqlgen v0.12.2
	github.com/Yamashou/gqlgenc v0.0.0-20200902035953-4dbef3551953
	github.com/antchfx/htmlquery v1.2.3
	github.com/bmatcuk/doublestar/v2 v2.0.1
	github.com/chromedp/cdproto v0.0.0-20200608134039-8a80cdaf865c
	github.com/chromedp/chromedp v0.5.3
	github.com/disintegration/imaging v1.6.0
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/gobuffalo/packr/v2 v2.0.2
	github.com/golang-migrate/migrate/v4 v4.3.1
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.0
	github.com/gorilla/websocket v1.4.2
	github.com/h2non/filetype v1.0.8
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/jmoiron/sqlx v1.2.0
	github.com/json-iterator/go v1.1.9
	github.com/mattn/go-sqlite3 v1.13.0
	github.com/natefinch/pie v0.0.0-20170715172608-9a0d72014007
	github.com/remeh/sizedwaitgroup v1.0.0
	github.com/rs/cors v1.6.0
	github.com/shurcooL/graphql v0.0.0-20181231061246-d48a9a75455f
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.5.1
	github.com/tidwall/gjson v1.6.0
	github.com/vektah/gqlparser/v2 v2.0.1
	github.com/vektra/mockery/v2 v2.2.1
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/image v0.0.0-20190802002840-cff245a6509b
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	golang.org/x/tools v0.0.0-20200915031644-64986481280e // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

go 1.13
