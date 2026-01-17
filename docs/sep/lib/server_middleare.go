package lib

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"



	"github.com/microcosm-cc/bluemonday"
)

func Compression(next http.Handler) http.Handler {
	fmt.Println("Compression Middleware...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Compression Middleware being returned...")

		// Check of the client accepts gzip encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Set the respone header
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Wrap the ResponseWriter
		w = &gzipResponseWriter{ResponseWriter: w, Writer: gz}

		next.ServeHTTP(w, r)
		fmt.Println("Sent response from Compression Middleware")
	})
}

// gzipResponseWriter wraps http.ResponseWriter to write gzipped responses
type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}



// api is hosted at www.myapi.com
// frontend server is at www.myfrontend.com

// Allowed origins
var allowedOrigins = []string{
	"https://my-origin-url.com",
	"https://www.myfrontend.com",
	"https://localhost:3000",
	"http://localhost:3000",
	"http://localhost:5173",
}

func Cors(next http.Handler) http.Handler {
	fmt.Println("Cors Middleware...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		// fmt.Println(origin)

		if isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			// Set other CORS headers
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")

		} else {
			http.Error(w, "Not allowed by CORS", http.StatusForbidden)
			return
		}


		// Handle preflight request
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isOriginAllowed(origin string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}


func XSSMiddleware(next http.Handler) http.Handler {
	fmt.Println("****** Intializing XSSMiddleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("++++++++++++ XSSMiddleware Ran")

		// Sanitize the URL Path
		sanitizedPath, err := clean(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Sanitize query params
		params := r.URL.Query()
		sanitizedQuery := make(map[string][]string)
		for key, values := range params {
			sanitizedKey, err := clean(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var sanitizedValues []string
			for _, value := range values {
				cleanValue, err := clean(value)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				sanitizedValues = append(sanitizedValues, cleanValue.(string))
			}
			sanitizedQuery[sanitizedKey.(string)] = sanitizedValues
		}

		r.URL.Path = sanitizedPath.(string)
		r.URL.RawQuery = url.Values(sanitizedQuery).Encode()

		// Sanitize request body
		if r.Header.Get("Content-Type") == "application/json" {
			if r.Body != nil {
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, ErrorHandler(err, "Error reading request body").Error(), http.StatusBadRequest)
					return
				}

				bodyString := strings.TrimSpace(string(bodyBytes))

				// Reset the request body
				r.Body = io.NopCloser(bytes.NewReader([]byte(bodyString)))

				if len(bodyString) > 0 {
					var inputData interface{}
					err := json.NewDecoder(bytes.NewReader([]byte(bodyString))).Decode(&inputData)
					if err != nil {
						http.Error(w, ErrorHandler(err, "Invalid JSON body").Error(), http.StatusBadRequest)
						return
					}

					// Sanitize the JSON body
					sanitizedData, err := clean(inputData)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					// Marshal the sanitized data back to the body
					sanitizedBody, err := json.Marshal(sanitizedData)
					if err != nil {
						http.Error(w, ErrorHandler(err, "Error sanitizing body").Error(), http.StatusBadRequest)
						return
					}

					r.Body = io.NopCloser(bytes.NewReader(sanitizedBody))
					fmt.Println("Sanitized body:", string(sanitizedBody))
				} else {
					log.Println("Request body is empty")
				}
			} else {
				log.Println("No body in the request")
			}
		} else if r.Header.Get("Content-Type") != "" {
			log.Printf("Received request with unsupported Content-Type: %s. Expected application/json.\n", r.Header.Get("Content-Type"))
			http.Error(w, "Unsupported Content-Type. Please use application/json.", http.StatusUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
		fmt.Println("Sending response from XSSMiddleware Ran")
	})
}

// Clean sanitizes input data to prevent XSS attacks
func clean(data interface{}) (interface{}, error) {

	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			v[key] = sanitizeValue(value)
		}
		return v, nil
	case []interface{}:
		for i, value := range v {
			v[i] = sanitizeValue(value)
		}
		return v, nil
	case string:
		return sanitizeString(v), nil
	default:
		return nil, ErrorHandler(fmt.Errorf("unsupported type: %T", data), fmt.Sprintf("unsupported type: %T", data))
	}
}

