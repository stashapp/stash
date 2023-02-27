package playground

import (
	"html/template"
	"net/http"
)

var altairPage = template.Must(template.New("altair").Parse(`<!doctype html>
<html>

<head>
  <meta charset="utf-8">
  <title>{{.title}}</title>
  <base href="https://cdn.jsdelivr.net/npm/altair-static@{{.version}}/build/dist/">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <link rel="icon" type="image/x-icon" href="favicon.ico">
  <link href="styles.css" rel="stylesheet" crossorigin="anonymous" integrity="{{.cssSRI}}"/>
</head>

<body>
  <app-root>
    <style>
      .loading-screen {
        display: none;
      }
    </style>
    <div class="loading-screen styled">
      <div class="loading-screen-inner">
        <div class="loading-screen-logo-container">
          <img src="assets/img/logo_350.svg" alt="Altair">
        </div>
        <div class="loading-screen-loading-indicator">
          <span class="loading-indicator-dot"></span>
          <span class="loading-indicator-dot"></span>
          <span class="loading-indicator-dot"></span>
        </div>
      </div>
    </div>
  </app-root>

  <script rel="preload" as="script" type="text/javascript" crossorigin="anonymous" integrity="{{.mainSRI}}" src="main.js"></script>
  <script rel="preload" as="script" type="text/javascript" crossorigin="anonymous" integrity="{{.polyfillsSRI}}" src="polyfills.js"></script>
  <script rel="preload" as="script" type="text/javascript" crossorigin="anonymous" integrity="{{.runtimeSRI}}" src="runtime.js"></script>

  <script>
{{- if .endpointIsAbsolute}}
	const url = {{.endpoint}};
	const subscriptionUrl = {{.subscriptionEndpoint}};
{{- else}}
	const url = location.protocol + '//' + location.host + {{.endpoint}};
	const wsProto = location.protocol == 'https:' ? 'wss:' : 'ws:';
	const subscriptionUrl = wsProto + '//' + location.host + {{.endpoint}};
{{- end}}
    var altairOptions = {
        endpointURL: url,
        subscriptionsEndpoint: subscriptionUrl,
    };
    window.addEventListener("load", function() {
      AltairGraphQL.init(altairOptions);
    });
  </script>
</body>

</html>`))

// AltairHandler responsible for setting up the altair playground
func AltairHandler(title, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := altairPage.Execute(w, map[string]interface{}{
			"title":                title,
			"endpoint":             endpoint,
			"endpointIsAbsolute":   endpointHasScheme(endpoint),
			"subscriptionEndpoint": getSubscriptionEndpoint(endpoint),
			"version":              "5.0.5",
			"cssSRI":               "sha256-kZ35e5mdMYN5ALEbnsrA2CLn85Oe4hBodfsih9BqNxs=",
			"mainSRI":              "sha256-nWdVTcGTlBDV1L04UQnqod+AJedzBCnKHv6Ct65liHE=",
			"polyfillsSRI":         "sha256-1aVEg2sROcCQ/RxU3AlcPaRZhZdIWA92q2M+mdd/R4c=",
			"runtimeSRI":           "sha256-cK2XhXqQr0WS1Z5eKNdac0rJxTD6miC3ubd+aEVMQDk=",
		})
		if err != nil {
			panic(err)
		}
	}
}
