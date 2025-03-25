### Key Points  
- Hệ thống backend cho DApp meeting với WebRTC có hai phần: Phần 1 lắng nghe sự kiện từ smart contract và gọi Phần 2 để xử lý; Phần 2 làm việc với Cloudflare và gọi lại smart contract để gửi sự kiện về frontend.  
- Khi người tham gia gia nhập phòng, Phần 1 nhận sự kiện "ParticipantJoined" từ smart contract, lấy thông tin track và sessionDescription, rồi gọi Phần 2 để tạo session và publish track trên Cloudflare.  
- Phần 2 sau khi xử lý với Cloudflare, yêu cầu Phần 1 gọi smart contract gửi sự kiện về frontend với session ID và phản hồi từ Cloudflare.  
- Hệ thống quản lý danh sách ví để đảm bảo mỗi ví chỉ thực hiện một giao dịch tại một thời điểm, với cơ chế phân luồng và hàng đợi nếu tất cả ví bận.  

---

### Cấu trúc Hệ thống Backend  

#### Tổng quan  
Hệ thống backend được thiết kế để hỗ trợ DApp meeting sử dụng WebRTC, với hai phần chính:  
- **Phần 1 (Event Listener và Transaction Manager)**: Lắng nghe sự kiện từ smart contract, quản lý danh sách ví để thực hiện giao dịch, và gọi Phần 2 để xử lý tương tác với Cloudflare.  
- **Phần 2 (Cloudflare Interaction)**: Làm việc với API Cloudflare để quản lý session và track, sau đó yêu cầu Phần 1 gọi smart contract gửi sự kiện về frontend.  

Hệ thống đảm bảo rằng mỗi ví chỉ thực hiện một giao dịch với smart contract tại một thời điểm, sử dụng cơ chế phân luồng và hàng đợi để xử lý nhiều yêu cầu đồng thời.  

#### Chi tiết Phần 1: Event Listener và Transaction Manager  
- **Kết nối Blockchain và Lắng nghe Sự kiện**:  
  - Sử dụng thư viện như go-ethereum để kết nối với blockchain (giả sử Ethereum hoặc chuỗi tương thích EVM).  
  - Lắng nghe sự kiện "ParticipantJoined" từ smart contract, bao gồm room ID, địa chỉ ví, và sessionDescription.  
  - Khi nhận sự kiện, gọi hàm `getParticipantInfo` để lấy danh sách track ban đầu của người tham gia, rồi gọi Phần 2 với dữ liệu (room ID, địa chỉ ví, track, sessionDescription).  

- **Quản lý Ví và Giao Dịch**:  
  - Duy trì danh sách ví, mỗi ví có:  
    - Địa chỉ ví.  
    - Khóa riêng tư (xử lý an toàn, không hardcode).  
    - Trạng thái: "available" (trống) hoặc "busy" (bận).  
  - Khi cần gọi smart contract (ví dụ: gửi sự kiện về frontend qua `forwardEventToFrontEnd`), Phần 1:  
    - Chọn ví trống (status = "available").  
    - Nếu không có ví trống, xếp yêu cầu vào hàng đợi.  
    - Chuẩn bị giao dịch: mã hóa gọi hàm smart contract, ký bằng khóa riêng tư của ví, và gửi lên blockchain.  
    - Đánh dấu ví là "busy" trong quá trình gửi, và trở lại "available" khi giao dịch được xác nhận hoặc hết thời gian chờ.  
  - Xử lý xác nhận giao dịch: Sử dụng cơ chế chờ (wait for receipt) để đảm bảo giao dịch thành công trước khi tái sử dụng ví.  

- **Hàng Đợi Yêu Cầu**:  
  - Nếu tất cả ví đang bận, yêu cầu được thêm vào hàng đợi (dạng slice hoặc channel).  
  - Khi có ví trống, lấy yêu cầu tiếp theo từ hàng đợi và xử lý.  

