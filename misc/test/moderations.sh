#!/bin/bash

BASE="${OPENAI_API_BASE:-https://api.openai.com/v1}"

curl -v "${BASE}/moderations" \
	-w "Total time: %{time_total} seconds\n" \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer ${OPENAI_API_KEY}" \
	-d '{"input": "你媽逼啊"}'
