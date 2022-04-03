package playground

import (
	"html/template"
	"net/http"
)

var page = template.Must(template.New("graphiql").Parse(`<!DOCTYPE html>
<html>
  <head>
    <title>{{.title}}</title>
    <link
		rel="stylesheet"
		href="https://cdn.jsdelivr.net/npm/graphiql@{{.version}}/graphiql.min.css"
		integrity="{{.cssSRI}}"
		crossorigin="anonymous"
	/>
  </head>
  <body style="margin: 0;">
    <div id="graphiql" style="height: 100vh;"></div>

	<script
		src="https://cdn.jsdelivr.net/npm/react@17.0.2/umd/react.production.min.js"
		integrity="{{.reactSRI}}"
		crossorigin="anonymous"
	></script>
	<script
		src="https://cdn.jsdelivr.net/npm/react-dom@17.0.2/umd/react-dom.production.min.js"
		integrity="{{.reactDOMSRI}}"
		crossorigin="anonymous"
	></script>
	<script
		src="https://cdn.jsdelivr.net/npm/graphiql@{{.version}}/graphiql.min.js"
		integrity="{{.jsSRI}}"
		crossorigin="anonymous"
	></script>

    <script>
      const url = location.protocol + '//' + location.host + '{{.endpoint}}';
      const wsProto = location.protocol == 'https:' ? 'wss:' : 'ws:';
      const subscriptionUrl = wsProto + '//' + location.host + '{{.endpoint}}';

      const fetcher = GraphiQL.createFetcher({ url, subscriptionUrl });
      ReactDOM.render(
        React.createElement(GraphiQL, {
          fetcher: fetcher,
          headerEditorEnabled: true,
          shouldPersistHeaders: true
        }),
        document.getElementById('graphiql'),
      );
    </script>
  </body>
</html>
`))

func Handler(title string, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		err := page.Execute(w, map[string]string{
			"title":       title,
			"endpoint":    endpoint,
			"version":     "1.5.16",
			"cssSRI":      "sha256-HADQowUuFum02+Ckkv5Yu5ygRoLllHZqg0TFZXY7NHI=",
			"jsSRI":       "sha256-uHp12yvpXC4PC9+6JmITxKuLYwjlW9crq9ywPE5Rxco=",
			"reactSRI":    "sha256-Ipu/TQ50iCCVZBUsZyNJfxrDk0E2yhaEIz0vqI+kFG8=",
			"reactDOMSRI": "sha256-nbMykgB6tsOFJ7OdVmPpdqMFVk4ZsqWocT6issAPUF0=",
		})
		if err != nil {
			panic(err)
		}
	}
}
