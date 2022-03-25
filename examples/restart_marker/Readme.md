# Dockerfile for restart markers test

Command should run from project root directory

# Build container
```
docker build -t rst-mrk -f ./examples/restart_marker/Dockerfile .
docker run --rm --mount type=bind,source="$(pwd)/examples/restart_marker/result",target="/restart_marker/result" -it rst-mrk
```

# Check
ffd[0..7] - Restart markers (https://dev.exiv2.org/projects/exiv2/wiki/The_Metadata_in_JPEG_files)
```dockerfile
xxd ./examples/restart_marker/result/rst-mrk-output-govips.jpeg > jpgraw
cat jpgraw | grep ffd0
```
