#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username admin --dbname "postgres" <<-EOSQL
    CREATE DATABASE auth_service;
EOSQL