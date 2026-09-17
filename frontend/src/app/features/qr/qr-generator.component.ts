import { Component, OnInit, signal, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import QRCode from 'qrcode';
import { ApiService } from '../../core/services/api.service';
import { ToastService } from '../../core/services/toast.service';

@Component({
  selector: 'app-qr-generator',
  standalone: true,
  imports: [CommonModule, FormsModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="space-y-5 fade-in max-w-5xl">
      <div>
        <h2 class="text-2xl font-bold text-primary-900">Generate QR Code Dokumen</h2>
        <p class="text-sm text-slate-500">
          Pilih transaksi untuk membuat QR code dokumen
        </p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- ==================== FORM ==================== -->
        <div class="bg-white rounded-2xl p-6 border border-primary-100 shadow-sm space-y-4">
          <!-- Pilih Transaksi -->
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">
              Pilih Transaksi <span class="text-red-500">*</span>
            </label>

            @if (loadingTx()) {
              <div class="w-full px-3 py-2 border border-primary-100 rounded-lg bg-slate-50 text-sm text-slate-400">
                Memuat daftar transaksi...
              </div>
            } @else {
              <select
                [(ngModel)]="selectedTxId"
                (change)="onTransactionChange()"
                class="w-full px-3 py-2 border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30">
                <option [ngValue]="null">-- Pilih Transaksi --</option>
                @for (tx of transactions(); track tx.id) {
                  <option [ngValue]="tx.id">
                    {{ tx.transactionCode }} — {{ tx.type }} — {{ tx.status }}
                  </option>
                }
              </select>

              @if (transactions().length === 0) {
                <p class="text-xs text-amber-600 mt-1">
                  Belum ada transaksi. Buat transaksi dulu di menu Transaksi.
                </p>
              }
            }
          </div>

          <!-- Tombol Refresh -->
          <div class="flex justify-end">
            <button (click)="loadTransactions()"
              class="px-3 py-1.5 text-xs rounded-lg bg-slate-100 hover:bg-slate-200 transition flex items-center gap-1">
              <span>⟳</span> Refresh Daftar
            </button>
          </div>

          <!-- Detail Transaksi Terpilih -->
          @if (selectedTransaction()) {
            <div class="p-4 bg-primary-50 rounded-xl space-y-2">
              <div class="text-xs text-slate-500 uppercase font-semibold">Detail Transaksi</div>
              <div class="grid grid-cols-2 gap-2 text-sm">
                <div>
                  <div class="text-xs text-slate-500">No Transaksi</div>
                  <div class="font-semibold text-primary-900">
                    {{ selectedTransaction()?.transactionCode }}
                  </div>
                </div>
                <div>
                  <div class="text-xs text-slate-500">Tipe</div>
                  <div class="font-semibold text-primary-900">
                    {{ selectedTransaction()?.type }}
                  </div>
                </div>
                <div>
                  <div class="text-xs text-slate-500">Total Item</div>
                  <div class="font-semibold text-primary-900">
                    {{ selectedTransaction()?.totalItems ?? 0 }}
                  </div>
                </div>
                <div>
                  <div class="text-xs text-slate-500">Total Value</div>
                  <div class="font-semibold text-primary-900">
                    {{ selectedTransaction()?.totalValue ?? 0 }}
                  </div>
                </div>
                <div>
                  <div class="text-xs text-slate-500">Status</div>
                  <div class="font-semibold text-primary-900 capitalize">
                    {{ selectedTransaction()?.status }}
                  </div>
                </div>
                <div>
                  <div class="text-xs text-slate-500">Tanggal</div>
                  <div class="font-semibold text-primary-900">
                    {{ selectedTransaction()?.transactionDate | date: 'dd/MM/yyyy' }}
                  </div>
                </div>
              </div>
            </div>
          }

          <!-- Deskripsi -->
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">Deskripsi</label>
            <textarea [(ngModel)]="description" rows="2"
              placeholder="Keterangan dokumen..."
              class="w-full px-3 py-2 border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30"></textarea>
          </div>

          <!-- Ukuran QR -->
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">Ukuran QR</label>
            <select [(ngModel)]="size"
              class="w-full px-3 py-2 border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30">
              <option [ngValue]="300">Sedang (300px)</option>
              <option [ngValue]="400">Besar (400px)</option>
              <option [ngValue]="500">Sangat Besar (500px)</option>
              <option [ngValue]="600">Extra (600px)</option>
            </select>
          </div>

          <!-- Info Payload -->
          @if (payloadPreview()) {
            <div class="p-3 bg-slate-50 rounded-lg">
              <div class="text-xs text-slate-500 font-mono break-all">
                <span class="font-semibold">Payload:</span> {{ payloadPreview() }}
              </div>
              <div class="text-xs text-slate-400 mt-1">
                Panjang: {{ payloadPreview().length }} karakter
                @if (payloadPreview().length > 100) {
                  <span class="text-amber-600 font-medium"> ⚠️ Terlalu panjang, QR akan padat</span>
                }
              </div>
            </div>
          }

          <button (click)="generate()" [disabled]="!selectedTxId"
            class="w-full py-2.5 rounded-lg bg-primary-600 text-white hover:bg-primary-700 disabled:opacity-60 disabled:cursor-not-allowed transition font-medium">
            Generate QR Code
          </button>
        </div>

        <!-- ==================== PREVIEW ==================== -->
        <div class="bg-white rounded-2xl p-6 border border-primary-100 shadow-sm">
          <h3 class="font-semibold text-primary-900 mb-4 text-center">Preview</h3>

          @if (qrDataUrl()) {
            <div class="flex flex-col items-center gap-4">
              <img [src]="qrDataUrl()" alt="QR Code"
                   class="rounded-xl border-4 border-primary-100 max-w-full" />

              <div class="text-center">
                <div class="text-sm font-semibold text-primary-900">
                  {{ selectedTransaction()?.transactionCode }}
                </div>
                @if (description) {
                  <div class="text-xs text-slate-500 mt-1">{{ description }}</div>
                }
              </div>

              <div class="flex gap-2 flex-wrap justify-center">
                <button (click)="downloadPng()"
                  class="px-4 py-2 text-sm rounded-lg bg-primary-600 text-white hover:bg-primary-700 transition">
                  Download PNG
                </button>
                <button (click)="downloadSvg()"
                  class="px-4 py-2 text-sm rounded-lg bg-white border border-primary-200 text-primary-700 hover:bg-primary-50 transition">
                  Download SVG
                </button>
                <button (click)="print()"
                  class="px-4 py-2 text-sm rounded-lg bg-white border border-primary-200 text-primary-700 hover:bg-primary-50 transition">
                  Print
                </button>
              </div>
            </div>
          } @else {
            <div class="flex flex-col items-center justify-center py-20 text-slate-400">
              <div class="text-5xl mb-3">▦</div>
              <div class="text-sm">Belum ada QR code</div>
              <div class="text-xs mt-1">Pilih transaksi lalu klik Generate</div>
            </div>
          }
        </div>
      </div>
    </div>
  `,
})
export class QrGeneratorComponent implements OnInit {
  transactions = signal<any[]>([]);
  loadingTx = signal(false);
  selectedTxId: number | null = null;

  description = '';
  size = 400;
  qrDataUrl = signal<string>('');
  payloadPreview = signal<string>('');

  constructor(
    private api: ApiService,
    private toast: ToastService,
  ) {}

  ngOnInit() {
    this.loadTransactions();
  }

  // ============================================================
  // LOAD TRANSACTIONS
  // ============================================================
  loadTransactions() {
    this.loadingTx.set(true);

    this.api.get<any>('transactions', { page: 1, size: 1000 }).subscribe({
      next: (res: any) => {
        this.loadingTx.set(false);
        console.log('[QrGenerator] Transactions response:', res);

        const root = res?.data ?? res;
        const inner = root?.data ?? root;
        const items: any[] = Array.isArray(inner)
          ? inner
          : Array.isArray(root)
            ? root
            : [];

        this.transactions.set(items);
        console.log('[QrGenerator] Loaded:', items.length, 'transactions');
      },
      error: (err) => {
        this.loadingTx.set(false);
        console.error('[QrGenerator] Load error:', err);
        this.toast.error('Gagal memuat daftar transaksi');
      },
    });
  }

  // ============================================================
  // ✅ BUILD PAYLOAD — SUPER COMPACT
  // Format: CODE|status|type|items|value|user|notes
  // Contoh: TRX-20260915-5500|approved|OUT|1|25000|superadmin|oke
  // Panjang: ~45 karakter (jauh lebih ringkas dari JSON)
  // ============================================================
  private buildPayload(tx: any): string {
    const safeNotes = String(tx.notes ?? '').replace(/[|\n\r]/g, ' ').trim();
    return [
      tx.transactionCode,                    // 1. Kode
      String(tx.status ?? '').toLowerCase(), // 2. Status (draft/pending/approved/...)
      String(tx.type ?? '').toUpperCase(),   // 3. Tipe (IN/OUT/RETURN/ADJUSTMENT)
      tx.totalItems ?? 0,                    // 4. Total items
      tx.totalValue ?? 0,                    // 5. Total value
      tx.user?.username ?? tx.user?.full_name ?? '-', // 6. User
      safeNotes,                             // 7. Notes
    ].join('|');
  }

  // ============================================================
  // HANDLE TRANSACTION CHANGE
  // ============================================================
  onTransactionChange() {
    this.qrDataUrl.set('');
    this.payloadPreview.set('');

    const tx = this.selectedTransaction();
    if (tx) {
      this.description =
        `${tx.type} - ${tx.totalItems ?? 0} item - Total ${tx.totalValue ?? 0}`;
      this.payloadPreview.set(this.buildPayload(tx));
    }
  }

  // ============================================================
  // GET SELECTED TRANSACTION
  // ============================================================
  selectedTransaction(): any | null {
    if (!this.selectedTxId) return null;
    return (
      this.transactions().find(
        (t: any) => Number(t.id) === Number(this.selectedTxId),
      ) ?? null
    );
  }

  // ============================================================
  // ✅ GENERATE QR — pakai payload COMPACT
  // ============================================================
  async generate() {
    const tx = this.selectedTransaction();
    if (!tx) {
      this.toast.warning('Pilih transaksi terlebih dahulu');
      return;
    }

    const qrText = this.buildPayload(tx);
    this.payloadPreview.set(qrText);
    console.log(
      '[QrGenerator] Payload:',
      qrText,
      '| Length:',
      qrText.length,
    );

    try {
      const url = await QRCode.toDataURL(qrText, {
        width: Number(this.size),
        margin: 4,                      // ✅ margin lebih lebar untuk kamera
        errorCorrectionLevel: 'M',      // ✅ M (medium) — cukup, tidak overkill
        color: { dark: '#000000', light: '#ffffff' }, // ✅ hitam pekat
      });

      this.qrDataUrl.set(url);
      this.toast.success(
        `QR code berhasil (${qrText.length} karakter)`,
      );
    } catch (err) {
      console.error('[QrGenerator] Error:', err);
      this.toast.error('Gagal generate QR code');
    }
  }

  // ============================================================
  // DOWNLOAD PNG
  // ============================================================
  downloadPng() {
    const url = this.qrDataUrl();
    const tx = this.selectedTransaction();
    if (!url || !tx) return;

    const a = document.createElement('a');
    a.href = url;
    a.download = `qr_${tx.transactionCode}.png`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
  }

  // ============================================================
  // DOWNLOAD SVG
  // ============================================================
  async downloadSvg() {
    const tx = this.selectedTransaction();
    if (!tx) return;

    const payload = this.buildPayload(tx);

    try {
      const svg = await QRCode.toString(payload, {
        type: 'svg',
        margin: 4,
        errorCorrectionLevel: 'M',
        color: { dark: '#000000', light: '#ffffff' },
      });

      const blob = new Blob([svg], { type: 'image/svg+xml' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `qr_${tx.transactionCode}.svg`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    } catch (err) {
      console.error('[QrGenerator] SVG error:', err);
      this.toast.error('Gagal download SVG');
    }
  }

  // ============================================================
  // PRINT
  // ============================================================
  print() {
    const url = this.qrDataUrl();
    const tx = this.selectedTransaction();
    if (!url || !tx) return;

    const w = window.open('', '_blank');
    if (!w) {
      this.toast.error('Popup diblokir. Izinkan popup untuk print.');
      return;
    }

    w.document.write(`
      <!DOCTYPE html>
      <html>
      <head>
        <title>Print QR - ${tx.transactionCode}</title>
        <style>
          body {
            text-align: center;
            font-family: sans-serif;
            padding: 40px;
          }
          .code { font-size: 14px; color: #666; margin-bottom: 20px; }
          .info { text-align: left; display: inline-block; margin-top: 20px; }
          .info div { margin: 4px 0; font-size: 14px; }
          .info strong { color: #1e3a8a; }
        </style>
      </head>
      <body>
        <h2>${tx.transactionCode}</h2>
        <p class="code">${this.description || ''}</p>
        <img src="${url}" style="width:${this.size}px;max-width:100%">
        <div class="info">
          <div><strong>Tipe:</strong> ${tx.type}</div>
          <div><strong>Status:</strong> ${tx.status}</div>
          <div><strong>Total Item:</strong> ${tx.totalItems ?? 0}</div>
          <div><strong>Total Value:</strong> ${tx.totalValue ?? 0}</div>
          <div><strong>Tanggal:</strong> ${tx.transactionDate ?? '-'}</div>
        </div>
        <script>window.onload=()=>{window.print();window.close()}<\/script>
      </body>
      </html>
    `);
    w.document.close();
  }
}