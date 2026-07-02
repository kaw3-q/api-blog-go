package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kaw3-q/api-blog-go/internal/auth"
	"github.com/kaw3-q/api-blog-go/internal/models"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	handler := SecurityHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Header().Get("X-Frame-Options") != "DENY" {
		t.Errorf("Expected X-Frame-Options DENY, got %s", rr.Header().Get("X-Frame-Options"))
	}
	if rr.Header().Get("X-XSS-Protection") != "1; mode=block" {
		t.Errorf("Expected X-XSS-Protection 1; mode=block, got %s", rr.Header().Get("X-XSS-Protection"))
	}
	if rr.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("Expected X-Content-Type-Options nosniff, got %s", rr.Header().Get("X-Content-Type-Options"))
	}
}

func TestRoleMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	roleAdminMiddleware := RoleMiddleware(models.RoleAdmin)(nextHandler)

	// Case 1: No claims in context -> Forbidden
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr1 := httptest.NewRecorder()
	roleAdminMiddleware.ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rr1.Code)
	}

	// Case 2: User role claims in context (not admin) -> Forbidden
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	claims2 := &auth.Claims{UserID: 1, Role: models.RoleUser}
	ctx2 := context.WithValue(req2.Context(), UserKey, claims2)
	req2 = req2.WithContext(ctx2)
	rr2 := httptest.NewRecorder()
	roleAdminMiddleware.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rr2.Code)
	}

	// Case 3: Admin role claims in context -> Success
	req3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	claims3 := &auth.Claims{UserID: 2, Role: models.RoleAdmin}
	ctx3 := context.WithValue(req3.Context(), UserKey, claims3)
	req3 = req3.WithContext(ctx3)
	rr3 := httptest.NewRecorder()
	roleAdminMiddleware.ServeHTTP(rr3, req3)
	if rr3.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr3.Code)
	}
}
