# Load balance Splunk HEC traffic

## Configure a load balancer for Splunk HEC

Splunk Connect for Kafka (SC4Kafka) includes client-side load balancing. The Splunk Distribution of OpenTelemetry Collector for Kafka delegates load balancing and high availability for Splunk HTTP Event Collector (HEC) endpoints to dedicated infrastructure components. This separates traffic management, health checks, and failover from the collector. In a multi-indexer environment, configure the collector with a single HEC endpoint that points to an external load balancer.

## Place a load balancer in front of the indexers

Place an external load balancer in front of the Splunk indexer pool to manage traffic, health checks, and failover.

## Configure Nginx as a load balancer

Nginx is one option for load balancing. For production deployments, follow Splunk's [guide to configuring Nginx to load balance HTTP Event Collector traffic](https://dev.splunk.com/enterprise/docs/devtools/httpeventcollector/confignginxloadhttp/).

The following `nginx.conf` example is for development and testing:
```
events {
    worker_connections 1024;
}

http {
    upstream hec {
        # List your Splunk indexers running HEC.
        server 10.236.10.37:8088;
        server 10.236.10.142:8088;
    }

    server {
        listen 8088 ssl;

        # --- IMPORTANT: REPLACE WITH YOUR CERTIFICATE PATHS ---
        ssl_certificate     /etc/nginx/ssl/your_domain.crt;
        ssl_certificate_key /etc/nginx/ssl/your_domain.key;

        location / {
            proxy_connect_timeout 1s;
            # Proxy requests to the 'hec' upstream group.
            # This assumes your backend HEC endpoints are using SSL (https).
            # If they use plain HTTP, change this to 'http://hec'.
            proxy_pass https://hec;
        }
    }
}
```

!!! note

    For development purposes, you can generate a self-signed certificate with the following command. Clients connecting to Nginx must trust this certificate or skip verification.

```bash
sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /etc/nginx/ssl/your_domain.key -out /etc/nginx/ssl/your_domain.crt
```
