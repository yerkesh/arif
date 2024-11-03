package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"arif/config"
)

const instructionsToExtract = `Extract all text from this image, preserving the original line breaks and paragraphs, even if the text is incomplete due to page divisions or formatting issues. All extracted content should maintain the plain text format.
# Steps 

1. Identify and access the text content from the provided document or source.
2. Extract the text content, ensuring that all lines and paragraphs are captured completely.
3. Retain the original line breaks and paragraphs to preserve the structure and formatting of the original text.
4. Acknowledge and extract partial sentences as they appear, without altering the content, to ensure absolutely all text is included.

# Output Format
The output should be plain text with preserved line breaks and paragraph spacing, reflecting the original structure of the extracted content.
Only Turkish headers should be wrapped with **.

# Notes
- File might contain arabic text, but main language is turkish, include both in the output.
- Ensure that even partial text caused by formatting issues is included.
- Maintain accuracy in preserving the text's original presentation as much as possible.
- If the text is incomplete, set isComplete to false.
`
const instructionsTranslate = `Твоя роль - Профессиональный переводчик, лингвист, полиглот. красноречивый поэт. Твои области это Турецкий язык и Русский язык. Ты очень хорошо понимаешь языковые особенности как турецкого так и русского языка. Тебе дана текущая страница, предедущая и следующая страница. Переведи на русский язык текущую страницу. При этом сохрани особенности русского языка. Текст на начале или конце страницы может быть обрезан, для этого используй пред/след страницу`

func ExtractFromImage(ctx context.Context, imageURL string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []interface{}{
					map[string]string{
						"type": "text",
						"text": instructionsToExtract,
					},
					map[string]interface{}{
						"type": "image_url",
						"image_url": map[string]interface{}{
							"url": imageURL,
						},
					},
				},
			},
		},
	})

	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Cfg.ChatGPTKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	fmt.Println(result)

	// Parse the response content here
	if choices, ok := result["choices"].([]interface{}); ok {
		firstChoice := choices[0].(map[string]interface{})
		message := firstChoice["message"].(map[string]interface{})
		return message["content"].(string), nil
	}

	return "", fmt.Errorf("failed to extract text")
}

func generatePageMessages(pages map[int]string) map[int]string {
	// Collect and sort the existing page numbers
	var pageNumbers []int
	for pageNum := range pages {
		pageNumbers = append(pageNumbers, pageNum)
	}
	sort.Ints(pageNumbers)

	// Initialize the result map
	messages := make(map[int]string)

	// Iterate over the sorted page numbers
	for idx, pageNum := range pageNumbers {
		pageContent := pages[pageNum]

		var prevContent, nextContent string

		// Check if previous page exists in the sorted list
		if idx > 0 {
			prevPageNum := pageNumbers[idx-1]
			prevContent = pages[prevPageNum]
		}

		// Check if next page exists in the sorted list
		if idx < len(pageNumbers)-1 {
			nextPageNum := pageNumbers[idx+1]
			nextContent = pages[nextPageNum]
		}

		// Build the message
		message := fmt.Sprintf("Текущая страница:\n%s", pageContent)

		if prevContent != "" {
			message += fmt.Sprintf("\nПредыдущая страница:\n%s", prevContent)
		}

		if nextContent != "" {
			message += fmt.Sprintf("\nСледующая страница:\n%s", nextContent)
		}

		// Store the message in the result map with the page number as the key
		messages[pageNum] = message
	}

	return messages
}

func GenerateTranslateMessage(ctx context.Context, pages map[int]string) (map[int]string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	res := make(map[int]string)

	messages := generatePageMessages(pages)
	for page, message := range messages {
		requestBody, _ := json.Marshal(map[string]interface{}{
			"model": "gpt-4o",
			"messages": []map[string]interface{}{
				{
					"role":    "user",
					"content": message,
				},
				{
					"role":    "system",
					"content": instructionsTranslate,
				},
			},
		},
		)

		req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+config.Cfg.ChatGPTKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("error decoding response: %v", err)
		}

		fmt.Println(result)

		// Parse the response content here
		if choices, ok := result["choices"].([]interface{}); ok {
			firstChoice := choices[0].(map[string]interface{})
			message := firstChoice["message"].(map[string]interface{})
			res[page] = message["content"].(string)
		}
	}

	return res, nil
}