#### Chi tiết Phần 2: Cloudflare Interaction  
- **Tương tác với Cloudflare**:  
  - Sử dụng `CloudflareService` đã cung cấp, với các API chính:  
    - `CreateSession`: Tạo session mới, trả về session ID.  
    - `PublishTracks`: Publish track vào session, với thông tin offer và danh sách track.  
    - `PullTracks`, `Renegotiate`, `CloseTracks`, `GetSessionState`: Các API bổ sung cho quản lý session.  
  - Khi nhận dữ liệu từ Phần 1 (room ID, địa chỉ ví, danh sách track, sessionDescription), Phần 2:  
    - Gọi `CreateSession` để tạo session ID.  
    - Chuẩn bị track để publish, định dạng theo yêu cầu API Cloudflare, sử dụng sessionDescription làm offer.  
    - Gọi `PublishTracks` để publish track vào session, nhận phản hồi.  

- **Gửi Sự Kiện Về Frontend**:  
  - Sau khi nhận phản hồi từ Cloudflare, chuẩn bị dữ liệu sự kiện dưới dạng chuỗi JSON, ví dụ:  
    - `{"sessionID": "abc123", "type": "publish-track","cloudflareResponse": {"answer": "someAnswer"}}`.  
  - Yêu cầu Phần 1 gọi hàm `forwardEventToFrontEnd` trong smart contract với:  
    - room ID.  
    - Địa chỉ ví người tham gia.  
    - Dữ liệu sự kiện (chuỗi JSON).  
  - Phần 1 sẽ chọn ví trống, gửi giao dịch, và đảm bảo sự kiện được phát về frontend.  

#### Luồng Xử Lý  
1. **Người tham gia gia nhập phòng**:  
   - Smart contract phát "ParticipantJoined" (room ID, địa chỉ ví, sessionDescription).  
   - Phần 1 lắng nghe, gọi `getParticipantInfo` để lấy danh sách track, rồi gọi Phần 2 với dữ liệu.  

2. **Phần 2 xử lý với Cloudflare**:  
   - Tạo session, publish track, nhận phản hồi.  
   - Chuẩn bị dữ liệu sự kiện (session ID, phản hồi Cloudflare), yêu cầu Phần 1 gọi `forwardEventToFrontEnd`.  

3. **Phần 1 quản lý giao dịch**:  
   - Chọn ví trống, gửi giao dịch, nếu bận thì xếp hàng đợi, đảm bảo mỗi ví chỉ xử lý một giao dịch tại một thời điểm.  

#### Ví Dụ Triển Khai  
- Người tham gia gọi `joinRoom`, smart contract phát "ParticipantJoined" với room ID, địa chỉ ví, và sessionDescription.  
- Phần 1 lấy track từ `getParticipantInfo`, gọi Phần 2 với dữ liệu.  
- Phần 2 tạo session ID "session123", publish track, nhận phản hồi, rồi yêu cầu Phần 1 gọi `forwardEventToFrontEnd` với dữ liệu:  
  - `{"sessionID": "session123", "type": "join-room","cloudflareResponse": {"answer": "someAnswer"}}`.  
- Phần 1 chọn ví trống, gửi giao dịch, frontend nhận sự kiện và cập nhật với session ID và phản hồi.  

#### Bảng Tổng Hợp  
| Thành phần          | Chức năng                                      | Công việc chính                                      | Giao tiếp với khác       |
|---------------------|-----------------------------------------------|-----------------------------------------------------|--------------------------|
| Phần 1 (EventListener) | Lắng nghe sự kiện, quản lý ví, gửi giao dịch | Kết nối blockchain, lấy track, chọn ví, xếp hàng đợi | Gọi Phần 2, xử lý Phần 2 |
| Phần 2 (CloudflareHandler) | Tương tác Cloudflare, gửi sự kiện về frontend | Tạo session, publish track, chuẩn bị dữ liệu sự kiện | Yêu cầu Phần 1 gửi giao dịch |

---
