# Peer 연결 명령어 설계 문서

## 1. 개요

PIXELZX POS EVM 체인에서 P2P 네트워크 연결 상태를 모니터링하고 관리하기 위한 새로운 CLI 명령어를 설계합니다. 이 명령어는 기존의 `admin` 명령어 그룹 아래에 추가되어 네트워크 피어 관리 기능을 제공합니다.

### 1.1 목적

- P2P 네트워크 연결 상태 확인
- 연결된 피어 정보 조회
- 피어 연결 관리 (추가, 제거 등)
- 네트워크 연결 문제 진단

### 1.2 기능 요구사항

- 연결된 피어 목록 조회
- 특정 피어의 상세 정보 확인
- 피어 연결 상태 모니터링
- 피어 연결 통계 정보 제공

### 1.3 사용자 대상

- 블록체인 노드 운영자
- 네트워크 관리자
- 기술 지원 팀

## 2. 아키텍처

### 2.1 시스템 구조

PIXELZX POS EVM 체인은 다음 계층으로 구성됩니다:

```
┌─────────────────────────────────────┐
│           CLI Interface             │
├─────────────────────────────────────┤
│         Application Layer           │
├─────────────────────────────────────┤
│          Network Layer              │
│  ┌───────────────────────────────┐  │
│  │      Network Manager          │  │
│  ├───────────────────────────────┤  │
│  │        Peer Management        │  │
│  ├───────────────────────────────┤  │
│  │     Peer Connection Logic     │  │
│  └───────────────────────────────┘  │
├─────────────────────────────────────┤
│          P2P Protocol               │
└─────────────────────────────────────┘
```

### 2.2 CLI 명령어 구조

기존의 `admin` 명령어 구조를 확장하여 `peer` 하위 명령어를 추가합니다:

```
pixelzx admin peer [command]
```

## 3. 명령어 설계

### 3.1 명령어 구조

```
pixelzx admin peer [subcommand] [flags]
```

### 3.2 하위 명령어

| 명령어 | 설명 | 사용 예시 |
|-------|------|----------|
| `list` | 연결된 피어 목록 조회 | `pixelzx admin peer list` |
| `info` | 특정 피어의 상세 정보 | `pixelzx admin peer info [peer-id]` |
| `stats` | 피어 연결 통계 정보 | `pixelzx admin peer stats` |
| `connect` | 새로운 피어에 연결 | `pixelzx admin peer connect [enode-url]` |
| `disconnect` | 특정 피어와 연결 해제 | `pixelzx admin peer disconnect [peer-id]` |

### 3.3 명령어 상세

#### 3.3.1 peer list

연결된 모든 피어의 목록을 표시합니다.

**옵션:**
- `--verbose`: 상세 정보 표시
- `--format`: 출력 형식 (table, json)

#### 3.3.2 peer info

특정 피어의 상세 정보를 표시합니다.

**인자:**
- `peer-id`: 조회할 피어의 ID

#### 3.3.3 peer stats

피어 연결 통계 정보를 표시합니다.

**옵션:**
- `--duration`: 통계 수집 기간 (기본: 1분)

#### 3.3.4 peer connect

새로운 피어에 연결을 시도합니다.

**인자:**
- `enode-url`: 연결할 피어의 ENODE URL

#### 3.3.5 peer disconnect

특정 피어와의 연결을 해제합니다.

**인자:**
- `peer-id`: 연결을 해제할 피어의 ID

## 4. 데이터 모델

### 4.1 Peer 정보 구조

```go
type PeerInfo struct {
    ID          string    // 피어 ID
    Enode       string    // ENODE URL
    IPAddress   string    // IP 주소
    Port        int       // 포트 번호
    Direction   string    // 연결 방향 (inbound/outbound)
    Latency     string    // 지연 시간
    LastSeen    time.Time // 마지막 통신 시간
    ConnectedAt time.Time // 연결 시간
    Capabilities []string // 지원하는 프로토콜
}
```

