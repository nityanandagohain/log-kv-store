
.PHONY: apigen 

apigen:
	jq -rs 'reduce .[] as $$item ({}; . * $$item)' ./api/*.json > apigen/result.json
	cat apigen/result.json | gojsontoyaml  > apigen/api.yaml
	oapi-codegen -package apigen apigen/api.yaml > apigen/server.gen.go
	# oapi-codegen -generate spec -package apigen apigen/api.yaml > apigen/spec.gen.go
	# oapi-codegen -generate types -package apigen  apigen/api.yaml > apigen/types.gen.go
	rm apigen/result.json