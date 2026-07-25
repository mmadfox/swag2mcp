#!/bin/bash
# scripts/add-license.sh — adds AGPL header after package declaration in all .go files

for f in $(find . -name '*.go' -not -path './vendor/*' -not -name 'mock_*' -not -name '*.pb.go'); do
    if grep -q "SPDX-License-Identifier" "$f"; then
        continue
    fi
    # Insert header after the package line, with a blank line before it
    sed -i '' '/^package /a\
\
// SPDX-License-Identifier: AGPL-3.0-only\
//\
// Use of this software is governed by the AGPL v3 license\
// included in the /LICENSE file.
' "$f"
    echo "added: $f"
done

# Fix formatting — goimports removes blank lines before package
goimports -w $(find . -name '*.go' -not -path './vendor/*' -not -name 'mock_*' -not -name '*.pb.go') 2>/dev/null