### 4.2 통계 정보 구조

```go
type PeerStats struct {
    TotalPeers     int       // 총 피어 수
    InboundPeers   int       // 수신 연결 피어 수
    OutboundPeers  int       // 송신 연결 피어 수
    AvgLatency     string    // 평균 지연 시간
    BytesReceived  int64     // 수신 데이터량
    BytesSent      int64     // 송신 데이터량
    StartTime      time.Time // 통계 시작 시간
}
```

## 5. 구현 계획

### 5.1 파일 구조

```
cmd/pixelzx/commands/
├── admin.go
├── admin_peer.go     <- 새로 추가할 파일
├── admin_status.go
├── ...
```

### 5.2 주요 함수

```go
// admin_peer.go
func adminPeerCmd() *cobra.Command
func adminPeerListCmd() *cobra.Command
func adminPeerInfoCmd() *cobra.Command
func adminPeerStatsCmd() *cobra.Command
func adminPeerConnectCmd() *cobra.Command
func adminPeerDisconnectCmd() *cobra.Command
```

### 5.3 네트워크 매니저 연동

`network.Manager`와 연동하여 다음 기능을 구현합니다:

1. 피어 목록 조회
2. 피어 상세 정보 조회
3. 피어 연결/해제
4. 통계 정보 수집

### 5.4 필요한 인터페이스

```go
type PeerManagerInterface interface {
    GetPeers() []*PeerInfo
    GetPeer(id string) (*PeerInfo, error)
    ConnectPeer(enodeURL string) error
    DisconnectPeer(id string) error
    GetPeerStats() *PeerStats
}
```

## 6. 기존 파일 수정

`admin.go` 파일에 peer 명령어를 추가해야 합니다:

```go
// AdminCmd creates the admin command group
func AdminCmd() *cobra.Command {
    cmd := &cobra.Command{
        // ... 기존 코드 ...
    }

    // 하위 명령어 추가
    cmd.AddCommand(
        adminStatusCmd(),
        adminResetCmd(),
        adminBackupCmd(),
        adminRestoreCmd(),
        adminConfigCmd(),
        adminDebugCmd(),
        adminPeerCmd(), // 새로 추가된 명령어
    )

    return cmd
}
```

## 7. 구현 예시 코드

다음은 `admin_peer.go` 파일의 구현 예시입니다:

### 6.2 주요 함수

```go
// admin_peer.go
func adminPeerCmd() *cobra.Command
func adminPeerListCmd() *cobra.Command
func adminPeerInfoCmd() *cobra.Command
func adminPeerStatsCmd() *cobra.Command
func adminPeerConnectCmd() *cobra.Command
func adminPeerDisconnectCmd() *cobra.Command
```

### 6.3 네트워크 매니저 연동

`network.Manager`와 연동하여 다음 기능을 구현합니다:

1. 피어 목록 조회
2. 피어 상세 정보 조회
3. 피어 연결/해제
4. 통계 정보 수집

#### 6.3.1 필요한 인터페이스

```go
type PeerManagerInterface interface {
    GetPeers() []*PeerInfo
    GetPeer(id string) (*PeerInfo, error)
    ConnectPeer(enodeURL string) error
    DisconnectPeer(id string) error
    GetPeerStats() *PeerStats
}
```

### 6.4 기존 파일 수정

`admin.go` 파일에 peer 명령어를 추가해야 합니다:

```go
// AdminCmd creates the admin command group
func AdminCmd() *cobra.Command {
    cmd := &cobra.Command{
        // ... 기존 코드 ...
    }

    // 하위 명령어 추가
    cmd.AddCommand(
        adminStatusCmd(),
        adminResetCmd(),
        adminBackupCmd(),
        adminRestoreCmd(),
        adminConfigCmd(),
        adminDebugCmd(),
        adminPeerCmd(), // 새로 추가된 명령어
    )

    return cmd
}
```

