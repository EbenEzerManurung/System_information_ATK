import { Component, OnInit, signal, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../../core/services/api.service';
import { ToastService } from '../../core/services/toast.service';
import { MasterItem, PageResult } from '../../core/models';
import {
  DataTableComponent,
  TableColumn,
} from '../../shared/components/data-table/data-table.component';
import { ExcelToolbarComponent } from '../../shared/components/excel-toolbar/excel-toolbar.component';
import { FilterBarComponent } from '../../shared/components/filter-bar/filter-bar.component';

@Component({
  selector: 'app-master',
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
          <h2 class="text-2xl font-bold text-primary-900">Master Data</h2>
          <p class="text-sm text-slate-500">Data referensi kategori, satuan, dan lokasi</p>
        </div>
        <div class="flex items-center gap-2">
          <app-excel-toolbar (import)="onImport($event)" (exportClick)="onExport()" />
          <button
            (click)="openForm()"
            class="px-4 py-2 text-sm rounded-lg bg-primary-600 text-white hover:bg-primary-700 transition"
          >
            + Tambah Master
          </button>
        </div>
      </div>

      <div
        class="bg-white rounded-2xl p-4 border border-primary-100 shadow-sm flex gap-3 flex-wrap"
      >
        <app-filter-bar placeholder="Cari master data..." (search)="onSearch($event)">
          <select
            [(ngModel)]="filterCategory"
            (change)="loadData(1)"
            class="px-3 py-2 text-sm border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
          >
            <option value="">Semua Kategori</option>
            @for (cat of categoryOptions; track cat) {
              <option [value]="cat">{{ cat }}</option>
            }
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
              {{ form.id ? 'Edit Master' : 'Tambah Master' }}
            </h3>
            <div class="space-y-3">
              <!-- KATEGORI — pakai input + datalist, bukan select -->
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">
                  Kategori <span class="text-red-500">*</span>
                </label>
                <input
                  [(ngModel)]="form.category"
                  list="category-options"
                  placeholder="Contoh: kertas, elektronik, perlengkapan, alat tulis"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                />
                <datalist id="category-options">
                  <option value="kertas"></option>
                  <option value="elektronik"></option>
                  <option value="perlengkapan"></option>
                  <option value="alat tulis"></option>
                </datalist>
                <p class="text-xs text-slate-400 mt-1">
                  Ketik atau pilih dari daftar
                </p>
              </div>

              <!-- KODE -->
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">
                  Kode <span class="text-red-500">*</span>
                </label>
                <input
                  [(ngModel)]="form.code"
                  placeholder="Contoh: AT001"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                />
              </div>

              <!-- NAMA -->
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">
                  Nama <span class="text-red-500">*</span>
                </label>
                <input
                  [(ngModel)]="form.name"
                  placeholder="Contoh: Pulpen Standard"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                />
              </div>

              <!-- UNIT (opsional) -->
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Satuan</label>
                <input
                  [(ngModel)]="form.unit"
                  list="unit-options"
                  placeholder="Contoh: pcs, rim, pack, box"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                />
                <datalist id="unit-options">
                  <option value="pcs"></option>
                  <option value="rim"></option>
                  <option value="pack"></option>
                  <option value="box"></option>
                  <option value="lusin"></option>
                </datalist>
              </div>

              <!-- DESKRIPSI -->
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Deskripsi</label>
                <textarea
                  [(ngModel)]="form.description"
                  rows="2"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                ></textarea>
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
export class MasterComponent implements OnInit {
  columns: TableColumn[] = [
    { key: 'category', label: 'Kategori', sortable: true },
    { key: 'code', label: 'Kode', sortable: true },
    { key: 'name', label: 'Nama' },
    { key: 'unit', label: 'Satuan' },
    { key: 'description', label: 'Deskripsi' },
  ];

  page = signal<PageResult<MasterItem> | null>(null);
  loading = signal(false);
  saving = signal(false);
  showForm = signal(false);
  form: any = {};
  search = '';
  filterCategory = '';
  currentPage = 1;
  trackById = (r: any) => r.id;

  // Daftar opsi kategori untuk dropdown filter
  categoryOptions = ['kertas', 'elektronik', 'perlengkapan', 'alat tulis'];

  constructor(
    private api: ApiService,
    private toast: ToastService,
  ) {}

  ngOnInit() {
    this.loadData(1);
  }

  // ============================================================
  // LOAD DATA
  // ============================================================
  loadData(page: number) {
    this.currentPage = page;
    this.loading.set(true);

    this.api
      .get<any>('master', {
        page,
        size: 10,
        limit: 10,
        search: this.search,
        category: this.filterCategory,
      })
      .subscribe({
        next: (res: any) => {
          this.loading.set(false);
          console.log('[Master] Raw response:', res);

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
          console.error('[Master] Error:', err);
          let msg = 'Gagal memuat data master';
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

  openForm(row?: MasterItem) {
    if (row) {
      this.form = { ...row };
    } else {
      this.form = {
        category: '',
        code: '',
        name: '',
        unit: '',
        description: '',
      };
    }
    this.showForm.set(true);
  }

  closeForm() {
    this.showForm.set(false);
    this.form = {};
  }

  save() {
    if (!this.form.category || !this.form.category.trim()) {
      this.toast.warning('Kategori wajib diisi');
      return;
    }
    if (!this.form.code || !this.form.code.trim()) {
      this.toast.warning('Kode wajib diisi');
      return;
    }
    if (!this.form.name || !this.form.name.trim()) {
      this.toast.warning('Nama wajib diisi');
      return;
    }

    this.saving.set(true);

    const payload = {
      code: this.form.code,
      name: this.form.name,
      category: this.form.category,
      unit: this.form.unit || '',
      description: this.form.description || '',
      minStock: this.form.minStock ?? 0,
      maxStock: this.form.maxStock ?? 0,
      price: this.form.price ?? 0,
    };

    console.log('[Master] Save payload:', payload);

    const req = this.form.id
      ? this.api.put<any>(`master/${this.form.id}`, payload)
      : this.api.post<any>('master', payload);

    req.subscribe({
      next: (res: any) => {
        this.saving.set(false);
        if (res?.success !== false) {
          this.toast.success('Data disimpan');
          this.closeForm();
          this.loadData(this.currentPage);
        } else {
          this.toast.error(res?.message || 'Gagal menyimpan');
        }
      },
      error: (err) => {
        this.saving.set(false);
        console.error('[Master] Save error:', err);
        this.toast.error(err?.error?.message || 'Gagal menyimpan');
      },
    });
  }

  onDelete(row: MasterItem) {
    if (!confirm(`Hapus "${row.name}"?`)) return;

    this.api.delete<any>(`master/${row.id}`).subscribe({
      next: (res: any) => {
        if (res?.success !== false) {
          this.toast.success('Data dihapus');
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
    this.api.upload<any>('master/import', file).subscribe({
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

  onExport() {
    this.api
      .exportFile('master/export', { category: this.filterCategory, search: this.search })
      .subscribe({
        next: (blob: Blob) => {
          if (!blob || blob.size === 0) {
            this.toast.error('File kosong dari server');
            return;
          }

          if (blob.type.includes('json') || blob.type.includes('html')) {
            blob.text().then(text => {
              console.error('[Master Export] Bukan Excel:', text.slice(0, 300));
              this.toast.error('Server tidak mengembalikan file Excel');
            });
            return;
          }

          const url = window.URL.createObjectURL(blob);
          const a = document.createElement('a');
          a.href = url;
          a.download = `master_${Date.now()}.xlsx`;
          document.body.appendChild(a);
          a.click();
          document.body.removeChild(a);
          window.URL.revokeObjectURL(url);

          this.toast.success('File berhasil diunduh');
        },
        error: (err) => {
          console.error('[Master Export] Error:', err);
          this.toast.error('Gagal export data');
        },
      });
  }
}