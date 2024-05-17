## Running oasdiff from docker
To run oasdiff from docker just replace the oasdiff command by `docker run --rm -t tufin/oasdiff` (or the image you prefer).  
Here are a few examples:

### Breaking changes with Docker
```bash
docker run --rm -t tufin/oasdiff breaking https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test1.yaml https://raw.githubusercontent.com/Tufin/oasdiff/main/data/openapi-test3.yaml
```

### Comparing local files with Docker
```bash
docker run --rm -t -v $(pwd)/data:/data:ro tufin/oasdiff diff /data/openapi-test1.yaml /data/openapi-test3.yaml
```

Replace `$(pwd)/data` by the path that contains your files.  
Note that the spec paths must begin with `/`.  