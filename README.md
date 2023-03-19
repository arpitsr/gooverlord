## Demo Setup

- Install vector from datadog, caddy, docker and docker-compose
- Build image Dockerfile.overlord
    ```$ docker build -t overlord_img:latest -f Dockerfile.overlord .```
- Build image Dockerfile.meili
    ```$ docker build -t meili-img:latest -f Dockerfile.meili .```
- Containers in Docker compose 
    - consul server
    - 2 index server instances - This will register indexer instances to consul service registry
    - 2 overlord instances to consul server

    ```$ docker-compose up```
- Overlord service instances are load balanced through local caddyserver
    ```$ caddy reverse-proxy --from :3000 --to localhost:3001 --to localhost:3002```

- Run vector pipeline (Replicating log scenario for demo purposes)
    - source - syslog (Demo purpose)
    - transform/parse
    - filter
    - sink (http sink to overlord instances through lb)
    ```$ vector -c vector-exp.toml```

