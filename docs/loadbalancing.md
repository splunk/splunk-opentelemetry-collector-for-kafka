## Load balancing

### Change in Load Balancing Strategy for Splunk HEC

Unlike the previous SC4Kafka connector, SOC4Kafka collector delegates the responsibility of load balancing and high availability for Splunk HEC endpoints to dedicated infrastructure components. This aligns with modern architectural best practices and provides a more scalable and resilient solution than client-side logic.

The collector should be configured with a single HEC endpoint. In a multi-indexer environment, this endpoint must be the address of a load balancer.

### Recommended architecture

The standard architecture involves placing an external load balancer in front of your Splunk indexer pool. This centralizes traffic management, health checks, and failover logic.

### Implementation Example: Using Nginx

Nginx is a lightweight, high-performance, and popular choice for this role. Splunk provides [an official, step-by-step guide for this exact use case](https://dev.splunk.com/enterprise/docs/devtools/httpeventcollector/confignginxloadhttp/), which we recommend following for production deployments.

Below is a minimal working example `nginx.conf` suitable for testing or development environments.
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

!!!note
    For development purposes, you can generate a self-signed certificate with the following command. Note that clients connecting to Nginx will need to be configured to trust this certificate or skip verification.
    `bash sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /etc/nginx/ssl/your_domain.key -out /etc/nginx/ssl/your_domain.crt`