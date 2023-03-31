#!/bin/bash

API_HOST="http://10.0.84.49:22035"
# API_HOST="http://localhost:22035"

curl "${API_HOST}/v1/engines/text-embedding-ada-002/embeddings" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d '{"input":[  "\u80fd\u91cf\u793c\u7269\u662f\u600e\u4e48\u56de\u4e8b\uff1f\u7528\u4e2d\u6587"], "encoding_format": "base64"}'
