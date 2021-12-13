
gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:pb

clean:
	rm pb/*.go

mongo-init:
	docker run -d --name mongo -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -e MONGO_INITDB_DATABASE=dev mongo

client-grpc-create:
	go run client/. create -a “Nico Vergauwen” -t “Learning Go” -c “The quick brown fox jumps over the lazy dog”

client-grpc-list:
	go run client/. list
