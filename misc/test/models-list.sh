#!/bin/bash

OPENAI_API_KEY="sk-rhjeVT1fkcuarBKnQR6ST$(cat ~/.config/openai)"

curl -s https://api.openai.com/v1/models \
	-H "Authorization: Bearer ${OPENAI_API_KEY}" \
	| jq . -
