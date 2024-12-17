#!/bin/bash
pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB" -h localhost || exit 1
