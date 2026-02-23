# Tài Liệu Yêu Cầu Sản Phẩm (PRD)
**Tên dự án:** [Tên Dự Án Của Bạn] - Nền tảng Light Novel Tương Tác
**Phiên bản:** 1.0 (MVP)
**Người tạo:** [Tên Của Bạn]
**Ngày tạo:** 23/02/2026

---

## 1. Tổng Quan Sản Phẩm (Product Overview)
**Mục tiêu:** Xây dựng một nền tảng web application cho phép người dùng sáng tạo, xuất bản và trải nghiệm các bộ light novel có cốt truyện rẽ nhánh (interactive fiction / visual novel dạng text). 
**Vấn đề giải quyết:** Các công cụ viết truyện truyền thống không hỗ trợ tốt việc quản lý cốt truyện phân nhánh phức tạp. Tác giả dễ bị rối logic, trong khi người đọc thiếu đi sự tương tác làm thay đổi kết cục câu chuyện.
**Giải pháp:** Cung cấp giao diện trực quan (Node-based Map) để tác giả thiết kế mạch truyện và một giao diện đọc mượt mà cho độc giả, được vận hành bởi hệ cơ sở dữ liệu đồ thị (Graph Database) tối ưu cho việc truy vấn rẽ nhánh.

---

## 2. Đối Tượng Người Dùng (Target Audience)
Hệ thống phân chia làm 2 nhóm người dùng chính:
1. **Creator (Tác giả):** Những người muốn viết và thiết kế game text-based, light novel tương tác. Cần công cụ trực quan để quản lý các nhánh truyện phức tạp mà không cần biết code.
2. **Reader (Người đọc):** Độc giả yêu thích thể loại truyện tương tác, nhập vai. Thích tự mình đưa ra quyết định để thay đổi diễn biến câu chuyện.

---

## 3. Kiến Trúc Kỹ Thuật (Tech Stack)
* **Frontend:** Vue.js 3 (Composition API)
  * Quản lý trạng thái: Pinia.
  * Giao diện sơ đồ truyện: Vue Flow.
* **Backend:** Golang (Gin hoặc Fiber framework)
  * Xử lý logic API với hiệu năng cao, concurrency tốt.
* **Database:** Neo4j (Graph Database)
  * Ngôn ngữ truy vấn: Cypher.
  * Tối ưu hóa việc lưu trữ Cảnh (Nodes) và Lựa chọn/Nhánh (Relationships).

---

## 4. Tính Năng Cốt Lõi (MVP Scope)

### 4.1. Creator Mode (Dành cho Tác giả)
| Tính năng | Mô tả chi tiết |
| :--- | :--- |
| **Story Management** | Tạo, sửa, xóa thông tin bộ truyện (Tên, Ảnh bìa, Tóm tắt). |
| **Visual Story Map** | Giao diện kéo thả sử dụng Vue Flow. Hiển thị toàn bộ cấu trúc rẽ nhánh của câu chuyện dưới dạng sơ đồ Node. |
| **Scene Editor** | Trình soạn thảo văn bản cho từng Node (Cảnh). Cho phép nhập text, upload ảnh minh họa (background/character). |
| **Choice Creator** | Tạo các liên kết (Edges) giữa các Node. Đặt tên cho lựa chọn (ví dụ: "Mở cửa"). |
| **Publishing** | Chuyển trạng thái truyện từ Draft (Bản nháp) sang Published (Công khai). |

### 4.2. Reader Mode (Dành cho Người đọc)
| Tính năng | Mô tả chi tiết |
| :--- | :--- |
| **Story Discovery** | Trang chủ hiển thị danh sách các bộ truyện đã được xuất bản để người dùng lựa chọn. |
| **Interactive Reading** | Giao diện đọc hiển thị nội dung của từng Scene. Cung cấp các nút bấm (Choices) ở cuối để người dùng quyết định hành động tiếp theo. |
| **Progress Tracking** | Lưu lại ID của Scene hiện tại (thông qua Pinia và Local Storage/Database) để người đọc có thể tiếp tục dang dở. |
| **History / Undo** | Cho phép xem lại lịch sử các lựa chọn đã đưa ra hoặc quay lại Node trước đó (Save/Load cơ bản). |

---

## 5. Yêu Cầu Dữ Liệu Cơ Bản (Database Schema - Neo4j)

**Nodes:**
* `(:User {id, username, email, password_hash, role})`
* `(:Story {id, title, summary, cover_image, status, created_at})`
* `(:Scene {id, content, image_url, is_start: boolean, is_end: boolean})`

**Relationships:**
* `(User)-[:CREATED]->(Story)`
* `(Story)-[:HAS_SCENE]->(Scene)`
* `(Story)-[:STARTS_AT]->(Scene)`
* `(Scene)-[:LEADS_TO {choice_text: "Text hiển thị"}]->(Scene)`

---

## 6. Yêu Cầu Phi Chức Năng (Non-Functional Requirements)
* **Hiệu năng:** Backend Golang cần đảm bảo thời gian phản hồi API dưới 200ms cho các truy vấn chuyển cảnh tiếp theo.
* **Trải nghiệm UX:** Chuyển cảnh ở chế độ Reader phải mượt mà, không giật lag, ưu tiên preload (tải trước) nội dung và hình ảnh của các nhánh tiếp theo liền kề.
* **Responsive:** Chế độ Reader phải hoạt động hoàn hảo trên thiết bị di động (Mobile-first approach). Chế độ Creator có thể ưu tiên Desktop/Tablet để có không gian thao tác với Story Map.

---

## 7. Các Tính Năng Cho Giai Đoạn 2 (Out of Scope for MVP)
* Hệ thống biến số (Variables) và Điều kiện (Conditions) cho các nhánh (ví dụ: cần có chìa khóa mới hiển thị lựa chọn mở cửa).
* Hệ thống Inventory (Túi đồ) và Stats (Chỉ số nhân vật).
* Tích hợp AI hỗ trợ tác giả viết tiếp cốt truyện.
* Đa ngôn ngữ.
* Nền tảng kiếm tiền (Monetization) cho tác giả.