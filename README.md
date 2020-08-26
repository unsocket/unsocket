<img align="left" alt="unsocket-logo" width="128" src="https://user-images.githubusercontent.com/198988/88155428-5cfb5000-cc08-11ea-9290-67259425178b.png" />


# unsocket

> consume websockets statelessly

To consume websockets you normally run a stateful app.
However, as it's becoming more and more common to build stateless systems hosted on AWS Lambda, Vercel or similar environments, consuming websockets isn't that easy anymore.
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
        |    1. HTTP request |                    |
        |    connection info |                    |
        | <----------------- |                    |
        |                    |                    |
        | 2. HTTP respond    | 3. establish web-  |
        | with ws://example  | socket connection  |
        | -----------------> | -----------------> |
        ⋮                     ⋮                    ⋮
        |   5. HTTP request  |    4. receive web- |
        |   received message |    socket message  |
        | <----------------- | <----------------- |
        |                    |                    |
        | 7. HTTP respond    | 6. pass reply as   |
        | with a reply       | websocket message  |
        | -----------------> | -----------------> |
        ⋮                     ⋮                    ⋮
        | 8. POST /message   | 9. pass message as |
        | spontaneously      | websocket message  |
        | -----------------> | -----------------> |
        ⋮                     ⋮                    ⋮
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
