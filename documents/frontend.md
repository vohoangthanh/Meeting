### Tài liệu cho frontend
#### Tổng quan về Frontend

Frontend của ứng dụng DApp Meeting là một ứng dụng web sử dụng HTML, CSS, và JavaScript thuần để tạo giao diện cho phép người dùng tham gia các cuộc họp video thông qua WebRTC. Frontend tương tác trực tiếp với Smart Contract trên blockchain và xử lý kết nối WebRTC thông qua các API chuẩn của trình duyệt.

#### Các thành phần chính của Frontend

1. **Xác thực và kết nối Blockchain**:
   - Sử dụng thư viện ethers.js để tương tác với blockchain
   - Xác thực người dùng thông qua private key (lưu trong localStorage)
   - Khởi tạo kết nối với smart contract thông qua ABI và địa chỉ hợp đồng

2. **Quản lý phòng họp**:
   - Tạo phòng mới với UUID ngẫu nhiên
   - Tham gia phòng hiện có thông qua room ID
   - Hiển thị danh sách người tham gia trong phòng

3. **Xử lý WebRTC**:
   - Tạo và quản lý RTCPeerConnection
   - Thu thập và quản lý local tracks (audio, video)
   - Tạo SDP offer/answer và xử lý trao đổi thông tin kết nối

4. **Tương tác với Smart Contract**:
   - Gọi các hàm trên smart contract (createRoom, joinRoom, leaveRoom)
   - Lắng nghe sự kiện từ smart contract (EventForwardedToFrontend)
   - Gửi và nhận dữ liệu qua cơ chế chuyển tiếp sự kiện

---

### Ghi Chú Chi Tiết: Frontend Làm Việc Với Smart Contract và WebRTC  

### Kiến trúc và Tương tác Frontend

#### Cấu trúc tệp và thành phần
Frontend được tổ chức với cấu trúc sau:

1. **Trang HTML chính**:
   - `index.html`: Trang chủ để tạo hoặc tham gia phòng
   - `join.html`: Trang tham gia phòng bằng ID
   - `room.html`: Giao diện của phòng họp
   - `check.html`: Trang kiểm tra trước khi vào phòng

2. **Mô-đun JavaScript**:
   - `auth.js`: Xử lý xác thực và kết nối ví
   - `smartContractIntegration.js`: Tương tác với smart contract
   - `CloudflareCalls.js`: Quản lý kết nối WebRTC và xử lý các cuộc gọi
   - `contractEventListener.js`: Lắng nghe các sự kiện từ smart contract
   - `room.js`: Logic cho phòng họp
   - `abi.js`: ABI của smart contract
   - `FaceMaskFilter.js` và `backgroundBlur.js`: Bộ lọc video

3. **Định nghĩa Lớp chính**:
   - `CloudflareCalls`: Quản lý kết nối WebRTC
   - `SmartContractConnector`: Tương tác với smart contract
   - `Auth`: Xử lý xác thực người dùng

#### Luồng xử lý dữ liệu

1. **Xác thực người dùng**:
   ```
   Người dùng → Private Key → auth.js → Kết nối blockchain → 
   Lưu địa chỉ ví → Kết nối smart contract
   ```

2. **Tạo phòng mới**:
   ```
   Người dùng → Nhập tên → Tạo UUID → smartContractIntegration.js → 
   Gọi createRoom() → Lưu room ID → Chuyển hướng đến check.html
   ```

3. **Tham gia phòng**:
   ```
   Người dùng → Nhập Room ID → smartContractIntegration.js → 
   Kiểm tra phòng tồn tại → Gọi joinRoom() → Chuyển hướng đến room.html
   ```

4. **Thiết lập WebRTC**:
   ```
   room.js → Tạo RTCPeerConnection → Yêu cầu MediaStream → 
   createOffer → setLocalDescription → Nén dữ liệu → 
   Gọi joinRoom() với offer → Nhận sự kiện EventForwardedToFrontend → 
   Giải nén answer → setRemoteDescription → Kết nối thành công
   ```

