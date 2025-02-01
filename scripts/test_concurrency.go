package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,min=3,max=200"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func updatePost(postId int, p UpdatePostPayload, wg *sync.WaitGroup) {
	defer wg.Done()

	// construct the url to update post endpoint
	url := fmt.Sprintf("http://localhost:1414/v1/posts/%d", postId)

	// Create the JSON payload
	b, _ := json.Marshal(p)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while sending request :", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Update post response : ", resp.Status)
}

func main() {
	var wg sync.WaitGroup

	// Assuming the postID to update is 2
	postId := 2

	// simulate UserA and UserB updating the same post
	wg.Add(2)
	content := "new Content from user B"
	title := "new Title from User A"

	go updatePost(postId, UpdatePostPayload{Content: &content, Title: &title}, &wg)
	go updatePost(postId, UpdatePostPayload{Content: &content, Title: &title}, &wg)
	wg.Wait()
}
