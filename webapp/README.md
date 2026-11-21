# SolidJS + TypeScript + Vite

Webapp for the personal finance tracker API. Built with [SolidJS 2.0](https://docs.solidjs.com),
[@solidjs/router](https://github.com/solidjs/solid-router) for routing, and
[Tailwind CSS v4](https://tailwindcss.com) for styling. It only talks to the Go API over HTTP
(default `http://localhost:8080`, override with `VITE_API_URL`).

Commands (from `webapp/`, uses pnpm):

- `pnpm dev` — dev server on :5173
- `pnpm test` — vitest run (jsdom)
- `pnpm build` — typecheck (tsc) + vite build
- `pnpm lint` — oxlint
