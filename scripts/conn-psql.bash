#!/bin/bash

echo "--------------"
echo "conn-psql.bash"
echo "--------------"

echo "Connecting to psql..."

psql "postgres://postgres:postgres@localhost:5432/blogator"
