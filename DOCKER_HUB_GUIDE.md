# Docker Hub Guide for PIXELZX POS EVM Chain

This guide explains how to use the PIXELZX POS EVM Chain Docker images hosted on Docker Hub.

## Available Images

The following Docker images are available on Docker Hub:

- `yuchanshin/pixelzx-evm:latest` - The latest stable release
- `yuchanshin/pixelzx-evm:develop` - The latest development build

## Quick Start

### Pull the image

```bash
docker pull yuchanshin/pixelzx-evm:latest
```

### Run a PIXELZX node

```bash
docker run -d \
  --name pixelzx-node \
  -p 8545:8545 \
  -p 8546:8546 \
  -p 30303:30303 \
  -v $(pwd)/data:/app/data \
  yuchanshin/pixelzx-evm:latest
```

### Initialize the network

```bash
docker run --rm \
  -v $(pwd)/data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx init
```

### Run in validator mode

```bash
docker run -d \
  --name pixelzx-validator \
  -p 8545:8545 \
  -p 8546:8546 \
  -p 30303:30303 \
  -v $(pwd)/data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --validator
```

## Configuration

You can customize the node configuration by mounting a custom configuration file:

```bash
docker run -d \
  --name pixelzx-node \
  -p 8545:8545 \
  -p 8546:8546 \
  -p 30303:30303 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/configs/production.yaml:/app/configs/production.yaml \
  yuchanshin/pixelzx-evm:latest
```

## Docker Compose

For more complex setups, you can use Docker Compose. See [docker-compose.yml](docker-compose.yml) for an example.

### Start a multi-node network

```bash
docker-compose up -d
```

### Start a development environment

```bash
docker-compose -f docker-compose.dev.yml up -d
```

## Ports

The following ports are used by PIXELZX nodes:

- `8545` - JSON-RPC
- `8546` - WebSocket
- `30303` - P2P (TCP and UDP)

## Volumes

The following volumes are recommended:

- `/app/data` - Node data directory
- `/app/configs` - Configuration files
- `/app/keystore` - Keystore files

## Environment Variables

- `CONFIG_ENV` - Set to "development" or "production" to use different config files

## Troubleshooting

### "exec format error"

This error occurs when running images on incompatible architectures. Make sure you're using the correct image for your platform.

### Connection issues

Ensure that the required ports are open and not blocked by a firewall.

### Permission issues

Make sure the data directory has the correct permissions for the container to read and write.