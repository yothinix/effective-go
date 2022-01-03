#!/bin/bash

echo "GET http://:8081/limitz" | vegeta attack -rate=10/s -duration=1s | vegeta report
