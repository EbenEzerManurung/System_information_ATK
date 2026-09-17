import { Component, OnInit, signal, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../../core/services/api.service';
import { ToastService } from '../../core/services/toast.service';
import {
  DataTableComponent,
  TableColumn,
} from '../../shared/components/data-table/data-table.component';
import { PageResult } from '../../core/models';

type ReportType = 'stock' | 'transaction';

@Component({
  selector: 'app-reports',
  standalone: true,
  imports: [CommonModule, FormsModule, DataTableComponent],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="space-y-5 fade-in">
      <!-- Header -->
      <div class="flex items-center justify-between flex-wrap gap-3">
        <div>
          <h2 class="text-2xl font-bold text-primary-900">Laporan</h2>
          <p class="text-sm text-slate-500">Laporan stock dan transaksi</p>
        </div>
        <!-- ✅ Tombol Import Excel dihilangkan, hanya Export Excel -->
        <button
          (click)="onExportExcel()"
          class="px-4 py-2 text-sm rounded-lg bg-white border border-primary-200 text-primary-700 hover:bg-primary-50 transition flex items-center gap-2 shadow-sm">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          Export Excel
        </button>
      </div>

      <!-- Tabs -->
      <div class="bg-white rounded-2xl p-1 border border-primary-100 shadow-sm inline-flex">
        @for (t of tabs; track t.value) {
          <button
            (click)="switchTab(t.value)"
            class="px-4 py-2 text-sm rounded-xl transition font-medium"
            [class.bg-primary-600]="activeTab() === t.value"
            [class.text-white]="activeTab() === t.value"
            [class.text-slate-600]="activeTab() !== t.value"
            [class.hover:bg-primary-50]="activeTab() !== t.value">
            {{ t.label }}
          </button>
        }
      </div>

      <!-- Filter Panel -->
      <div class="bg-white rounded-2xl p-4 border border-primary-100 shadow-sm">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-3">
          <div>
            <label class="block text-xs font-medium text-slate-600 mb-1">Dari Tanggal</label>
            <input type="date" [(ngModel)]="dateFrom"
                   class="w-full px-3 py-2 text-sm border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none" />
          </div>
          <div>
            <label class="block text-xs font-medium text-slate-600 mb-1">Sampai Tanggal</label>
            <input type="date" [(ngModel)]="dateTo"
                   class="w-full px-3 py-2 text-sm border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none" />
          </div>
          <div>
            <label class="block text-xs font-medium text-slate-600 mb-1">
              {{ activeTab() === 'stock' ? 'Kategori' : 'Tipe Transaksi' }}
            </label>
            <select [(ngModel)]="filterType"
                    class="w-full px-3 py-2 text-sm border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30">
              @if (activeTab() === 'stock') {
                <option value="">Semua Kategori</option>
                <option value="kertas">Kertas</option>
                <option value="elektronik">Elektronik</option>
                <option value="perlengkapan">Perlengkapan</option>
                <option value="alat tulis">Alat Tulis</option>
              } @else {
                <option value="">Semua Tipe</option>
                <option value="IN">Masuk</option>
                <option value="OUT">Keluar</option>
              }
            </select>
          </div>
          <div class="flex items-end gap-2">
            <button (click)="loadData(1)"
                    class="flex-1 px-4 py-2 text-sm rounded-lg bg-primary-600 text-white hover:bg-primary-700 transition">
              Terapkan Filter
            </button>
            <button (click)="resetFilter()"
                    class="px-4 py-2 text-sm rounded-lg bg-slate-100 hover:bg-slate-200 transition">
              Reset
            </button>
          </div>
        </div>
      </div>

      <!-- ========================================== -->
      <!-- SUMMARY CARDS: LAPORAN STOCK               -->
      <!-- ========================================== -->
      @if (activeTab() === 'stock') {
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 items-start">
          <!-- TOTAL ITEM -->
          <div class="bg-white rounded-2xl p-5 border border-primary-100 shadow-sm">
            <div class="text-xs tracking-wide text-slate-500 uppercase">Total Item</div>
            <div class="text-3xl font-bold text-primary-900 mt-2">{{ totalItem() }}</div>
          </div>

          <!-- TOTAL QTY -->
          <div class="bg-white rounded-2xl p-5 border border-primary-100 shadow-sm">
            <div class="text-xs tracking-wide text-slate-500 uppercase">Total Qty</div>
            <div class="text-3xl font-bold text-primary-900 mt-2">{{ totalQty() | number }}</div>
          </div>

          <!-- STOCK RENDAH (Ascending paling rendah) -->
          <div class="bg-white rounded-2xl p-5 border border-amber-200 shadow-sm">
            <div class="flex items-center justify-between">
              <div class="text-xs tracking-wide text-slate-500 uppercase">Stock Lowest</div>
              <span class="px-2 py-0.5 text-xs rounded-full bg-amber-100 text-amber-700 font-medium">
                {{ lowestStockItems().length }} item
              </span>
            </div>
            @if (lowestStockItems().length === 0) {
              <div class="text-sm text-slate-400 mt-2">Tidak ada data stock</div>
            } @else {
              <div class="mt-3 space-y-2 max-h-44 overflow-y-auto pr-1">
                @for (item of lowestStockItems(); track item.id) {
                  <div class="p-2.5 rounded-lg bg-amber-50 border border-amber-100">
                    <div class="flex items-center justify-between gap-2">
                      <span class="text-xs font-semibold text-amber-800">{{ item.master?.code }}</span>
                      <span class="text-xs font-bold text-amber-700">Qty {{ item.quantity | number }}</span>
                    </div>
                    <div class="text-sm font-medium text-slate-800 truncate" [title]="item.master?.name">
                      {{ item.master?.name }}
                    </div>
                    <div class="text-xs text-slate-500">{{ item.master?.category || '-' }}</div>
                  </div>
                }
              </div>
            }
          </div>

          <!-- STOCK HABIS (Qty = 0 atau < 0) -->
          <div class="bg-white rounded-2xl p-5 border border-red-200 shadow-sm">
            <div class="flex items-center justify-between">
              <div class="text-xs tracking-wide text-slate-500 uppercase">Stock Habis</div>
              <span class="px-2 py-0.5 text-xs rounded-full bg-red-100 text-red-700 font-medium">
                {{ emptyStockItems().length }} item
              </span>
            </div>
            @if (emptyStockItems().length === 0) {
              <div class="text-sm text-slate-400 mt-2">Tidak ada item yang habis</div>
            } @else {
              <div class="mt-3 space-y-2 max-h-44 overflow-y-auto pr-1">
                @for (item of emptyStockItems(); track item.id) {
                  <div class="p-2.5 rounded-lg bg-red-50 border border-red-100">
                    <div class="flex items-center justify-between gap-2">
                      <span class="text-xs font-semibold text-red-800">{{ item.master?.code }}</span>
                      <span class="text-xs font-bold text-red-700">Qty {{ item.quantity | number }}</span>
                    </div>
                    <div class="text-sm font-medium text-slate-800 truncate" [title]="item.master?.name">
                      {{ item.master?.name }}
                    </div>
                    <div class="text-xs text-slate-500">{{ item.master?.category || '-' }}</div>
                  </div>
                }
              </div>
            }
          </div>
        </div>
      }

      <!-- ========================================== -->
      <!-- SUMMARY CARDS: LAPORAN TRANSAKSI            -->
      <!-- ========================================== -->
      @if (activeTab() === 'transaction') {
        @if (summary().length > 0) {
          <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
            @for (s of summary(); track s.label) {
              <div class="bg-white rounded-xl p-4 border border-primary-100 shadow-sm">
                <div class="text-xs text-slate-500 uppercase">{{ s.label }}</div>
                <div class="text-xl font-bold text-primary-900 mt-1">{{ s.value | number }}</div>
              </div>
            }
          </div>
        }
      }

      <!-- Table dengan PDF per row -->
      <app-data-table
        [columns]="columns()"
        [page]="page()"
        [loading]="loading()"
        [trackBy]="trackById"
        [showPdf]="true"
        [showEdit]="false"
        [showDelete]="false"
        [showDetail]="false"
        (pageChange)="loadData($event)"
        (pdf)="onExportPdfRow($event)"
      />
    </div>
  `,
})
export class ReportsComponent implements OnInit {
  tabs = [
    { label: 'Laporan Stock', value: 'stock' as ReportType },
    { label: 'Laporan Transaksi', value: 'transaction' as ReportType },
  ];

  activeTab = signal<ReportType>('stock');
  page = signal<PageResult<any> | null>(null);
  loading = signal(false);
  summary = signal<{ label: string; value: number }[]>([]);
  columns = signal<TableColumn[]>([]);

  // ✅ Signals untuk Laporan Stock
  totalItem = signal(0);
  totalQty = signal(0);
  lowestStockItems = signal<any[]>([]);
  emptyStockItems = signal<any[]>([]);

  dateFrom = '';
  dateTo = '';
  filterType = '';
  currentPage = 1;

  trackById = (r: any) => r.id;

  constructor(
    private api: ApiService,
    private toast: ToastService,
  ) {}

  ngOnInit() {
    const today = new Date();
    const firstDay = new Date(today.getFullYear(), today.getMonth(), 1);
    this.dateFrom = firstDay.toISOString().slice(0, 10);
    this.dateTo = today.toISOString().slice(0, 10);
    this.updateColumns();
    this.loadData(1);
  }

  // ============================================================
  // SWITCH TAB
  // ============================================================
  switchTab(t: ReportType) {
    this.activeTab.set(t);
    this.filterType = '';
    this.updateColumns();
    this.loadData(1);
  }

  // ============================================================
  // UPDATE COLUMNS — berdasarkan tab aktif
  // ============================================================
  updateColumns() {
    if (this.activeTab() === 'stock') {
      this.columns.set([
        { key: 'master.code', label: 'Kode', sortable: true },
        { key: 'master.name', label: 'Nama Item', sortable: true },
        { key: 'master.category', label: 'Kategori' },
        { key: 'master.unit', label: 'Satuan' },
        { key: 'quantity', label: 'Qty', sortable: true },
        { key: 'location', label: 'Lokasi' },
        { key: 'status', label: 'Status' },
      ]);
    } else {
      this.columns.set([
        { key: 'transactionCode', label: 'No Transaksi', sortable: true },
        { key: 'type', label: 'Tipe' },
        { key: 'totalItems', label: 'Total Item' },
        { key: 'totalValue', label: 'Total Value' },
        { key: 'status', label: 'Status' },
        { key: 'user.username', label: 'User' },
        { key: 'transactionDate', label: 'Tanggal', sortable: true },
      ]);
    }
  }

  // ============================================================
  // LOAD DATA — langsung ke endpoint stock atau transactions
  // ============================================================
  loadData(page: number) {
    this.currentPage = page;
    this.loading.set(true);

    const endpoint = this.activeTab() === 'stock' ? 'stock' : 'transactions';

    const params: any = {
      page,
      size: 10,
    };

    if (this.activeTab() === 'stock') {
      if (this.filterType) params.category = this.filterType;
      if (this.dateFrom) params.dateFrom = this.dateFrom;
      if (this.dateTo) params.dateTo = this.dateTo;
    } else {
      if (this.filterType) params.type = this.filterType;
      if (this.dateFrom) params.dateFrom = this.dateFrom;
      if (this.dateTo) params.dateTo = this.dateTo;
    }

    console.log(`[Reports] Loading ${endpoint} with params:`, params);

    this.api.get<any>(endpoint, params).subscribe({
      next: (res: any) => {
        this.loading.set(false);
        console.log(`[Reports] Response:`, res);

        const root = res?.data ?? res;
        const inner = root?.data ?? root;
        const items: any[] = Array.isArray(inner)
          ? inner
          : Array.isArray(root)
            ? root
            : [];

        const meta = res?.meta ?? root?.meta ?? root ?? {};

        this.page.set({
          data: items,
          page: Number(meta?.page ?? page),
          size: Number(meta?.size ?? meta?.limit ?? 10),
          total: Number(meta?.total ?? items.length),
          totalPages: Number(meta?.totalPages ?? meta?.totalPage ?? 1),
        });

        this.buildSummary(items);
      },
      error: (err) => {
        this.loading.set(false);
        console.error('[Reports] Error:', err);
        this.toast.error('Gagal memuat laporan');
        this.page.set({ data: [], total: 0, page: 1, size: 10, totalPages: 0 });
        this.summary.set([]);
        this.totalItem.set(0);
        this.totalQty.set(0);
        this.lowestStockItems.set([]);
        this.emptyStockItems.set([]);
      },
    });
  }

  // ============================================================
  // BUILD SUMMARY
  // ============================================================
  private buildSummary(items: any[]) {
    if (this.activeTab() === 'stock') {
      const totalItems = items.length;
      const totalQty = items.reduce((s, x) => s + Number(x.quantity ?? 0), 0);
      
      // Mencari Qty terendah (Ascending) yang masih > 0
      const availableItems = items.filter(x => Number(x.quantity) > 0);
      let lowestQty = 0;
      let lowestItems: any[] = [];
      
      if (availableItems.length > 0) {
        lowestQty = Math.min(...availableItems.map(x => Number(x.quantity)));
        lowestItems = availableItems.filter(x => Number(x.quantity) === lowestQty);
      }

      // Mencari Stock Habis (Qty <= 0)
      const outOfStockItems = items.filter(x => Number(x.quantity) <= 0);

      this.totalItem.set(totalItems);
      this.totalQty.set(totalQty);
      this.lowestStockItems.set(lowestItems);
      this.emptyStockItems.set(outOfStockItems);
      
      // Kosongkan summary transaksi
      this.summary.set([]);
    } else {
      // Logika untuk Tab Transaksi
      const totalTx = items.length;
      const totalIn = items.filter((x) => x.type === 'IN').length;
      const totalOut = items.filter((x) => x.type === 'OUT').length;
      const totalValue = items.reduce((s, x) => s + Number(x.totalValue ?? 0), 0);

      this.summary.set([
        { label: 'Total Transaksi', value: totalTx },
        { label: 'Masuk', value: totalIn },
        { label: 'Keluar', value: totalOut },
        { label: 'Total Value', value: totalValue },
      ]);

      // Kosongkan summary stock
      this.totalItem.set(0);
      this.totalQty.set(0);
      this.lowestStockItems.set([]);
      this.emptyStockItems.set([]);
    }
  }

  // ============================================================
  // RESET FILTER
  // ============================================================
  resetFilter() {
    const today = new Date();
    const firstDay = new Date(today.getFullYear(), today.getMonth(), 1);
    this.dateFrom = firstDay.toISOString().slice(0, 10);
    this.dateTo = today.toISOString().slice(0, 10);
    this.filterType = '';
    this.loadData(1);
  }

  // ============================================================
  // EXPORT EXCEL — semua data
  // ============================================================
  onExportExcel() {
    const endpoint = this.activeTab() === 'stock'
      ? 'stock/export'
      : 'transactions/export';

    const params: any = {};
    if (this.filterType) {
      params[this.activeTab() === 'stock' ? 'category' : 'type'] = this.filterType;
    }
    if (this.dateFrom) params.dateFrom = this.dateFrom;
    if (this.dateTo) params.dateTo = this.dateTo;

    this.api.exportFile(endpoint, params).subscribe({
      next: (blob: Blob) => {
        if (!blob || blob.size === 0) {
          this.toast.error('File kosong dari server');
          return;
        }
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `report_${this.activeTab()}_${Date.now()}.xlsx`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);
        this.toast.success('File Excel berhasil diunduh');
      },
      error: (err) => {
        console.error('[Reports] Export error:', err);
        this.toast.error('Gagal export Excel');
      },
    });
  }

  // ============================================================
  // EXPORT PDF PER ROW — pakai window.print
  // ============================================================
  onExportPdfRow(row: any) {
    console.log('[Reports] Export PDF for row:', row);

    const title = this.activeTab() === 'stock'
      ? `Laporan Stock - ${row.master?.code ?? row.id}`
      : `Laporan Transaksi - ${row.transactionCode ?? row.id}`;

    const rows: { label: string; value: string }[] = [];

    if (this.activeTab() === 'stock') {
      rows.push(
        { label: 'Kode Item', value: String(row.master?.code ?? '-') },
        { label: 'Nama Item', value: String(row.master?.name ?? '-') },
        { label: 'Kategori', value: String(row.master?.category ?? '-') },
        { label: 'Satuan', value: String(row.master?.unit ?? '-') },
        { label: 'Quantity', value: String(row.quantity ?? 0) },
        { label: 'Min Stock', value: String(row.master?.minStock ?? 0) },
        { label: 'Lokasi', value: String(row.location ?? '-') },
        { label: 'Status', value: String(row.status ?? '-') },
      );
    } else {
      rows.push(
        { label: 'No Transaksi', value: String(row.transactionCode ?? '-') },
        { label: 'Tipe', value: String(row.type ?? '-') },
        { label: 'Total Item', value: String(row.totalItems ?? 0) },
        { label: 'Total Value', value: String(row.totalValue ?? 0) },
        { label: 'Status', value: String(row.status ?? '-') },
        { label: 'User', value: String(row.user?.username ?? row.user?.full_name ?? '-') },
        { label: 'Tanggal', value: String(row.transactionDate ?? '-') },
        { label: 'Catatan', value: String(row.notes ?? '-') },
      );
    }

    const html = `
      <!DOCTYPE html>
      <html>
      <head>
        <meta charset="utf-8">
        <title>${title}</title>
        <style>
          @page { size: A4; margin: 15mm; }
          * { box-sizing: border-box; }
          body {
            font-family: 'Segoe UI', Arial, sans-serif;
            padding: 20px;
            color: #0f172a;
            margin: 0;
          }
          .header {
            text-align: center;
            border-bottom: 3px solid #1e3a8a;
            padding-bottom: 15px;
            margin-bottom: 20px;
          }
          .header h1 {
            margin: 0;
            color: #1e3a8a;
            font-size: 22px;
            font-weight: 700;
          }
          .header .subtitle {
            color: #64748b;
            font-size: 12px;
            margin-top: 4px;
          }
          .meta {
            display: flex;
            justify-content: space-between;
            font-size: 11px;
            color: #64748b;
            margin-bottom: 15px;
          }
          table {
            width: 100%;
            border-collapse: collapse;
            font-size: 12px;
          }
          table th {
            background: #1e3a8a;
            color: #ffffff;
            text-align: left;
            padding: 10px 12px;
            border: 1px solid #1e40af;
            font-weight: 600;
          }
          table td {
            padding: 10px 12px;
            border: 1px solid #cbd5e1;
          }
          table tr:nth-child(even) td {
            background: #f8fafc;
          }
          .label-col { width: 35%; font-weight: 600; background: #eff6ff !important; color: #1e3a8a; }
          .footer {
            margin-top: 30px;
            text-align: center;
            font-size: 10px;
            color: #94a3b8;
            border-top: 1px solid #e2e8f0;
            padding-top: 10px;
          }
        </style>
      </head>
      <body>
        <div class="header">
          <h1>System Information ATK</h1>
          <div class="subtitle">${title}</div>
        </div>

        <div class="meta">
          <div>Dicetak: ${new Date().toLocaleString('id-ID')}</div>
          <div>Oleh: ${this.getCurrentUser()}</div>
        </div>

        <table>
          <tbody>
            ${rows.map((r) => `
              <tr>
                <td class="label-col">${this.escapeHtml(r.label)}</td>
                <td>${this.escapeHtml(r.value)}</td>
              </tr>
            `).join('')}
          </tbody>
        </table>

        <div class="footer">
          Dokumen ini digenerate otomatis oleh System Information ATK
        </div>

        <script>
          window.onload = function() {
            window.print();
            setTimeout(() => window.close(), 500);
          };
        <\/script>
      </body>
      </html>
    `;

    const w = window.open('', '_blank', 'width=900,height=700');
    if (!w) {
      this.toast.error('Popup diblokir. Izinkan popup untuk export PDF');
      return;
    }
    w.document.write(html);
    w.document.close();
  }

  private escapeHtml(str: string): string {
    return String(str ?? '')
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#39;');
  }

  private getCurrentUser(): string {
    try {
      const raw = localStorage.getItem('atk_user');
      if (raw) {
        const u = JSON.parse(raw);
        return u?.full_name ?? u?.username ?? 'System';
      }
    } catch {}
    return 'System';
  }
}