package upload

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

    const maxUploadSize = 10 << 20 // 10 MB
    r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

    if err := r.ParseMultipartForm(maxUploadSize); err != nil {
        http.Error(w, "File too large or invalid form data", http.StatusBadRequest)
        return
    }

    file, fileHeader, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "Unable to retrieve the file from form data", http.StatusBadRequest)
        return
    }
    defer file.Close()

    fileName := fileHeader.Filename
    mimeType := fileHeader.Header.Get("Content-Type")

    log.Printf("Received file: %s (%s)", fileName, mimeType)

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"message":"Uploaded file: %s (%s)"}`, fileName, mimeType)

    fileBytes, err := io.ReadAll(file)
    if err != nil {
        http.Error(w, "Failed to read file content", http.StatusInternalServerError)
        return
    }

    encoded := base64.StdEncoding.EncodeToString(fileBytes)

    url := "http://127.0.0.1:2024/runs/wait"

    // Escape the encoded string properly for JSON
    encodedJSON, _ := json.Marshal(encoded)

    payloadStr := fmt.Sprintf(`{
		"assistant_id": "fe096781-5601-53d2-b2f6-0d3403f7e9ca",
		"input": {
			"file_content": %s
		}
	}`, string(encodedJSON))

    req, _ := http.NewRequest("POST", url, strings.NewReader(payloadStr))
    req.Header.Add("Content-Type", "application/json")

    res, _ := http.DefaultClient.Do(req)
    defer res.Body.Close()

    body, _ := io.ReadAll(res.Body)
    fmt.Println(res)
    fmt.Println(string(body))
	fmt.Println("\n\n\n\n\n")
}

//   "metadata": {},
//   "config": {
//     "tags": [""],
//     "recursion_limit": 1,
//     "configurable": {}
//   },
//   "interrupt_before": "*",
//   "interrupt_after": "*",
//   "stream_mode": ["values"],
//   "feedback_keys": [""],
//   "stream_subgraphs": false,
//   "on_completion": "delete",
//   "on_disconnect": "cancel",
//   "after_seconds": 1,
//   "checkpoint_during": false
//   "command": {
//     "update": null,
//     "resume": null
//   },
// "goto": {
//   "node": "chatbot",
//   "input": null
// }
