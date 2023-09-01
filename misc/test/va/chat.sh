#!/bin/bash

URL="http://localhost:22035/va/chat"

curl -s "$URL" \
	-H "Content-Type: application/json" \
	-H "VA-TOKEN: ${ORCA_VA_TOKEN}" \
	-d '{
	"system":"翻译下列语言为中文：",
	"user":"Hello, world!",
	"debug":true
}' | tee tmp-chat.json
echo

jq . tmp-chat.json
