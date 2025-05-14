#!/bin/bash
    echo "Starting Auth Service..."
    ./tmp/auth &
    sleep 2
    
    echo "Starting Equipment Service..."
    ./tmp/equipment &
    sleep 2
    
    echo "Starting Kafka Service..."
    ./tmp/kafka &
    sleep 2
    
    echo "Starting Main Server..."
    ./tmp/server

    wait
