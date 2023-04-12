# !/usr/bin/bash
# Clean
go clean
# Remove cache directories
rm -rf 'bin'
rm -rf 'log'
rm -rf 'export'
rm -rf 'training'
# Build in platform
go build -o ./bin/caos