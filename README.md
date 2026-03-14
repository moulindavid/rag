#wip


docker run -d --name pgvector -e POSTGRES_PASSWORD=rag -e POSTGRES_DB=rag -p 5432:5432 pgvector/pgvector:pg16
ollama pull nomic-embed-text
ollama pull llama3


go run compile
