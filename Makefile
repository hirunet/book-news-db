run:
	go build
	./book-news-db

db:
	psql $(DATABASE_URL)
