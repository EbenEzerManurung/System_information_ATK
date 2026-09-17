import { Component, OnInit, signal, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../../core/services/api.service';
import { ToastService } from '../../core/services/toast.service';
import { PageResult, Role } from '../../core/models';
import {
  DataTableComponent,
  TableColumn,
} from '../../shared/components/data-table/data-table.component';
import { ExcelToolbarComponent } from '../../shared/components/excel-toolbar/excel-toolbar.component';
import { FilterBarComponent } from '../../shared/components/filter-bar/filter-bar.component';

@Component({
  selector: 'app-roles',
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
          <h2 class="text-2xl font-bold text-primary-900">Manajemen Role</h2>
          <p class="text-sm text-slate-500">Kelola role dan hak akses</p>
        </div>
        <div class="flex items-center gap-2">
          <app-excel-toolbar (import)="onImport($event)" (exportClick)="onExport()" />
          <button
            (click)="openForm()"
            class="px-4 py-2 text-sm rounded-lg bg-primary-600 text-white hover:bg-primary-700 transition"
          >
            + Tambah Role
          </button>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-primary-100 shadow-sm">
        <app-filter-bar placeholder="Cari role..." (search)="onSearch($event)" />
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
              {{ form.id ? 'Edit Role' : 'Tambah Role' }}
            </h3>
            <div class="space-y-3">
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">
                  Nama Role <span class="text-red-500">*</span>
                </label>
                <input
                  [(ngModel)]="form.name"
                  placeholder="Contoh: admin"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg focus:ring-2 focus:ring-primary-500/30 outline-none"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Deskripsi</label>
                <textarea
                  [(ngModel)]="form.description"
                  rows="3"
                  placeholder="Deskripsi role..."
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
export class RolesComponent implements OnInit {
  columns: TableColumn[] = [
    { key: 'name', label: 'Nama Role', sortable: true },
    { key: 'description', label: 'Deskripsi' },
  ];

  page = signal<PageResult<Role> | null>(null);
  loading = signal(false);
  saving = signal(false);
  showForm = signal(false);
  form: any = {};
  search = '';
  currentPage = 1;
  trackById = (r: any) => r.id;

  constructor(
    private api: ApiService,
    private toast: ToastService,
  ) {}

  ngOnInit() {
    this.loadData(1);
  }

  // ============================================================
  // LOAD DATA — normalisasi berbagai format response
  // ============================================================
  loadData(page: number) {
    this.currentPage = page;
    this.loading.set(true);

    this.api.get<any>('roles', { page, size: 10, search: this.search }).subscribe({
      next: (res: any) => {
        this.loading.set(false);
        console.log('[Roles] Raw response:', res);

        // ==== NORMALISASI ====
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
        console.error('[Roles] Error:', err);
        let msg = 'Gagal memuat data role';
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

  openForm(row?: Role) {
    this.form = row ? { ...row } : {};
    this.showForm.set(true);
  }

  closeForm() {
    this.showForm.set(false);
    this.form = {};
  }

  save() {
    if (!this.form.name || !this.form.name.trim()) {
      this.toast.warning('Nama role wajib diisi');
      return;
    }

    this.saving.set(true);

    const req = this.form.id
      ? this.api.put<any>(`roles/${this.form.id}`, this.form)
      : this.api.post<any>('roles', this.form);

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
        console.error('[Roles] Save error:', err);
        this.toast.error(err?.error?.message || 'Gagal menyimpan');
      },
    });
  }

  onDelete(row: Role) {
    if (!confirm(`Hapus role "${row.name}"?`)) return;

    this.api.delete<any>(`roles/${row.id}`).subscribe({
      next: (res: any) => {
        if (res?.success !== false) {
          this.toast.success('Role dihapus');
          this.loadData(this.currentPage);
        } else {
          this.toast.error(res?.message || 'Gagal menghapus');
        }
      },
      error: (err) => {
        console.error('[Roles] Delete error:', err);
        this.toast.error(err?.error?.message || 'Gagal menghapus');
      },
    });
  }

  onImport(file: File) {
    this.api.upload<any>('roles/import', file).subscribe({
      next: (res: any) => {
        if (res?.success !== false) {
          this.toast.success('Import berhasil');
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
  // EXPORT — pola sama dengan Users & Stock
  // ============================================================
  onExport() {
    this.api.exportFile('roles/export', { search: this.search }).subscribe({
      next: (blob: Blob) => {
        // Validasi blob
        if (!blob || blob.size === 0) {
          this.toast.error('File kosong dari server');
          return;
        }

        // Cek content-type
        if (blob.type.includes('json') || blob.type.includes('html')) {
          blob.text().then(text => {
            console.error('[Roles Export] Server return bukan Excel:', text.slice(0, 300));
            this.toast.error('Server tidak mengembalikan file Excel');
          });
          return;
        }

        // Simpan file
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `roles_${Date.now()}.xlsx`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);

        this.toast.success('File berhasil diunduh');
      },
      error: (err) => {
        console.error('[Roles Export] Error:', err);
        this.toast.error('Gagal export data');
      },
    });
  }
}