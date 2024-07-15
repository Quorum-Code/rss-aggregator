#!/bin/bash

echo "---------------"
echo "goose-down.bash"
echo "---------------"

echo "Down migrating..."

( cd sql/schema ; goose postgres "postgres://postgres:postgres@localhost:5432/blogator" down )