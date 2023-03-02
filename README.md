# caddy-cfkv-adapter
Using Cloudflare KV to store Caddy server configuration


## Using

using /etc/caddy/Caddyfile for init

### config.json

```json
{
  "caddyfile": "caddyfile path, if the value is empty, use the inner Caddyfile.",
  "id": "The caddy ID, if the value is empty, use the public IP as the value of ID.",
  "api_token": "api token, required",
  "account_id": "account id, required",
  "namespace_id": "kv namespace id, required"
}
```


## 

ID: public IP

```
CADDY_${ID}_SRV_1: {vhost1 caddyfile} 
CADDY_${ID}_SRV_2: {vhost2 caddyfile}
CADDY_${ID}_CADDYFILE: {Caddyfile}
CADDY_${ID}_UDPATEAT: [0]
```

### How to add VHost3

```
CADDY_${ID}_SRV_3: {vhost3 caddyfile}
CADDY_${ID}_UDPATEAT: [now timestamp]
```

### How to update Vhost3

```
CADDY_${ID}_SRV_3: {vhost3 caddyfile}
CADDY_${ID}_UDPATEAT: [now timestamp]
```

### How to delete VHost3

```
DEL -> CADDY_${ID}_SRV_3
CADDY_${ID}_UDPATEAT: [now timestamp]
```

## Sequence diagram

```mermaid
sequenceDiagram
    autonumber
    Note right of caddy: Running

    loop Every 5 minutes
        caddy->>+adapter: get config
        adapter->>+cfkv: cloudflare kv
        cfkv->>-adapter: return
        adapter->>-caddy: convert to json
    end
```
