# 📌 Giới thiệu sản phẩm

## 1️⃣ Thông tin chung
- **Tên sản phẩm:** DABE Backend System
- **Chức năng chính:** Hệ thống backend hỗ trợ xử lý dữ liệu, giao tiếp với frontend và cung cấp API.
- **Product Owner (PO):** Lý Vĩnh Thái

## 2️⃣ Công nghệ sử dụng 🚀

### 🌐 Backend
- ![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white) - Ngôn ngữ lập trình chính.
- ![Gin](https://img.shields.io/badge/Gin-00ADD8?style=for-the-badge&logo=go&logoColor=white) - Xây dựng API mạnh mẽ.
- ![gRPC](https://img.shields.io/badge/gRPC-4285F4?style=for-the-badge&logo=google-cloud&logoColor=white) - Giao tiếp hiệu quả giữa các microservices.
- ![GraphQL](https://img.shields.io/badge/GraphQL-E10098?style=for-the-badge&logo=graphql&logoColor=white) - Truy vấn dữ liệu linh hoạt.
- ![GORM](https://img.shields.io/badge/GORM-512BD4?style=for-the-badge&logo=go&logoColor=white) - ORM cho Golang.

### 🗄️ Cơ sở dữ liệu
- ![MySQL](https://img.shields.io/badge/MySQL-4479A1?style=for-the-badge&logo=mysql&logoColor=white) - Hệ quản trị cơ sở dữ liệu quan hệ.
- ![Firestore](https://img.shields.io/badge/Firestore-FFCA28?style=for-the-badge&logo=firebase&logoColor=black) - Lưu trữ dữ liệu phi quan hệ.

### ☁️ Cloud & DevOps
- ![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white) - Đóng gói ứng dụng.
- ![Google Cloud](https://img.shields.io/badge/GCP-4285F4?style=for-the-badge&logo=google-cloud&logoColor=white) - Deploy backend trên cloud.
- ![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=github-actions&logoColor=white) - CI/CD tự động.

### 🔒 Bảo mật & Xác thực
- ![JWT](https://img.shields.io/badge/JWT-000000?style=for-the-badge&logo=json-web-tokens&logoColor=white) - Xác thực và bảo mật API.
- ![OAuth 2.0](https://img.shields.io/badge/OAuth_2.0-EC2025?style=for-the-badge&logo=oauth&logoColor=white) - Xác thực với các dịch vụ bên thứ ba.

## 3️⃣ Kiến trúc hệ thống 🏗️
Hệ thống được thiết kế theo mô hình **microservices**, với các dịch vụ riêng biệt giao tiếp qua gRPC.

```
┌───────────────────┐       ┌───────────────────┐
│  Client (FE)      │ <-->  │  API Gateway      │
└───────────────────┘       └───────────────────┘
           |                        |       
           |                        |       
 ┌────────────────┐       ┌────────────────┐
 │ Service Auth   │  <--> │ Service Data   │
 └────────────────┘       └────────────────┘
```

## 4️⃣ Chức năng chính 🔥
- **Quản lý người dùng** (Đăng nhập, đăng ký, phân quyền,...)
- **Xử lý dữ liệu GraphQL**
- **Tương tác với cơ sở dữ liệu MySQL**
- **Triển khai hệ thống với Docker & GCP**
- **Tích hợp gRPC giữa các microservices**
- **Xác thực và bảo mật với JWT/OAuth 2.0**

## 5️⃣ Kết luận ✨
Dự án **DABE Backend System** là một hệ thống backend mạnh mẽ, áp dụng các công nghệ hiện đại để đảm bảo hiệu suất và bảo mật cao. Sản phẩm này được phát triển nhằm hỗ trợ các ứng dụng frontend hoạt động mượt mà và hiệu quả.

📌 **Liên hệ PO:** *Lý Vĩnh Thái*