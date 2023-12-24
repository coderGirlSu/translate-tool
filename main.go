package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func main() {
	// read all environment variables in .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	http.HandleFunc("/translate", translate)
	http.HandleFunc("/writing", writing)
	http.HandleFunc("/grammar", grammar)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func translate(w http.ResponseWriter, r *http.Request) {
	if !authorized(w, r) {
		return
	}

	// request from apple Shortcuts
	clientReq := r.FormValue("input")
	log.Println(clientReq)
	// clientReq, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	w.WriteHeader(400) // 400 means a problem with client
	// 	w.Write([]byte(err.Error()))
	// }
	prompt := "你是一个短信翻译人员，当你收到英语时，把英语翻译成中文，当你收到中文时，把中文翻译成英语。Your replies should not be wrapped in quotes."

	resp, err := callOpenAI(r.Context(), prompt, clientReq)
	if err != nil {
		sendErrorResponse(err, w)
	}
	sendResponse(resp, w)
}

func grammar(w http.ResponseWriter, r *http.Request) {
	if !authorized(w, r) {
		return
	}
	clientReq := r.FormValue("input")
	fmt.Println(clientReq)

	prompt := "Improve and fix the text grammar I provided to you, your replies should not be wrapped in quotes"

	resp, err := callOpenAI(r.Context(), prompt, clientReq)
	if err != nil {
		sendErrorResponse(err, w)
		return
	}
	sendResponse(resp, w)
}

func writing(w http.ResponseWriter, r *http.Request) {
	if !authorized(w, r) {
		return
	}
	clientReq := r.FormValue("input")
	fmt.Println(clientReq)

	prompt := "Improve the writing, make it sounds more natural and friendly, your replies should not be wrapped in quotes"
	resp, err := callOpenAI(r.Context(), prompt, clientReq)
	if err != nil {
		sendErrorResponse(err, w)
	}
	sendResponse(resp, w)
}

func callOpenAI(ctx context.Context, prompt string, clientReq string) (openai.ChatCompletionResponse, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY")) // make a instance of client

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem, // System is instruction
			Content: prompt,
		},
		{
			Role:    openai.ChatMessageRoleUser, // User is what you want send to chatgpt
			Content: `"` + clientReq + `"`,
		},
	}

	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    messages,
		MaxTokens:   256, // the maximum length of the response
		Temperature: 1,   // how imaginative
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Println(err)
	}
	return resp, nil
}

func sendResponse(resp openai.ChatCompletionResponse, w http.ResponseWriter) {
	// send respond back to Client
	w.WriteHeader(200)
	_, err := w.Write([]byte(resp.Choices[0].Message.Content))
	if err != nil {
		log.Println(err)
		return
	}
}

func sendErrorResponse(errIn error, w http.ResponseWriter) {
	w.WriteHeader(500)
	_, err := w.Write([]byte(errIn.Error()))
	if err != nil {
		log.Println(err)
		return
	}
}

func authorized(w http.ResponseWriter, r *http.Request) bool {
	apikey := r.Header.Get("Authorization")
	if apikey != "Bearer "+os.Getenv("TRANSLATE_API_KEY") {
		sendErrorResponse(errors.New("unauthorised"), w)
		return false
	}
	return true
}
