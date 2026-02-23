# Phase 01 — Authentication + Users (JWT + Roles)

**Date:** 2026-02-23
**Priority:** High
**Status:** ⬜ TODO
**Estimate:** M (1–3 days)
**Depends on:** Phase 00

---

## Overview
Implement user registration/login with JWT-based auth using HttpOnly cookies and role-based access control.

## Key Insights
- HttpOnly cookie prevents XSS token theft (safer than LocalStorage).
- Simple access token only for MVP (no refresh token rotation — YAGNI).
- Two roles: `creator` and `reader`.

## Requirements

### Functional
- User can register with username, email, password.
- User can login and receive JWT in HttpOnly cookie.
- User can logout (clear cookie).
- Authenticated user can fetch their profile (`GET /api/me`).
- Creator-only endpoints reject unauthenticated/unauthorized requests.

### Non-Functional
- Password hashed with bcrypt.
- JWT expiry: 7 days (keep simple for MVP).

## Architecture

```
Backend:
  internal/
    auth/
      jwt.go           # JWT generation, validation, cookie helpers
    http/
      middleware/
        auth.go         # Extract JWT from cookie, set user context
      handlers/
        auth.go         # register, login, logout, me
    domain/
      user.go           # User struct
    services/
      user_service.go   # business logic
    repo/
      user_repo.go      # Neo4j user queries

Frontend:
  src/
    api/
      auth.ts           # API client for auth endpoints
    stores/
      auth.ts           # Pinia auth store (user profile)
    pages/
      LoginPage.vue
      RegisterPage.vue
```

## Related Code Files

### Files to Create
- `backend/internal/domain/user.go`
- `backend/internal/repo/user_repo.go`
- `backend/internal/services/user_service.go`
- `backend/internal/auth/jwt.go`
- `backend/internal/http/middleware/auth.go`
- `backend/internal/http/handlers/auth.go`
- `backend/migrations/001_user_constraints.cypher`
- `frontend/src/api/auth.ts`
- `frontend/src/stores/auth.ts`

### Files to Modify
- `backend/internal/http/router.go` — add auth routes + middleware
- `frontend/src/pages/LoginPage.vue` — implement form
- `frontend/src/pages/RegisterPage.vue` — implement form
- `frontend/src/router/index.ts` — add route guards

## Implementation Steps

### 1. Neo4j constraints + indexes (S)
Create `backend/migrations/001_user_constraints.cypher`:
```cypher
CREATE CONSTRAINT user_id_unique IF NOT EXISTS FOR (u:User) REQUIRE u.id IS UNIQUE;
CREATE CONSTRAINT user_email_unique IF NOT EXISTS FOR (u:User) REQUIRE u.email IS UNIQUE;
CREATE CONSTRAINT user_username_unique IF NOT EXISTS FOR (u:User) REQUIRE u.username IS UNIQUE;
```

### 2. User domain model (S)
`backend/internal/domain/user.go`:
```go
type User struct {
    ID           string `json:"id"`
    Username     string `json:"username"`
    Email        string `json:"email"`
    PasswordHash string `json:"-"`
    Role         string `json:"role"` // "creator" or "reader"
}
```

### 3. User repository (M)
`backend/internal/repo/user_repo.go`:

**Create user:**
```cypher
CREATE (u:User {
  id: $id,
  username: $username,
  email: $email,
  password_hash: $password_hash,
  role: $role
})
RETURN u { .id, .username, .email, .role } AS user;
```

**Find user by email:**
```cypher
MATCH (u:User {email: $email})
RETURN u { .id, .username, .email, .password_hash, .role } AS user;
```

**Find user by ID:**
```cypher
MATCH (u:User {id: $id})
RETURN u { .id, .username, .email, .role } AS user;
```

### 4. Password hashing (S)
- Use `golang.org/x/crypto/bcrypt`
- Hash on register, compare on login.

### 5. JWT implementation (M)
`backend/internal/auth/jwt.go`:
- Generate token with claims: `user_id`, `role`, `exp` (7 days).
- Validate token and extract claims.
- Cookie helpers:
  - `SetAuthCookie(c *gin.Context, token string)` — HttpOnly, Secure, SameSite=Lax, Path=/
  - `ClearAuthCookie(c *gin.Context)`

### 6. Auth middleware (M)
`backend/internal/http/middleware/auth.go`:
- Extract `access_token` cookie.
- Validate JWT → set `userID` and `role` in Gin context.
- Return 401 if missing/invalid.

### 7. Auth endpoints (M)
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/auth/register` | POST | Create user, set cookie |
| `/api/auth/login` | POST | Validate credentials, set cookie |
| `/api/auth/logout` | POST | Clear cookie |
| `/api/me` | GET | Return current user profile |

**Register request:**
```json
{ "username": "string", "email": "string", "password": "string", "role": "creator|reader" }
```

**Login request:**
```json
{ "email": "string", "password": "string" }
```

### 8. Frontend auth pages (M)
- Login/Register forms with validation.
- Pinia `auth` store:
  - `user: User | null`
  - `isAuthenticated: boolean` (computed)
  - `fetchMe()` — call `GET /api/me` on app init
  - `login(email, password)` — call POST, then fetchMe
  - `register(...)` — call POST, then fetchMe
  - `logout()` — call POST, clear user
- Router guard: redirect to `/login` if unauthenticated on creator routes.

## Todo List
- [ ] Create Neo4j user constraints migration
- [ ] Implement User domain model
- [ ] Implement user repository (create, findByEmail, findByID)
- [ ] Implement password hashing with bcrypt
- [ ] Implement JWT generation and validation
- [ ] Implement cookie helpers (set/clear)
- [ ] Implement auth middleware
- [ ] Implement register endpoint
- [ ] Implement login endpoint
- [ ] Implement logout endpoint
- [ ] Implement /api/me endpoint
- [ ] Update router.go with auth routes
- [ ] Implement frontend auth API client
- [ ] Implement Pinia auth store
- [ ] Build Login page UI
- [ ] Build Register page UI
- [ ] Add router navigation guards
- [ ] Test: register → login → /api/me → logout flow

## Success Criteria
- User can register and login successfully.
- Auth cookie is set (HttpOnly) and required endpoints reject 401 without it.
- Creator-only endpoints enforce auth middleware.
- Frontend redirects unauthenticated users to login.

## Risk Assessment
| Risk | Impact | Mitigation |
|------|--------|------------|
| Token stored in LocalStorage (XSS) | High | Use HttpOnly cookie only |
| JWT secret leaked | High | Load from env var, never commit |
| Brute-force login | Medium | Add rate limiting later if abuse appears (YAGNI for MVP) |

## Security Considerations
- Passwords hashed with bcrypt (cost ≥ 10).
- JWT in HttpOnly cookie, Secure flag in production.
- SameSite=Lax to prevent basic CSRF.
- No tokens in LocalStorage or response body.

## Next Steps
- → Phase 02: Creator Story CRUD
