# Newsletter

This repo is WIP.

> This repo contains the source code I use to self host my personal newsletter. The subscription form integrates with [listmonk](https://listmonk.app).

Subscribe to my newsletter at [news.mrkaran.dev](news.mrkaran.dev)

## Setup

### Frontend

```sh
$ cd frontend
$ yarn # to add the required dependencies
$ yarn run start:dev # starts a local dev server
$ yarn watch # looks for changes in `css` directory and uses `tailwind` css to build the final minified css.
```

### Backend

```sh
$ cd backend
$ make deps
$ make dist
$ make run # starts a webserver with static assets packed into binary. API is available at `/api` and static assets at `/static`.
```