# doitlive

Window reloader in Go, with websockets and filewatchers: watches changed files and signals (via websockets) that a web app should reload (without cache)

## Usage

1. Build/download the binary:

```sh
go build -o doitlive .
```

2. Move it to a directory where you're building a web app

3. Insert this script into your web app, somehow:

```html
<script src="http://localhost:35729/doitlive.js"></script>
```

<details>
<summary>Example Django middleware</summary>

some middleware.py:

```python
class DoItLiveMiddleware:
    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        response = self.get_response(request)

        # Only inject into HTML pages (not JSON, CSS, etc.)
        content_type = response.get("Content-Type", "")
        if "text/html" in content_type and hasattr(response, "content"):
            snippet = """
            <script src="http://localhost:35729/doitlive.js"></script>
            """

            # Insert before closing </body>, fallback append at end
            content = response.content.decode("utf-8")
            if "</body>" in content:
                content = content.replace("</body>", snippet + "</body>")
            else:
                content += snippet

            response.content = content.encode("utf-8")
            response["Content-Length"] = len(response.content)

        return response
```

in settings.py:

```python
if DEBUG:
    MIDDLEWARE += ["config.middleware.DoItLiveMiddleware"]
```

</details>

4. Run it (in the background) alongside your app:

```sh
./doitlive
[doitlive] v1.0.0
[doitlive] JS url: http://localhost:35729/doitlive.js
[doitlive] WS Endpoint: http://localhost:35729/ws/reload
```

5. Run your app (somehow), then change files to see the page reload

## Assumptions

- `.git` and `node_modules` should be ignored
- it will only be run on `localhost`
- consumer can inject script (such as Django middleware)
- ~~consumer is happy with hard reload (no caching)~~

<details>
<summary>Using soft reload in Django</summary>

In context_processors.py:

```python

def static_key (_request) :
	if settings.DEBUG:
		import random
		import string

		def random_string(length=5):
				letters = string.ascii_letters
				return ''.join(random.choice(letters) for _ in range(length))

		return {
			'static_key': random_string()
		}
	else:
		return {
			# TODO: needs to pull version
			'static_key': 'v1'
		}
```

Then in the template:

```html
<link
  rel="stylesheet"
  type="text/css"
  href="{% static 'css/style.css' %}?{{ static_key }}"
/>
```

</details>