5. **Xử lý Track Media**:
   ```
   Khi track local → addTrack vào PeerConnection → 
   RTCPeerConnection.getSenders() → Gọi addTrack() hoặc forwardEventToBackend() →
   Smart Contract → Backend → Cloudflare → Backend → Smart Contract →
   Phát sự kiện EventForwardedToFrontend → Frontend xử lý
   ```

#### Công nghệ và thư viện chính

1. **Blockchain và Smart Contract**:
   - ethers.js: Tương tác với blockchain Ethereum
   - pako.js: Nén và giải nén dữ liệu lớn (SDP)

2. **WebRTC**:
   - APIs chuẩn: RTCPeerConnection, MediaStream
   - adapter.js: Đảm bảo tương thích giữa các trình duyệt

3. **UI và đồ họa**:
   - Material Icons: Icons cho UI
   - CSS tùy chỉnh: Định nghĩa giao diện người dùng

### Triển khai chi tiết Frontend

#### Khởi tạo kết nối

Frontend sử dụng thư viện ethers.js để tương tác với blockchain thay vì Web3.js như sau:

```javascript
// Khởi tạo SmartContractConnector
class SmartContractConnector {
    constructor() {
        this.contract = null;
        this.provider = null;
        this.signer = null;
        this.isInitialized = false;
    }

    async initialize() {
        try {
            // Lấy private key từ localStorage
            const privateKey = localStorage.getItem('privateKey');
            if (!privateKey) throw new Error('No private key found');

            // Khởi tạo provider cho mạng
            this.provider = new ethers.providers.JsonRpcProvider("https://rpc-endpoint");
            
            // Tạo ví và signer từ private key
            const wallet = new ethers.Wallet(privateKey, this.provider);
            this.signer = wallet;
            
            // Địa chỉ smart contract
            const contractAddress = 'CONTRACT_ADDRESS_HERE';
            
            // Khởi tạo contract với ABI và signer
            this.contract = new ethers.Contract(
                contractAddress,
                CONTRACT_ABI,
                this.signer
            );
            
            this.isInitialized = true;
            console.log('Smart contract initialized');
            return true;
        } catch (error) {
            console.error('Error initializing smart contract:', error);
            return false;
        }
    }
}
```

#### Quản lý kết nối WebRTC

Lớp `CloudflareCalls` quản lý kết nối WebRTC với các phương thức chính:

```javascript
class CloudflareCalls {
    constructor(config) {
        this.config = config;
        this.peerConnection = null;
        this.localStream = null;
        this.roomId = null;
        this.sessionId = null;
        this.token = null;
        // Queue để xử lý Pull Track
        this.pullTrackQueue = [];
        this.isProcessingPullTrack = false;
        // Maps để theo dõi tracks
        this.midToSessionId = new Map();
        this.midToTrackName = new Map();
    }
    
    // Tạo RTCPeerConnection với cấu hình ICE servers
    async _createPeerConnection() {
        const pc = new RTCPeerConnection({
            iceServers: [
                { urls: 'stun:stun.cloudflare.com:3478' },
                // Các máy chủ TURN nếu cần
            ]
        });
        
        // Xử lý các sự kiện kết nối
        pc.ontrack = this._handleTrack.bind(this);
        pc.onicecandidate = this._handleIceCandidate.bind(this);
        pc.onconnectionstatechange = this._handleConnectionStateChange.bind(this);
        
        return pc;
    }
    
    // Quá trình tham gia phòng họp
    async joinRoom(roomId, metadata = {}) {
        this.roomId = roomId;
        
        // Tạo PeerConnection
        this.peerConnection = await this._createPeerConnection();
        
        // Lấy local media streams
        if (!this.localStream) {
            this.localStream = await navigator.mediaDevices.getUserMedia({ 
                video: true, 
                audio: true 
            });
        }
        
        // Thêm local tracks vào PeerConnection
        this.localStream.getTracks().forEach(track => {
            this.peerConnection.addTrack(track, this.localStream);
        });
        
        // Tạo offer và set local description
        const offer = await this.peerConnection.createOffer();
        await this.peerConnection.setLocalDescription(offer);
        
        // Nén SDP offer
        const offerData = JSON.stringify(offer);
        const encoder = new TextEncoder();
        const data = encoder.encode(offerData);
        const compressedOffer = pako.deflate(data);
        
        // Gọi Smart Contract để tham gia phòng với offer
        return await SmartContractConnector.joinRoom(
            roomId,
            metadata.name || 'User',
            [], // initial tracks được xử lý riêng
            compressedOffer
        );
    }
}

#### Xử lý tham gia phòng và lắng nghe sự kiện

Quá trình tham gia phòng và lắng nghe sự kiện từ smart contract được thực hiện như sau:

```javascript
// Từ file room.js - Hàm tham gia phòng họp
async function joinRoom() {
    // 1. Kiểm tra kết nối với smart contract
    if (!SmartContractConnector.contract) {
        await SmartContractConnector.initialize();
    }

    // 2. Thiết lập listeners trước khi join
    const removeAnswerListener = await SmartContractConnector.listenForTrackPublishedAnswer(
        roomId,
        async (answerData) => {
            if (answerData && answerData.sessionDescription) {
                // Thiết lập remote description từ answer
                if (calls.peerConnection && calls.peerConnection.signalingState === 'have-local-offer') {
                    const remoteDesc = new RTCSessionDescription(answerData.sessionDescription);
                    await calls.peerConnection.setRemoteDescription(remoteDesc);
                    console.log('Remote description set successfully');
                }
            }
        }
    );

    // 3. Thiết lập các callbacks cho WebRTC
    await setupCallbacks();

    // 4. Lấy danh sách người tham gia hiện tại từ smart contract
    const participants = await calls.getParticipantsFromContract(roomId);
    
    // 5. Tìm thông tin của bản thân trong danh sách người tham gia
    const wallet = new ethers.Wallet(localStorage.getItem('privateKey'));
    const myWalletAddress = wallet.address;
    const participant = participants.find(p => 
        p.walletAddress.toLowerCase() === myWalletAddress.toLowerCase()
    );
    
    if (participant) {
        calls.sessionId = participant.sessionId;
    }

    // 6. Xử lý kết nối và pull tracks từ người tham gia khác
    for (const participant of participants) {
        // Bỏ qua nếu là bản thân
        if (participant.walletAddress.toLowerCase() === myWalletAddress.toLowerCase()) continue;
        
        // Pull các tracks từ người tham gia này
        for (const trackName of participant.publishedTracks) {
            await calls._pullTracks(participant.sessionId, trackName);
        }
    }
}
```

#### Lắng nghe sự kiện từ smart contract

Frontend lắng nghe các sự kiện từ smart contract thông qua module ContractEventListener:

```javascript
// Từ file contractEventListener.js
class ContractEventListener {
    constructor() {
        this.eventListeners = {};
        this.currentSessionId = null;
    }
    
    setCurrentSessionId(sessionId) {
        this.currentSessionId = sessionId;
    }
    
    // Lắng nghe sự kiện từ smart contract
    listenForEventsToFrontend(roomId, callback) {
        if (!SmartContractConnector.contract) return;
        
        // Tạo filter để lắng nghe sự kiện EventForwardedToFrontend
        const filter = SmartContractConnector.contract.filters.EventForwardedToFrontend(
            roomId,
            null,  // Lắng nghe cho mọi participant
            null   // Lắng nghe mọi dữ liệu sự kiện
        );
        
        // Xử lý khi nhận được sự kiện
        const handler = (roomId, participant, eventData) => {
            try {
                // Chỉ xử lý sự kiện dành cho mình
                const userAddress = auth.userAddress.toLowerCase();
                if (participant.toLowerCase() !== userAddress) return;
                
                // Giải nén dữ liệu sự kiện
                const decompressedData = this._decompressEventData(eventData);
                if (decompressedData) {
                    callback(decompressedData);
                }
            } catch (error) {
                console.error('Error processing frontend event:', error);
            }
        };
        
        // Đăng ký listener với smart contract
        SmartContractConnector.contract.on(filter, handler);
        
        // Trả về hàm để hủy đăng ký
        return () => {
            SmartContractConnector.contract.off(filter, handler);
        };
    }
    
