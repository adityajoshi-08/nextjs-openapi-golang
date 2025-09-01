build:
	go build -o nextjs-to-openapi cmd/main.go

run: build
	./nextjs-to-openapi --api-dir ../course-viewer/app/api --model phi3:mini --output openapi.json
