package ratelimit

import (
	"RateLimiter/fileReader"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)



func RateApiLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    limits := fileReader.GetLimits();
		route := mux.CurrentRoute(r)
        pathTemplate, err := route.GetPathTemplate()
		
        if err != nil {
            fmt.Println("No route matched.")
            next.ServeHTTP(w, r)
            return
        }

        endpoint := fmt.Sprintf("%s %s", r.Method, pathTemplate)
        limitInfo, ok := limits[endpoint];

        if !ok || limitInfo == nil {
            fmt.Println("No rate limit info found for", endpoint)
            next.ServeHTTP(w, r)
            return
        }

        fileReader.RefillToken(limitInfo)
		
        if limitInfo.Remaining > 0 {
			limitInfo.Remaining -= 1;
			rb := &fileReader.ResponseBuffer{ResponseWriter: w}
			next.ServeHTTP(rb, r)
			modifiedContent := fmt.Sprintf("Remaining: %d. %s", limitInfo.Remaining, rb.Buffer.String())
			rb.Buffer.Reset()
			rb.Buffer.WriteString(modifiedContent)
			w.Write(rb.Buffer.Bytes())
            
        }else{
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
		}
        
    })
}

