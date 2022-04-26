# use :

`docker run -d --name=gpstrans \
-p 8006:8006 \
-e traccarserveraddr= your traccar server addr\
--restart=unless-stopped \
gpstrans:latest \
`

# build:

`
docker build -t gpstrans:latest -f Dockerfile .`