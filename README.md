# Concurrent Manga Data Pipeline & Streaming API

## 1. Context & Motivation

Traditional manga reading applications often rely on slow, monolithic architectures that link out to fragmented third party sources,
resulting in broken images and high latency.

This project solves that by implementing a concurrent data ingestion pipeline written in Go. Taking the core mechanics of concurrent
image fetching and scaling them into a continuous backend service, this system aggregates chapter metadata across multiple domains,
normalizes the data and streams image bytes directly to a lightweight React frontend. It demonstrates a clear shift away from heavy
JVM based frameworks toward idiomatic Go concurrency.

## 2. Goals 

* High Throughput Ingestion: Implement bounded worker pools and Fan-in / Fan-out channel patterns to concureently scrape 3+ external
sources without triggering rate limits or memory leaks.

* Resilient Error Handling: Utilize context.Context for strict network timeouts and dead letter queues to ensure a single failed external
request does not crash the pipeline.

* Optimizd I/O Streaming: Fetch and pipe large image byte streams directly to the client using Go's io.Reader and io.Writer interfaces,
minimizing RAM overhead.

* Database Efficiency: Perform bulk insert operations into PostgreSQL to handle rapid concurrent writes from the ingestion workers safely.

## 3. Non Goals

* No Recommendation Algorithms: The system will display chronological releases and explicit search results, not predictive suggestions.
* No Advanced Anti-Bot Bypassing: We will target sources with accesible API's or simple DOM structures.

## 4. Tech Stack

* Backend: Go (Goroutines, Channels, standard library net/http)
* Database: PostgreSQL
* Frontend: React
* Infrastructure: Docker