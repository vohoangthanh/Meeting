## Tài liệu: Smart Constract cho DApp Meeting với WebRTC

### Phân tích yêu cầu
Yêu cầu xây dựng một hợp đồng thông minh để quản lý phòng họp cho DApp sử dụng WebRTC, với các đặc điểm sau:
- **Thông tin phòng (Room)**: Lưu trữ `roomID` (string), danh sách người tham gia (`Participant[]`), và thời gian tạo (`uint256`).
- **Thông tin người tham gia (Participant)**: Bao gồm địa chỉ ví (`address`), tên (`string`), `sessionID` (string, do Cloudflare Calls cung cấp, không tạo trong hợp đồng), và danh sách track (`Track[]`).
- **Thông tin track (Track)**: Gồm `trackName` (string), `mid` (string), `location` (string), `isPublished` (bool), và `trackData` (bytes).
- **Yêu cầu chức năng**:
  - Một hàm cho client gọi để phát sự kiện đến backend (ví dụ: tham gia phòng).
  - Một hàm cho backend gọi để cập nhật thông tin ngược lại (ví dụ: thiết lập `sessionID`).
  - **Yêu cầu mới**: Khi người tham gia gia nhập phòng, lưu track ban đầu của họ và phát sự kiện `TrackAdded`.

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
    string trackName;   // Tên của track
    string mid;         // Định danh media trong WebRTC
    string location;    // Vị trí của track
    bool isPublished;   // Trạng thái công khai
    bytes trackData;    // Dữ liệu nhị phân của track
}

struct Participant {
    address walletAddress;  // Địa chỉ ví của người tham gia
    string name;            // Tên người tham gia
    string sessionID;       // ID phiên từ Cloudflare Calls
    Track[] tracks;         // Danh sách các track
}

