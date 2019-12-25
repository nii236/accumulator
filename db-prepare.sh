#!/bin/bash

./bin/migrate -path ./migrations/  -database sqlite3://accumulator.db up
./bin/sqlboiler ./bin/sqlboiler-sqlite3