clear && printf '\e[3J'
GOCACHE=off go test -cover -coverprofile cover.out -v ./
go tool cover -html=cover.out -o cover.html
open cover.html
