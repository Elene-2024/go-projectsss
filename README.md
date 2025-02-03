HTTP Caching Middleware in Go
This Go project implements a simple HTTP caching middleware. It caches responses for GET requests with a 200 OK status code. Cached responses are returned for subsequent requests until the cache expires (based on a configurable TTL).

Features:
Caches GET requests with 200 OK status code.
Supports cache expiration via TTL (Time-to-Live).
Thread-safe cache with sync.Map.
Usage:
Clone or download the project.

Run the program with:

bash
Copy
go run main.go
The program will make three requests:

First request: Fetches from the server.
Second request: Fetches from cache.
Third request: Fetches from the server again (after cache TTL expires).
Cache Configuration:
TTL: 5 seconds (for testing).
Cacheable Methods: GET.
Cacheable Status Codes: 200 OK.

