import {
  Component,
  Input,
  Output,
  EventEmitter,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { PageResult } from '../../../core/models';

export interface TableColumn {
  key: string;
  label: string;
  sortable?: boolean;
  width?: string;
}

@Component({
  selector: 'app-data-table',
  standalone: true,
  imports: [CommonModule, FormsModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="bg-white rounded-2xl shadow-sm border border-primary-100 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead class="bg-primary-50 border-b border-primary-100">
            <tr>
              <th class="px-4 py-3 text-left font-semibold text-primary-800 w-12">#</th>
              @for (col of columns; track col.key) {
                <th
                  class="px-4 py-3 text-left font-semibold text-primary-800"
                  [style.width]="col.width"
                  [class.cursor-pointer]="col.sortable"
                  (click)="col.sortable && toggleSort(col.key)"
                >
                  <div class="flex items-center gap-1">
                    {{ col.label }}
                    @if (col.sortable) {
                      <span class="text-xs opacity-60">
                        @if (sortBy === col.key) {
                          {{ sortDir === 'asc' ? '▲' : '▼' }}
                        } @else {
                          ⇅
                        }
                      </span>
                    }
                  </div>
                </th>
              }
              <th class="px-4 py-3 text-right font-semibold text-primary-800"
                  [class.w-32]="actionColumnWidth === 'sm'"
                  [class.w-40]="actionColumnWidth === 'md'"
                  [class.w-48]="actionColumnWidth === 'lg'"
                  [class.w-56]="actionColumnWidth === 'xl'">
                Aksi
              </th>
            </tr>
          </thead>
          <tbody>
            @if (loading) {
              @for (i of [1, 2, 3, 4, 5]; track i) {
                <tr class="border-b border-primary-50">
                  @for (c of skeletonCols; track $index) {
                    <td class="px-4 py-3">
                      <div class="h-4 bg-primary-100 rounded animate-pulse"></div>
                    </td>
                  }
                </tr>
              }
            } @else if (rows.length > 0) {
              @for (row of rows; track trackBy ? trackBy(row) : $index; let i = $index) {
                <tr
                  class="border-b border-primary-50 hover:bg-primary-50/50 transition-colors fade-in"
                >
                  <td class="px-4 py-3 text-slate-500">
                    {{ startIndex + i + 1 }}
                  </td>
                  @for (col of columns; track col.key) {
                    <td class="px-4 py-3 text-slate-700">{{ getValue(row, col.key) }}</td>
                  }
                  <td class="px-4 py-3 text-right">
                    <div class="flex justify-end gap-1 flex-wrap">
                      <!-- ✅ Tombol PDF -->
                      @if (showPdf) {
                        <button
                          (click)="pdf.emit(row)"
                          class="px-3 py-1 text-xs rounded-lg bg-purple-100 text-purple-700 hover:bg-purple-200 transition"
                          title="Export PDF"
                        >
                          PDF
                        </button>
                      }

                      <!-- ✅ Tombol Detail -->
                      @if (showDetail) {
                        <button
                          (click)="detail.emit(row)"
                          class="px-3 py-1 text-xs rounded-lg bg-blue-100 text-blue-700 hover:bg-blue-200 transition"
                        >
                          Detail
                        </button>
                      }

                      <!-- ✅ Tombol Edit -->
                      @if (showEdit) {
                        <button
                          (click)="edit.emit(row)"
                          class="px-3 py-1 text-xs rounded-lg bg-primary-100 text-primary-700 hover:bg-primary-200 transition"
                        >
                          Edit
                        </button>
                      }

                      <!-- ✅ Tombol Delete -->
                      @if (showDelete) {
                        <button
                          (click)="delete.emit(row)"
                          class="px-3 py-1 text-xs rounded-lg bg-red-100 text-red-700 hover:bg-red-200 transition"
                        >
                          Hapus
                        </button>
                      }
                    </div>
                  </td>
                </tr>
              }
            } @else {
              <tr>
                <td
                  [attr.colspan]="columns.length + 2"
                  class="px-4 py-12 text-center text-slate-400"
                >
                  Tidak ada data
                </td>
              </tr>
            }
          </tbody>
        </table>
      </div>

      @if (total > 0) {
        <div
          class="px-4 py-3 border-t border-primary-100 flex items-center justify-between flex-wrap gap-3 bg-primary-50/50"
        >
          <div class="text-sm text-slate-600">
            Menampilkan
            <span class="font-semibold text-primary-700">{{ startIndex + 1 }}</span>
            –
            <span class="font-semibold text-primary-700">{{ endIndex }}</span>
            dari <span class="font-semibold text-primary-700">{{ total }}</span> data
          </div>
          <div class="flex items-center gap-1">
            <button
              [disabled]="currentPage === 1"
              (click)="goTo(1)"
              class="px-3 py-1.5 text-sm rounded-lg bg-white border border-primary-200 text-primary-700 hover:bg-primary-50 disabled:opacity-40 disabled:cursor-not-allowed transition"
            >
              «
            </button>
            <button
              [disabled]="currentPage === 1"
              (click)="goTo(currentPage - 1)"
              class="px-3 py-1.5 text-sm rounded-lg bg-white border border-primary-200 text-primary-700 hover:bg-primary-50 disabled:opacity-40 disabled:cursor-not-allowed transition"
            >
              ‹
            </button>
            @for (p of pages; track p) {
              <button
                (click)="goTo(p)"
                class="px-3 py-1.5 text-sm rounded-lg border transition"
                [class.bg-primary-600]="p === currentPage"
                [class.text-white]="p === currentPage"
                [class.border-primary-600]="p === currentPage"
                [class.bg-white]="p !== currentPage"
                [class.border-primary-200]="p !== currentPage"
                [class.text-primary-700]="p !== currentPage"
              >
                {{ p }}
              </button>
            }
            <button
              [disabled]="currentPage === totalPages"
              (click)="goTo(currentPage + 1)"
              class="px-3 py-1.5 text-sm rounded-lg bg-white border border-primary-200 text-primary-700 hover:bg-primary-50 disabled:opacity-40 disabled:cursor-not-allowed transition"
            >
              ›
            </button>
            <button
              [disabled]="currentPage === totalPages"
              (click)="goTo(totalPages)"
              class="px-3 py-1.5 text-sm rounded-lg bg-white border border-primary-200 text-primary-700 hover:bg-primary-50 disabled:opacity-40 disabled:cursor-not-allowed transition"
            >
              »
            </button>
          </div>
        </div>
      }
    </div>
  `,
})
export class DataTableComponent {
  @Input() columns: TableColumn[] = [];
  @Input() loading = false;
  @Input() sortBy = '';
  @Input() sortDir: 'asc' | 'desc' = 'asc';
  @Input() trackBy?: (row: any) => any;

  // ✅ Toggle visibility tombol Aksi
  @Input() showPdf = false;       // default: sembunyikan tombol PDF
  @Input() showDetail = false;    // default: sembunyikan tombol Detail
  @Input() showEdit = true;       // default: tampilkan Edit
  @Input() showDelete = true;     // default: tampilkan Hapus

  // ✅ Auto width kolom Aksi berdasarkan tombol yang aktif
  get actionColumnWidth(): 'sm' | 'md' | 'lg' | 'xl' {
    const count =
      (this.showPdf ? 1 : 0) +
      (this.showDetail ? 1 : 0) +
      (this.showEdit ? 1 : 0) +
      (this.showDelete ? 1 : 0);

    if (count <= 1) return 'sm';
    if (count === 2) return 'md';
    if (count === 3) return 'lg';
    return 'xl';
  }

  // ✅ Pakai backing field + getter/setter agar aman dari null/undefined
  private _page: PageResult<any> | null = null;

  @Input()
  set page(value: PageResult<any> | null | undefined) {
    this._page = value ?? null;
  }
  get page(): PageResult<any> | null {
    return this._page;
  }

  @Output() pageChange = new EventEmitter<number>();
  @Output() sortChange = new EventEmitter<{ sortBy: string; sortDir: 'asc' | 'desc' }>();
  @Output() edit = new EventEmitter<any>();
  @Output() delete = new EventEmitter<any>();
  @Output() detail = new EventEmitter<any>();
  @Output() pdf = new EventEmitter<any>();

  Math = Math;

  get rows(): any[] {
    const data = this._page?.data;
    return Array.isArray(data) ? data : [];
  }

  get total(): number {
    return Number(this._page?.total ?? 0);
  }

  get currentPage(): number {
    return Number(this._page?.page ?? 1);
  }

  get totalPages(): number {
    const t = Number(this._page?.totalPages ?? 1);
    return t < 1 ? 1 : t;
  }

  get pageSize(): number {
    return Number(this._page?.size ?? 10);
  }

  get startIndex(): number {
    return (this.currentPage - 1) * this.pageSize;
  }

  get endIndex(): number {
    return Math.min(this.startIndex + this.rows.length, this.total);
  }

  get skeletonCols(): any[] {
    const len = this.columns?.length ?? 0;
    return new Array(Math.max(1, len)).fill(0);
  }

  get pages(): number[] {
    if (!this._page) return [];
    const total = this.totalPages;
    const current = this.currentPage;
    const arr: number[] = [];
    const start = Math.max(1, current - 2);
    const end = Math.min(total, current + 2);
    for (let i = start; i <= end; i++) arr.push(i);
    return arr;
  }

  getValue(row: any, key: string): any {
    if (!row || !key) return '-';
    return key.split('.').reduce((o, k) => (o == null ? undefined : o[k]), row) ?? '-';
  }

  goTo(p: number) {
    if (p >= 1 && p <= this.totalPages && p !== this.currentPage) {
      this.pageChange.emit(p);
    }
  }

  toggleSort(key: string) {
    const dir = this.sortBy === key && this.sortDir === 'asc' ? 'desc' : 'asc';
    this.sortChange.emit({ sortBy: key, sortDir: dir });
  }
}