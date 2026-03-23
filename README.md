#wip


docker run -d --name pgvector -e POSTGRES_PASSWORD=rag -e POSTGRES_DB=rag -p 5432:5432 pgvector/pgvector:pg16
ollama pull nomic-embed-text
ollama pull llama3


go build ./...

curl -X POST http://localhost:8080/documents -F "file=@test2.txt"

docker exec pgvector psql -U postgres -d rag -c "SELECT d.filename, count(c.id) as chunks FROM documents d JOIN chunks c ON c.document_id = d.id GROUP BY d.filename;"
  
docker exec pgvector psql -U postgres -d rag -c "DELETE FROM chunks; DELETE FROM documents;"
