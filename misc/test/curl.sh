#!/bin/bash

curl http://localhost:21035/v1/engines/text-embedding-ada-002/embeddings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d '{"input": ["\u80fd\u91cf\u793c\u7269\u662f\u600e\u4e48\u56de\u4e8b\uff1f\u7528\u4e2d\u6587"], "encoding_format": "base64"}'
