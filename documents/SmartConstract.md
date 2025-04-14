## Tài liệu: Smart Constract cho DApp Meeting với WebRTC

### Phân tích yêu cầu
Hợp đồng thông minh (Smart Contract) cần quản lý phòng họp video phi tập trung kết hợp với WebRTC, với các đặc điểm sau:
- **Thông tin phòng (Room)**: Lưu trữ `roomID` (string), danh sách người tham gia (`Participant[]`), và thời gian tạo (`uint256`).
- **Thông tin người tham gia (Participant)**: Bao gồm địa chỉ ví (`address`), tên hiển thị (`string`), `sessionID` (string, do Cloudflare Calls cung cấp, được cập nhật bởi backend), và danh sách track (`Track[]`).
- **Thông tin track (Track)**: Gồm `trackName` (string, như "audio", "video", "screen"), `mid` (string, định danh media trong WebRTC), `location` (string, vị trí "local" hoặc "remote"), `isPublished` (bool, trạng thái đã publish hay chưa), `sessionId` (string, liên kết với session), `roomId` (string, phòng chứa track) và dữ liệu bổ sung.
- **Yêu cầu chức năng**:
  - Cơ chế tạo phòng và quản lý người tham gia.
  - Khả năng tạo và lưu trữ thông tin về các track media.
  - Cơ chế để frontend gửi và nhận thông tin WebRTC qua backend.
  - Hệ thống sự kiện để thông báo khi người dùng tham gia/rời phòng và thêm track mới.

### Thiết kế cấu trúc dữ liệu
- **Track**: Lưu trữ thông tin chi tiết của một track media.
- **Participant**: Đại diện cho một người tham gia trong phòng, bao gồm thông tin cá nhân và danh sách track.
- **Room**: Quản lý thông tin phòng họp, bao gồm danh sách người tham gia.
- **Mappings**:
  - `mapping(string => Room) public rooms`: Lưu trữ thông tin phòng theo `roomId`.
  - `mapping(string => mapping(address => uint)) public participantIndices`: Theo dõi chỉ số của người tham gia trong mảng `participants`.
  - `mapping(string => mapping(address => bool)) public participantsInRoom`: Kiểm tra người tham gia có trong phòng không.

### Cấu trúc dữ liệu
```solidity
struct Track {
    string trackName;   // Tên track (audio, video, screen)
    string mid;         // Định danh media trong WebRTC
    string location;    // Vị trí của track (local/remote)
    bool isPublished;   // Trạng thái đã publish chưa
    string sessionId;   // SessionID liên kết với track
    string roomId;      // ID phòng chứa track
}

struct Participant {
    address walletAddress;  // Địa chỉ ví người tham gia
    string name;            // Tên hiển thị
    string sessionID;       // Session ID từ Cloudflare Calls 
    Track[] tracks;         // Danh sách track media
}

struct Room {
    string roomId;          // ID phòng họp
    uint256 creationTime;   // Thời gian tạo (timestamp)
    Participant[] participants; // Danh sách người tham gia
}

// Cấu trúc ParticipantDetails sử dụng để trả về chi tiết người tham gia
struct ParticipantDetails {
    address walletAddress;  // Địa chỉ ví
    string name;            // Tên hiển thị
    string sessionID;       // Session ID
    Track[] tracks;         // Danh sách track
}
```

### Các hàm chính

1. **`createRoom(string memory _roomId)`**  
   - **Mô tả**: Tạo phòng họp mới với ID và timestamp tạo.
   - **Tham số**: `_roomId` - ID duy nhất cho phòng (thường là UUID).
   - **Xử lý**: Khởi tạo cấu trúc Room và lưu vào mapping rooms.
   - **Điều kiện**: Phòng chưa tồn tại.

2. **`joinRoom(string memory _roomId, string memory _name, Track[] memory _initialTracks, bytes sessionDescription)`**  
   - **Mô tả**: Đăng ký người tham gia vào phòng họp với thông tin ban đầu.
   - **Tham số**:
     - `_roomId`: ID phòng họp 
     - `_name`: Tên hiển thị
     - `_initialTracks`: Danh sách track ban đầu (audio, video)
     - `sessionDescription`: Dữ liệu SDP offer từ WebRTC client
   - **Xử lý**: Tạo người tham gia mới, phát sự kiện ParticipantJoined để backend biết và tạo session.
   - **Sự kiện**: Phát sự kiện ParticipantJoined và TrackAdded cho mỗi track.

