# Changelog

## v1.0.8 — 2026-08-21

### Security
- Fixed stored XSS in message links (attribute injection)
- WebSocket CheckOrigin now restricted to localhost/LAN
- WS token: Authorization header preferred over query param
- Login rate limiting: 5 attempts / 5 min per IP
- Setup race condition fixed (mutex)
- Settings routes now require admin role
- GenerateKey no longer panics on rand failure

### Features
- Avatar upload (POST /api/users/avatar) — new backend
- Profile save (PUT /api/users/me) — new backend
- Message search pagination (limit + offset params)
- Max message length: 10,000 chars
- Poll option/question length validation (200/500 chars)
- Room name max 100 chars + XSS filter
- Room description/topic fields with validation
- Keyboard shortcuts: Cmd+K search, Esc close
- Push notifications via Service Worker
- 12-theme system (6 base × light/dark)
- Admin setup moved to /panel/setup
- Expired sessions cleaned hourly

### Fixes
- Password change now invalidates all sessions
- Poll vote requires room membership
- Poll close restricted to creator/admin/supervisor
- Logout clears all session localStorage keys
- CSS minified (21% smaller)

