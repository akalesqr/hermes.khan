# Khan v1.0.8

## 🔴 Security
- Fixed stored XSS in message links (attribute injection)
- WebSocket CheckOrigin restricted to localhost/LAN
- WS token: Authorization header preferred
- Login rate limiting: 5 attempts / 5 min per IP
- Setup race condition fixed (mutex)
- Settings routes require admin role
- GenerateKey no longer panics

## 🟠 Features & Fixes
- Avatar upload + profile save (new backend routes)
- Message search pagination (limit/offset)
- Max message length: 10,000 chars
- Poll option/question validation
- Room name max 100 chars + XSS filter
- Password change invalidates all sessions
- Poll vote requires membership; close restricted to creator/admin
- Expired sessions cleaned hourly
- Admin setup moved to /panel/setup
- 12-theme system, keyboard shortcuts, push notifications (SW)
- CSS minified (-21%)

**Tested:** full integration test passed on live server (setup → login → avatar → profile → rate limit → password change → poll → message limits)

## 📦 Downloads
| File | Platform |
|------|----------|
| khan-linux-amd64 | Linux x64 |
| khan-linux-arm64 | Linux ARM64 |
| khan-darwin-amd64 | macOS Intel |
| khan-darwin-arm64 | macOS Apple Silicon |
| khan-windows-amd64.exe | Windows x64 |
