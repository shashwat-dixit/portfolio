package middleware

// TODO: implement Redis response caching middleware
// - generate cache key from request path + query params
// - check Redis for cached response
// - if hit: write cached response and return
// - if miss: wrap ResponseWriter to capture body, call next, cache response
// - skip caching for non-GET methods and error responses