3. **`leaveRoom(string memory _roomId)`**  
   - **Mô tả**: Cho phép người tham gia rời phòng họp.
   - **Tham số**: `_roomId` - ID phòng họp.
   - **Xử lý**: Xóa người tham gia và track ra khỏi phòng, cập nhật các mapping liên quan.
   - **Sự kiện**: Phát sự kiện ParticipantLeft để thông báo cho backend và người tham gia khác.

4. **`setParticipantSessionID(string memory _roomId, address _participantAddress, string memory _sessionID)`**  
   - **Mô tả**: Cập nhật sessionID cho người tham gia (gọi bởi backend sau khi tạo session).
   - **Tham số**:
     - `_roomId`: ID phòng họp
     - `_participantAddress`: Địa chỉ ví của người tham gia
     - `_sessionID`: Session ID từ Cloudflare Calls
   - **Điều kiện**: Chỉ những địa chỉ backend được ủy quyền mới gọi được (onlyAuthorizedBackend).

5. **`addTrack(string memory _roomId, Track memory _newTrack)`**  
   - **Mô tả**: Thêm track mới vào danh sách track của người tham gia.
   - **Tham số**:
     - `_roomId`: ID phòng họp
     - `_newTrack`: Thông tin track cần thêm
   - **Sự kiện**: Phát sự kiện TrackAdded để thông báo cho backend và người tham gia khác.

6. **`getRoomParticipantsDetails(string memory _roomId)`**  
   - **Mô tả**: Lấy thông tin chi tiết về tất cả người tham gia trong phòng.
   - **Tham số**: `_roomId` - ID phòng họp.
   - **Kết quả trả về**: Mảng ParticipantDetails chứa thông tin chi tiết của mỗi người tham gia.

7. **`forwardEventToBackend(string memory _roomId, bytes memory _eventData)`**  
   - **Mô tả**: Cho phép frontend gửi dữ liệu sự kiện đến backend qua smart contract.
   - **Tham số**:
     - `_roomId`: ID phòng họp
     - `_eventData`: Dữ liệu sự kiện dạng bytes (thường là dữ liệu nén)
   - **Sự kiện**: Phát EventForwardedToBackend để backend biết và xử lý.
   - **Ứng dụng**: Dùng để gửi thông tin publish/pull track, renegotiate, v.v.

8. **`forwardEventToFrontend(string memory _roomId, address _participant, bytes memory _eventData)`**  
   - **Mô tả**: Cho phép backend gửi dữ liệu sự kiện đến frontend cụ thể.
   - **Tham số**:
     - `_roomId`: ID phòng họp
     - `_participant`: Địa chỉ ví của người tham gia cần nhận thông báo
     - `_eventData`: Dữ liệu sự kiện dạng bytes (thường là dữ liệu nén)
   - **Điều kiện**: Chỉ những địa chỉ backend được ủy quyền mới gọi được.
   - **Sự kiện**: Phát EventForwardedToFrontend để frontend nhận và xử lý.
   - **Ứng dụng**: Dùng để gửi SDP answer, thông tin session, kết quả publish/pull track.

### Sự kiện tương tác

- **`ParticipantJoined(string roomId, address participant, bytes sessionDescription, Track[] memory initialTracks)`**  
  - **Mô tả**: Phát khi người dùng tham gia phòng họp qua hàm `joinRoom`.
  - **Tham số**: 
    - `roomId`: ID phòng họp
    - `participant`: Địa chỉ ví người tham gia
    - `sessionDescription`: Dữ liệu SDP offer từ WebRTC client
    - `initialTracks`: Danh sách track ban đầu người dùng đăng ký
  - **Người nhận**: Backend lắng nghe sự kiện này để tạo session mới trên Cloudflare Calls và publish các track ban đầu.

