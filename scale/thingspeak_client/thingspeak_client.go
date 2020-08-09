package thingspeak_client

import (
	"encoding/json"
	"net/http"
	"bytes"
)

type ChannelWriter struct {
	key    string
	fields map[int]string
}

func NewChannelWriter(key string) *ChannelWriter {
	w := new(ChannelWriter)
	w.key = key
	w.fields = make(map[int]string)
	return w
}

func (w *ChannelWriter) AddField(n int, value string) {
	w.fields[n] = value
}

func (w *ChannelWriter) Update() (resp *http.Response, err error) {
	requestBody, err := json.Marshal(map[string]string{
		"api_key": w.key,
		"field1":  w.fields[0],
	})
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	r, err := http.NewRequest("POST", "https://api.thingspeak.com/update.json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-type", "application/json")
	resp, err = client.Do(r)
	return resp, err
}