### 6.5 구현 예시 코드

다음은 `admin_peer.go` 파일의 구현 예시입니다:

```go
package commands

import (
    "fmt"
    "time"

    "github.com/ethereum/go-ethereum/p2p/enode"
    "github.com/spf13/cobra"
)

// adminPeerCmd creates the peer command group
func adminPeerCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "peer",
        Short: "P2P 네트워크 피어 관리",
        Long: `P2P 네트워크에 연결된 피어를 관리합니다.

연결된 피어 목록을 조회하고, 특정 피어의 상세 정보를 확인하며,
피어와의 연결을 관리할 수 있습니다.`,
    }

    // 하위 명령어 추가
    cmd.AddCommand(
        adminPeerListCmd(),
        adminPeerInfoCmd(),
        adminPeerStatsCmd(),
        adminPeerConnectCmd(),
        adminPeerDisconnectCmd(),
    )

    return cmd
}

// adminPeerListCmd lists connected peers
func adminPeerListCmd() *cobra.Command {
    var (
        verbose bool
        format  string
    )

    cmd := &cobra.Command{
        Use:   "list",
        Short: "연결된 피어 목록 조회",
        Long:  "현재 P2P 네트워크에 연결된 모든 피어의 목록을 표시합니다.",
        RunE: func(cmd *cobra.Command, args []string) error {
            // 실제 구현에서는 네트워크 매니저에서 피어 목록을 가져옴
            // 여기서는 예시 데이터 표시
            return nil
        },
    }

    cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "상세 정보 표시")
    cmd.Flags().StringVar(&format, "format", "table", "출력 형식 (table, json)")

    return cmd
}

// 다른 함수들도 유사한 방식으로 구현...
```

## 7. 보안 고려사항

1. 피어 연결/해제 명령어는 관리자 권한이 필요
2. 외부 노드 연결 시 보안 검증 수행
3. 잘못된 ENODE URL 처리
4. 연결 실패 시 적절한 오류 메시지 제공

## 8. 테스트 계획

### 8.1 단위 테스트

1. 명령어 파싱 테스트
2. 피어 정보 표시 테스트
3. 통계 정보 계산 테스트
4. 연결/해제 기능 테스트

### 8.2 통합 테스트

1. 실제 네트워크 연결 상태에서 명령어 동작 확인
2. 다양한 피어 연결 시나리오 테스트
3. 에러 상황 처리 테스트

## 9. 사용 예시

### 9.1 피어 목록 조회

```bash
$ pixelzx admin peer list
ID                                    IP 주소          포트   방향    지연시간
16Uiu2HAm9Q8R7nxZ8GWqK9vK1b2X3h4mL5nF8dY3...  192.168.1.100    30303  수신    45ms
16Uiu2HAm5P7M6nxY9FVqH8vJ1a2W2g3hK4mE7cX2...  203.123.45.67    30303  송신    120ms
```

### 9.2 피어 상세 정보 조회

```bash
$ pixelzx admin peer info 16Uiu2HAm9Q8R7nxZ8GWqK9vK1b2X3h4mL5nF8dY3...
피어 ID: 16Uiu2HAm9Q8R7nxZ8GWqK9vK1b2X3h4mL5nF8dY3...
ENODE: enode://9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7@192.168.1.100:30303
IP 주소: 192.168.1.100
포트: 30303
연결 방향: 수신
지연 시간: 45ms
마지막 통신: 2024-01-25 10:30:45 UTC
연결 시간: 2024-01-25 09:15:22 UTC
지원 프로토콜: pixelzx/1.0
```

### 9.3 통계 정보 조회

```bash
$ pixelzx admin peer stats
총 피어 수: 24개
수신 연결: 12개
송신 연결: 12개
평균 지연 시간: 122ms
수신 데이터: 2.3 GB
송신 데이터: 1.8 GB
```


















































































































































