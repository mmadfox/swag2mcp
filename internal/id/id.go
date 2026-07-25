package id

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// Domain generates an MD5 hash ID for a domain name.
func Domain(domainName string) string {
	return hash(domainName)
}

// Collection generates an MD5 hash ID for a collection (domain + spec name).
func Collection(domainName, specName string) string {
	return hash(domainName, specName)
}

// Tag generates an MD5 hash ID for a tag (domain + collection + tag).
func Tag(domainName, collectionName, tag string) string {
	return hash(domainName, collectionName, tag)
}

// Method generates an MD5 hash ID for a method endpoint.
// It combines domainName, collectionName, method, path, and opID into a single hash.
func Method(domainName, collectionName, _ string, method, path, opID string) string {
	return hash(domainName, collectionName, method, path, opID)
}

func hash(parts ...string) string {
	for i, p := range parts {
		parts[i] = trimVal(p)
	}
	h := md5.Sum([]byte(strings.Join(parts, ":")))
	return hex.EncodeToString(h[:])
}

func trimVal(val string) string {
	val = strings.TrimSpace(val)
	val = strings.ReplaceAll(val, "/", "_")
	val = strings.ReplaceAll(val, " ", "_")
	return strings.ToLower(val)
}
