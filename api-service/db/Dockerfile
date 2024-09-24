# Use the official PostgreSQL image as the base
FROM postgres:15

# Install necessary build dependencies
RUN apt-get update && apt-get install -y \
  build-essential \
  git \
  curl \
  postgresql-server-dev-all \
  && rm -rf /var/lib/apt/lists/*

# Clone pgvector repository
RUN git clone https://github.com/pgvector/pgvector.git /tmp/pgvector

# Build and install pgvector extension
RUN cd /tmp/pgvector && make && make install

# Clean up after installation
RUN rm -rf /tmp/pgvector

