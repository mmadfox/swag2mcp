/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

// Package auth provides authentication methods for API specifications.
// It supports 9 methods: none, basic, bearer, digest, hmac, oauth2-cc, oauth2-pwd, api-key, and script.
// Each method implements the Authenticator interface and can be applied to HTTP requests.
package auth
