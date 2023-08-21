#!/bin/bash -ex

BASE="${OPENAI_API_BASE:-https://api.openai.com/v1}"

curl -v "${BASE}/chat/completions" \
	-w "Total time: %{time_total} seconds\n" \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer ${OPENAI_API_KEY}" \
	-d '{
	"model": "gpt-3.5-turbo",
	"messages": [{"role": "user", "content": "Hello!"}]
}'