- **`ParticipantLeft(string roomId, address participant, string sessionId)`**  
  - **Mô tả**: Phát khi người dùng rời phòng họp qua hàm `leaveRoom`.
  - **Tham số**:
    - `roomId`: ID phòng họp
    - `participant`: Địa chỉ ví người tham gia rời phòng
    - `sessionId`: Session ID của người tham gia
  - **Người nhận**: Backend (để đóng session) và các frontend khác (để cập nhật UI và xóa video element).

- **`TrackAdded(string roomId, address participant, string trackName, string sessionId)`**  
  - **Mô tả**: Phát khi một track mới được thêm vào thông qua `joinRoom` hoặc `addTrack`.
  - **Tham số**:
    - `roomId`: ID phòng họp
    - `participant`: Địa chỉ ví người thêm track
    - `trackName`: Tên track ("audio", "video", "screen")
    - `sessionId`: Session ID của người thêm track
  - **Người nhận**: Backend (để cập nhật track trong session) và các frontend khác (để biết có track mới có thể pull).

- **`EventForwardedToBackend(string roomId, address sender, bytes eventData)`**  
  - **Mô tả**: Được phát khi frontend gọi `forwardEventToBackend` để gửi thông tin đến backend.
  - **Tham số**:
    - `roomId`: ID phòng họp
    - `sender`: Địa chỉ ví người gửi yêu cầu
    - `eventData`: Dữ liệu nén dạng bytes chứa thông tin sự kiện (publish track, pull track, renegotiate...)
  - **Người nhận**: Backend lắng nghe để nhận và xử lý các yêu cầu từ frontend.

- **`EventForwardedToFrontend(string roomId, address participant, bytes eventData)`**  
  - **Mô tả**: Được phát khi backend gọi `forwardEventToFrontend` để gửi phản hồi đến frontend.
  - **Tham số**: 
    - `roomId`: ID phòng họp
    - `participant`: Địa chỉ ví người nhận thông báo
    - `eventData`: Dữ liệu nén dạng bytes chứa thông tin phản hồi (SDP answer, kết quả publish/pull track...)
  - **Người nhận**: Frontend cụ thể có địa chỉ ví trùng với `participant`.

### Thiết kế chi tiết chức năng

1. **`createRoom(string memory _roomId)`**
   - **Kiểm tra điều kiện**:
     - Phòng chưa tồn tại: `require(rooms[_roomId].roomId == "", "Room already exists")`
   - **Thực hiện**:
     - Khởi tạo Room mới với roomId và timestamp tạo: `rooms[_roomId] = Room(_roomId, block.timestamp, new Participant[](0))`
   - **Quy trình**:
     - Frontend tạo UUID làm roomId
     - Gọi `createRoom` với roomId này
     - Lưu roomId vào localStorage để sử dụng sau

2. **`joinRoom(string memory _roomId, string memory _name, Track[] memory _initialTracks, bytes memory sessionDescription)`**
   - **Kiểm tra điều kiện**:
     - Phòng tồn tại: `require(rooms[_roomId].roomId != "", "Room does not exist")`
     - Người dùng chưa tham gia: `require(!participantsInRoom[_roomId][msg.sender], "Already in room")`
   - **Thực hiện**:
     - Tạo participant mới với địa chỉ ví, tên và sessionID rỗng
     - Thêm participant vào danh sách participants của phòng
     - Cập nhật mappings `participantIndices` và `participantsInRoom`
     - Phát sự kiện `ParticipantJoined` với roomId, địa chỉ ví, sessionDescription và danh sách track ban đầu
     - Với mỗi track trong _initialTracks, thêm vào danh sách tracks và phát sự kiện `TrackAdded`
   - **Quy trình WebRTC**:
     - Frontend tạo RTCPeerConnection và thêm local tracks
     - Frontend tạo offer và set local description
     - Frontend nén offer thành bytes và gọi `joinRoom`
     - Backend lắng nghe sự kiện `ParticipantJoined`, tạo session Cloudflare
     - Backend gọi `setParticipantSessionID` và `forwardEventToFrontend` với SDP answer

