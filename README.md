# x-ui

xray panel supporting multi-protocol multi-user

# Features

System Status Monitoring- 
- Support multi-user multi-protocol, web page visualization operation
- Supported protocols: vmess, vless, trojan, shadowsocks, dokodemo-door, socks, http
- Support for configuring more transport configurations
- Traffic statistics, limit traffic, limit expiration time
- Customizable xray configuration templates
- Support https access panel (self-provided domain name + ssl certificate)
- Support one-click SSL certificate application and automatic renewal
- For more advanced configuration items, please refer to the panel

# Install & Upgrade

```
bash <(curl -Ls https://raw.githubusercontent.com/blackestwhite/x-ui/master/install.sh)
```

## Manual install & upgrade

1. First download the latest compressed package from https://github.com/blackestwhite/x-ui/releases
2. Then upload the compressed package to the server's `/root/`, then login to the server as `root` user

> If your server cpu architecture is not  `amd64`，replace `amd64` with another arch

```
cd /root/
rm x-ui/ /usr/local/x-ui/ /usr/bin/x-ui -rf
tar zxvf x-ui-linux-amd64.tar.gz
chmod +x x-ui/x-ui x-ui/bin/xray-linux-* x-ui/x-ui.sh
cp x-ui/x-ui.sh /usr/bin/x-ui
cp -f x-ui/x-ui.service /etc/systemd/system/
mv x-ui/ /usr/local/
systemctl daemon-reload
systemctl enable x-ui
systemctl restart x-ui
```

## install using docker

> This docker tutorial and docker image are provided by [Chasing66](https://github.com/Chasing66)提供

1. install docker

```shell
curl -fsSL https://get.docker.com | sh
```

2. install x-ui

```shell
mkdir x-ui && cd x-ui
docker run -itd --network=host \
    -v $PWD/db/:/etc/x-ui/ \
    -v $PWD/cert/:/root/cert/ \
    --name x-ui --restart=unless-stopped \
    enwaiax/x-ui:latest
```

> Build image

```shell
docker build -t x-ui .
```

## SSL certificate

> This feature and tutorial are provided by [FranzKafkaYu](https://github.com/FranzKafkaYu)提供

The script has a built-in SSL certificate application function. To use this script to apply for a certificate, the following conditions must be met:

- Know the Cloudflare registered email
- Know the Cloudflare Global API Key
- The domain name has been resolved to the current server through cloudflare

How to get the Cloudflare Global API Key:
    ![](media/bda84fbc2ede834deaba1c173a932223.png)
    ![](media/d13ffd6a73f938d1037d0708e31433bf.png)

When using, just enter `域名`, `邮箱`, `API KEY` and the schematic diagram is as follows:
        ![](media/2022-04-04_141259.png)

Precautions:

- The script uses DNS API for certificate request
- By default, Let'sEncrypt is used as the CA party
- The certificate installation directory is the /root/cert directory
- The certificates applied for by this script are all generic domain name certificates

## suggessted distro

- CentOS 7+
- Ubuntu 16+
- Debian 8+

# common problems

## issue closed

All kinds of small white problems see high blood pressure

## API guide

pass username in `x-api-username` and password in `x-api-password` headers.

available routes:
- POST /xui/api/inbound/add
- POST /xui/api/inbound/list
- POST /xui/api/inbound/del/:id
- POST /xui/api/inbound/update/:id

### /xui/api/inbound/add
post body:
```json
{
    "up": 0,
    "down": 0,
    "total": 0,
    "remark": "customer name",
    "enable": true,
    "expiryTime": 0,
    "listen": null,
    "port": 6942,
    "protocol": "vmess",
    "settings": "{
        \"clients\": [
            {
                \"id\": \"uuid-v4-id-preferably\",
                \"alterId\": 0,
            }
        ],
        \"disableInsecureEncryption\": false
    }",
    "streamSettings": "{
        \"network\": \"ws\",
        \"security\": \"none\",
        \"wsSettings\": {
            \"path\": \"/\",
            \"headers\": {}
        }
    }",
    "sniffing": "{
        \"enabled\": true,
        \"destOverride\": [
            \"http\",
            \"tls\"
        ]
    }"
}
```