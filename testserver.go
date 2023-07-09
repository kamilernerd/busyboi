package main

import (
	"net/http"
	"os"
)

func testserver() {
	dom := `
<html lang="en">
<head>
		<meta charset="UTF-8">
		<title>Title</title>
</head>
<body>
		<div class="images">
			<img class="product_image" src="./testdata/randomimage.jpg" alt="Image hi">
			<img class="product_image" src="./testdata/randomimage.jpg" alt="Image hi">
		</div>
		<h1 class="product_title">Some random item title</h1>
		<p class="product_price">5555.12 kr</p>
		<a class="product_url" href="www.google.com">Some link</a>
</body>
</html>`

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(dom))
	})

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		f, _ := os.ReadFile("./testdata/robots.txt")
		w.Write([]byte(f))
	})

	http.ListenAndServe("localhost:8080", nil)
}
