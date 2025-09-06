# PIXELZX POS EVM Chain

PIXELZX POS EVM Chain is a Proof of Stake (PoS) based Ethereum Virtual Machine (EVM) compatible blockchain network with PIXELZX as its native token.

## Table of Contents

- [Key Features](#key-features)
- [Token Specifications](#token-specifications)
- [Network Parameters](#network-parameters)
- [Architecture](#architecture)
- [Directory Structure](#directory-structure)
- [Build and Run](#build-and-run)
  - [Dependencies](#dependencies)
  - [Local Build](#local-build)
  - [Docker Build](#docker-build)
  - [Execution](#execution)
- [**Docker Quick Start**](#docker-quick-start) 🚀
  - [Basic Docker Commands](#basic-docker-commands)
  - [Environment Variable Setup](#environment-variable-setup)
  - [Volume Mount Guide](#volume-mount-guide)
  - [Health Check and Status Verification](#health-check-and-status-verification)
- [**CLI Commands Reference**](#cli-commands-reference) ⌨️
  - [Global Options](#global-options)
  - [Core Commands](#core-commands)
  - [Docker Commands](#docker-commands)
  - [Account Management](#account-management)
  - [Node Management](#node-management)
  - [Network Commands](#network-commands)
  - [Staking Commands](#staking-commands)
  - [Validator Commands](#validator-commands)
  - [Governance Commands](#governance-commands)
  - [Configuration Commands](#configuration-commands)
  - [Admin Commands](#admin-commands)
- [**P2P Network Connection**](#p2p-network-connection) 🌐
  - [P2P Port Setup](#p2p-port-setup)
  - [Bootnode Connection](#bootnode-connection)
  - [Network Status Monitoring](#network-status-monitoring)
  - [P2P Connection Troubleshooting](#p2p-connection-troubleshooting)
- [**Node Initialization and Configuration**](#node-initialization-and-configuration) ⚙️
  - [Genesis File Initialization](#genesis-file-initialization)
  - [Data Directory Setup](#data-directory-setup)
  - [Customizing Configuration Files](#customizing-configuration-files)
  - [Keystore Management](#keystore-management)
  - [Initialization Verification](#initialization-verification)
- [API Endpoints](#api-endpoints)
- [**Troubleshooting**](#troubleshooting) 🚑
  - [Permission Issues Resolution](#permission-issues-resolution)
  - [Docker Related Issues](#docker-related-issues)
  - [P2P Network Issues](#p2p-network-issues)
  - [API Connection Issues](#api-connection-issues)
  - [Performance Issues](#performance-issues)
  - [Log Analysis](#log-analysis)
- [Docker Hub](#docker-hub)
- [License](#license)

## Key Features

- **Native Token**: PIXELZX (PXZ)
- **Consensus Mechanism**: Proof of Stake (PoS)
- **EVM Compatibility**: Full Ethereum smart contract support
- **High Performance**: 3-second block time, 1000+ TPS
- **Low Fees**: Gas fee optimization
- **Multi-Architecture**: linux/amd64, linux/arm64, linux/arm/v7 support

## Token Specifications

| Attribute | Value |
|-----------|-------|
| Token Name | PIXELZX |
| Symbol | PXZ |
| Total Supply | 10,000,000,000,000,000 PXZ |
| Decimal Places | 18 |
| Token Type | Native Token |

## Network Parameters

| Parameter | Value |
|-----------|-------|
| Block Time | 3 seconds |
| Block Size Limit | 30MB |
| Gas Limit | 30,000,000 |
| Max Validators | 125 |
| Unbonding Period | 21 days |

## Architecture

### Layered Structure

1. **Application Layer**: DApp interface, API endpoints
2. **EVM Layer**: Ethereum Virtual Machine, smart contract execution
3. **Consensus Layer**: PoS consensus algorithm, block creation/validation
4. **Network Layer**: P2P communication, block propagation
5. **Storage Layer**: State storage, blockchain database

## Directory Structure

```
pos/
├── cmd/                    # Executable binaries
├── consensus/              # PoS consensus mechanism
├── core/                   # Core blockchain logic
├── evm/                    # EVM integration and execution environment
├── network/                # P2P networking
├── api/                    # JSON-RPC, WebSocket API
├── staking/                # Staking and validator management
├── governance/             # Governance system
├── storage/                # Data storage and state management
├── crypto/                 # Encryption and security functions
├── tests/                  # Test code
├── docs/                   # Documentation
└── scripts/                # Utility scripts
```

## Build and Run

### Dependencies

- Go 1.21+
- Git
- Docker (optional)
- Docker Buildx (for multi-platform builds)

### Local Build

```bash
go mod tidy
go build -o bin/pixelzx ./cmd/pixelzx
```

### Docker Build

#### Single Platform Build
```bash
# Build image for current platform
make docker-build-local

# Or build directly
docker build -t pixelzx-pos:latest .
```

#### Multi-Platform Build
```bash
# Docker Buildx setup
make buildx-setup

# Build and deploy for all platforms
make docker-push-multi

# Platform-specific testing
make docker-test-multi
```

#### Supported Platforms
- **linux/amd64**: Intel/AMD 64-bit processors
- **linux/arm64**: ARM 64-bit processors (Apple Silicon, ARM servers)
- **linux/arm/v7**: ARM 32-bit processors (Raspberry Pi, etc.)

### Execution

#### Local Execution
```bash
# Initialize genesis file
./bin/pixelzx init

# Start node
./bin/pixelzx start
```

#### Docker Execution
```bash
# Production environment
docker-compose -f docker-compose.production.yml up -d

# Development environment
docker-compose -f docker-compose.dev.yml up -d

# Direct execution
docker run -d --name pixelzx-node \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  yuchanshin/pixelzx-evm:latest
```

## Docker Quick Start

Learn how to quickly and easily start a PixelZX node with Docker.

### Basic Docker Commands

#### 1. Image Download
```bash
# Download latest image
docker pull yuchanshin/pixelzx-evm:latest

# Download specific version
docker pull yuchanshin/pixelzx-evm:v1.0.0
```

#### 2. Node Initialization (Optional)
```bash
# Mainnet initialization
docker run --rm \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx init --datadir /app/data

# Testnet initialization
docker run --rm \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx init --datadir /app/data --chain-id 8889
```

## CLI Commands Reference

The PIXELZX CLI provides a comprehensive set of commands for managing your blockchain node, similar to Ethereum's geth client.

### Global Options

| Option | Description | Default |
|--------|-------------|---------|
| `--config` | Path to config file | |
| `--datadir` | Data directory path | `./data` |
| `--log-level` | Log level (debug, info, warn, error) | `info` |
| `--testnet` | Run in testnet mode | `false` |
| `--help`, `-h` | Show help | |
| `--version`, `-v` | Show version | |

### Core Commands

| Command | Description |
|---------|-------------|
| `pixelzx account` | Manage accounts |
| `pixelzx init` | Bootstrap and initialize a new genesis block |
| `pixelzx start` | Start PIXELZX node |
| `pixelzx config` | Manage configuration |
| `pixelzx staking` | Manage staking features |
| `pixelzx validator` | Manage validator features |
| `pixelzx governance` | Manage governance features |
| `pixelzx admin` | Admin node commands |
| `pixelzx version` | Show version information |
| `pixelzx help` | Show help for commands |

### Docker Commands

For Docker-based deployments, you can use these commands:

| Command | Description |
|---------|-------------|
| `docker run pixelzx/pixelzx-evm pixelzx [command]` | Run PIXELZX CLI commands in a Docker container |
| `docker-compose -f docker-compose.production.yml up` | Start production node with Docker Compose |
| `docker exec -it pixelzx-node pixelzx [command]` | Execute CLI commands in a running container |

### Account Management

Manage accounts, including creating new accounts, listing existing accounts, importing private keys, 
and updating accounts.

```bash
# Create a new account
pixelzx account new

# List all accounts
pixelzx account list

# Get account balance
pixelzx account balance [address]

# Import a private key
pixelzx account import --private-key [key]

# Export account private key
pixelzx account export [address]

# Update an existing account
pixelzx account update [address]
```

### Node Management

Initialize and start your blockchain node.

```bash
# Initialize genesis block
pixelzx init [network-name]

# Start node
pixelzx start

# Show version information
pixelzx version

# Show help
pixelzx help
```

### Network Commands

Manage network connections and peer interactions.

```bash
# Attach to a running node
pixelzx attach [endpoint]

# Start interactive JavaScript console
pixelzx console

# List connected peers
pixelzx admin peer list

# Show local node enode information
pixelzx admin peer self

# Show detailed information about a specific peer
pixelzx admin peer info [peer-id]

# Connect to a new peer
pixelzx admin peer connect [enode-url]

# Disconnect from a specific peer
pixelzx admin peer disconnect [peer-id]

# Show peer connection statistics
pixelzx admin peer stats
```

#### P2P Network Connection

To connect to the PIXELZX network, you need to connect to at least one peer. You can use the following 
bootnode addresses to establish initial connections:

```bash
# Connect to mainnet bootnode
pixelzx admin peer connect enode://[mainnet-bootnode-id]@[mainnet-bootnode-ip]:30303

# Connect to testnet bootnode
pixelzx admin peer connect enode://[testnet-bootnode-id]@[testnet-bootnode-ip]:30303 --testnet
```

Once connected to a peer, the node will automatically discover and connect to other peers in the network.

#### Peer Management

Manage peer connections and view detailed information about connected nodes.

| Command | Description |
|---------|-------------|
| `pixelzx admin peer list` | List connected peers |
| `pixelzx admin peer self` | Show local node enode information |
| `pixelzx admin peer info [peer-id]` | Show detailed information about a specific peer |
| `pixelzx admin peer connect [enode-url]` | Connect to a new peer |
| `pixelzx admin peer disconnect [peer-id]` | Disconnect from a specific peer |
| `pixelzx admin peer stats` | Show peer connection statistics |

##### List Peers (`admin peer list`)

Display a table of all currently connected peers with their basic information.

```bash
# List peers in table format (default)
pixelzx admin peer list

# List peers in JSON format
pixelzx admin peer list --format json

# Show verbose information
pixelzx admin peer list --verbose
```

##### Local Node Information (`admin peer self`)

Show the enode URL and related information for the local node. This information can be shared with other nodes to establish connections.

```bash
# Show local node information in text format (default)
pixelzx admin peer self

# Show local node information in JSON format
pixelzx admin peer self --format json
```

##### Peer Information (`admin peer info`)

Show detailed information about a specific peer including network addresses, connection status, and capabilities.

```bash
# Show information for a specific peer
pixelzx admin peer info [peer-id]
```

##### Connect to Peer (`admin peer connect`)

Establish a connection to a new peer using its enode URL.

```bash
# Connect to a peer
pixelzx admin peer connect enode://[node-id]@[ip-address]:[port]
```

##### Disconnect from Peer (`admin peer disconnect`)

Terminate connection with a specific peer.

```bash
# Disconnect from a peer
pixelzx admin peer disconnect [peer-id]
```

##### Peer Statistics (`admin peer stats`)

Display network statistics including connection counts, data transfer rates, and protocol information.

```bash
# Show peer statistics
pixelzx admin peer stats
```

### Staking Commands

Manage staking, unstaking, delegating, and viewing rewards.

```bash
# Stake tokens to a validator
pixelzx staking stake [validator-address] --amount [amount]

# Unstake tokens
pixelzx staking unstake [validator-address] --amount [amount]

# Delegate tokens to a validator
pixelzx staking delegate [validator-address] --amount [amount]

# Undelegate tokens
pixelzx staking undelegate [validator-address] --amount [amount]

# View staking rewards
pixelzx staking rewards [address]

# View staking status
pixelzx staking status [address]
```

### Validator Commands

Manage validator registration, status checking, and configuration changes.

```bash
# List validators
pixelzx validator list

# Register a new validator
pixelzx validator register --address [address] --pubkey [pubkey]

# Show validator information
pixelzx validator info [validator-address]

# Update validator information
pixelzx validator update [validator-address] --commission [rate]
```

### Governance Commands

Manage governance proposals, voting, and results.

```bash
# List governance proposals
pixelzx governance list

# Show proposal details
pixelzx governance info [proposal-id]

# Submit a new proposal
pixelzx governance submit --title [title] --description [description]

# Vote on a proposal
pixelzx governance vote [proposal-id] --vote [yes|no|abstain]

# Show proposal result
pixelzx governance result [proposal-id]
```

### Configuration Commands

Manage node configuration including viewing, setting, and validating configurations.

```bash
# Show current configuration
pixelzx config show

# Set configuration value
pixelzx config set [key] [value]

# Reset configuration to defaults
pixelzx config reset --confirm

# Validate configuration
pixelzx config validate
```

### Admin Commands

Advanced administration features for managing and monitoring your PIXELZX node.

| Command | Description |
|---------|-------------|
| `pixelzx admin status` | Node status monitoring |
| `pixelzx admin backup` | Backup important data |
| `pixelzx admin restore` | Restore data from backup |
| `pixelzx admin config` | Advanced configuration management |
| `pixelzx admin debug` | Debugging and diagnostic tools |
| `pixelzx admin peer` | P2P network peer management |
| `pixelzx admin reset` | Node data and configuration reset |
| `pixelzx admin metrics` | 노드 성능 메트릭스 수집 |
| `pixelzx admin snapshot` | 블록체인 스냅샷 관리 |

#### Node Status (`admin status`)

Monitor various aspects of your node's current status including basic information, network connections, staking status, and validator information.

```bash
# Show node basic information and status
pixelzx admin status node

# Show P2P network connection status
pixelzx admin status network

# Show staking pool status
pixelzx admin status staking

# Show validator set information
pixelzx admin status validators
```

#### Data Backup and Restore (`admin backup` / `admin restore`)

Create backups of your node's important data and restore from previous backups.

```bash
# Backup node data
pixelzx admin backup

# Restore node data from backup
pixelzx admin restore
```

#### Configuration Management (`admin config`)

Advanced configuration file management and validation.

```bash
# Show current configuration
pixelzx admin config show

# Validate configuration file
pixelzx admin config validate

# Reset configuration to defaults
pixelzx admin config reset
```

#### Debugging Tools (`admin debug`)

Diagnostic tools for troubleshooting and analyzing node performance.

```bash
# Show detailed logs
pixelzx admin debug logs

# Analyze node performance
pixelzx admin debug profile

# Check system resources
pixelzx admin debug system
```

#### Node Reset (`admin reset`)

Reset node data and configuration to initial state.

```bash
# Reset node data
pixelzx admin reset
```

### Environment Variable Setup

| Variable Name | Default Value | Description | Example |
|---------------|---------------|-------------|---------|
| PIXELZX_CHAIN_ID | 8888 | Chain ID | 8888 (mainnet), 8889 (testnet) |
| PIXELZX_NETWORK | mainnet | Network type | mainnet, testnet, devnet |
| PIXELZX_P2P_PORT | 30303 | P2P communication port | 30303 |
| PIXELZX_RPC_PORT | 8545 | JSON-RPC API port | 8545 |
| PIXELZX_WS_PORT | 8546 | WebSocket API port | 8546 |
| PIXELZX_DATA_DIR | /app/data | Data directory | /app/data |
| PIXELZX_KEYSTORE_DIR | /app/keystore | Keystore directory | /app/keystore |

### Volume Mount Guide

#### Docker Volume Creation
```bash
# Create data and keystore volumes
docker volume create pixelzx-data
docker volume create pixelzx-keystore

# Inspect volume locations
docker volume inspect pixelzx-data
docker volume inspect pixelzx-keystore
```

#### Host Directory Mount
```bash
# Create host directories
mkdir -p $HOME/pixelzx/{data,keystore}

# Mount to host directories
docker run -d \
  --name pixelzx-node \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v $HOME/pixelzx/data:/app/data \
  -v $HOME/pixelzx/keystore:/app/keystore \
  yuchanshin/pixelzx-evm:latest
```

### Health Check and Status Verification

```bash
# Check container status
docker ps | grep pixelzx-node

# Check logs
docker logs pixelzx-node

# Check real-time logs
docker logs -f pixelzx-node

# Access container shell
docker exec -it pixelzx-node /bin/sh

# Check node version
docker exec pixelzx-node pixelzx version

# Check block height
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545
```

## P2P Network Connection

Learn how to set up P2P connections between your PixelZX node and other nodes in the network.

### P2P Port Setup

#### Firewall Setup
```bash
# Ubuntu/Debian firewall setup
sudo ufw allow 30303/tcp
sudo ufw allow 30303/udp

# CentOS/RHEL firewall setup
sudo firewall-cmd --permanent --add-port=30303/tcp
sudo firewall-cmd --permanent --add-port=30303/udp
sudo firewall-cmd --reload
```

#### Docker Port Verification
```bash
# Verify P2P port
docker exec pixelzx-node netstat -tulpn | grep 30303

# Verify port binding
docker port pixelzx-node
```

### Bootnode Connection

#### Network Information Verification
```bash
# Verify current node information
docker exec pixelzx-node pixelzx admin nodeInfo

# Verify connected peer list
docker exec pixelzx-node pixelzx admin peers

# Verify peer count
docker exec pixelzx-node pixelzx admin peerCount
```

#### Passive Peer Addition
```bash
# Connect to a specific peer
docker exec pixelzx-node pixelzx admin addPeer "enode://[PEER_ID]@[IP]:[PORT]"

# Example: Connect to bootnode
docker exec pixelzx-node pixelzx admin addPeer "enode://abcd1234@52.123.45.67:30303"
```

### Network Status Monitoring

#### Sync Status Verification
```bash
# Verify block sync status
docker exec pixelzx-node pixelzx eth syncing

# Verify current block number
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545

# Verify network ID
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' \
  http://localhost:8545
```

#### Peer Connection Status Verification
```bash
# Verify connected peer count
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}' \
  http://localhost:8545

# Verify listening status
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"net_listening","params":[],"id":1}' \
  http://localhost:8545
```

### P2P Connection Troubleshooting

| Issue | Symptoms | Cause | Solution |
|-------|----------|-------|----------|
| Peer connection failure | peerCount is 0 | Firewall blocking | Open port 30303 |
| Slow synchronization | Block height not increasing | Bootnode unresponsive | Try different bootnode |
| NAT issue | Inbound connection impossible | No public IP | Use --nat option |
| Port conflict | Node start failure | Port already in use | Use different port |

#### Detailed Debugging
```bash
# Verify network connection status
docker exec pixelzx-node ss -tulpn | grep 30303

# Test external port access
telnet [YOUR_PUBLIC_IP] 30303

# Verify Docker network settings
docker inspect pixelzx-node | grep -A 10 "NetworkSettings"

# Verify firewall status (Ubuntu)
sudo ufw status verbose

# Restart node with NAT settings
docker run -d \
  --name pixelzx-node-nat \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --nat extip:[YOUR_PUBLIC_IP]
```

## Node Initialization and Configuration

Learn how to initialize and configure your node when starting it for the first time.

### Genesis File Initialization

#### Basic Initialization
```bash
# Mainnet genesis initialization
docker run --rm \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx init --datadir /app/data

# Testnet genesis initialization
docker run --rm \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx init --datadir /app/data --chain-id 8889 --network testnet
```

#### Custom Genesis File Usage
```bash
# Prepare custom genesis file
cat > custom-genesis.json << EOF
{
  "chainId": 8888,
  "homesteadBlock": 0,
  "eip150Block": 0,
  "eip155Block": 0,
  "eip158Block": 0,
  "byzantiumBlock": 0,
  "constantinopleBlock": 0,
  "petersburgBlock": 0,
  "istanbulBlock": 0,
  "berlinBlock": 0,
  "londonBlock": 0,
  "alloc": {
    "0x742d35cc6672c0532925a3b8d6f7b71b47c0062f": {
      "balance": "1000000000000000000000000"
    }
  },
  "difficulty": "0x1",
  "gasLimit": "0x1c9c380"
}
EOF

# Initialize with custom genesis
docker run --rm \
  -v $(pwd)/custom-genesis.json:/app/genesis.json \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx init /app/genesis.json --datadir /app/data
```

### Data Directory Setup

#### Volume Management
```bash
# Create data volume
docker volume create pixelzx-data
docker volume create pixelzx-keystore

# Backup volume
docker run --rm \
  -v pixelzx-data:/source \
  -v $(pwd):/backup \
  alpine tar czf /backup/pixelzx-data-backup.tar.gz -C /source .

# Restore volume
docker run --rm \
  -v pixelzx-data:/target \
  -v $(pwd):/backup \
  alpine tar xzf /backup/pixelzx-data-backup.tar.gz -C /target

# Verify volume contents
docker run --rm \
  -v pixelzx-data:/data \
  alpine ls -la /data
```

#### Directory Structure
```
pixelzx-data/
├── chaindata/          # Blockchain data
├── nodes/              # Node information
├── trie/               # State trie
└── ancient/            # Archive data

pixelzx-keystore/
├── UTC--[timestamp]--[address]  # Key files
└── ...
```

### Customizing Configuration Files

#### Extract Default Configuration Files
```bash
# Verify config files
docker run --rm yuchanshin/pixelzx-evm:latest ls -la /app/configs/

# Extract production config file
docker run --rm \
  -v $(pwd):/backup \
  yuchanshin/pixelzx-evm:latest \
  cp /app/configs/production.yaml /backup/

# Extract development config file
docker run --rm \
  -v $(pwd):/backup \
  yuchanshin/pixelzx-evm:latest \
  cp /app/configs/development.yaml /backup/
```

#### Run with Custom Configuration
```bash
# Modify config file (example)
cat > custom-config.yaml << EOF
chain_id: 8888
network_id: 8888
data_dir: "/app/data"
keystore_dir: "/app/keystore"

rpc:
  enabled: true
  host: "0.0.0.0"
  port: 8545
  cors: ["*"]
  api: ["eth", "net", "web3", "personal", "admin"]

ws:
  enabled: true
  host: "0.0.0.0"
  port: 8546
  origins: ["*"]
  api: ["eth", "net", "web3"]

p2p:
  enabled: true
  host: "0.0.0.0"
  port: 30303
  max_peers: 50
  
logging:
  level: "info"
  format: "json"
EOF

# Run node with custom config
docker run -d \
  --name pixelzx-custom \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v $(pwd)/custom-config.yaml:/app/configs/production.yaml \
  -v pixelzx-data:/app/data \
  -v pixelzx-keystore:/app/keystore \
  yuchanshin/pixelzx-evm:latest
```

### Keystore Management

#### Account Creation
```bash
# Create new account
docker exec -it pixelzx-node pixelzx account new

# List accounts
docker exec pixelzx-node pixelzx account list

# Account info
docker exec pixelzx-node pixelzx account info [ADDRESS]
```

#### Keystore File Management
```bash
# Verify keystore files
docker exec pixelzx-node ls -la /app/keystore/

# Backup keystore files
docker cp pixelzx-node:/app/keystore/ ./keystore-backup/

# Restore keystore files
docker cp ./keystore-backup/ pixelzx-node:/app/keystore/
```

### Initialization Verification

#### System Status Verification
```bash
# Verify node version
docker exec pixelzx-node pixelzx version

# Verify chain ID
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
  http://localhost:8545

# Verify genesis block
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x0",true],"id":1}' \
  http://localhost:8545

# Verify account balance
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x742d35cc6672c0532925a3b8d6f7b71b47c0062f","latest"],"id":1}' \
  http://localhost:8545
```

#### Initialization Troubleshooting
```bash
# Data directory initialization (Caution: All data will be deleted)
docker volume rm pixelzx-data
docker volume create pixelzx-data

# Resolve permission issues
docker exec pixelzx-node chown -R 1000:1000 /app/data
docker exec pixelzx-node chown -R 1000:1000 /app/keystore

# Verify errors in logs
docker logs pixelzx-node | grep -i error
docker logs pixelzx-node | grep -i fatal
```

## API Endpoints

### JSON-RPC API

- **Port**: 8545
- **URL**: http://localhost:8545

### WebSocket API

- **Port**: 8546
- **URL**: ws://localhost:8546

### P2P Network

- **Port**: 30303
- **Protocol**: TCP/UDP

## Troubleshooting

### Permission Issues Resolution

Learn how to resolve permission issues that may occur during PIXELZX node initialization.

#### Common Permission Errors

The following errors may occur:
```bash
Error: Keystore directory creation failed: mkdir data/keystore: permission denied
```

#### Automatic Permission Verification and Resolution Guide

Starting from PIXELZX v2.0, automatic permission verification and detailed resolution guides are provided:

```bash
# Initialization with permission verification
./pixelzx init

# Automatic resolution guide when permission error occurs
```

#### Local Environment Solutions

**1. Use Admin Privileges**
```bash
# Run entire process with admin privileges
sudo ./pixelzx init
sudo ./pixelzx start
```

**2. Use Home Directory (Recommended)**
```bash
# Store data in home directory
./pixelzx init --datadir ~/pixelzx-data
./pixelzx start --datadir ~/pixelzx-data
```

**3. Use Temporary Directory (For Testing)**
```bash
# Store data in temporary directory
./pixelzx init --datadir /tmp/pixelzx-data
./pixelzx start --datadir /tmp/pixelzx-data
```

**4. Change Directory Ownership**
```bash
# Change ownership to current user
sudo chown -R $USER:$USER ./data
chmod -R 755 ./data

# Can then run as regular user
./pixelzx init
./pixelzx start
```

#### Docker Environment Solutions

**1. Use Docker Helper Script (Recommended)**
```bash
# Tool for automatic permission resolution
./docker-helper.sh check  # Check permission status
./docker-helper.sh fix    # Automatically fix permission issues
./docker-helper.sh init   # Chain initialization
./docker-helper.sh start  # Node start
```

**2. Manual Permission Setup**
```bash
# Set host volume permissions
sudo chown -R 1001:1001 ./data ./keystore ./logs
chmod -R 755 ./data ./keystore ./logs

# Start with Docker Compose
docker-compose up -d
```

**3. Use Development Environment**
```bash
# Development Docker Compose (Minimizes permission issues)
export UID=$(id -u)
export GID=$(id -g)
docker-compose -f docker-compose.dev.yml up -d
```

**4. Use Container Internal Paths**
```bash
# Run without host volume mounts
docker run -it yuchanshin/pixelzx-evm:latest init
docker run -d yuchanshin/pixelzx-evm:latest start
```

#### Troubleshooting Steps

**1st Step: Simple Resolution Attempts**
```bash
# Move to home directory and retry
cd ~ && pixelzx init --datadir ~/pixelzx-data
```

**2nd Step: Detailed Diagnosis**
```bash
# Verify current user and permissions
whoami
id
pwd
ls -la

# Verify target directory permissions
ls -la ./
ls -la ./data 2>/dev/null || echo "Data directory does not exist"
```

**3rd Step: Scalable Solutions**
```bash
# Method 1: Change ownership
sudo chown -R $USER:$USER .

# Method 2: Use admin privileges
sudo pixelzx init

# Method 3: Use different location
pixelzx init --datadir /tmp/pixelzx-test
```

#### Prevention Methods

**Set Permissions During Installation**
```bash
# Set appropriate permissions for operating system during build
make install-with-permissions

# Or manual installation
sudo cp bin/pixelzx /usr/local/bin/
sudo chmod +x /usr/local/bin/pixelzx
sudo mkdir -p /etc/pixelzx /var/lib/pixelzx
sudo chown $USER:$USER /var/lib/pixelzx
```

**Environment Variable Setup**
```bash
# Add to .bashrc or .zshrc
export PIXELZX_HOME=$HOME/.pixelzx
export PIXELZX_DATA_DIR=$PIXELZX_HOME/data

# Automatically set directories on use
pixelzx init  # Automatically uses $PIXELZX_DATA_DIR
```

### Docker Related Issues

#### Exec Format Error

When `exec format error` occurs while running Docker container:

1. **Use Multi-Platform Image**: 
   ```bash
   docker run --rm yuchanshin/pixelzx-evm:latest /usr/local/bin/pixelzx version
   ```

2. **Explicit Platform Specification**:
   ```bash
   docker run --rm --platform linux/amd64 yuchanshin/pixelzx-evm:latest /usr/local/bin/pixelzx version
   ```

3. **Use Local Build**:
   ```bash
   make docker-build-local
   docker run --rm yuchanshin/pixelzx-evm:local /usr/local/bin/pixelzx version
   ```

For more details, refer to the [EXEC_FORMAT_ERROR_SOLUTION.md](./EXEC_FORMAT_ERROR_SOLUTION.md) document.

#### Container Start Failure
```bash
# Check container logs
docker logs pixelzx-node

# Check container status
docker ps -a | grep pixelzx

# Verify port conflicts
sudo netstat -tulpn | grep -E '(8545|8546|30303)'

# Restart container
docker restart pixelzx-node

# Completely recreate container
docker stop pixelzx-node
docker rm pixelzx-node
docker run -d --name pixelzx-node \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest
```

#### Volume Permission Issues
```bash
# Verify volume permissions
docker exec pixelzx-node ls -la /app/

# Change permissions
docker exec pixelzx-node chown -R 1000:1000 /app/data
docker exec pixelzx-node chown -R 1000:1000 /app/keystore

# Volume mount issues in SELinux environment
docker run -d --name pixelzx-node \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data:Z \
  -v pixelzx-keystore:/app/keystore:Z \
  yuchanshin/pixelzx-evm:latest
```

### P2P Network Issues

#### Peer Connection Issues
```bash
# Verify firewall status
sudo ufw status
sudo firewall-cmd --list-ports

# Verify port forwarding in NAT environment
# Forward port 30303 to node IP in router settings

# Verify network connection
telnet [REMOTE_NODE_IP] 30303

# Start node in P2P debugging mode
docker run -d --name pixelzx-debug \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --verbosity 5
```

#### Synchronization Issues
```bash
# Verify detailed block sync status
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' \
  http://localhost:8545

# Compare current node block with network latest block
# 1. Current node block height
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545

# 2. Verify latest block from other node
# Use official block explorer or other node API

# Restart synchronization
docker restart pixelzx-node

# Fast synchronization mode (Snapshot usage)
docker run -d --name pixelzx-fast \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --syncmode fast
```

### API Connection Issues

#### JSON-RPC API Connection Failure
```bash
# Verify API service status
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":1}' \
  http://localhost:8545

# Verify port listening status
docker exec pixelzx-node netstat -tulpn | grep 8545

# Allow API ports in firewall
sudo ufw allow 8545/tcp
sudo ufw allow 8546/tcp

# Resolve CORS issues
docker run -d --name pixelzx-cors \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --http.corsdomain "*" --ws.origins "*"
```

#### WebSocket Connection Issues
```bash
# Test WebSocket connection
wscat ws://localhost:8546

# Or test with JavaScript
node -e "
  const WebSocket = require('ws');
  const ws = new WebSocket('ws://localhost:8546');
  ws.on('open', () => {
    console.log('WebSocket connection successful');
    ws.close();
  });
  ws.on('error', (err) => {
    console.log('WebSocket connection failed:', err.message);
  });
"

# Verify WebSocket service status
docker exec pixelzx-node netstat -tulpn | grep 8546
```

### Performance Issues

#### Insufficient Memory
```bash
# Verify container resource usage
docker stats pixelzx-node

# Set memory limits
docker run -d --name pixelzx-limited \
  --memory="2g" --memory-swap="4g" \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest

# Adjust garbage collection settings
docker run -d --name pixelzx-gc \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --cache 1024 --gcmode archive
```

#### Slow Response Times
```bash
# Increase cache size
docker run -d --name pixelzx-cache \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --cache 2048

# SSD usage recommended (When using host directory mounts)
mkdir -p /fast-ssd/pixelzx/{data,keystore}
docker run -d --name pixelzx-ssd \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v /fast-ssd/pixelzx/data:/app/data \
  -v /fast-ssd/pixelzx/keystore:/app/keystore \
  yuchanshin/pixelzx-evm:latest
```

### Log Analysis

#### Key Log Patterns
```bash
# Verify error logs
docker logs pixelzx-node 2>&1 | grep -i error
docker logs pixelzx-node 2>&1 | grep -i fatal
docker logs pixelzx-node 2>&1 | grep -i panic

# P2P connection logs
docker logs pixelzx-node 2>&1 | grep -i peer
docker logs pixelzx-node 2>&1 | grep -i "connection"

# Synchronization logs
docker logs pixelzx-node 2>&1 | grep -i sync
docker logs pixelzx-node 2>&1 | grep -i "block"

# API request logs
docker logs pixelzx-node 2>&1 | grep -i "rpc"
docker logs pixelzx-node 2>&1 | grep -i "http"
```

#### Log Level Adjustment
```bash
# Debug log mode
docker run -d --name pixelzx-debug \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --verbosity 5

# Quiet log mode
docker run -d --name pixelzx-quiet \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --verbosity 1

# JSON formatted logs
docker run -d --name pixelzx-json \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --log.json
```

## Docker Hub

The PIXELZX POS EVM Chain is now available on Docker Hub with multi-architecture support:

- **Repository**: [yuchanshin/pixelzx-evm](https://hub.docker.com/r/yuchanshin/pixelzx-evm)
- **Supported Architectures**: linux/amd64, linux/arm64, linux/arm/v7
- **Latest Version**: 1.0.0

### Docker Pull Commands

```bash
# Pull the latest image
docker pull yuchanshin/pixelzx-evm:latest

# Pull a specific version
docker pull yuchanshin/pixelzx-evm:1.0.0

# Pull image for a specific architecture
docker pull --platform linux/arm64 yuchanshin/pixelzx-evm:latest
```

### Docker Run Commands

```bash
# Run a PIXELZX node
docker run -d --name pixelzx-node \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  yuchanshin/pixelzx-evm:latest

# Run with volume mounts for data persistence
docker run -d --name pixelzx-node \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  -v pixelzx-keystore:/app/keystore \
  yuchanshin/pixelzx-evm:latest

# Initialize genesis file
docker run --rm \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx init
```

## License

MIT License