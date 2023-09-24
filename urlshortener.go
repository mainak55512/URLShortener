package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	// "log"
	"math/rand"
	"net/http"
)

type PostPayload struct {
	Short_key string `json:"short_key"`
	Url       string `json:"url"`
}

type GetPayload struct {
	Short_key string `json:"short_key"`
}

type GetPayloadResponse struct {
	OriginalURL string `json:"originalURL"`
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handleForm)
	http.HandleFunc("/shorten", HandleShorten)
	http.HandleFunc("/short/", HandleRedirect)

	fmt.Println("URL Shortener is running on :"+port)
	http.ListenAndServe(":"+port, nil)

}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 10

	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/shorten", http.StatusSeeOther)
		return
	}

	// Serve the HTML form
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>URL Shortener</title>
		</head>
		<body>
			<h2>URL Shortener</h2>
			<form method="post" action="/shorten">
				<input type="url" name="url" placeholder="Enter a URL" required>
				<input type="submit" value="Shorten">
			</form>
		</body>
		</html>
	`)
}

func HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	// Generate a unique shortened key for the original URL
	shortKey := generateShortKey()

	//payload creation
	payload := PostPayload{
		Short_key: shortKey,
		Url:       originalURL,
	}

	payloadJSON, _ := json.MarshalIndent(payload, "", " ")

	_, err := http.Post("https://dev174161.service-now.com/api/692302/shorturl_store/store", "application/json", bytes.NewBuffer(payloadJSON))

	if err != nil {
		fmt.Print(err)
	}

	// Construct the full shortened URL
	shortenedURL := fmt.Sprintf("https://urlshortener-guli.onrender.com/short/%s", shortKey)

	// Render the HTML response with the shortened URL
	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
        <h2>URL Shortener</h2>
        <p>Original URL: %s</p>
        <p>Shortened URL: <a href="%s">%s</a></p>
        <form method="post" action="/shorten">
            <input type="text" name="url" placeholder="Enter a URL">
            <input type="submit" value="Shorten">
        </form>
    `, originalURL, shortenedURL, shortenedURL)
	fmt.Fprintf(w, responseHTML)
}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	shortKey := r.URL.Path[len("/short/"):]
	fmt.Println(shortKey)
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	// Retrieve the original URL from the `urls` map using the shortened key
	payload := GetPayload{
		Short_key: shortKey,
	}
	payloadJSON, _ := json.MarshalIndent(payload, "", "	")
	resp, err := http.Post("https://dev174161.service-now.com/api/692302/shorturl_store/fetch", "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		fmt.Print(err)
	}

	stResp, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	var responsee GetPayloadResponse
	json.Unmarshal(stResp, &responsee)

	if responsee.OriginalURL == "" {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}
	fmt.Println(responsee.OriginalURL)
	// Redirect the user to the original URL
	http.Redirect(w, r, responsee.OriginalURL, http.StatusMovedPermanently)

}
