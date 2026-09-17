# 📦 ATK Inventory Management System

### Enterprise Office Stationery & Inventory Management Platform

**A production-grade, full-stack web application for managing the complete office stationery (ATK) lifecycle — from procurement requests and multi-level approval to QR Code-based verification and real-time analytics.**

[![Go](https://img.shields.io/badge/Go-1.26.0-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Gin](https://img.shields.io/badge/Gin-1.11-00ADD8?style=for-the-badge&logo=gin&logoColor=white)](https://gin-gonic.com)
[![GORM](https://img.shields.io/badge/GORM-1.30-00ADD8?style=for-the-badge)](https://gorm.io)
[![Angular](https://img.shields.io/badge/Angular-22.1.6-DD0031?style=for-the-badge&logo=angular&logoColor=white)](https://angular.dev)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.9-3178C6?style=for-the-badge&logo=typescript&logoColor=white)](https://www.typescriptlang.org)
[![Tailwind](https://img.shields.io/badge/Tailwind-4.3.3-06B6D4?style=for-the-badge&logo=tailwindcss&logoColor=white)](https://tailwindcss.com)
[![MySQL](https://img.shields.io/badge/MySQL-8.4-4479A1?style=for-the-badge&logo=mysql&logoColor=white)](https://www.mysql.com)
[![JWT](https://img.shields.io/badge/JWT-Secure-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white)](https://jwt.io)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](LICENSE)
![Status](https://img.shields.io/badge/status-active--development-brightgreen?style=flat-square)

[Overview](#-overview) •
[Features](#-key-features) •
[Tech Stack](#-technology-stack) •
[Architecture](#-project-architecture) •
[Getting Started](#-getting-started) •
[Screenshots](#-screenshots)

---

## 📖 Overview

**ATK Inventory Management System** is an enterprise-grade platform designed to replace manual, spreadsheet-driven office stationery tracking with a secure, auditable, and fully digital workflow.

Built on a scalable RESTful architecture with a **Go (Gin) backend** and a modern **Angular 22 + Signals** frontend, the system enables organizations to manage inventory in real time — from request submission, multi-level approval, stock-in/stock-out transactions, to QR Code-based item verification and analytical reporting.

The project was built to demonstrate a real-world, production-style implementation of:

- **Scalable RESTful API design** with role-based authorization
- **Multi-level approval workflows** modeled after real corporate processes
- **QR Code generation and scanning** for tamper-resistant item tracking
- **Real-time stock analytics** with dynamic charts and trend monitoring
- **Bulk Excel import/export** for efficient data migration
- Clean separation of concerns across **frontend, backend, and data layers**

> 💡 This system reflects how modern enterprises manage shared office resources — combining governance, traceability, and operational efficiency in a single platform.

---

## ✨ Key Features

### 📊 Dashboard & Analytics
- Real-time monthly transaction counter with quick-glance KPIs
- 7-day transaction trend (Stock IN vs OUT) with interactive bar charts
- Category-based stock distribution with visual progress bars
- Role-aware dashboard view with responsive layout
- Recent activity feed with audit trail

### 📦 Item & Stock Management
- Full CRUD for inventory items with category classification
- Real-time stock tracking (available, minimum threshold, unit of measure)
- Stock movement history (IN / OUT) with audit trail
- Category management (Alat Tulis, Elektronik, Kertas, Perlengkapan)
- Low-stock alerts and automated threshold monitoring
- Export stock data to Excel

### 🔄 Transaction Management
- Request submission for stationery procurement
- Stock-IN & Stock-OUT recording with mandatory documentation
- Automatic stock recalculation on every transaction
- Full transaction history with pagination and advanced filtering
- Export transaction records to Excel

### ✅ Multi-Level Approval Workflow
- Structured review process — **Manager → Admin → Super Admin**
- Approve / Reject with mandatory comments for accountability
- Real-time status tracking with approval timeline
- Centralized pending-approvals dashboard
- Notification on status changes

### 📱 QR Code Integration
- Unique QR Code generated automatically for every item
- Live camera-based QR scanning for instant verification
- Downloadable & printable QR Code labels
- Batch QR generation for multiple items
- Instant item detail retrieval from scan

### 📈 Reports & Analytics
- Transaction reports by date range, category, and status
- Stock valuation and usage statistics
- Export-ready reports (Excel) with styled headers
- Historical trend analysis with chart visualizations
- Custom date-range filtering

### 👥 User & Role Management
- Secure authentication via **JWT**
- Fine-grained **Role-Based Access Control (RBAC)**
- Self-service profile management & password changes
- Full activity/audit logging with timestamp
- Dynamic menu rendering based on user role

**Supported roles:**

| Role | Level | Responsibilities |
|---|---|---|
| **Super Admin** | 5 | Full system administration, user & role management, system configuration, audit oversight |
| **Admin** | 4 | User management, master data, stock & transaction management |
| **Manager** | 3 | Approve transactions, monitor stock, generate reports |
| **Staff** | 2 | Create transactions, manage inventory, submit approval requests |
| **Auditor** | 1 | Read-only access for internal audit, compliance & verification |

### 🎨 Modern UI/UX
- Responsive layout (Desktop / Tablet / Mobile)
- Angular **Signals** for reactive state management
- **Tailwind CSS v4** utility-first styling with CSS-first config
- Toast notifications for real-time user feedback
- Collapsible sidebar with role-filtered navigation
- Elegant modal dialogs with keyboard interactions

---

## 🛠 Technology Stack

| Layer | Technology |
|---|---|
| **Backend Language** | Go 1.26.0 |
| **Backend Framework** | Gin (HTTP router) |
| **ORM** | GORM |
| **Database** | MySQL 8.4 |
| **Authentication** | JWT + bcrypt |
| **Frontend Framework** | Angular 22.1.6 (Standalone Components) |
| **Reactive State** | Angular Signals |
| **Language** | TypeScript 5.9 |
| **Styling** | Tailwind CSS v4.3.3 |
| **HTTP Client** | Angular HttpClient + RxJS |
| **Charts** | Custom SVG Charts |
| **QR Code** | qrcode (generation) + jsQR (scanning) |
| **Excel Handling** | excelize (Go) |
| **API Style** | RESTful (JSON over HTTP) |
| **Architecture** | Clean Layered (Handler → Service → Repository) |

Responsive across **Desktop**, **Tablet**, and **Mobile**.

---

## 🏗 Project Architecture