3. **`leaveRoom(string memory _roomId)`**
   - **Kiểm tra điều kiện**:
     - Phòng tồn tại: `require(rooms[_roomId].roomId != "", "Room does not exist")`
     - Người dùng đang trong phòng: `require(participantsInRoom[_roomId][msg.sender], "Not in room")`
   - **Thực hiện**:
     - Xóa participant khỏi danh sách, dùng `participantIndices` để biết vị trí
     - Cập nhật lại các mappings và array
     - Phát sự kiện `ParticipantLeft` với roomId, địa chỉ ví và sessionId
   - **Quy trình đồng bộ**:
     - Frontend đóng RTCPeerConnection trước hoặc sau khi gọi `leaveRoom`
     - Backend nhận sự kiện ParticipantLeft và cập nhật trạng thái
     - Các frontend khác nhận sự kiện ParticipantLeft và xóa video element

4. **`setParticipantSessionID(string memory _roomId, address _participantAddress, string memory _sessionID)`**
   - **Kiểm tra điều kiện**:
     - Chỉ backend được ủy quyền: `modifier onlyAuthorizedBackend`
     - Phòng tồn tại và người tham gia hợp lệ
   - **Thực hiện**:
     - Tìm participant với địa chỉ _participantAddress
     - Cập nhật sessionID: `participants[participantIdx].sessionID = _sessionID`
   - **Quy trình**:
     - Backend nhận sự kiện `ParticipantJoined`
     - Backend tạo session Cloudflare và nhận sessionID
     - Backend gọi hàm này để cập nhật sessionID

5. **`getRoomParticipantsDetails(string memory _roomId)`**
   - **Kiểm tra điều kiện**:
     - Phòng tồn tại: `require(rooms[_roomId].roomId != "", "Room does not exist")`
   - **Thực hiện**:
     - Tạo mảng ParticipantDetails với kích thước bằng số người tham gia
     - Với mỗi participant, tạo ParticipantDetails chứa thông tin và track
     - Trả về mảng ParticipantDetails
   - **Quy trình**:
     - Frontend gọi hàm này khi vào phòng để lấy danh sách người tham gia và track
     - Frontend tạo các video element và bắt đầu pull track từ người khác

6. **`addTrack(string memory _roomId, Track memory _newTrack)`**
   - **Kiểm tra điều kiện**:
     - Phòng tồn tại và người gọi là người tham gia
   - **Thực hiện**:
     - Thêm track mới vào danh sách tracks của participant
     - Phát sự kiện `TrackAdded` với roomId, địa chỉ ví, tên track và sessionId
   - **Quy trình**:
     - Frontend thêm track mới (ví dụ: bật camera sau khi vào phòng với mic)
     - Backend nhận sự kiện TrackAdded và cập nhật danh sách track
     - Các frontend khác biết có track mới và có thể pull

### Quản lý sự kiện tương tác
- **Client → Backend**:  
  - Gọi `joinRoom` với danh sách track ban đầu → Phát `ParticipantJoined` và `TrackAdded` cho từng track → Backend xử lý (tạo session trên Cloudflare Calls).  
  - Gọi `addTrack` → Phát `TrackAdded` → Backend cập nhật thông tin WebRTC.  
- **Backend → Hợp đồng**:  
  - Backend gọi `setParticipantSessionID` để cập nhật `sessionID` sau khi xử lý sự kiện từ client.

### Ví dụ luồng hoạt động chi tiết

