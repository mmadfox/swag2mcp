#!/bin/bash
# scripts/add-license.sh — adds Apache 2.0 SPDX header to all .go files

HEADER='/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */'

for f in $(find . -name '*.go' -not -path './vendor/*' -not -name 'mock_*' -not -name '*.pb.go'); do
    if grep -q "SPDX-FileCopyrightText" "$f"; then
        continue
    fi
    printf '%s\n\n' "$HEADER" | cat - "$f" > "${f}.tmp" && mv "${f}.tmp" "$f"
    echo "added: $f"
done
