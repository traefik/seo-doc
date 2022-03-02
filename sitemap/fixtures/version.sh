#!/bin/sh

# Clean data

rm -rf {traefik,traefik-enterprise,traefik-mesh}

# Generate data

mkdir -p traefik/{v1.0,v1.1,v1.2,v1.3,v1.4,v1.5,v1.6,v1.7,v2.0,v2.1,v2.2,v2.3,v2.4,v2.5,v2.6,master}
touch traefik/index.html
touch traefik/{v1.0,v1.1,v1.2,v1.3,v1.4,v1.5,v1.6,v1.7,v2.0,v2.1,v2.2,v2.3,v2.4,v2.5,v2.6,master}/index.html

mkdir -p traefik-enterprise/{v1.0,v1.1,v1.2,v1.3,v2.0,v2.1,v2.2,v2.3,v2.4,v2.5,master}
touch traefik-enterprise/index.html
touch traefik-enterprise/{v1.0,v1.1,v1.2,v1.3,v2.0,v2.1,v2.2,v2.3,v2.4,v2.5,master}/index.html

mkdir -p traefik-mesh/{v1.0,v1.1,v1.2,v1.3,v1.4,master}
touch traefik-mesh/index.html
touch traefik-mesh/{v1.0,v1.1,v1.2,v1.3,v1.4,master}/index.html
