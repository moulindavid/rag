# Database Migrations

## 001 — Documents table

Enables the pgvector extension (so PostgreSQL can store and search vectors), then creates the `documents` table. Each row represents an uploaded file with a unique ID, its filename, and when it was added.

## 002 — Chunks table

Creates the `chunks` table. When a document is ingested, its text gets split into smaller pieces (chunks). Each chunk stores:
- A link back to its parent document
- The text content
- A 768-dimension vector embedding (the numeric representation used for similarity search)

If a document is deleted, all its chunks are automatically removed too (cascade delete).

Also creates an IVFFlat index on the embeddings to make similarity searches fast. Think of it as organizing vectors into 100 buckets so PostgreSQL doesn't have to compare against every single row.
