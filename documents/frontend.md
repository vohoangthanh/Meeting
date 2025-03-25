### Tài liệu cho frontend
#### Kết Nối Với Blockchain và Smart Contract  
Frontend cần kết nối với blockchain bằng thư viện như Web3.js, sử dụng ví  để xác thực. Sau đó, khởi tạo smart contract với ABI và địa chỉ đã triển khai.  

#### Tham Gia Phòng và Xử Lý WebRTC  
Khi tham gia phòng, frontend tạo offer WebRTC, gọi `joinRoom` trên smart contract, và nghe sự kiện `EventForwardedToFrontEnd` để nhận answer từ backend. Sử dụng answer này để thiết lập kết nối WebRTC.  

#### Thêm Track và Rời Phòng  
Để thêm track, gọi `forwardEventToBackend` trên smart contract với dữ liệu track, backend sẽ xử lý publish track. Để rời phòng, gọi `leaveRoom` và đóng kết nối WebRTC.  

---

### Ghi Chú Chi Tiết: Frontend Làm Việc Với Smart Contract và WebRTC  

#### Giới Thiệu  
Tài liệu này trình bày cách frontend của ứng dụng họp DApp tương tác với smart contract và WebRTC để hỗ trợ giao tiếp thời gian thực. Smart contract quản lý thông tin phòng và người tham gia, trong khi backend xử lý tương tác với dịch vụ WebRTC của Cloudflare. Frontend sử dụng API WebRTC chuẩn để quản lý luồng media và kết nối dựa trên thông tin từ smart contract.  

#### Phân Tích Yêu Cầu  
- Frontend cần kết nối với smart contract để thực hiện các hành động như tham gia phòng, thêm track, và nhận cập nhật qua sự kiện.  
- Khi tham gia phòng, frontend tạo offer WebRTC, gửi qua smart contract, và nhận answer từ backend qua sự kiện để hoàn tất kết nối.  
- Backend xử lý tương tác với Cloudflare, cung cấp thông tin như session ID và answer qua smart contract.  
- Frontend cần quản lý WebRTC bằng API chuẩn, không phụ thuộc vào thư viện Cloudflare Calls, dựa trên thông tin từ smart contract.  
- Đặc biệt, publish track và pull track sẽ được xử lý qua hàm `forwardEventToBackend` để gửi sự kiện tới backend, và backend sẽ gọi `forwardEventToFrontend` để gửi lại thông tin cho frontend.  

#### Thiết Kế Hệ Thống  
- **Smart Contract**: Quản lý phòng, người tham gia, và track, với các hàm như `joinRoom`, `addTrack`, `leaveRoom`, và phát sự kiện như `ParticipantJoined`, `EventForwardedToFrontEnd`, `EventForwardedToBackend`.  
- **Backend**: Lắng nghe sự kiện từ smart contract, làm việc với Cloudflare để tạo session và quản lý track, rồi gọi smart contract gửi sự kiện về frontend.  
- **Frontend**: Kết nối blockchain, gọi hàm smart contract, nghe sự kiện, và sử dụng API WebRTC chuẩn để thiết lập kết nối.  

#### Các Bước Cho Frontend  

##### 1. Kết Nối Với Blockchain và Smart Contract  
Frontend cần kết nối với blockchain và khởi tạo smart contract. Sử dụng thư viện Web3.js để tương tác.

```javascript
import Web3 from 'web3';

const web3 = new Web3(window.ethereum);
const contractAddress = '0x...'; // Địa chỉ hợp đồng đã triển khai
const contractABI = [...]; // ABI của hợp đồng

const contract = new web3.eth.Contract(contractABI, contractAddress);
```

##### 2. Tạo Offer Cho Phiên WebRTC  
Khi tham gia phòng, frontend tạo offer cho phiên WebRTC bằng RTCPeerConnection.  

```javascript
const pc = new RTCPeerConnection({
  iceServers: [{ urls: 'stun:stun.cloudflare.com:3478' }], // Ví dụ server ICE
});

// Thêm track địa phương, ví dụ:
const localStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
localStream.getTracks().forEach(track => pc.addTrack(track, localStream));

async function generateOffer() {
  const offer = await pc.createOffer();
  await pc.setLocalDescription(offer);
  return offer;
}
```

##### 3. Tham Gia Phòng Qua Smart Contract  
Gọi hàm `joinRoom` trên smart contract, cung cấp room ID, tên, danh sách track ban đầu, và offer đã tạo.  

```javascript
const roomId = 'room1';
const name = 'Participant Name';
const initialTracks = [
  {
    trackName: 'audio',
    mid: 'm0',
    location: 'local',
    isPublished: true,
    trackData: '0x...' // Dữ liệu track mã hóa
  }
];
const offer = await generateOffer();
const sessionDescription = JSON.stringify(offer); // Chuyển offer thành chuỗi

contract.methods.joinRoom(roomId, name, initialTracks, sessionDescription).send({ from: userAddress })
  .on('transactionHash', (hash) => {
    console.log('Hash giao dịch:', hash);
  })
  .on('error', (error) => {
    console.error('Lỗi:', error);
  });
```

##### 4. Nghe Sự Kiện Từ Smart Contract  
Frontend nghe sự kiện từ smart contract để nhận thông tin cập nhật, đặc biệt là `EventForwardedToFrontEnd`, chứa answer từ Cloudflare.  

```javascript
contract.events.EventForwardedToFrontEnd({}, (error, event) => {
  if (error) {
    console.error(error);
  } else {
    const { roomId, participant, eventData } = event.returnValues;
    const data = JSON.parse(eventData);
    handleSessionInformation(data);
  }
});
```

##### 5. Xử Lý Thông Tin Phiên  (Offer)
Khi nhận thông tin phiên từ sự kiện, frontend sử dụng answer từ Cloudflare để thiết lập kết nối WebRTC.  

```javascript
function handleSessionInformation(data) {
  const { sessionID, cloudflareResponse } = data;
  const answer = JSON.parse(cloudflareResponse.answer);

  pc.setRemoteDescription(new RTCSessionDescription(answer))
    .then(() => {
      console.log('Đã thiết lập mô tả từ xa thành công');
    })
    .catch((error) => {
      console.error('Lỗi thiết lập mô tả từ xa:', error);
    });
}
```

##### 6. Publish Track  
Để publish track mới, frontend gọi `forwardEventToBackend` trên smart contract với dữ liệu track và room ID.  

```javascript
const trackData = JSON.stringify({
  trackName: 'video',
  mid: 'm1',
  location: 'local',
  isPublished: true,
  trackData: '0x...' // Dữ liệu track mã hóa
});

contract.methods.forwardEventToBackend(roomId, trackData).send({ from: userAddress })
  .on('transactionHash', (hash) => {
    console.log('Hash giao dịch publish track:', hash);
  })
  .on('error', (error) => {
    console.error('Lỗi publish track:', error);
  });
```

##### 7. Pull Track  
Tương tự, để pull track từ người khác, frontend gọi `forwardEventToBackend` với dữ liệu pull track.  

```javascript
const pullData = JSON.stringify({
  action: 'pullTrack',
  roomId: roomId,
  participant: userAddress,
  trackInfo: '...' // Dữ liệu pull track
});

contract.methods.forwardEventToBackend(roomId, pullData).send({ from: userAddress })
  .on('transactionHash', (hash) => {
    console.log('Hash giao dịch pull track:', hash);
  })
  .on('error', (error) => {
    console.error('Lỗi pull track:', error);
  });
```

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
