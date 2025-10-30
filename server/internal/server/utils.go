package server

import (
	"io"
	"fmt"
	"log"
	"strings"
	"unicode"
	"net/http"
	"unicode/utf8"

	"github.com/JValtteri/qure/server/internal/utils"
	ware "github.com/JValtteri/qure/server/internal/middleware"
)

func sendJsonResponse[O any](w http.ResponseWriter, event O) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%v", utils.UnloadJSON(event))
}

func sanitize(input string) string {
	var result strings.Builder
	for i := 0; i < len(input); {
		r, size := utf8.DecodeRuneInString(input[i:])
		if unicode.IsSpace(r) || unicode.IsLetter(r) || unicode.IsDigit(r) || r=='-' {
			result.WriteRune(r)
			i += size
		} else {
			i++
		}
	}
	return strings.ToLower(result.String())
}

func loadRequestBody [R ware.Request](request *http.Request, obj R) (R, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("Error closing request body %v\n", err)
	} else {
		defer close(request.Body)
		utils.LoadJSON(body, &obj)
	}
	return obj, err
}

func close(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		log.Printf("Error closing request body %v\n", err)
	}
}
