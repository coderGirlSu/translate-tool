package main

import (
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
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func translate(w http.ResponseWriter, r *http.Request) {
	// request from apple Shortcuts
	clientReq := r.FormValue("input")
	fmt.Print(clientReq)
	// clientReq, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	fmt.Print(err)
	// 	w.WriteHeader(400) // 400 means a problem with client
	// 	w.Write([]byte(err.Error()))
	// }

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY")) // make a instance of client
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "你是一个短信翻译人员，当你收到英语时，把英语翻译成中文，然后我回复中文，你把中文翻译成英语发送出去",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: clientReq,
		},
	}
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    messages,
		MaxTokens:   256, // the maximum length of the response
		Temperature: 1,   // how imaginative
	}

	resp, err := client.CreateChatCompletion(r.Context(), req)
	if err != nil {
		fmt.Print(err)
		return
	}

	// send respond back to Client
	w.WriteHeader(200)
	_, err = w.Write([]byte(resp.Choices[0].Message.Content))

	if err != nil {
		fmt.Print(err)
		return
	}
}