#### 1. Tạo và tham gia phòng họp mới
```
Frontend                     Smart Contract                 Backend                      Cloudflare
   |                              |                            |                              |
   |-- createRoom("room123") ---->|                            |                              |
   |                              |                            |                              |
   |<------ Transaction confirmed-|                            |                              |
   |                              |                            |                              |
   |-- Tạo RTCPeerConnection -----|                            |                              |
   |-- getUserMedia() ------------|                            |                              |
   |-- createOffer() -------------|                            |                              |
   |-- setLocalDescription() -----|                            |                              |
   |                              |                            |                              |
   |-- joinRoom("room123",        |                            |                              |
   |   "Manh", [{audio,video}],   |                            |                              |
   |   compressedOfferSDP) ------>|                            |                              |
   |                              |-- ParticipantJoined ------>|                              |
   |                              |                            |-- CreateSession() ---------->|
   |                              |                            |<-- sessionId: "sess123" -----|
   |                              |                            |                              |
   |                              |                            |-- PublishTracks(offerSDP) -->|
   |                              |                            |<-- answerSDP ----------------|
   |                              |                            |                              |
   |                              |<- setParticipantSessionID -|                              |
   |                              |   ("room123", 0x123,       |                              |
   |                              |    "sess123")              |                              |
   |                              |                            |                              |
   |                              |<- forwardEventToFrontend --|                              |
   |                              |   ("room123", 0x123,       |                              |
   |                              |    answerSDP)              |                              |
   |                              |                            |                              |
   |<-- EventForwardedToFrontend -|                            |                              |
   |      (answerSDP)             |                            |                              |
   |-- setRemoteDescription() ----|                            |                              |
   |-- Kết nối WebRTC thiết lập --|                            |                              |
```

#### 2. Thêm track mới và pull track từ người tham gia khác
```
Frontend A                 Smart Contract                 Backend                  Frontend B
   |                            |                             |                         |
   |-- addTrack("room123", --->|                             |                         |
   |   {video, mid:"m2"...})   |                             |                         |
   |                            |-- TrackAdded ------------->|                         |
   |                            |   ("room123", 0xA, "video")|                         |
   |                            |                             |-- Cập nhật track list ->|
   |                            |                             |                         |
   |<---- Transaction confirmed-|                             |                         |
   |                            |                             |                         |
   |                            |                             |                         |
   |                            |<--- forwardEventToBackend --|                         |
   |                            |     ("room123", 0xB,        |                         |
   |                            |      pullTrackData)         |                         |
   |                            |                             |                         |
   |                            |-- EventForwardedToBackend ->|                         |
   |                            |                             |-- Pull track từ        |
   |                            |                             |   Cloudflare ---------->|
   |                            |                             |                         |
   |                            |<- forwardEventToFrontend ---|                         |
   |                            |   ("room123", 0xB,          |                         |
   |                            |    pullResultData)          |                         |
   |                            |                             |                         |
   |                            |                             |                         |
   |                            |                             |                         |
```

#### 3. Rời phòng họp
```
Frontend                     Smart Contract                 Backend                   
   |                              |                            |                             
   |-- leaveRoom("room123") ----->|                            |
   |                              |-- ParticipantLeft -------->|
   |                              |   ("room123", 0x123,       |
   |                              |    "sess123")              |
   |<----- Transaction confirmed -|                            |
   |                              |                            |
   |-- close RTCPeerConnection ---|                            |
   |-- Chuyển hướng -------------|                            |
   |   về trang chủ              |                            |
   |                              |                            |
```

#### 4. Luồng xử lý sự kiện và giao tiếp giữa các thành phần

- **Frontend đến Smart Contract**: Frontend gọi hàm trong Smart Contract như `createRoom`, `joinRoom`, `leaveRoom`
- **Smart Contract đến Backend**: Phát sự kiện `ParticipantJoined`, `ParticipantLeft`, `TrackAdded`, `EventForwardedToBackend`
- **Backend đến Smart Contract**: Gọi hàm `setParticipantSessionID`, `forwardEventToFrontend`
- **Smart Contract đến Frontend**: Phát sự kiện `EventForwardedToFrontend`
- **Backend đến Cloudflare**: Gọi API `CreateSession`, `PublishTracks`, `PullTracks`
- **Cloudflare đến Backend**: Trả về SDP answer và thông tin kết nối

Mỗi khi frontend cần thực hiện tác vụ như publish track mới hoặc pull track từ người khác:
1. Frontend nén dữ liệu yêu cầu và gọi `forwardEventToBackend`
2. Backend nhận sự kiện, giải nén và xử lý với Cloudflare
3. Backend nén kết quả và gọi `forwardEventToFrontend` 
4. Frontend nhận sự kiện, giải nén và cập nhật kết nối WebRTC

### Bảng tổng hợp chức năng và sự kiện

