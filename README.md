# ⏹ unsocket

> consume websockets statelessly

To consume websockets you normally run a stateful app.
However, as it's becoming more and more common to build stateless systems hosted on AWS Lambda, Vercel or similar environments, consuming webhooks isn't that easy anymore.
`unsocket` is a proxy service that takes care of your stateful websocket connections and talks to your app through stateless HTTP calls.

* `unsocket` takes care of all websocket connection handling (#1, #2, #3)
* `unsocket` calls your app's HTTP endpoint on incoming websocket messages (#4, #5, #6, #7)
* `unsocket` turns your app's HTTP calls into outgoing websocket messages (#8, #9)

```
                     +----------------+
+----------------+   |                |   +----------------+
| my.app/webhook |   |    unsocket    |   |  ws://example  | 
+----------------+   |                |   +----------------+
        |            +----------------+           |
        |                    |                    |
        |           1. READY |                    |
        | <----------------- |                    |
        |                    |                    |
        | 2. CONNECT to      | 3. establish web-  |
        | ws://example       | socket connection  |
        | -----------------> | -----------------> |
        |                    |                    |
        |        ...         |        ...         |
        |                    |                    |
        |    5. pass message |    4. receive web- |
        |    in HTTP request |    socket message  |
        | <----------------- | <----------------- |
        |                    |                    |
        | 7. respond with    | 6. pass reply as   |
        | reply message      | websocket message  |
        | -----------------> | -----------------> |
        |                    |                    |
        |        ...         |        ...         |
        |                    |                    |
        | 8. POST /message   | 9. pass message as |
        | spontaneously      | websocket message  |
        | -----------------> | -----------------> |
        |                    |                    |
```

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
