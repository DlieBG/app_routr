# app-routr

## run on console
``` 
./routr config.yml
```

> Keep in mind to provide the builded frontend files on ./frontend/dist or ../frontend/dist of the executable \
> The dist folder already contains all files of the latest build

## run on docker-compose
```
sudo docker-compose up
```

> Keep in mind to bind the right config file to /app/config.yml on the container in the docker-compose.yml
 
<br />
<br />

# Configuration
The configuration files defines on the one hand services, which are containing routes, having paths with their associated endpoints.\
On the other hand there are apps, which manaage the assignemend of a service to a socket. These apps differentiate between production and development systems.

## services
``` yml
services:
    - name: sample name
      description: sample description
      host: one.host.com
      hosts: ["list", "of", "hosts"]
      port: 8080
      dev_port: 1234    # port for the dev server; leave empty for only production usage
      app: sample_app   # name of the app, defined in the app section
```

## apps
``` yml
apps:
    - name: sample name
      description: sample description
      routes:
        - name: sample route name
          description: sample route description
          path: /
          endpoint: http://sample.endpoint.com  # don't mind about trailing slashes
```