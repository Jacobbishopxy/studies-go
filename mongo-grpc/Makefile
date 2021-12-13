
proto-gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:pb

proto-clean:
	rm pb/*.go

mongo-init:
	docker run -d --name mongo -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -e MONGO_INITDB_DATABASE=dev mongo

server-start:
	cd server && go run .

client-create:
	cd client && go run . create -a "Nico Vergauwen" -t "Learning Go" -c "The quick brown fox jumps over the lazy dog"

client-list:
	cd client && go run . list
