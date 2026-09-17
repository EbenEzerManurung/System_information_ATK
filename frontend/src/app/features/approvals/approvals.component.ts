import {
  Component,
  ElementRef,
  OnInit,
  ViewChild,
  signal,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../../core/services/api.service';
import { ToastService } from '../../core/services/toast.service';
import { PageResult, Transaction } from '../../core/models';
import {
  DataTableComponent,
  TableColumn,
} from '../../shared/components/data-table/data-table.component';
import { FilterBarComponent } from '../../shared/components/filter-bar/filter-bar.component';

@Component({
  selector: 'app-approvals',
  standalone: true,
  imports: [CommonModule, FormsModule, DataTableComponent, FilterBarComponent],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="space-y-5 fade-in">
      <div>
        <h2 class="text-2xl font-bold text-primary-900">Approval & Signature Digital</h2>
        <p class="text-sm text-slate-500">
          Review dan setujui transaksi dengan tanda tangan digital
        </p>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-primary-100 shadow-sm">
        <app-filter-bar placeholder="Cari transaksi..." (search)="onSearch($event)">
          <select
            [(ngModel)]="filterStatus"
            (change)="loadData(1)"
            class="px-3 py-2 text-sm border border-primary-200 rounded-lg outline-none"
          >
            <option value="draft">Draft (Menunggu Approval)</option>
            <option value="pending">Pending (Sedang Ditinjau)</option>
            <option value="approved">Approved</option>
            <option value="rejected">Rejected</option>
            <option value="">Semua Status</option>
          </select>
        </app-filter-bar>
      </div>

      <app-data-table
        [columns]="columns"
        [page]="page()"
        [loading]="loading()"
        [trackBy]="trackById"
        (pageChange)="loadData($event)"
        (edit)="openApproval($event)"
        (delete)="onDelete($event)"
      />

      @if (showModal()) {
        <div
          class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4 overflow-y-auto"
        >
          <div class="bg-white rounded-2xl shadow-2xl w-full max-w-2xl p-6 my-8 fade-in">
            <h3 class="text-lg font-bold text-primary-900 mb-4">Approval Transaksi</h3>

            <div class="grid grid-cols-2 gap-3 mb-4 p-4 bg-primary-50 rounded-xl">
              <div>
                <div class="text-xs text-slate-500">No Transaksi</div>
                <div class="font-semibold text-primary-900">
                  {{ current?.transactionCode ?? '-' }}
                </div>
              </div>
              <div>
                <div class="text-xs text-slate-500">Tipe</div>
                <div class="font-semibold text-primary-900">{{ current?.type ?? '-' }}</div>
              </div>
              <div>
                <div class="text-xs text-slate-500">Total Item</div>
                <div class="font-semibold text-primary-900">
                  {{ current?.totalItems ?? 0 }}
                </div>
              </div>
              <div>
                <div class="text-xs text-slate-500">Total Value</div>
                <div class="font-semibold text-primary-900">
                  {{ current?.totalValue ?? 0 }}
                </div>
              </div>
              <div class="col-span-2">
                <div class="text-xs text-slate-500">Catatan</div>
                <div class="font-semibold text-primary-900">
                  {{ current?.notes ?? '-' }}
                </div>
              </div>
            </div>

            @if (canApprove()) {
              <div class="mb-4">
                <label class="block text-sm font-medium text-slate-700 mb-2">
                  Tanda Tangan Digital
                </label>
                <div class="border-2 border-dashed border-primary-200 rounded-xl p-2 bg-primary-50/30">
                  <canvas
                    #canvas
                    width="600"
                    height="200"
                    class="w-full bg-white rounded-lg cursor-crosshair touch-none"
                  ></canvas>
                </div>
                <div class="flex justify-end gap-2 mt-2">
                  <button
                    (click)="clearSignature()"
                    class="px-3 py-1.5 text-xs rounded-lg bg-slate-100 hover:bg-slate-200 transition"
                  >
                    Hapus Tanda Tangan
                  </button>
                </div>
              </div>

              <div class="mb-4">
                <label class="block text-sm font-medium text-slate-700 mb-1">Catatan</label>
                <textarea
                  [(ngModel)]="notes"
                  rows="2"
                  placeholder="Catatan approval..."
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30"
                ></textarea>
              </div>
            } @else {
              <div class="mb-4 p-3 rounded-lg bg-slate-50 border border-slate-200 text-sm text-slate-600">
                Transaksi berstatus <strong>{{ current?.status }}</strong>.
                Tidak dapat di-approve lagi.
              </div>
            }

            <div class="flex justify-end gap-2 flex-wrap">
              <button
                (click)="closeModal()"
                class="px-4 py-2 text-sm rounded-lg bg-slate-100 hover:bg-slate-200 transition"
              >
                Tutup
              </button>

              @if (canApprove()) {
                <!-- ✅ Tombol PENDING — user masih mempertimbangkan -->
                <button
                  (click)="setPending()"
                  [disabled]="processing()"
                  class="px-4 py-2 text-sm rounded-lg bg-amber-500 text-white hover:bg-amber-600 disabled:opacity-60 transition"
                  title="Tandai sebagai pending untuk ditinjau lebih lanjut"
                >
                  {{ processing() ? 'Memproses...' : 'Pending' }}
                </button>

                <button
                  (click)="reject()"
                  [disabled]="processing()"
                  class="px-4 py-2 text-sm rounded-lg bg-red-600 text-white hover:bg-red-700 disabled:opacity-60 transition"
                >
                  {{ processing() ? 'Memproses...' : 'Tolak' }}
                </button>
                <button
                  (click)="approve()"
                  [disabled]="processing()"
                  class="px-4 py-2 text-sm rounded-lg bg-primary-600 text-white hover:bg-primary-700 disabled:opacity-60 transition"
                >
                  {{ processing() ? 'Memproses...' : 'Setujui & Tanda Tangan' }}
                </button>
              }
            </div>

            @if (canApprove()) {
              <p class="text-xs text-slate-400 mt-3 text-center">
                <strong>Pending</strong> = tandai untuk ditinjau lebih lanjut
                (tetap bisa di-approve nanti)
              </p>
            }
          </div>
        </div>
      }
    </div>
  `,
})
export class ApprovalsComponent implements OnInit {
  @ViewChild('canvas') canvasRef?: ElementRef<HTMLCanvasElement>;

  columns: TableColumn[] = [
    { key: 'transactionCode', label: 'No Transaksi', sortable: true },
    { key: 'type', label: 'Tipe' },
    { key: 'totalItems', label: 'Total Item' },
    { key: 'totalValue', label: 'Total Value' },
    { key: 'status', label: 'Status' },
    { key: 'user.username', label: 'Diminta Oleh' },
    { key: 'transactionDate', label: 'Tanggal' },
  ];

  page = signal<PageResult<Transaction> | null>(null);
  loading = signal(false);
  processing = signal(false);
  showModal = signal(false);
  search = '';
  filterStatus = '';
  currentPage = 1;
  current: Transaction | null = null;
  notes = '';

  private ctx?: CanvasRenderingContext2D;
  private isDrawing = false;
  private hasSignature = false;

  trackById = (r: any) => r.id;

  constructor(
    private api: ApiService,
    private toast: ToastService,
  ) {}

  ngOnInit() {
    this.loadData(1);
  }

  // ============================================================
  // Helper: cek apakah transaksi bisa di-approve
  // ============================================================
  canApprove(): boolean {
    const status = String(this.current?.status ?? '').toLowerCase().trim();
    return status === 'draft' || status === 'pending';
  }

  // ============================================================
  // LOAD DATA
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
        status: this.filterStatus,
      })
      .subscribe({
        next: (res: any) => {
          this.loading.set(false);
          console.log('[Approvals] Raw response:', res);

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
          console.error('[Approvals] Error:', err);
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

  openApproval(row: Transaction) {
    this.current = row;
    this.notes = row.notes ?? '';
    this.hasSignature = false;
    this.showModal.set(true);

    setTimeout(() => this.initCanvas(), 100);
  }

  closeModal() {
    this.showModal.set(false);
    this.current = null;
    this.hasSignature = false;
    this.notes = '';
  }

  // ============================================================
  // CANVAS SIGNATURE
  // ============================================================
  private initCanvas() {
    const canvas = this.canvasRef?.nativeElement;
    if (!canvas) {
      console.error('[Approvals] Canvas not found');
      return;
    }

    canvas.width = canvas.width;

    this.ctx = canvas.getContext('2d')!;
    this.ctx.strokeStyle = '#1e3a8a';
    this.ctx.lineWidth = 2;
    this.ctx.lineCap = 'round';
    this.ctx.lineJoin = 'round';

    const getPos = (e: MouseEvent | TouchEvent) => {
      const rect = canvas.getBoundingClientRect();
      const scaleX = canvas.width / rect.width;
      const scaleY = canvas.height / rect.height;
      const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX;
      const clientY = 'touches' in e ? e.touches[0].clientY : e.clientY;
      return {
        x: (clientX - rect.left) * scaleX,
        y: (clientY - rect.top) * scaleY,
      };
    };

    const start = (e: MouseEvent | TouchEvent) => {
      e.preventDefault();
      this.isDrawing = true;
      this.hasSignature = true;
      const p = getPos(e);
      this.ctx!.beginPath();
      this.ctx!.moveTo(p.x, p.y);
    };

    const draw = (e: MouseEvent | TouchEvent) => {
      if (!this.isDrawing) return;
      e.preventDefault();
      const p = getPos(e);
      this.ctx!.lineTo(p.x, p.y);
      this.ctx!.stroke();
    };

    const end = () => {
      this.isDrawing = false;
    };

    canvas.addEventListener('mousedown', start);
    canvas.addEventListener('mousemove', draw);
    canvas.addEventListener('mouseup', end);
    canvas.addEventListener('mouseleave', end);
    canvas.addEventListener('touchstart', start, { passive: false });
    canvas.addEventListener('touchmove', draw, { passive: false });
    canvas.addEventListener('touchend', end);
  }

  clearSignature() {
    const canvas = this.canvasRef?.nativeElement;
    if (canvas && this.ctx) {
      this.ctx.clearRect(0, 0, canvas.width, canvas.height);
      this.hasSignature = false;
    }
  }

  getSignature(): string {
    return this.canvasRef?.nativeElement.toDataURL('image/png') ?? '';
  }

  // ============================================================
  // ✅ SET PENDING — PUT /transactions/:id dengan status=pending
  // User masih mempertimbangkan, jadi hanya ubah status
  // ============================================================
  setPending() {
    if (!this.current) return;

    if (!confirm(`Tandai transaksi "${this.current.transactionCode}" sebagai PENDING?`)) {
      return;
    }

    this.processing.set(true);

    // Panggil 2 endpoint sekaligus sebagai fallback
    // 1. POST /transactions/:id/submit (ubah draft → pending)
    // 2. PUT /transactions/:id dengan status=pending

    const id = this.current.id;

    this.api.post<any>(`transactions/${id}/submit`, {}).subscribe({
      next: (res: any) => {
        this.processing.set(false);
        if (res?.success !== false) {
          this.toast.success('Transaksi ditandai sebagai PENDING');
          this.closeModal();
          this.loadData(this.currentPage);
        } else {
          this.toast.error(res?.message || 'Gagal menandai pending');
        }
      },
      error: (err) => {
        console.error('[Approvals] Submit error, fallback to PUT:', err);

        // Fallback: PUT dengan status
        this.api.put<any>(`transactions/${id}`, { status: 'pending' }).subscribe({
          next: (res: any) => {
            this.processing.set(false);
            if (res?.success !== false) {
              this.toast.success('Transaksi ditandai sebagai PENDING');
              this.closeModal();
              this.loadData(this.currentPage);
            } else {
              this.toast.error(res?.message || 'Gagal menandai pending');
            }
          },
          error: (err2) => {
            this.processing.set(false);
            console.error('[Approvals] PUT error:', err2);
            this.toast.error(err2?.error?.message || 'Gagal menandai pending');
          },
        });
      },
    });
  }

  // ============================================================
  // APPROVE
  // ============================================================
  approve() {
    if (!this.current) return;
    if (!this.hasSignature) {
      this.toast.warning('Tanda tangan masih kosong');
      return;
    }

    this.processing.set(true);

    const sig = this.getSignature();

    const payload = {
      signature: sig,
      signatureData: sig,
      notes: this.notes,
    };

    console.log('[Approvals] Approve transaction ID:', this.current.id);
    console.log('[Approvals] Signature length:', sig.length);

    this.api
      .post<any>(`transactions/${this.current.id}/approve`, payload)
      .subscribe({
        next: (res: any) => {
          this.processing.set(false);
          if (res?.success !== false) {
            this.toast.success('Transaksi berhasil disetujui');
            this.closeModal();
            this.loadData(this.currentPage);
          } else {
            this.toast.error(res?.message || 'Gagal menyetujui');
          }
        },
        error: (err) => {
          this.processing.set(false);
          console.error('[Approvals] Approve error:', err);
          const msg =
            err?.error?.message ||
            err?.error?.error ||
            err?.message ||
            'Gagal memproses';
          this.toast.error(msg);
        },
      });
  }

  // ============================================================
  // REJECT
  // ============================================================
  reject() {
    if (!this.current) return;
    if (!confirm(`Tolak transaksi "${this.current.transactionCode}"?`)) return;

    this.processing.set(true);

    const payload = {
      notes: this.notes,
      reason: this.notes,
    };

    this.api
      .post<any>(`transactions/${this.current.id}/reject`, payload)
      .subscribe({
        next: (res: any) => {
          this.processing.set(false);
          if (res?.success !== false) {
            this.toast.success('Transaksi ditolak');
            this.closeModal();
            this.loadData(this.currentPage);
          } else {
            this.toast.error(res?.message || 'Gagal menolak');
          }
        },
        error: (err) => {
          this.processing.set(false);
          console.error('[Approvals] Reject error:', err);
          this.toast.error(err?.error?.message || 'Gagal memproses');
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
}