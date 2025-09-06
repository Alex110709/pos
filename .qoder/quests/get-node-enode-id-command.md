# PIXELZX 노드 enode ID 조회 명령어 설계

## 1. 개요

PIXELZX POS EVM 체인 노드의 고유 식별자인 enode ID를 조회하는 새로운 CLI 명령어를 설계합니다. 이 명령어는 노드 운영자가 자신의 노드 정보를 확인하고 다른 노드와 연결 설정을 할 때 사용할 수 있습니다.

## 2. 요구사항

### 2.1 기능 요구사항
- 노드의 enode URL을 조회하여 표시
- enode URL은 `enode://[public-key]@[ip]:[port]` 형식을 따름
- JSON 형식 출력 옵션 제공
- 오류 발생 시 적절한 오류 메시지 출력

### 2.2 비기능 요구사항
- 기존 코드 구조와 일관성 유지
- Cobra CLI 프레임워크 사용
- 기존 네트워크 관리자 인터페이스와 통합

## 3. 설계

### 3.1 아키텍처

```mermaid
graph TD
    A[CLI Command] --> B[Cobra Command Handler]
    B --> C[Network Manager]
    C --> D[P2P Server]
    D --> E[Self() 메서드]
    E --> F[enode.Node 정보]
    F --> B
    B --> G[Output Formatter]
    G --> H[Console 출력]
```

### 3.2 컴포넌트 설계

#### 3.2.1 CLI 명령어 구조
새로운 명령어는 `admin` 명령어 그룹 아래에 위치하며, 다음과 같은 구조를 가집니다:
```
pixelzx admin peer self
```

#### 3.2.2 Network Manager 인터페이스 확장
기존 Network Manager에 로컬 노드의 enode 정보를 조회하는 새로운 메서드를 추가합니다:
```go
// GetLocalEnode returns the local node's enode information
func (nm *Manager) GetLocalEnode() (*enode.Node, error) {
    if nm.server == nil {
        return nil, fmt.Errorf("P2P server not initialized")
    }
    return nm.server.Self(), nil
}
```

#### 3.2.3 Cobra Command Handler
`admin_peer.go` 파일에 새로운 명령어 핸들러를 추가합니다:
```go
func adminPeerSelfCmd() *cobra.Command {
    var format string
    
    cmd := &cobra.Command{
        Use:   "self",
        Short: "로컬 노드의 enode 정보 조회",
        Long:  "현재 PIXELZX 노드의 enode URL과 관련 정보를 표시합니다.",
        RunE: func(cmd *cobra.Command, args []string) error {
            return showLocalEnode(format)
        },
    }
    
    cmd.Flags().StringVar(&format, "format", "text", "출력 형식 (text, json)")
    
    return cmd
}
```

## 4. 구현 계획

### 4.1 파일 수정 계획

1. `network/manager.go`:
   - `GetLocalEnode()` 메서드 추가

2. `cmd/pixelzx/commands/admin_peer.go`:
   - `adminPeerSelfCmd()` 함수 추가
   - `showLocalEnode()` 함수 추가
   - `adminPeerCmd()` 함수에 새로운 하위 명령어 등록

### 4.2 상세 구현

#### 4.2.1 Network Manager 수정
```go
// GetLocalEnode returns the local node's enode information
func (nm *Manager) GetLocalEnode() (*enode.Node, error) {
    nm.mu.RLock()
    defer nm.mu.RUnlock()
    
    if nm.server == nil {
        return nil, fmt.Errorf("P2P server not initialized")
    }
    
    return nm.server.Self(), nil
}
```

#### 4.2.2 CLI 명령어 추가
```go
// adminPeerSelfCmd shows local node enode information
func adminPeerSelfCmd() *cobra.Command {
    var format string
    
    cmd := &cobra.Command{
        Use:   "self",
        Short: "로컬 노드의 enode 정보 조회",
        Long:  "현재 PIXELZX 노드의 enode URL과 관련 정보를 표시합니다.",
        RunE: func(cmd *cobra.Command, args []string) error {
            return showLocalEnode(format)
        },
    }
    
    cmd.Flags().StringVar(&format, "format", "text", "출력 형식 (text, json)")
    
    return cmd
}

// showLocalEnode displays local node enode information
func showLocalEnode(format string) error {
    // 네트워크 관리자에서 로컬 enode 정보 조회
    // 실제 구현에서는 네트워크 관리자 인스턴스를 통해 조회해야 함
    node, err := getNetworkManager().GetLocalEnode()
    if err != nil {
        return fmt.Errorf("enode 정보 조회 실패: %w", err)
    }
    
    if format == "json" {
        return showLocalEnodeJSON(node)
    }
    
    return showLocalEnodeText(node)
}

// showLocalEnodeText displays enode information in text format
func showLocalEnodeText(node *enode.Node) error {
    fmt.Printf("로컬 노드 enode 정보:\n")
    fmt.Printf("========================\n")
    fmt.Printf("enode URL: %s\n", node.String())
    fmt.Printf("Node ID: %s\n", node.ID().String())
    fmt.Printf("IP 주소: %s\n", node.IP().String())
    fmt.Printf("TCP 포트: %d\n", node.TCP())
    fmt.Printf("UDP 포트: %d\n", node.UDP())
    
    return nil
}

// showLocalEnodeJSON displays enode information in JSON format
func showLocalEnodeJSON(node *enode.Node) error {
    info := struct {
        EnodeURL string `json:"enode_url"`
        NodeID   string `json:"node_id"`
        IP       string `json:"ip"`
        TCPPort  int    `json:"tcp_port"`
        UDPPort  int    `json:"udp_port"`
    }{
        EnodeURL: node.String(),
        NodeID:   node.ID().String(),
        IP:       node.IP().String(),
        TCPPort:  node.TCP(),
        UDPPort:  node.UDP(),
    }
    
    output, err := json.MarshalIndent(info, "", "  ")
    if err != nil {
        return fmt.Errorf("JSON 직렬화 실패: %w", err)
    }
    
    fmt.Println(string(output))
    return nil
}
```

#### 4.2.3 명령어 등록
```go
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
        adminPeerSelfCmd(), // 새로 추가된 명령어
    )

    return cmd
}
```

## 5. 테스트 계획

### 5.1 단위 테스트
- `showLocalEnodeText` 함수 테스트
- `showLocalEnodeJSON` 함수 테스트
- `validateEnodeURL` 함수에 대한 테스트 케이스 추가 (관련된 경우)

### 5.2 통합 테스트
- CLI 명령어 실행 테스트
- 네트워크 관리자와의 통합 테스트
- 다양한 출력 형식 테스트

## 6. 사용 예시

### 6.1 기본 사용법
```bash
# 기본 텍스트 형식으로 enode 정보 조회
pixelzx admin peer self

# JSON 형식으로 enode 정보 조회
pixelzx admin peer self --format json
```

### 6.2 예상 출력
```
로컬 노드 enode 정보:
========================
enode URL: enode://a1b2c3d4e5f67890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef12345@192.168.1.100:30303
Node ID: a1b2c3d4e5f67890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
IP 주소: 192.168.1.100
TCP 포트: 30303
UDP 포트: 30303
```

## 7. 보안 고려사항

- enode URL은 네트워크 연결을 위한 공개 정보이므로 보안상의 문제가 없음
- 노드의 IP 주소가 노출될 수 있으나, 이는 P2P 네트워크 운영에 필수적인 정보임