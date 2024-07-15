#!/bin/bash

echo "-------------"
echo "goose-up.bash"
echo "-------------"

echo "Up migrating..."

( cd ../sql/schema ; goose postgres "postgres://postgres:postgres@localhost:5432/blogator" up )