#### Bảng 1: Tổng hợp các hàm và phương thức
| Hàm                       | Mô tả                                      | Ai gọi       | Phát sự kiện                       | Ứng dụng thực tế                   |
|---------------------------|--------------------------------------------|--------------|------------------------------------|-----------------------------------|
| `createRoom`              | Tạo phòng mới với ID                       | Frontend     | Không                              | Khởi tạo phòng họp mới            |
| `joinRoom`                | Tham gia phòng với tracks ban đầu          | Frontend     | `ParticipantJoined`, `TrackAdded`  | Bắt đầu kết nối WebRTC           |
| `leaveRoom`               | Rời khỏi phòng                             | Frontend     | `ParticipantLeft`                  | Kết thúc cuộc họp                |
| `setParticipantSessionID` | Cập nhật sessionID từ Cloudflare           | Backend      | Không                              | Liên kết participant với session  |
| `addTrack`                | Thêm track mới (video, audio, screen)      | Frontend     | `TrackAdded`                       | Bật camera/mic sau khi vào phòng  |
| `getRoomParticipantsDetails` | Lấy thông tin người tham gia và tracks  | Frontend     | Không                              | Hiển thị danh sách người tham gia |
| `forwardEventToBackend`   | Gửi yêu cầu từ frontend đến backend        | Frontend     | `EventForwardedToBackend`          | Publish/Pull tracks, Renegotiate  |
| `forwardEventToFrontend`  | Gửi phản hồi từ backend đến frontend       | Backend      | `EventForwardedToFrontend`         | Trả SDP answer, kết quả xử lý    |
| `addNewTrackAfterPublish` | Thêm track mới sau khi publish thành công  | Backend      | Không                              | Cập nhật track đã publish lên contract |

#### Bảng 2: Tổng hợp các sự kiện và xử lý
| Sự kiện                    | Phát ra khi                             | Người lắng nghe | Xử lý                              | Dữ liệu chính                 |
|----------------------------|----------------------------------------|----------------|-----------------------------------|-------------------------------|
| `ParticipantJoined`        | Người dùng tham gia phòng               | Backend        | Tạo session Cloudflare, publish tracks | sessionDescription (SDP offer) |
| `ParticipantLeft`          | Người dùng rời phòng                    | Backend & Frontend | Đóng session, cập nhật UI      | roomId, địa chỉ ví, sessionId |
| `TrackAdded`               | Thêm track mới                          | Backend & Frontend | Cập nhật danh sách track, thông báo người khác | trackName, sessionId |
| `EventForwardedToBackend`  | Frontend gửi yêu cầu qua smart contract | Backend        | Xử lý với Cloudflare (publish, pull, v.v.) | eventData (nén) |
| `EventForwardedToFrontend` | Backend gửi phản hồi qua smart contract | Frontend       | Cập nhật kết nối WebRTC           | eventData (nén) |

#### Bảng 3: Quy trình xử lý các tác vụ chính
| Tác vụ                | Frontend                                   | Smart Contract                       | Backend                               | Cloudflare              |
|-----------------------|-------------------------------------------|-------------------------------------|---------------------------------------|-------------------------|
| Tạo phòng họp         | Tạo UUID, gọi createRoom                  | Lưu thông tin phòng                 | -                                     | -                       |
| Tham gia phòng        | Tạo offer, gọi joinRoom                   | Phát ParticipantJoined             | Tạo session, publish tracks, gửi answer | Tạo session, xử lý offer |
| Thêm track mới        | Tạo MediaStreamTrack, gọi addTrack        | Phát TrackAdded                    | Cập nhật danh sách track              | -                       |
| Publish track         | Tạo offer, gọi forwardEventToBackend      | Phát EventForwardedToBackend       | Publish track trên Cloudflare, gửi answer | Xử lý publish request  |
| Pull track            | Gọi forwardEventToBackend với track info  | Phát EventForwardedToBackend       | Pull track từ Cloudflare, gửi kết quả  | Thiết lập kết nối giữa sessions |
| Rời phòng             | Đóng RTCPeerConnection, gọi leaveRoom     | Phát ParticipantLeft               | Đóng session Cloudflare               | Xóa session            |
