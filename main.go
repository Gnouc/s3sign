package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"
)

var defaultS3URL = "https://sgp1.digitaloceanspaces.com"

func main() {
	bucket := os.Getenv("BUCKET")
	key := os.Getenv("KEY")
	if bucket == "" || key == "" {
		fmt.Fprintln(os.Stderr, "bucket or key is empty, set BUCKET=xxx and KEY=xxx")
		os.Exit(1)
	}
	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	method := os.Getenv("METHOD")
	if method == "" {
		method = "GET"
	}

	s3URL := os.Getenv("S3_URL")
	if s3URL == "" {
		s3URL = defaultS3URL
	}

	expiredIn, err := strconv.Atoi(os.Getenv("EXPIRED_IN"))
	if err != nil {
		expiredIn = 3600
	}
	timestamp := time.Now().Unix() + int64(expiredIn)
	message := buildMessage(method, bucket, key, timestamp)

	fmt.Printf("%s/%s/%s?AWSAccessKeyId=%s&Expires=%d&Signature=%s\n", s3URL, bucket, key, accessKey, timestamp, sign(secretKey, message))
}

func buildMessage(method, bucket, key string, timestamp int64) string {
	return fmt.Sprintf("%s\n\n\n%d\n/%s/%s", method, timestamp, bucket, key)
}

func sign(secretKey, message string) string {
	h := hmac.New(sha1.New, []byte(secretKey))
	_, _ = h.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
