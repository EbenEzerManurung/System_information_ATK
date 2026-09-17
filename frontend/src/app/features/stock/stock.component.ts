import { Component, OnInit, signal, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../../core/services/api.service';
import { ToastService } from '../../core/services/toast.service';
import { PageResult, Stock } from '../../core/models';
import {
  DataTableComponent,
  TableColumn,
} from '../../shared/components/data-table/data-table.component';
import { ExcelToolbarComponent } from '../../shared/components/excel-toolbar/excel-toolbar.component';
import { FilterBarComponent } from '../../shared/components/filter-bar/filter-bar.component';

@Component({
  selector: 'app-stock',
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
          <h2 class="text-2xl font-bold text-primary-900">Manajemen Stock</h2>
          <p class="text-sm text-slate-500">Kelola persediaan ATK</p>
        </div>
        <div class="flex items-center gap-2">
          <app-excel-toolbar (import)="onImport($event)" (exportClick)="onExport()" />
          <button
            (click)="openForm()"
            class="px-4 py-2 text-sm rounded-lg bg-primary-600 text-white hover:bg-primary-700 transition"
          >
            + Tambah Stock
          </button>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-primary-100 shadow-sm flex gap-3 flex-wrap">
        <app-filter-bar placeholder="Cari item..." (search)="onSearch($event)">
          <select
            [(ngModel)]="filterStatus"
            (change)="loadData(1)"
            class="px-3 py-2 text-sm border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
          >
            <option value="">Semua Status</option>
            <option value="available">Available</option>
            <option value="low">Low</option>
            <option value="out">Out of Stock</option>
          </select>
        </app-filter-bar>
      </div>

      <app-data-table
        [columns]="columns"
        [page]="page()"
        [loading]="loading()"
        [trackBy]="trackById"
        (pageChange)="loadData($event)"
        (edit)="openForm($event)"
        (delete)="onDelete($event)"
      />

      @if (showForm()) {
        <div
          class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4"
          (click)="closeForm()"
        >
          <div
            class="bg-white rounded-2xl shadow-2xl w-full max-w-lg p-6 fade-in"
            (click)="$event.stopPropagation()"
          >
            <h3 class="text-lg font-bold text-primary-900 mb-4">
              {{ form.id ? 'Edit Stock' : 'Tambah Stock' }}
            </h3>
            <div class="grid grid-cols-2 gap-3">
              <!-- Master Item -->
              <div class="col-span-2">
                <label class="block text-sm font-medium text-slate-700 mb-1">
                  Master Item <span class="text-red-500">*</span>
                </label>
                @if (form.id) {
                  <input
                    [value]="form.master?.name || form.master?.code || '-'"
                    disabled
                    class="w-full px-3 py-2 border border-primary-200 rounded-lg bg-slate-50 outline-none"
                  />
                } @else {
                  <select
                    [(ngModel)]="form.masterId"
                    class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                  >
                    <option [ngValue]="null">-- Pilih Master Item --</option>
                    @for (m of masterList(); track m.id) {
                      <option [ngValue]="m.id">{{ m.code }} - {{ m.name }}</option>
                    }
                  </select>
                  @if (masterList().length === 0) {
                    <p class="text-xs text-amber-600 mt-1">
                      Belum ada master item. Tambah di menu Master Data dulu.
                    </p>
                  }
                }
              </div>

              <!-- Lokasi -->
              <div class="col-span-2">
                <label class="block text-sm font-medium text-slate-700 mb-1">
                  Lokasi <span class="text-red-500">*</span>
                </label>
                <input
                  [(ngModel)]="form.location"
                  placeholder="Contoh: Gudang Utama"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                />
              </div>

              <!-- Status -->
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Status</label>
                <select
                  [(ngModel)]="form.status"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                >
                  <option value="available">Available</option>
                  <option value="low">Low</option>
                  <option value="out">Out of Stock</option>
                </select>
              </div>

              <!-- Quantity -->
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Quantity</label>
                <input
                  type="number"
                  min="0"
                  [(ngModel)]="form.quantity"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                />
              </div>
            </div>
            <div class="flex justify-end gap-2 mt-6">
              <button
                (click)="closeForm()"
                class="px-4 py-2 text-sm rounded-lg bg-slate-100 hover:bg-slate-200 transition"
              >
                Batal
              </button>
              <button
                (click)="save()"
                [disabled]="saving()"
                class="px-4 py-2 text-sm rounded-lg bg-primary-600 text-white hover:bg-primary-700 disabled:opacity-60 transition"
              >
                {{ saving() ? 'Menyimpan...' : 'Simpan' }}
              </button>
            </div>
          </div>
        </div>
      }
    </div>
  `,
})
export class StockComponent implements OnInit {
  columns: TableColumn[] = [
    { key: 'master.code', label: 'Kode', sortable: false },
    { key: 'master.name', label: 'Nama Item', sortable: false },
    { key: 'master.category', label: 'Kategori', sortable: false },
    { key: 'quantity', label: 'Qty', sortable: true },
    { key: 'master.unit', label: 'Satuan', sortable: false },
    { key: 'location', label: 'Lokasi' },
    { key: 'status', label: 'Status' },
  ];

  page = signal<PageResult<Stock> | null>(null);
  loading = signal(false);
  saving = signal(false);
  showForm = signal(false);
  masterList = signal<any[]>([]);
  form: any = {};
  search = '';
  filterStatus = '';
  currentPage = 1;
  trackById = (r: any) => r.id;

  constructor(
    private api: ApiService,
    private toast: ToastService,
  ) {}

  ngOnInit() {
    this.loadData(1);
    this.loadMasterList();
  }

  // ============================================================
  // LOAD MASTER LIST
  // ============================================================
  loadMasterList() {
    this.api.get<any>('master', { page: 1, size: 1000 }).subscribe({
      next: (res: any) => {
        const root = res?.data ?? res;
        const inner = root?.data ?? root;
        const items: any[] = Array.isArray(inner)
          ? inner
          : Array.isArray(root)
            ? root
            : Array.isArray(res)
              ? res
              : [];

        this.masterList.set(items);
      },
      error: () => this.masterList.set([]),
    });
  }

  // ============================================================
  // LOAD STOCK DATA
  // ============================================================
  loadData(page: number) {
    this.currentPage = page;
    this.loading.set(true);

    this.api
      .get<any>('stock', {
        page,
        size: 10,
        limit: 10,
        search: this.search,
        status: this.filterStatus,
      })
      .subscribe({
        next: (res: any) => {
          this.loading.set(false);

          const root = res?.data ?? res;
          const inner = root?.data ?? root;
          const items: any[] = Array.isArray(inner)
            ? inner
            : Array.isArray(root)
              ? root
              : Array.isArray(res)
                ? res
                : [];

          const meta =
            res?.meta ??
            root?.meta ??
            (typeof root === 'object' && !Array.isArray(root) ? root : {}) ??
            {};

          const total = Number(meta?.total ?? meta?.totalCount ?? items.length);
          const size = Number(meta?.size ?? meta?.limit ?? meta?.perPage ?? 10);
          const currentPage = Number(meta?.page ?? meta?.currentPage ?? page);
          const totalPages = Number(
            meta?.totalPages ??
              meta?.totalPage ??
              meta?.lastPage ??
              Math.max(1, Math.ceil(total / size)),
          );

          this.page.set({
            data: items,
            page: currentPage,
            size,
            total,
            totalPages,
          });
        },
        error: (err) => {
          this.loading.set(false);
          let msg = 'Gagal memuat data stock';
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

  openForm(row?: Stock) {
    if (row) {
      this.form = { ...row };
    } else {
      this.form = {
        masterId: null,
        quantity: 0,
        location: 'Gudang Utama',
        status: 'available',
      };
    }
    this.showForm.set(true);
  }

  closeForm() {
    this.showForm.set(false);
    this.form = {};
  }

  save() {
    // Validasi CREATE
    if (!this.form.id) {
      if (!this.form.masterId) {
        this.toast.warning('Pilih Master Item terlebih dahulu');
        return;
      }
      if (!this.form.location || !this.form.location.trim()) {
        this.toast.warning('Lokasi wajib diisi');
        return;
      }
    }

    this.saving.set(true);

    // CREATE
    if (!this.form.id) {
      const payload = {
        masterId: Number(this.form.masterId),
        quantity: Number(this.form.quantity) || 0,
        location: this.form.location,
        status: this.form.status || 'available',
      };

      this.api.post<any>('stock', payload).subscribe({
        next: (res: any) => {
          this.saving.set(false);
          if (res?.success !== false) {
            this.toast.success('Stock berhasil dibuat');
            this.closeForm();
            this.loadData(this.currentPage);
          } else {
            this.toast.error(res?.message || 'Gagal membuat stock');
          }
        },
        error: (err) => {
          this.saving.set(false);
          this.toast.error(err?.error?.message || 'Gagal membuat stock');
        },
      });
      return;
    }

    // UPDATE
    const payload = {
      location: this.form.location,
      status: this.form.status,
      quantity: Number(this.form.quantity) || 0,
    };

    this.api.put<any>(`stock/${this.form.id}`, payload).subscribe({
      next: (res: any) => {
        this.saving.set(false);
        if (res?.success !== false) {
          this.toast.success('Stock berhasil diupdate');
          this.closeForm();
          this.loadData(this.currentPage);
        } else {
          this.toast.error(res?.message || 'Gagal update stock');
        }
      },
      error: (err) => {
        this.saving.set(false);
        this.toast.error(err?.error?.message || 'Gagal update stock');
      },
    });
  }

  onDelete(row: Stock) {
    const label = (row as any).master?.name ?? row.location;
    if (!confirm(`Hapus stock "${label}"?`)) return;

    this.api.delete<any>(`stock/${row.id}`).subscribe({
      next: (res: any) => {
        if (res?.success !== false) {
          this.toast.success('Stock berhasil dihapus');
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
    this.api.upload<any>('stock/import', file).subscribe({
      next: (res: any) => {
        const data = res?.data ?? {};
        if (res?.success !== false) {
          this.toast.success(
            `Import: ${data?.imported ?? 0} baru, ${data?.updated ?? 0} update, ${data?.failed ?? 0} gagal`
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

  // ============================================================
  // EXPORT — pola SAMA dengan users.component.ts (terbukti berhasil)
  // ============================================================
  onExport() {
    this.api.exportFile('stock/export', { search: this.search }).subscribe({
      next: (blob: Blob) => {
        // Cek blob valid
        if (!blob || blob.size === 0) {
          this.toast.error('File kosong dari server');
          return;
        }

        // Cek content-type — kalau JSON/HTML, berarti server error
        if (blob.type.includes('json') || blob.type.includes('html')) {
          blob.text().then(text => {
            console.error('[Export] Server return bukan Excel:', text.slice(0, 300));
            this.toast.error('Server tidak mengembalikan file Excel');
          });
          return;
        }

        // Simpan file
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `stock_${Date.now()}.xlsx`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);

        this.toast.success('File berhasil diunduh');
      },
      error: (err) => {
        console.error('[Export] Error:', err);
        this.toast.error('Gagal export data');
      },
    });
  }
}