    // Giải nén dữ liệu sự kiện
    _decompressEventData(eventData) {
        try {
            // Chuyển từ hex sang bytes
            const bytes = ethers.utils.arrayify(eventData);
            // Giải nén với pako
            const decompressed = pako.inflate(bytes);
            // Chuyển bytes thành text
            const decoder = new TextDecoder();
            const jsonString = decoder.decode(decompressed);
            // Parse JSON
            return JSON.parse(jsonString);
        } catch (error) {
            console.error('Error decompressing event data:', error);
            return null;
        }
    }
}

#### Xử lý tracks và renegotiation

Frontend xử lý publish và pull tracks thông qua lớp CloudflareCalls với cơ chế nén/giải nén dữ liệu để tối ưu tương tác với blockchain:

```javascript
// Method trong CloudflareCalls.js để publish tracks
async _publishTracks() {
    if (!this.peerConnection || !this.localStream) {
        throw new Error('No peer connection or local stream available');
    }

    try {
        console.log('Publishing tracks to session...');
        
        // Tạo offer cho tracks mới
        const offer = await this.peerConnection.createOffer();
        await this.peerConnection.setLocalDescription(offer);
        
        // Nén dữ liệu SDP offer
        const offerData = JSON.stringify({
            type: 'publish',
            sdp: offer.sdp,
            roomId: this.roomId,
            sessionId: this.sessionId,
            timestamp: Date.now()
        });
        
        // Nén dữ liệu với pako
        const encoder = new TextEncoder();
        const data = encoder.encode(offerData);
        const compressedData = pako.deflate(data);
        const eventData = ethers.utils.hexlify(compressedData);
        
        // Gọi smart contract để forward event đến backend
        await SmartContractConnector.forwardEventToBackend(
            this.roomId,
            eventData
        );
        
        console.log('Published tracks successfully');
    } catch (error) {
        console.error('Error publishing tracks:', error);
        throw error;
    }
}

// Method trong CloudflareCalls.js để pull track từ người khác
async pullTrack(remoteSessionId, trackName) {
    if (!this.peerConnection || this.peerConnection.connectionState !== 'connected') {
        throw new Error('PeerConnection not ready');
    }
    
    console.log(`Pulling track '${trackName}' from session ${remoteSessionId}`);
    
    // Queue yêu cầu pull track để xử lý tuần tự
    return new Promise((resolve, reject) => {
        this.pullTrackQueue.push({
            remoteSessionId,
            trackName,
            resolve,
            reject
        });
        
        // Bắt đầu xử lý queue nếu chưa có quá trình đang chạy
        if (!this.isProcessingPullTrack) {
            this._processPullTrackQueue();
        }
    });
}

// Method xử lý queue pull track
async _processPullTrackQueue() {
    if (this.pullTrackQueue.length === 0) {
        this.isProcessingPullTrack = false;
        return;
    }
    
    this.isProcessingPullTrack = true;
    const { remoteSessionId, trackName, resolve, reject } = this.pullTrackQueue.shift();
    
    try {
        // Nén dữ liệu yêu cầu pull track
        const pullData = JSON.stringify({
            type: 'pullTrack',
            trackName,
            remoteSessionId,
            sessionId: this.sessionId,
            roomId: this.roomId,
            timestamp: Date.now()
        });
        
        const encoder = new TextEncoder();
        const data = encoder.encode(pullData);
        const compressedData = pako.deflate(data);
        const eventData = ethers.utils.hexlify(compressedData);
        
        // Thiết lập listener cho pull answer
        const removeListener = SmartContractConnector.listenForEventsToFrontend(
            this.roomId,
            async (pullAnswerData) => {
                if (pullAnswerData && pullAnswerData.type === 'pullAnswer' && 
                    pullAnswerData.sessionDescription) {
                    
                    // Cập nhật mapping từ mid sang session và track name
                    if (pullAnswerData.tracks) {
                        pullAnswerData.tracks.forEach(track => {
                            const { mid, sessionId, trackName } = track;
                            this.midToSessionId.set(mid, sessionId);
                            this.midToTrackName.set(mid, trackName);
                        });
                    }
                    
                    // Set remote description từ answer
                    await this.peerConnection.setRemoteDescription(
                        pullAnswerData.sessionDescription
                    );
                    
                    // Tạo và set local answer cho renegotiation
                    const localAnswer = await this.peerConnection.createAnswer();
                    await this.peerConnection.setLocalDescription(localAnswer);
                    
                    // Gửi renegotiation request
                    await SmartContractConnector.callToRenegotiate(
                        this.sessionId,
                        localAnswer.sdp,
                        localAnswer.type
                    );
                    
                    resolve();
                }
            }
        );
        
        // Gửi yêu cầu pull track
        await SmartContractConnector.forwardEventToBackend(
            this.roomId,
            eventData
        );
    } catch (error) {
        console.error('Error processing pull track:', error);
        reject(error);
    } finally {
        // Xử lý yêu cầu tiếp theo trong queue
        this._processPullTrackQueue();
    }
}
```

