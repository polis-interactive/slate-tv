
.PHONY: compile
compile:

	protoc ./api/v1/control.proto \
		--proto_path=./api/v1 \
		--go_out=./api/v1/go \
		--go-grpc_out=./api/v1/go \
		--go_opt=paths=source_relative \
        --go-grpc_opt=paths=source_relative

	protoc ./api/v1/control.proto \
		--proto_path=./api/v1 \
		--js_out=import_style=commonjs:./api/v1/web  \
		--grpc-web_out=import_style=typescript,mode=grpcwebtext:./api/v1/web

