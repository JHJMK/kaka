.PHONY: api
api:
	cd kaka/cmd/api && goctl api go -api *.api -dir ../  --style=goZero