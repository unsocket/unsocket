# ⏹ unsocket

> consume websockets statelessly

To consume websockets you normally run a stateful app.
However, as it's becoming more and more common to build stateless systems hosted on AWS Lambda, Vercel or similar environments, consuming webhooks isn't that easy anymore.
`unsocket` is a proxy service that takes care of your stateful websocket connections and talks to your app through stateless HTTP calls.

* `unsocket` calls your app's HTTP endpoint on incoming websocket messages
* `unsocket` turns your app's HTTP calls into outgoing websocket messages
* `unsocket` takes care of all websocket connection handling

## Run

Call `unsocket` with only the url to your app's HTTP endpoint:

```
unsocket http://localhost:3000/messages
```

It's expected that your endpoint returns proper connection data during initialization:

```
POST http://localhost:3000/messages
[{"init":true}]

HTTP 200 OK
[{"type":"connect","url":"wss://example.com"},{"type":"message"}]
```

## License

MIT
