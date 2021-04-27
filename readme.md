# Plan 9 Crypto Monorepo

While all of the API's are in development there is intermittent and very
neccessary changes that break everything, everything is being piled into this
repository in order to simplify several aspects of the work.

## Building

```bash
go run ./pod/podbuild/. builder
podbuild install
```

For other commands

```bash
podbuild help
```

A one-liner to run the node:

```bash
go get github.com/p9c/monorepo/pod/podbuild; go run github.com/p9c/monorepo/pod/podbuild node
```

Note that the `podbuild` isn't essential to use but without it, the build will not split file paths correctly,
and your logs will be very ugly. 

- @l0k18
