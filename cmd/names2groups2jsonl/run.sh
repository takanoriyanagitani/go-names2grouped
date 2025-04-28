#!/bin/sh

export ENV_COMMON_PREFIX_LEN=4

printf '%s\n' \
	'helo, wrld' \
	'hello, world' \
	'hell, world' \
	'unsorted sample 1' \
	'unsorted sample 2' \
	'pre-sort required' \
	'unsorted sample 3' \
	'unsorted sample 4' |
	./names2groups2jsonl |
	dasel --read=json --write=yaml |
	bat --language=yaml
