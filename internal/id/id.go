package id

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

// Domain generates an MD5 hash ID for a domain name.
func Domain(domainName string) string {
	domainName = trimVal(domainName)
	hash := md5.Sum([]byte(domainName))
	return hex.EncodeToString(hash[:])
}

// Collection generates an MD5 hash ID for a collection (domain + spec name).
func Collection(domainName, specName string) string {
	domainName = trimVal(domainName)
	specName = trimVal(specName)
	input := fmt.Sprintf("%s:%s", domainName, specName)
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// Tag generates an MD5 hash ID for a tag (domain + collection + tag).
func Tag(domainName, collectionName, tag string) string {
	domainName = trimVal(domainName)
	collectionName = trimVal(collectionName)
	tag = trimVal(tag)
	input := fmt.Sprintf("%s:%s:%s", domainName, collectionName, tag)
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// Method generates an MD5 hash ID for a method endpoint.
// It combines domainName, collectionName, tagName, method, path, and opID into a single hash.
func Method(domainName, collectionName, _ string, method, path, opID string) string {
	domainName = trimVal(domainName)
	collectionName = trimVal(collectionName)
	method = trimVal(method)
	path = trimVal(path)
	opID = trimVal(opID)
	input := fmt.Sprintf("%s:%s:%s:%s:%s", domainName, collectionName, method, path, opID)
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func trimVal(val string) string {
	val = strings.TrimSpace(val)
	val = strings.ReplaceAll(val, "/", "_")
	val = strings.ReplaceAll(val, " ", "_")
	val = strings.ToLower(val)
	return val
}
