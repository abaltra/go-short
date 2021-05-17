## A URL shortener written in Golang with a Svelte frontend!

### Why?

I wanted to learn some Golang. And in the process of it I read up a little on Svelte and it sounded fun, so I threw it in there

### How to run

Just start both frontend and backend services. They are currently using the harcoded ports `5000` and `8000` respectively

```shell
cd server
go get
go run .
```

```shell
cd client
npm i
npm run dev
```

### Is this prod ready?

Oh God, not by a long shot. The code can be way cleaner, it needs a persistent backend and several optimizations. This is no more than a toy used as a learning exercise

### Contributing

PRs are welcome!