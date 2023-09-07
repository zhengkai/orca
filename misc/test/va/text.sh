#!/bin/bash

URL="http://localhost:22035/va/text"

curl -s "$URL" \
	-H "Content-Type: application/json" \
	-H "VA-TOKEN: ${ORCA_VA_TOKEN}" \
	-d '{
	"prompt":"翻译下列语言为中文：\n\n3 2 1 Hello, world!",
	"debug":true
}' | tee tmp-text.json
echo

jq . tmp-text.json
