import { Component, OnInit, signal, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../../core/services/api.service';
import { ToastService } from '../../core/services/toast.service';
import { PageResult, Transaction } from '../../core/models';
import {
  DataTableComponent,
  TableColumn,
} from '../../shared/components/data-table/data-table.component';
import { ExcelToolbarComponent } from '../../shared/components/excel-toolbar/excel-toolbar.component';
import { FilterBarComponent } from '../../shared/components/filter-bar/filter-bar.component';

@Component({
  selector: 'app-transactions',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    DataTableComponent,
    ExcelToolbarComponent,
    FilterBarComponent,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="space-y-5 fade-in">
      <div class="flex items-center justify-between flex-wrap gap-3">
        <div>
          <h2 class="text-2xl font-bold text-primary-900">Transaksi</h2>
          <p class="text-sm text-slate-500">Transaksi masuk dan keluar ATK</p>
        </div>
        <div class="flex items-center gap-2">
          <app-excel-toolbar (import)="onImport($event)" (exportClick)="onExport()" />
          <button
            (click)="openForm()"
            class="px-4 py-2 text-sm rounded-lg bg-primary-600 text-white hover:bg-primary-700 transition"
          >
            + Transaksi Baru
          </button>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-primary-100 shadow-sm flex gap-3 flex-wrap">
        <app-filter-bar placeholder="Cari transaksi..." (search)="onSearch($event)">
          <select [(ngModel)]="filterType" (change)="loadData(1)"
            class="px-3 py-2 text-sm border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30">
            <option value="">Semua Tipe</option>
            <option value="IN">Masuk</option>
            <option value="OUT">Keluar</option>
          </select>
          <select [(ngModel)]="filterStatus" (change)="loadData(1)"
            class="px-3 py-2 text-sm border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30">
            <option value="">Semua Status</option>
            <option value="draft">Draft</option>
            <option value="pending">Pending</option>
            <option value="approved">Approved</option>
            <option value="rejected">Rejected</option>
        
    
          </select>
        </app-filter-bar>
      </div>

      <app-data-table
        [columns]="columns"
        [page]="page()"
        [loading]="loading()"
        [trackBy]="trackById"
        [showDetail]="true"
        (pageChange)="loadData($event)"
        (edit)="openForm($event)"
        (delete)="onDelete($event)"
        (detail)="openDetail($event)"
      />

      <!-- ==================== MODAL FORM ==================== -->
      @if (showForm()) {
        <div class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4"
             (click)="closeForm()">
          <div class="bg-white rounded-2xl shadow-2xl w-full max-w-lg p-6 fade-in max-h-[90vh] overflow-y-auto"
               (click)="$event.stopPropagation()">
            <h3 class="text-lg font-bold text-primary-900 mb-4">
              {{ form.id ? 'Edit Transaksi' : 'Transaksi Baru' }}
            </h3>
            <div class="grid grid-cols-2 gap-3">
              <!-- Tipe -->
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">
                  Tipe <span class="text-red-500">*</span>
                </label>
                <select [(ngModel)]="form.type"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30">
                  <option value="IN">Masuk</option>
                  <option value="OUT">Keluar</option>
                  
            
                </select>
              </div>

              <!-- Tanggal -->
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Tanggal</label>
                <input type="date" [(ngModel)]="form.transactionDate"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30" />
              </div>

              <!-- ✅ Master Item — pakai ngModel + compareWith, TANPA index -->
              <div class="col-span-2">
                <label class="block text-sm font-medium text-slate-700 mb-1">
                  Master Item <span class="text-red-500">*</span>
                </label>
                <select
                  [(ngModel)]="form.masterId"
                  [compareWith]="compareById"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30">
                  <option [ngValue]="null">-- Pilih Master Item --</option>
                  @for (m of masterList(); track m.id) {
                    <option [ngValue]="m.id">{{ m.code }} - {{ m.name }}</option>
                  }
                </select>
                @if (masterList().length === 0) {
                  <p class="text-xs text-amber-600 mt-1">Belum ada master item.</p>
                }
                @if (form.masterId) {
                  <p class="text-xs text-slate-400 mt-1">
                    ID Master: {{ form.masterId }}
                  </p>
                }
              </div>

              <!-- Quantity -->
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">
                  Quantity <span class="text-red-500">*</span>
                </label>
                <input type="number" min="1" [(ngModel)]="form.quantity"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30" />
              </div>

              <!-- Status — DISABLED -->
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Status</label>
                <input [value]="form.status || 'draft'" disabled
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg bg-slate-100 outline-none text-slate-600 capitalize cursor-not-allowed" />
                <p class="text-xs text-slate-400 mt-1">Berubah via approval</p>
              </div>

              <!-- Catatan -->
              <div class="col-span-2">
                <label class="block text-sm font-medium text-slate-700 mb-1">Catatan</label>
                <textarea [(ngModel)]="form.notes" rows="3" placeholder="Catatan transaksi..."
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30"></textarea>
              </div>
            </div>

            <div class="flex justify-end gap-2 mt-6">
              <button (click)="closeForm()"
                class="px-4 py-2 text-sm rounded-lg bg-slate-100 hover:bg-slate-200 transition">
                Batal
              </button>
              <button (click)="save()" [disabled]="saving()"
                class="px-4 py-2 text-sm rounded-lg bg-primary-600 text-white hover:bg-primary-700 disabled:opacity-60 transition">
                {{ saving() ? 'Menyimpan...' : 'Simpan' }}
              </button>
            </div>
          </div>
        </div>
      }

      <!-- ==================== MODAL DETAIL ==================== -->
      @if (showDetailModal() && detailData()) {
        <div class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4 overflow-y-auto"
             (click)="closeDetail()">
          <div class="bg-white rounded-2xl shadow-2xl w-full max-w-2xl p-6 my-8 fade-in"
               (click)="$event.stopPropagation()">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-lg font-bold text-primary-900">Detail Transaksi</h3>
              <span class="px-3 py-1 text-xs rounded-full font-medium capitalize"
                    [ngClass]="statusClass(detailData()?.status)">
                {{ detailData()?.status }}
              </span>
            </div>

            <div class="grid grid-cols-2 gap-3 mb-4 p-4 bg-primary-50 rounded-xl">
              <div>
                <div class="text-xs text-slate-500">No Transaksi</div>
                <div class="font-semibold text-primary-900">{{ detailData()?.transactionCode }}</div>
              </div>
              <div>
                <div class="text-xs text-slate-500">Tipe</div>
                <div class="font-semibold text-primary-900">{{ detailData()?.type }}</div>
              </div>
              <div>
                <div class="text-xs text-slate-500">Tanggal</div>
                <div class="font-semibold text-primary-900">{{ detailData()?.transactionDate }}</div>
              </div>
              <div>
                <div class="text-xs text-slate-500">User</div>
                <div class="font-semibold text-primary-900">
                  {{ detailData()?.user?.full_name ?? detailData()?.user?.username ?? '-' }}
                </div>
              </div>
            </div>

            @if ((detailData()?.items ?? []).length > 0) {
              <div class="mb-4">
                <div class="text-sm font-semibold text-slate-700 mb-2">Item Transaksi</div>
                <div class="border border-primary-100 rounded-xl overflow-hidden">
                  <table class="w-full text-sm">
                    <thead class="bg-primary-50">
                      <tr>
                        <th class="px-3 py-2 text-left text-xs font-semibold text-primary-800">#</th>
                        <th class="px-3 py-2 text-left text-xs font-semibold text-primary-800">Kode</th>
                        <th class="px-3 py-2 text-left text-xs font-semibold text-primary-800">Nama Item</th>
                        <th class="px-3 py-2 text-right text-xs font-semibold text-primary-800">Qty</th>
                        <th class="px-3 py-2 text-right text-xs font-semibold text-primary-800">Harga</th>
                        <th class="px-3 py-2 text-right text-xs font-semibold text-primary-800">Subtotal</th>
                      </tr>
                    </thead>
                    <tbody>
                      @for (item of detailData()?.items; track item.id; let i = $index) {
                        <tr class="border-t border-primary-50">
                          <td class="px-3 py-2 text-slate-500">{{ i + 1 }}</td>
                          <td class="px-3 py-2 font-medium">{{ item.master?.code ?? '-' }}</td>
                          <td class="px-3 py-2">{{ item.master?.name ?? '-' }}</td>
                          <td class="px-3 py-2 text-right">{{ item.quantity }}</td>
                          <td class="px-3 py-2 text-right">{{ item.price }}</td>
                          <td class="px-3 py-2 text-right font-medium">{{ item.subTotal }}</td>
                        </tr>
                      }
                    </tbody>
                    <tfoot class="bg-primary-50">
                      <tr>
                        <td colspan="3" class="px-3 py-2 text-right font-semibold text-primary-800">Total</td>
                        <td class="px-3 py-2 text-right font-semibold text-primary-800">{{ detailData()?.totalItems }}</td>
                        <td colspan="2" class="px-3 py-2 text-right font-bold text-primary-900">{{ detailData()?.totalValue }}</td>
                      </tr>
                    </tfoot>
                  </table>
                </div>
              </div>
            }

            <div class="mb-4">
              <div class="text-sm font-semibold text-slate-700 mb-1">Catatan</div>
              <div class="p-3 bg-slate-50 rounded-lg text-sm text-slate-700">
                {{ detailData()?.notes || '(tidak ada catatan)' }}
              </div>
            </div>

            @if (detailData()?.approvedBy || detailData()?.approvedAt) {
              <div class="mb-4 p-3 bg-green-50 border border-green-200 rounded-lg">
                <div class="text-xs text-green-700 font-semibold mb-1">APPROVAL</div>
                <div class="text-sm text-slate-700">
                  Disetujui pada: {{ detailData()?.approvedAt }}
                </div>
              </div>
            }

            <div class="flex justify-end gap-2">
              <button (click)="closeDetail()"
                class="px-4 py-2 text-sm rounded-lg bg-slate-100 hover:bg-slate-200 transition">
                Tutup
              </button>
            </div>
          </div>
        </div>
      }
    </div>
  `,
})
export class TransactionsComponent implements OnInit {
  columns: TableColumn[] = [
    { key: 'transactionCode', label: 'No Transaksi', sortable: true },
    { key: 'type', label: 'Tipe' },
    { key: 'totalItems', label: 'Total Item' },
    { key: 'totalValue', label: 'Total Value' },
    { key: 'status', label: 'Status' },
    { key: 'user.username', label: 'User' },
    { key: 'transactionDate', label: 'Tanggal' },
  ];

  page = signal<PageResult<Transaction> | null>(null);
  loading = signal(false);
  saving = signal(false);
  showForm = signal(false);
  showDetailModal = signal(false);
  detailData = signal<any>(null);
  masterList = signal<any[]>([]);
  form: any = {};
  search = '';
  filterType = '';
  filterStatus = '';
  currentPage = 1;
  trackById = (r: any) => r.id;

  constructor(
    private api: ApiService,
    private toast: ToastService,
  ) {}

  ngOnInit() {
    this.loadMasterList();
    this.loadData(1);
  }

  // ============================================================
  // ✅ compareWith — agar ngModel cocokkan ID angka
  // ============================================================
  compareById(a: any, b: any): boolean {
    if (a == null || b == null) return a === b;
    return Number(a) === Number(b);
  }

  // ============================================================
  // LOAD MASTER LIST
  // ============================================================
  loadMasterList() {
    this.api.get<any>('master', { page: 1, size: 1000 }).subscribe({
      next: (res: any) => {
        const root = res?.data ?? res;
        const inner = root?.data ?? root;
        const rawItems: any[] = Array.isArray(inner)
          ? inner
          : Array.isArray(root)
            ? root
            : [];

        // ✅ Normalisasi: paksa id jadi Number
        const items = rawItems
          .map((m: any) => ({
            ...m,
            id: Number(m?.id ?? m?.ID ?? m?.Id ?? 0),
          }))
          .filter((m: any) => m.id > 0);

        console.log('[Transactions] Master loaded:', items.length, 'items');
        if (items.length > 0) {
          console.log('[Transactions] Sample master[0]:', items[0]);
        }

        this.masterList.set(items);
      },
      error: (err) => {
        console.error('[Transactions] Master load error:', err);
        this.masterList.set([]);
      },
    });
  }

  // ============================================================
  // LOAD DATA transaksi
  // ============================================================
  loadData(page: number) {
    this.currentPage = page;
    this.loading.set(true);

    this.api
      .get<any>('transactions', {
        page,
        size: 10,
        limit: 10,
        search: this.search,
        type: this.filterType,
        status: this.filterStatus,
      })
      .subscribe({
        next: (res: any) => {
          this.loading.set(false);
          console.log('[Transactions] Raw response:', res);

          const root = res?.data ?? res;
          const inner = root?.data ?? root;
          const items: any[] = Array.isArray(inner)
            ? inner
            : Array.isArray(root)
              ? root
              : Array.isArray(res)
                ? res
                : [];

          const meta = res?.meta ?? root?.meta ?? root ?? {};

          this.page.set({
            data: items,
            page: Number(meta?.page ?? page),
            size: Number(meta?.size ?? meta?.limit ?? 10),
            total: Number(meta?.total ?? items.length),
            totalPages: Number(meta?.totalPages ?? meta?.totalPage ?? 1),
          });
        },
        error: (err) => {
          this.loading.set(false);
          console.error('[Transactions] Error:', err);
          let msg = 'Gagal memuat data transaksi';
          if (err?.status === 0) msg = 'Tidak dapat terhubung ke server';
          else if (err?.status === 401) msg = 'Sesi habis, silakan login ulang';
          else if (err?.error?.message) msg = err.error.message;
          this.toast.error(msg);
        },
      });
  }

  onSearch(q: string) {
    this.search = q;
    this.loadData(1);
  }

  // ============================================================
  // ✅ OPEN FORM — prefill dari items[0] (support berbagai field name)
  // ============================================================
  openForm(row?: Transaction) {
    if (row) {
      const items = (row as any).items ?? [];
      const firstItem = items[0] ?? {};

      // ✅ Coba berbagai kemungkinan field name untuk masterId
      const masterIdRaw =
        (firstItem as any).masterId ??
        (firstItem as any).master_id ??
        (firstItem as any).MasterID ??
        (firstItem as any).MasterId ??
        (firstItem as any).master?.id ??
        (firstItem as any).master?.ID ??
        0;

      const masterId = Number(masterIdRaw);
      const quantity = Number((firstItem as any).quantity ?? 1);

      console.log('[Transactions] openForm - items:', items);
      console.log('[Transactions] openForm - firstItem:', firstItem);
      console.log('[Transactions] openForm - extracted masterId:', masterId);
      console.log('[Transactions] openForm - masterList size:', this.masterList().length);

      this.form = {
        id: row.id,
        type: (row as any).type ?? 'IN',
        status: (row as any).status ?? 'draft',
        transactionDate: (row as any).transactionDate
          ? String((row as any).transactionDate).slice(0, 10)
          : new Date().toISOString().slice(0, 10),
        masterId: masterId > 0 ? masterId : null,
        quantity: quantity > 0 ? quantity : 1,
        notes: (row as any).notes ?? '',
      };

      console.log('[Transactions] Edit form prefilled:', this.form);

      // ✅ Kalau masterList belum ada, load dulu
      if (this.masterList().length === 0) {
        console.warn('[Transactions] masterList kosong, reload...');
        this.loadMasterList();
      }
    } else {
      this.form = {
        type: 'IN',
        status: 'draft',
        transactionDate: new Date().toISOString().slice(0, 10),
        quantity: 1,
        masterId: null,
        notes: '',
      };
    }

    this.showForm.set(true);
  }

  closeForm() {
    this.showForm.set(false);
    this.form = {};
  }

  // ============================================================
  // OPEN DETAIL
  // ============================================================
  openDetail(row: Transaction) {
    console.log('[Transactions] Detail row:', row);
    this.detailData.set(row);
    this.showDetailModal.set(true);
  }

  closeDetail() {
    this.showDetailModal.set(false);
    this.detailData.set(null);
  }

  statusClass(status: string): string {
    const map: Record<string, string> = {
      draft: 'bg-slate-100 text-slate-700',
      pending: 'bg-amber-100 text-amber-700',
      approved: 'bg-green-100 text-green-700',
      rejected: 'bg-red-100 text-red-700',
      completed: 'bg-blue-100 text-blue-700',
      cancelled: 'bg-gray-200 text-gray-600',
    };
    return map[(status ?? '').toLowerCase()] ?? 'bg-slate-100 text-slate-700';
  }

  // ============================================================
  // SAVE — create / update
  // ============================================================
  save() {
    const masterId = Number(this.form.masterId);
    const quantity = Number(this.form.quantity);

    if (!masterId || isNaN(masterId) || masterId <= 0) {
      this.toast.warning('Pilih Master Item terlebih dahulu');
      return;
    }
    if (!quantity || isNaN(quantity) || quantity <= 0) {
      this.toast.warning('Quantity harus lebih dari 0');
      return;
    }

    this.saving.set(true);

    const tpe = String(this.form.type || 'IN').trim();
    const nts = String(this.form.notes || '').trim();

    // ===== CREATE =====
    if (!this.form.id) {
      const payload = {
        type: tpe,
        notes: nts,
        status: 'draft',
        items: [
          {
            masterId: masterId,
            quantity: quantity,
            notes: nts,
          },
        ],
      };

      console.log('[Transactions] Create payload:', JSON.stringify(payload));

      this.api.post<any>('transactions', payload).subscribe({
        next: (res: any) => {
          this.saving.set(false);
          if (res?.success !== false) {
            this.toast.success('Transaksi berhasil dibuat');
            this.closeForm();
            this.loadData(this.currentPage);
          } else {
            this.toast.error(res?.message || 'Gagal membuat transaksi');
          }
        },
        error: (err) => {
          this.saving.set(false);
          console.error('[Transactions] Create error:', err);
          this.toast.error(
            err?.error?.message || err?.error?.error || 'Gagal membuat transaksi',
          );
        },
      });
      return;
    }

    // ===== UPDATE =====
    const payload = {
      type: tpe,
      notes: nts,
      items: [
        {
          masterId: masterId,
          quantity: quantity,
          notes: nts,
        },
      ],
    };

    console.log('[Transactions] Update payload:', JSON.stringify(payload));

    this.api.put<any>(`transactions/${this.form.id}`, payload).subscribe({
      next: (res: any) => {
        this.saving.set(false);
        if (res?.success !== false) {
          this.toast.success('Transaksi berhasil diupdate');
          this.closeForm();
          this.loadData(this.currentPage);
        } else {
          this.toast.error(res?.message || 'Gagal update transaksi');
        }
      },
      error: (err) => {
        this.saving.set(false);
        console.error('[Transactions] Update error:', err);
        this.toast.error(err?.error?.message || 'Gagal update transaksi');
      },
    });
  }

  onDelete(row: Transaction) {
    if (!confirm(`Hapus transaksi "${row.transactionCode}"?`)) return;

    this.api.delete<any>(`transactions/${row.id}`).subscribe({
      next: (res: any) => {
        if (res?.success !== false) {
          this.toast.success('Transaksi dihapus');
          this.loadData(this.currentPage);
        } else {
          this.toast.error(res?.message || 'Gagal menghapus');
        }
      },
      error: (err) => {
        this.toast.error(err?.error?.message || 'Gagal menghapus');
      },
    });
  }

  onImport(file: File) {
    this.api.upload<any>('transactions/import', file).subscribe({
      next: (res: any) => {
        const data = res?.data ?? {};
        if (res?.success !== false) {
          this.toast.success(
            `Import: ${data?.imported ?? 0} baru, ${data?.failed ?? 0} gagal`,
          );
          this.loadData(1);
        } else {
          this.toast.error(res?.message || 'Import gagal');
        }
      },
      error: (err) => {
        this.toast.error(err?.error?.message || 'Import gagal');
      },
    });
  }

  onExport() {
    this.api
      .exportFile('transactions/export', {
        search: this.search,
        type: this.filterType,
        status: this.filterStatus,
      })
      .subscribe({
        next: (blob: Blob) => {
          if (!blob || blob.size === 0) {
            this.toast.error('File kosong dari server');
            return;
          }

          if (blob.type.includes('json') || blob.type.includes('html')) {
            blob.text().then((text) => {
              console.error('[Transactions Export] Bukan Excel:', text.slice(0, 300));
              this.toast.error('Server tidak mengembalikan file Excel');
            });
            return;
          }

          const url = window.URL.createObjectURL(blob);
          const a = document.createElement('a');
          a.href = url;
          a.download = `transactions_${Date.now()}.xlsx`;
          document.body.appendChild(a);
          a.click();
          document.body.removeChild(a);
          window.URL.revokeObjectURL(url);
          this.toast.success('File berhasil diunduh');
        },
        error: (err) => {
          console.error('[Transactions Export] Error:', err);
          this.toast.error('Gagal export data');
        },
      });
  }
}