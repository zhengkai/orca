#!/bin/bash -ex

BASE="${OPENAI_API_BASE:-https://api.openai.com/v1}"
curl "${BASE}/completions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${OPENAI_API_KEY}" \
  -d '{
    "model": "text-davinci-003",
    "prompt": "Say this is a test",
    "max_tokens": 7,
    "temperature": 0
  }'