struct Room {
    string roomId;          // ID của phòng họp
    uint256 creationTime;   // Thời gian tạo (timestamp)
    Participant[] participants; // Danh sách người tham gia
}
```

### Các hàm chính
1. **`createRoom(string memory _roomId)`**  
   - Tạo phòng mới với `roomId` và thời gian tạo (`block.timestamp`).

2. **`joinRoom(string memory _roomId, string memory _name, Track[] memory _initialTracks, bytes sessionDescription)`**  
   - Thêm người tham gia vào phòng, lưu danh sách track ban đầu và phát sự kiện đến backend.

3. **`leaveRoom(string memory _roomId)`**  
   - Cho phép người tham gia rời khỏi phòng, xóa thông tin liên quan và phát sự kiện thông báo.

4. **`setParticipantSessionID(string memory _roomId, address _participantAddress, string memory _sessionID)`**  
   - Cập nhật `sessionID` cho người tham gia, chỉ backend được phép gọi.

5. **`addTrack(string memory _roomId, Track memory _newTrack)`**  
   - Thêm track mới vào danh sách track của người tham gia.

6. **`getParticipantInfo(string memory _roomId)`**  
   - Trả về thông tin của người tham gia (bao gồm `sessionID`).

7. **`forwardEventToBackend(string memory _roomId, string memory _eventData)`**  
   - **Mô tả**: Cho phép frontend gửi sự kiện tới backend thông qua smart contract bằng cách sử dụng dữ liệu có sẵn (room ID và dữ liệu sự kiện).  
   - **Ai gọi**: Frontend (bất kỳ người tham gia nào trong phòng).  
   - **Hành động**: Phát sự kiện để backend nhận và xử lý.

8. **`forwardEventToFrontend(string memory _roomId, address _participant, string memory _eventData)`**  
   - **Mô tả**: Cho phép backend gửi sự kiện tới frontend thông qua smart contract bằng cách sử dụng dữ liệu có sẵn (room ID, địa chỉ người tham gia, và dữ liệu sự kiện).  
   - **Ai gọi**: Backend (yêu cầu ủy quyền).  
   - **Hành động**: Phát sự kiện để frontend nhận và cập nhật.

### Sự kiện tương tác
- **`ParticipantJoined(string roomId, address participant, bytes Trackdata, Track[] memory _initialTracks, bytes sessionDescription)`**  
  - Phát khi người tham gia gia nhập phòng, thông báo cho backend.
- **`ParticipantLeft(string roomId, address participant)`**  
  - Phát khi người tham gia rời khỏi phòng, thông báo cho backend và các người tham gia khác.
- **`TrackAdded(string roomId, address participant, string trackName)`**  
  - Phát khi người tham gia thêm track mới (bao gồm khi gia nhập phòng với track ban đầu).
- **`EventForwardedToBackend(string roomId, address sender, string eventData)`**  
  - **Mô tả**: Được phát khi frontend gọi `forwardEventToBackend`, gửi thông tin phòng, địa chỉ người gửi và dữ liệu sự kiện tới backend.
- **`EventForwardedToFrontend(string roomId, address participant, string eventData)`**  
  - **Mô tả**: Được phát khi backend gọi `forwardEventToFrontend`, gửi thông tin phòng, địa chỉ người tham gia và dữ liệu sự kiện tới frontend.

### Thiết kế chức năng
1. **`createRoom(string memory _roomId)`**  
   - **Điều kiện**: Phòng chưa tồn tại (`rooms[_roomId].roomId == ""`).  
   - **Hành động**: Khởi tạo `Room` với `roomId` và `creationTime`.

2. **`joinRoom(string memory _roomId, string memory _name, Track[] memory _initialTracks, bytes sessionDescription)`**  
   - **Điều kiện**: Phòng tồn tại, người tham gia chưa có trong phòng.  
   - **Hành động**:  
     - Thêm `Participant` với `walletAddress = msg.sender`, `name = _name`, `sessionID = ""`, và danh sách track ban đầu `_initialTracks`.  
     - Phát sự kiện `ParticipantJoined`.  
     - Đối với mỗi track trong `_initialTracks`, thêm track vào danh sách và phát sự kiện `TrackAdded`.

3. **`leaveRoom(string memory _roomId)`**  
   - **Điều kiện**: Phòng tồn tại, người gọi đang là người tham gia trong phòng.  
   - **Hành động**:  
     - Xóa thông tin người tham gia khỏi danh sách participants trong phòng.  
     - Cập nhật mappings `participantIndices` và `participantsInRoom`.  
     - Phát sự kiện `ParticipantLeft` để thông báo cho backend và những người tham gia khác.

4. **`setParticipantSessionID(string memory _roomId, address _participantAddress, string memory _sessionID)`**  
   - **Điều kiện**: Chỉ backend được ủy quyền gọi (modifier `onlyAuthorized`), phòng tồn tại, người tham gia hợp lệ.  
   - **Hành động**: Cập nhật `sessionID` cho người tham gia.

5. **`addTrack(string memory _roomId, Track memory _newTrack)`**  
   - **Điều kiện**: Phòng tồn tại, người gọi là người tham gia hợp lệ.  
   - **Hành động**: Thêm `_newTrack` vào danh sách track, phát sự kiện `TrackAdded`.

6. **`getParticipantInfo(string memory _roomId)`**  
   - **Điều kiện**: Phòng tồn tại, người gọi là người tham gia.  
   - **Hành động**: Trả về thông tin `Participant` của `msg.sender`.

### Quản lý sự kiện tương tác
- **Client → Backend**:  
  - Gọi `joinRoom` với danh sách track ban đầu → Phát `ParticipantJoined` và `TrackAdded` cho từng track → Backend xử lý (tạo session trên Cloudflare Calls).  
  - Gọi `addTrack` → Phát `TrackAdded` → Backend cập nhật thông tin WebRTC.  
- **Backend → Hợp đồng**:  
  - Backend gọi `setParticipantSessionID` để cập nhật `sessionID` sau khi xử lý sự kiện từ client.

### Ví dụ triển khai
1. **Client tham gia phòng với track ban đầu**:  
   - Gọi `joinRoom("room1", "Nguyen Van A", [{trackName: "audio", mid: "m0", location: "local", isPublished: true, trackData: 0x...}])`.  
   - Phát `ParticipantJoined("room1", 0x123...)`.  
   - Phát `TrackAdded("room1", 0x123..., "audio")`.  
   - Backend nhận sự kiện, tạo `sessionID = "session123"`.

2. **Backend cập nhật sessionID**:  
   - Gọi `setParticipantSessionID("room1", 0x123..., "session123")`.  
   - `sessionID` của người tham gia được cập nhật.

3. **Client kiểm tra thông tin**:  
   - Gọi `getParticipantInfo("room1")`.  
   - Nhận thông tin bao gồm `sessionID = "session123"` và danh sách track `[{trackName: "audio", ...}]`.

4. **Thêm track**:  
   - Gọi `addTrack("room1", {trackName: "video", mid: "m1", location: "local", isPublished: true, trackData: 0x...})`.  
   - Phát `TrackAdded("room1", 0x123..., "video")`.

5. **Frontend gửi sự kiện tới backend**:  
   - Frontend gọi hàm với dữ liệu như `"room1"` và `"start_meeting"`.  
   - Smart contract phát sự kiện chứa thông tin này.  
   - Backend nhận sự kiện và thực hiện hành động (ví dụ: bắt đầu phiên họp).

6. **Backend gửi sự kiện tới frontend**:  
   - Backend gọi hàm với dữ liệu như `"room1"`, địa chỉ người tham gia, và `"meeting_started"`.  
   - Smart contract phát sự kiện chứa thông tin này.  
   - Frontend nhận sự kiện và cập nhật giao diện (ví dụ: hiển thị trạng thái họp).

### Bảng tổng hợp
| Hàm                     | Mô tả                                      | Ai gọi       | Phát sự kiện                       | Kết quả chính                     |
|-------------------------|--------------------------------------------|--------------|------------------------------------|-----------------------------------|
| `createRoom`            | Tạo phòng mới                              | Client       | Không                              | Phòng mới được khởi tạo           |
| `joinRoom`              | Tham gia phòng, thông báo backend          | Client       | `ParticipantJoined`, `TrackAdded`  | Người tham gia và track được thêm |
| `leaveRoom`             | Rời khỏi phòng, thông báo backend          | Client       | `ParticipantLeft`                  | Người tham gia bị xóa khỏi phòng  |
| `setParticipantSessionID`| Cập nhật `sessionID` từ backend           | Backend      | Không                              | `sessionID` được thiết lập        |
| `addTrack`              | Thêm track mới, thông báo backend          | Client       | `TrackAdded`                       | Track được thêm vào danh sách     |
| `getParticipantInfo`    | Lấy thông tin người tham gia               | Client       | Không                              | Trả về thông tin `Participant`    |
| `forwardEventToBackend` | Gửi sự kiện từ frontend tới backend        | Frontend     | `EventForwardedToBackend`          | Sự kiện được chuyển tới backend   |
| `forwardEventToFrontend`| Gửi sự kiện từ backend tới frontend        | Backend      | `EventForwardedToFrontend`         | Sự kiện được chuyển tới frontend  |
