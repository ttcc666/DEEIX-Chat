# 10 — Files, RAG, vectors and Memory

Required chain: `Upload → Validation → Storage → Extraction/OCR → Chunking → Embedding → Vector Persistence → Knowledge Base → Retrieval → Context Injection`.

Audit source areas: `backend/internal/application/embedding/`, `conversation/`, `knowledgebase/`, `memory/`; `backend/internal/infra/extraction/`, `embedding/`, `objectstorage/`, `persistence/sqlitevec/`; model definitions in `models/chat.go`, `knowledgebase.go`, `memory.go`.

Preserve file hashes, ownership, quotas, object keys, MIME detection, processing states, extraction/OCR adapters, chunking, embedding dimensions, vector-space signatures, retrieval filtering/ranking, RAG context injection and memory lifecycle. Embeddings support 64–4096 dimensions. PostgreSQL stores native dimensions and uses 4000-dimension halfvec HNSW candidate retrieval plus full-vector re-ranking; SQLite sqlite-vec uses 4096-dimensional slots. Rebuild tasks need leases, retries, idempotency and stale-worker fencing.

G6 tests: file/hash/object-key readback, OCR/extraction fixtures, dimension mismatch, vector serialization/ranking, signature isolation, quota, task lease/retry/concurrency and memory/context equivalence. Not implemented in .NET.
