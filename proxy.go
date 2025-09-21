package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"text/template"
)

var scriptTemplate = template.Must(template.New("").Parse(
	`<script>
(function() {
	const socket = new WebSocket("ws://{{.Host}}{{.Path}}");
	socket.onmessage = (event) => {
		if (event.data === "reload") {
			console.log('[doitlive]', 'Reloading page');
			// true = hard refresh
			window.location.reload({{.HardRefresh}});
		}
	};
})();
</script></body>`,
))

func runProxy() (proxy *httputil.ReverseProxy) {
	host, err := url.Parse(fmt.Sprintf("http://localhost:%d", *host_port))

	if err != nil {
		msg := fmt.Sprintf("Can't parse host url with port: %d", *host_port)
		panic(msg)
	}

	var injectedJS bytes.Buffer

	scriptTemplate.Execute(&injectedJS, map[string]any{
		"Host":        fmt.Sprintf("localhost:%d", *proxy_port),
		"Path":        wsEndpoint,
		"HardRefresh": *is_hard_refresh,
	})

	proxy = httputil.NewSingleHostReverseProxy(host)

	proxy.ModifyResponse = func(r *http.Response) error {
		ctype := r.Header.Get("Content-Type")

		debug("[proxy]", ctype)

		if strings.HasPrefix(ctype, "text/html") {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Println("Failed to read body")
				return nil
			}
			r.Body.Close()

			modified := bytes.Replace(body, []byte("</body>"), injectedJS.Bytes(), 1)

			r.Body = io.NopCloser(bytes.NewReader(modified))
			r.ContentLength = int64(len(modified))
			r.Header.Set("Content-Length", fmt.Sprint(len(modified)))
		}

		return nil
	}

	return
}