func sanitizeValue(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		return sanitizeString(v)
	case map[string]interface{}:
		for k, value := range v {
			v[k] = sanitizeValue(value)
		}
		return v
	case []interface{}:
		for i, value := range v {
			v[i] = sanitizeValue(value)
		}
		return v
	default:
		return v // Return v as it is unsupported
	}
}

func sanitizeString(value string) string {
	return bluemonday.UGCPolicy().Sanitize(value)
}


func SecurityHeaders(next http.Handler) http.Handler {
	fmt.Println("SecurityHeaders Middleware...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("SecurityHeaders Middleware being returned...")
		w.Header().Set("X-DNS-Prefetch-Control", "off")

		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Powered-By", "Django")
		w.Header().Set("Server", "")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Permissions-Policy", "geolocation=(self), microphone=()")

		next.ServeHTTP(w, r)
		fmt.Println("SecurityHeaders Middleware ends...")
	})
}

// BASIC MIDDLEWARE SKELETON
// func securityHeaders(next http.Handler) http.Handler {
// 	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
// next.ServeHTTP(w, r)
// 	})
// }



type HPPOptions struct {
	CheckQuery                  bool
	CheckBody                   bool
	CheckBodyOnlyForContentType string
	Whitelist                   []string
}

func Hpp(options HPPOptions) func(http.Handler) http.Handler {
	fmt.Println("HPP Middleware...")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("HPP Middleware being returned...")
			if options.CheckBody && r.Method == http.MethodPost && isCorrectContentType(r, options.CheckBodyOnlyForContentType) {
				// filter the body params
				filterBodyParams(r, options.Whitelist)
			}
			if options.CheckQuery && r.URL.Query() != nil {
				// filter the query params
				filterQueryParams(r, options.Whitelist)
			}
			next.ServeHTTP(w, r)
			fmt.Println("HPP Middleware ends...")
		})
	}
}

func isCorrectContentType(r *http.Request, contentType string) bool {
	return strings.Contains(r.Header.Get("Content-Type"), contentType)
}

func filterBodyParams(r *http.Request, whitelist []string) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, v := range r.Form {
		if len(v) > 1 {
			r.Form.Set(k, v[0]) // first value
			// r.Form.Set(k, v[len(v)-1]) // last value
		}
		if !isWhiteListed(k, whitelist) {
			delete(r.Form, k)
		}
	}
}

func filterQueryParams(r *http.Request, whitelist []string) {
	query := r.URL.Query()

	for k, v := range query {
		if len(v) > 1 {
			query.Set(k, v[0]) // first value
			// query.Set(k, v[len(v)-1]) // last value
		}
		if !isWhiteListed(k, whitelist) {
			query.Del(k)
		}
	}
	r.URL.RawQuery = query.Encode()
}

func isWhiteListed(param string, whitelist []string) bool {
	for _, v := range whitelist {
		if param == v {
			return true
		}
	}
	return false
}


type rateLimiter struct {
	mu        sync.Mutex
	visitors  map[string]int
	limit     int
	resetTime time.Duration
}

func NewRateLimiter(limit int, resetTime time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors:  make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}
	// start the reset routine
	go rl.resetVisitorCount()
	return rl
}

func (rl *rateLimiter) resetVisitorCount() {
	for {
		time.Sleep(rl.resetTime)
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) RLMiddleware(next http.Handler) http.Handler {
	fmt.Println("Rate Limiter Middleware...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Rate Limiter Middleware being returned...")
		rl.mu.Lock()
		defer rl.mu.Unlock()

		visitorIP := r.RemoteAddr // You might want to extract the IP in a more sophisticated way
		rl.visitors[visitorIP]++
		// fmt.Printf("Vistor count from %v is %v\n", visitorIP, rl.visitors[visitorIP])

		if rl.visitors[visitorIP] > rl.limit {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
		fmt.Println("Rate Limiter ends...")
	})
}

