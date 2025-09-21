package main

import (
	"html/template"
	"net/http"
)

var jsTemplate = template.Must(template.New("").Parse(
	`(function() {
    const socket = new WebSocket("ws://{{.Host}}{{.Path}}");
    socket.onmessage = (event) => {
        if (event.data === "reload") {
            console.log('[doitlive]', 'Reloading page');
            // true = hard refresh
            window.location.reload({{.HardRefresh}});
        }
    };
})();`))

func handleJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")

	// create a websocket that points to the connect endpoint
	jsTemplate.Execute(w, map[string]string{
		"Host":        wsHost,
		"Path":        wsEndpoint,
		"HardRefresh": HardRefresh,
	})
}
