// SPDX-License-Identifier: AGPL-3.0-or-later
// SPDX-FileCopyrightText: 2026 GOTUNIX Networks <code@gotunix.net>
// SPDX-FileCopyrightText: 2026 Justin Ovens <code@gotunix.net>

package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-discord-notifications/config"
)

func TestVerifyGitHubSignature(t *testing.T) {
	secret := "my-secret-key"
	payload := []byte(`{"zen": "Keep it simple"}`)

	// Temporarily set/restore config.WebhookSecret
	oldSecret := config.WebhookSecret
	config.WebhookSecret = secret
	defer func() { config.WebhookSecret = oldSecret }()

	// Compute signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	t.Run("Valid Signature", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/webhook/github", bytes.NewReader(payload))
		req.Header.Set("X-Hub-Signature-256", "sha256="+expectedSig)
		req.Header.Set("X-GitHub-Event", "ping")

		body, err := verifyGitHubSignature(req, secret)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !bytes.Equal(body, payload) {
			t.Errorf("expected body %s, got %s", payload, body)
		}
	})

	t.Run("Invalid Signature", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/webhook/github", bytes.NewReader(payload))
		req.Header.Set("X-Hub-Signature-256", "sha256=invalid-signature-hex")
		req.Header.Set("X-GitHub-Event", "ping")

		_, err := verifyGitHubSignature(req, secret)
		if err == nil {
			t.Fatal("expected signature mismatch error, got nil")
		}
	})

	t.Run("Missing Signature Header", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/webhook/github", bytes.NewReader(payload))
		req.Header.Set("X-GitHub-Event", "ping")

		_, err := verifyGitHubSignature(req, secret)
		if err == nil {
			t.Fatal("expected error due to missing signature header, got nil")
		}
	})
}

func TestGitHubHandler(t *testing.T) {
	// Temporarily clear or set config webhook secret to simplify or control tests
	oldSecret := config.WebhookSecret
	defer func() { config.WebhookSecret = oldSecret }()
	config.WebhookSecret = "" // No secret to skip signature checking for testing handlers

	t.Run("Ping Event", func(t *testing.T) {
		payload := []byte(`{"zen": "Avoid placeholders", "repository": {"full_name": "test/repo", "html_url": "https://github.com/test/repo"}}`)
		req := httptest.NewRequest("POST", "/webhook/github", bytes.NewReader(payload))
		req.Header.Set("X-GitHub-Event", "ping")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(githubHandler)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Push Event", func(t *testing.T) {
		payload := []byte(`{
			"ref": "refs/heads/main",
			"compare": "https://github.com/test/repo/compare/a...b",
			"repository": {"full_name": "test/repo"},
			"pusher": {"name": "test-user"},
			"commits": [
				{"id": "sha123456", "message": "First commit\nDetail info", "url": "https://github.com/c1", "author": {"name": "test-user"}}
			]
		}`)
		req := httptest.NewRequest("POST", "/webhook/github", bytes.NewReader(payload))
		req.Header.Set("X-GitHub-Event", "push")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(githubHandler)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Pull Request Event", func(t *testing.T) {
		payload := []byte(`{
			"action": "opened",
			"number": 42,
			"pull_request": {
				"number": 42,
				"html_url": "https://github.com/test/repo/pull/42",
				"title": "Amazing feature",
				"body": "PR description",
				"user": {"login": "pr-author"},
				"merged": false
			},
			"repository": {"full_name": "test/repo"},
			"sender": {"login": "pr-author"}
		}`)
		req := httptest.NewRequest("POST", "/webhook/github", bytes.NewReader(payload))
		req.Header.Set("X-GitHub-Event", "pull_request")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(githubHandler)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Invalid Payload JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/webhook/github", bytes.NewReader([]byte(`{invalid-json`)))
		req.Header.Set("X-GitHub-Event", "push")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(githubHandler)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("Form Urlencoded Event", func(t *testing.T) {
		formBody := []byte(`payload=%7B%22zen%22%3A%22Form+encoded+Zen%22%2C%22repository%22%3A%7B%22full_name%22%3A%22test%2Frepo%22%7D%7D`)
		req := httptest.NewRequest("POST", "/webhook/github", bytes.NewReader(formBody))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("X-GitHub-Event", "ping")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(githubHandler)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}
	})
}