#### Xử lý renegotiation kết nối WebRTC

Quá trình renegotiation rất quan trọng khi thêm hoặc xóa tracks:

```javascript
// Method xử lý renegotiation trong CloudflareCalls
async _renegotiate() {
    if (!this.peerConnection) {
        console.log('No peer connection available for renegotiation');
        return;
    }
    
    try {
        console.log('Creating renegotiation answer...');
        const answer = await this.peerConnection.createAnswer();
        
        if (this.peerConnection.signalingState === 'have-remote-offer') {
            console.log('Setting local description for renegotiation');
            await this.peerConnection.setLocalDescription(answer);
            
            // Thiết lập listener cho renegotiation response
            const removeListener = SmartContractConnector.ListenToAnswerRenegotiate(
                this.sessionId,
                async (answerData) => {
                    try {
                        if (answerData.sessionDescription) {
                            const remoteDesc = new RTCSessionDescription(answerData.sessionDescription);
                            await this.peerConnection.setRemoteDescription(remoteDesc);
                            console.log('Renegotiation completed successfully');
                        }
                    } catch (error) {
                        console.error('Error handling renegotiation response:', error);
                    } finally {
                        if (removeListener) removeListener();
                    }
                }
            );
            
            // Gửi renegotiation request
            await SmartContractConnector.callToRenegotiate(
                this.sessionId,
                answer.sdp,
                answer.type
            );
        } else {
            console.log('Peer connection not in correct state for renegotiation');
        }
    } catch (error) {
        console.error('Error during renegotiation:', error);
    }
}

##### 8. Rời Phòng  
Để rời phòng, gọi `leaveRoom` trên smart contract và đóng PeerConnection.  

```javascript
contract.methods.leaveRoom(roomId).send({ from: userAddress })
  .on('transactionHash', (hash) => {
    console.log('Hash giao dịch rời phòng:', hash);
    pc.close();
  })
  .on('error', (error) => {
    console.error('Lỗi rời phòng:', error);
  });
``` 

#### Bảng Tổng Hợp  
| Bước                     | Mô Tả                                      | Công Việc Chính                                      | Giao Tiếp Với Khác       |
|--------------------------|--------------------------------------------|-----------------------------------------------------|--------------------------|
| Kết Nối Blockchain        | Kết nối với blockchain và smart contract   | Sử dụng Web3.js, khởi tạo hợp đồng                  | -                        |
| Tạo Offer                 | Tạo offer cho phiên WebRTC                 | Sử dụng RTCPeerConnection, thêm track địa phương    | -                        |
| Tham Gia Phòng            | Gọi `joinRoom` trên smart contract         | Gửi room ID, tên, track, và offer                   | Gửi sự kiện đến backend  |
| Nghe Sự Kiện              | Nghe sự kiện từ smart contract             | Nhận answer từ Cloudflare qua `EventForwardedToFrontEnd` | -                        |
| Xử Lý Thông Tin Phiên     | Thiết lập kết nối WebRTC                   | Sử dụng answer để gọi `setRemoteDescription`        | -                        |
| Publish Track             | Gửi sự kiện publish track tới backend      | Gọi `forwardEventToBackend` với dữ liệu track       | Gửi sự kiện đến backend  |
| Pull Track                | Gửi sự kiện pull track tới backend         | Gọi `forwardEventToBackend` với dữ liệu pull track  | Gửi sự kiện đến backend  |
| Rời Phòng                 | Gọi `leaveRoom` và đóng kết nối            | Gọi `leaveRoom`, đóng PeerConnection                | Gửi sự kiện đến backend  |
