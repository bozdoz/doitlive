# doitlive

Window reloader in Go, with websockets and filewatchers: watches changed files and signals (via websockets) that a web app should reload

## Usage

1. Build/download the binary:

```sh
go build -o doitlive .
```

2. Move it to the root directory where you're building a web app (and want to watch files)

3. Run it (in the background) alongside your app (change `--host` or `--proxy` if you want with `--host=8080 --port=80`):

```sh
./doitlive
[doitlive] v2.0.0
[doitlive] Host: http://localhost:8000
[doitlive] Proxied: http://localhost:4000
```

5. Run your app (somehow), visit the proxy URL, then change files to see the page reload

## Assumptions

- `.git` and `node_modules` should be ignored
- it will only be run on `localhost`
- ~~consumer can inject script (such as Django middleware)~~
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
