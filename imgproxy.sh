 #!/bin/bash

docker run -p 8000:8080 \
    -v "$PWD/web/images:/images" \
    -e IMGPROXY_LOCAL_FILESYSTEM_ROOT=/images \
    -e IMGPROXY_ENABLE_WEBP_DETECTION=true \
    -e IMGPROXY_ENABLE_AVIF_DETECTION=false \
    -e IMGPROXY_DEVELOPMENT_ERRORS_MODE=false \
    darthsim/imgproxy
