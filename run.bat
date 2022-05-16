go build -o resvbooking cmd/web/main.go mailchannel.go middlewares.go routes.go
./resvbooking -dbname=postgres -dbuser=postgres -inproduction=false -usecache=false