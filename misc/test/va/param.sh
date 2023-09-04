#!/bin/bash

URL="http://localhost:22035/va/chat"

curl -s "$URL" \
	-H "Content-Type: application/json" \
	-H "VA-TOKEN: ${ORCA_VA_TOKEN}" \
	-d '{
	"system":"翻译下列语言为中文：",
	"user":"Hello, world!",
	"param":{
		"temperature":0.3,
		"maxOutputTokens": 100,
		"topP": 1,
		"topK": 40
	},
	"debug":true
}' | tee tmp-chat.json
echo

jq . tmp-chat.json
