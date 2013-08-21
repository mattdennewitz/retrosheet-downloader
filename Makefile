rs:
	go build rs.go

clean:
	go clean .

all:
	rm -f rs
	make clean
	make rs
