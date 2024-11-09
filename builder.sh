go build -o dial ./main/client/main.go
go build -o serve ./main/server/main.go


echo run ./dial ./serve