rs:
	go build rs.go

clean:
	go clean .

all:
	make clean
	make rs
