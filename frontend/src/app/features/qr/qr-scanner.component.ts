import {
  Component,
  OnInit,
  OnDestroy,
  signal,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ZXingScannerModule } from '@zxing/ngx-scanner';
import { BarcodeFormat, BrowserMultiFormatReader } from '@zxing/library';
import { ApiService } from '../../core/services/api.service';
import { ToastService } from '../../core/services/toast.service';

interface ParsedQRData {
  documentCode?: string;
  transactionCode?: string;
  type?: string;
  status?: string;
  totalItems?: number;
  totalValue?: number;
  transactionDate?: string;
  user?: string;
  notes?: string;
  docNo?: string;
}

@Component({
  selector: 'app-qr-scanner',
  standalone: true,
  imports: [CommonModule, FormsModule, ZXingScannerModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="space-y-5 fade-in max-w-5xl">
      <div>
        <h2 class="text-2xl font-bold text-primary-900">Scan QR Code Dokumen</h2>
        <p class="text-sm text-slate-500">Validasi dokumen dengan scan QR</p>
      </div>

      <!-- Tab -->
      <div class="bg-white rounded-2xl p-1 border border-primary-100 shadow-sm inline-flex">
        <button (click)="switchMode('camera')"
          class="px-4 py-2 text-sm rounded-xl font-medium transition"
          [class.bg-primary-600]="mode() === 'camera'"
          [class.text-white]="mode() === 'camera'"
          [class.text-slate-600]="mode() !== 'camera'">
          📷 Kamera
        </button>
        <button (click)="switchMode('upload')"
          class="px-4 py-2 text-sm rounded-xl font-medium transition"
          [class.bg-primary-600]="mode() === 'upload'"
          [class.text-white]="mode() === 'upload'"
          [class.text-slate-600]="mode() !== 'upload'">
          🖼️ Upload Gambar
        </button>
        <button (click)="switchMode('manual')"
          class="px-4 py-2 text-sm rounded-xl font-medium transition"
          [class.bg-primary-600]="mode() === 'manual'"
          [class.text-white]="mode() === 'manual'"
          [class.text-slate-600]="mode() !== 'manual'">
          ⌨️ Input Manual
        </button>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- INPUT -->
        <div class="bg-white rounded-2xl p-4 border border-primary-100 shadow-sm">

          <!-- ===== MODE: KAMERA ===== -->
          @if (mode() === 'camera') {
            <div class="flex items-center justify-between mb-3">
              <div class="text-sm font-medium text-primary-900">
                Kamera
                @if (cameraReady()) {
                  <span class="ml-2 inline-flex items-center gap-1 text-xs text-green-600">
                    <span class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
                    Aktif
                  </span>
                }
              </div>
              <button (click)="resetScanner()"
                class="px-3 py-1 text-xs rounded-lg bg-slate-100 hover:bg-slate-200 transition">
                ⟳ Reset
              </button>
            </div>

            <div class="aspect-square bg-primary-900 rounded-xl overflow-hidden relative">
              <zxing-scanner
                [formats]="allowedFormats"
                [tryHarder]="true"
                [autofocusEnabled]="true"
                (scanSuccess)="onScan($event)"
                (camerasFound)="onCamerasFound($event)"
                (permissionResponse)="onPermission($event)"
                class="w-full h-full">
              </zxing-scanner>
            </div>

            <div class="text-center text-xs text-slate-500 mt-3">
              💡 Arahkan QR ke kamera, jarak 10-20 cm dengan cahaya cukup
            </div>

            @if (!hasPermission()) {
              <div class="mt-3 p-3 rounded-lg bg-amber-50 border border-amber-200 text-xs text-amber-700">
                ⚠️ Izin kamera diperlukan
              </div>
            }
          }

          <!-- ===== MODE: UPLOAD ===== -->
          @if (mode() === 'upload') {
            <div class="text-sm font-medium text-primary-900 mb-3">Upload Gambar QR</div>

            <div class="border-2 border-dashed border-primary-200 rounded-xl p-8 text-center bg-primary-50/30">
              <input type="file" accept="image/*" #fileInput
                (change)="onFileSelected($event)" class="hidden" />
              <div class="text-4xl mb-3">🖼️</div>
              <div class="text-sm text-slate-600 mb-3">Pilih gambar QR</div>
              <button (click)="fileInput.click()"
                class="px-4 py-2 text-sm rounded-lg bg-primary-600 text-white hover:bg-primary-700 transition">
                Pilih Gambar
              </button>
            </div>

            @if (uploadedPreview()) {
              <div class="mt-4 text-center">
                <img [src]="uploadedPreview()" alt="Preview"
                     class="max-w-full max-h-64 mx-auto rounded-lg border border-primary-200" />
                @if (decodeStatus()) {
                  <div class="text-xs mt-2"
                       [class.text-green-600]="decodeOk()"
                       [class.text-red-600]="!decodeOk()">
                    {{ decodeStatus() }}
                  </div>
                }
              </div>
            }

            <div class="text-xs text-slate-500 mt-3 text-center">
              💡 Cocok kalau kamera tidak bisa fokus ke QR
            </div>
          }

          <!-- ===== MODE: MANUAL ===== -->
          @if (mode() === 'manual') {
            <div class="text-sm font-medium text-primary-900 mb-3">Input Kode Manual</div>
            <div class="space-y-3">
              <div>
                <label class="block text-sm text-slate-700 mb-1">Kode Dokumen / Transaksi</label>
                <input [(ngModel)]="manualCode"
                  placeholder="Contoh: TRX-20260915-5500"
                  (keyup.enter)="submitManual()"
                  class="w-full px-3 py-2 border border-primary-200 rounded-lg outline-none focus:ring-2 focus:ring-primary-500/30" />
              </div>
              <button (click)="submitManual()" [disabled]="!manualCode.trim()"
                class="w-full py-2.5 rounded-lg bg-primary-600 text-white hover:bg-primary-700 disabled:opacity-60 transition font-medium">
                Validasi Kode
              </button>
            </div>
            <div class="text-xs text-slate-500 mt-3 text-center">
              💡 Test validasi tanpa kamera
            </div>
          }
        </div>

        <!-- HASIL -->
        <div class="bg-white rounded-2xl p-6 border border-primary-100 shadow-sm">
          <div class="flex items-center justify-between mb-4">
            <h3 class="font-semibold text-primary-900">Hasil Scan</h3>
            @if (scanResult()) {
              <button (click)="clearResult()"
                class="px-2 py-1 text-xs rounded-lg bg-slate-100 hover:bg-slate-200 transition">
                Clear
              </button>
            }
          </div>

          @if (scanResult()) {
            <div class="space-y-3 fade-in">
              @if (validationStatus() === 'checking') {
                <div class="p-4 rounded-xl bg-blue-50 border border-blue-200">
                  <div class="flex items-center gap-2">
                    <div class="animate-spin rounded-full h-4 w-4 border-2 border-blue-600 border-t-transparent"></div>
                    <span class="text-sm text-blue-700 font-medium">Memvalidasi...</span>
                  </div>
                </div>
              }

              @if (validationStatus() === 'valid') {
                <div class="p-4 rounded-xl border" [ngClass]="statusBadgeClass(currentStatus())">
                  <div class="flex items-center gap-2">
                    <span class="text-2xl">{{ statusIcon(currentStatus()) }}</span>
                    <div>
                      <div class="text-sm font-semibold uppercase" [ngClass]="statusTextClass(currentStatus())">
                        {{ statusLabel(currentStatus()) }}
                      </div>
                      <div class="text-xs" [ngClass]="statusSubTextClass(currentStatus())">
                        Dokumen terverifikasi di sistem
                      </div>
                    </div>
                  </div>
                </div>
              }

              @if (validationStatus() === 'invalid') {
                <div class="p-4 rounded-xl bg-red-50 border border-red-200">
                  <div class="flex items-center gap-2">
                    <span class="text-2xl">❌</span>
                    <div>
                      <div class="text-sm font-semibold text-red-700">DOKUMEN TIDAK VALID</div>
                      <div class="text-xs text-red-600">
                        {{ validationMessage() || 'Dokumen tidak ditemukan' }}
                      </div>
                    </div>
                  </div>
                </div>
              }

              @if (parsedData(); as data) {
                <div class="p-4 bg-primary-50 rounded-xl">
                  <div class="text-xs text-slate-500 uppercase font-semibold mb-2">Data Dokumen</div>
                  <div class="grid grid-cols-2 gap-2 text-sm">
                    @if (data.documentCode || data.transactionCode) {
                      <div class="col-span-2">
                        <div class="text-xs text-slate-500">No Dokumen</div>
                        <div class="font-semibold text-primary-900">
                          {{ data.documentCode || data.transactionCode }}
                        </div>
                      </div>
                    }
                    @if (data.type) {
                      <div>
                        <div class="text-xs text-slate-500">Tipe</div>
                        <div class="font-semibold text-primary-900">{{ data.type }}</div>
                      </div>
                    }
                    @if (data.status) {
                      <div>
                        <div class="text-xs text-slate-500">Status</div>
                        <div class="font-semibold text-primary-900 capitalize">{{ data.status }}</div>
                      </div>
                    }
                    @if (data.totalItems !== undefined && data.totalItems !== null) {
                      <div>
                        <div class="text-xs text-slate-500">Total Item</div>
                        <div class="font-semibold text-primary-900">{{ data.totalItems }}</div>
                      </div>
                    }
                    @if (data.totalValue !== undefined && data.totalValue !== null) {
                      <div>
                        <div class="text-xs text-slate-500">Total Value</div>
                        <div class="font-semibold text-primary-900">{{ data.totalValue }}</div>
                      </div>
                    }
                    @if (data.user) {
                      <div>
                        <div class="text-xs text-slate-500">User</div>
                        <div class="font-semibold text-primary-900">{{ data.user }}</div>
                      </div>
                    }
                    @if (data.notes) {
                      <div class="col-span-2">
                        <div class="text-xs text-slate-500">Catatan</div>
                        <div class="font-medium text-primary-900">{{ data.notes }}</div>
                      </div>
                    }
                  </div>
                </div>
              }

              <details>
                <summary class="cursor-pointer text-xs text-slate-500 hover:text-slate-700 select-none">
                  ▼ Lihat data mentah
                </summary>
                <div class="mt-2 p-3 bg-slate-50 rounded-lg text-xs font-mono break-all text-slate-600 max-h-40 overflow-y-auto">
                  {{ scanResult() }}
                </div>
              </details>
            </div>
          } @else {
            <div class="flex flex-col items-center justify-center py-20 text-slate-400">
              <div class="text-5xl mb-3">📷</div>
              <div class="text-sm">Menunggu scan...</div>
            </div>
          }

          @if (history().length > 0) {
            <div class="mt-6 pt-4 border-t border-slate-100">
              <div class="flex items-center justify-between mb-2">
                <div class="text-xs font-semibold text-slate-500 uppercase">
                  Riwayat ({{ history().length }})
                </div>
                <button (click)="clearHistory()" class="text-xs text-slate-400 hover:text-red-500 transition">
                  Hapus
                </button>
              </div>
              <div class="space-y-1 max-h-40 overflow-y-auto">
                @for (h of history(); track $index) {
                  <div class="flex items-center gap-2 p-2 rounded-lg bg-slate-50 text-xs">
                    <span [class.text-green-600]="h.valid" [class.text-red-600]="!h.valid">
                      {{ h.valid ? '✅' : '❌' }}
                    </span>
                    <span class="flex-1 truncate text-slate-600">{{ h.code }}</span>
                    <span class="text-slate-400">{{ h.time }}</span>
                  </div>
                }
              </div>
            </div>
          }
        </div>
      </div>
    </div>
  `,
})
export class QrScannerComponent implements OnInit, OnDestroy {
  allowedFormats = [BarcodeFormat.QR_CODE];

  mode = signal<'camera' | 'upload' | 'manual'>('camera');
  cameraReady = signal(false);
  hasPermission = signal(true);

  scanResult = signal<string>('');
  parsedData = signal<ParsedQRData | null>(null);
  validationStatus = signal<'idle' | 'checking' | 'valid' | 'invalid'>('idle');
  validationMessage = signal<string>('');

  history = signal<{ code: string; valid: boolean; time: string }[]>([]);

  manualCode = '';
  uploadedPreview = signal<string>('');
  decodeStatus = signal<string>('');
  decodeOk = signal(false);

  private lastScanned = '';
  private scanTimeout?: any;

  constructor(
    private api: ApiService,
    private toast: ToastService,
  ) {}

  ngOnInit() {
    this.suppressZXingNoise();
  }

  ngOnDestroy() {
    if (this.scanTimeout) clearTimeout(this.scanTimeout);
  }

  // ============================================================
  // SUPPRESS ZXING NOISE
  // ============================================================
  private suppressZXingNoise() {
    const keywords = [
      'MultiFormatReader',
      'non-ReaderException',
      'ChecksumException',
      'NotFoundException',
      'FormatException',
      'Dimensions could be not found',
    ];

    const isNoisy = (args: any[]) => {
      const first = args[0];
      if (!first) return false;
      const msg = typeof first === 'string' ? first : (first.message ?? String(first));
      return keywords.some((k) => msg.includes(k));
    };

    const origError = console.error;
    console.error = (...args: any[]) => {
      if (isNoisy(args)) return;
      origError.apply(console, args);
    };

    const origWarn = console.warn;
    console.warn = (...args: any[]) => {
      if (isNoisy(args)) return;
      origWarn.apply(console, args);
    };
  }

  // ============================================================
  // MODE SWITCH
  // ============================================================
  switchMode(m: 'camera' | 'upload' | 'manual') {
    this.mode.set(m);
    this.clearResult();
    this.manualCode = '';
    this.uploadedPreview.set('');
    this.decodeStatus.set('');
    this.decodeOk.set(false);
  }

  // ============================================================
  // CAMERA EVENTS
  // ============================================================
  onCamerasFound(cameras: MediaDeviceInfo[]) {
    console.log('[QR Scanner] Cameras found:', cameras.length);
    this.cameraReady.set(cameras.length > 0);
  }

  onPermission(response: boolean) {
    this.hasPermission.set(response);
    if (!response) this.toast.warning('Izin kamera ditolak');
  }

  resetScanner() {
    this.clearResult();
    this.lastScanned = '';
    this.toast.info('Scanner direset');
  }

  // ============================================================
  // ✅ UPLOAD GAMBAR
  // ============================================================
  async onFileSelected(event: Event) {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;

    this.decodeStatus.set('Memproses gambar...');
    this.decodeOk.set(false);

    // Preview
    const reader = new FileReader();
    reader.onload = (e) => {
      this.uploadedPreview.set(e.target?.result as string);
    };
    reader.readAsDataURL(file);

    try {
      const decodedText = await this.decodeImageFile(file);

      if (!decodedText) {
        throw new Error('QR tidak terdeteksi');
      }

      this.decodeStatus.set('✅ QR berhasil dibaca');
      this.decodeOk.set(true);
      this.scanResult.set(decodedText);
      this.parsedData.set(this.parseQRData(decodedText));

      const code = this.extractCode(decodedText);
      this.validateCode(code);
    } catch (err: any) {
      console.error('[QR Scanner] Image decode error:', err);
      this.decodeStatus.set('❌ ' + (err?.message || 'QR tidak terdeteksi'));
      this.decodeOk.set(false);
      this.toast.error('Tidak dapat membaca QR dari gambar');
    }
  }

  // ============================================================
  // ✅ DECODE IMAGE FILE — pakai decodeFromImageUrl + canvas.toDataURL
  // (decodeFromCanvas tidak ada di type definition)
  // ============================================================
  private async decodeImageFile(file: File): Promise<string> {
    const codeReader = new BrowserMultiFormatReader();
    const img = await this.loadImage(file);

    const attempts = [
      { scale: 1, rotate: 0 },
      { scale: 1.5, rotate: 0 },
      { scale: 2, rotate: 0 },
      { scale: 0.75, rotate: 0 },
      { scale: 1, rotate: 90 },
      { scale: 1, rotate: 180 },
      { scale: 1, rotate: 270 },
    ];

    for (const attempt of attempts) {
      try {
        const canvas = this.createCanvas(img, attempt.scale, attempt.rotate);

        // ✅ Pakai decodeFromImageUrl dengan Data URL
        const dataUrl = canvas.toDataURL('image/png');
        const result = await codeReader.decodeFromImageUrl(dataUrl);
        const text = result.getText();

        console.log(
          `[QR Scanner] ✅ Decoded: scale=${attempt.scale}, rotate=${attempt.rotate}°`,
          text,
        );
        return text;
      } catch (e) {
        console.log(
          `[QR Scanner] ❌ Failed: scale=${attempt.scale}, rotate=${attempt.rotate}°`,
        );
      }
    }

    return '';
  }

  private createCanvas(
    img: HTMLImageElement,
    scale: number,
    rotateDeg: number,
  ): HTMLCanvasElement {
    const canvas = document.createElement('canvas');
    const rotated = rotateDeg === 90 || rotateDeg === 270;
    const w = rotated ? img.height : img.width;
    const h = rotated ? img.width : img.height;

    canvas.width = Math.round(w * scale);
    canvas.height = Math.round(h * scale);

    const ctx = canvas.getContext('2d')!;
    ctx.fillStyle = '#ffffff';
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    ctx.translate(canvas.width / 2, canvas.height / 2);
    ctx.rotate((rotateDeg * Math.PI) / 180);
    ctx.scale(scale, scale);
    ctx.drawImage(img, -img.width / 2, -img.height / 2);

    return canvas;
  }

  private loadImage(file: File): Promise<HTMLImageElement> {
    return new Promise((resolve, reject) => {
      const img = new Image();
      const url = URL.createObjectURL(file);
      img.onload = () => {
        URL.revokeObjectURL(url);
        resolve(img);
      };
      img.onerror = () => {
        URL.revokeObjectURL(url);
        reject(new Error('Gagal load image'));
      };
      img.src = url;
    });
  }

  // ============================================================
  // MANUAL INPUT
  // ============================================================
  submitManual() {
    const code = this.manualCode.trim();
    if (!code) return;

    this.scanResult.set(code);
    this.parsedData.set(null);
    this.validationStatus.set('checking');
    this.validationMessage.set('');

    const cleaned = this.extractCode(code);
    this.validateCode(cleaned);
  }

  // ============================================================
  // EXTRACT CODE
  // ============================================================
  private extractCode(raw: string): string {
    console.log('[QR Scanner] Raw input:', raw);

    let code = raw.trim();

    // 1. Try parse JSON
    try {
      const parsed = JSON.parse(raw);
      const candidates = [
        parsed.c, parsed.documentCode,
        parsed.transactionCode, parsed.transaction_code,
        parsed.document_code, parsed.transactionNo,
        parsed.transaction_no, parsed.code,
        parsed.kode, parsed.docNo, parsed.doc_no,
      ];
      for (const c of candidates) {
        if (c && typeof c === 'string' && c.trim() !== '') {
          code = c.trim();
          console.log('[QR Scanner] ✅ Code dari JSON:', code);
          break;
        }
      }
    } catch {
      // Bukan JSON
    }

    // 2. Pipe-delimited format: TRX-xxx|status|type|...
    if (code.includes('|')) {
      code = code.split('|')[0].trim();
      console.log('[QR Scanner] ✅ Code dari pipe format:', code);
    }

    // 3. Auto-clean: hapus prefix & extension
    code = code
      .replace(/^qr_/i, '')
      .replace(/^qr-/i, '')
      .replace(/\.(png|jpg|jpeg|gif|svg|webp|bmp)$/i, '')
      .trim();

    // 4. Regex: extract TRX-/SIG-/DOC- pattern
    const match = code.match(/(TRX|SIG|DOC)-\d{8}-\d{4}/i);
    if (match) {
      console.log('[QR Scanner] ✅ Final code (regex):', match[0]);
      return match[0].toUpperCase();
    }

    console.log('[QR Scanner] ✅ Final code:', code);
    return code;
  }

  // ============================================================
  // SCAN HANDLER (dari kamera)
  // ============================================================
  onScan(result: string) {
    if (!result || result === this.lastScanned) return;
    this.lastScanned = result;

    console.log('[QR Scanner] ✅ QR dari kamera:', result);

    if ('vibrate' in navigator) navigator.vibrate(50);

    this.scanResult.set(result);
    this.parsedData.set(this.parseQRData(result));
    this.validationStatus.set('checking');
    this.validationMessage.set('');

    const code = this.extractCode(result);
    this.validateCode(code);

    if (this.scanTimeout) clearTimeout(this.scanTimeout);
    this.scanTimeout = setTimeout(() => {
      this.lastScanned = '';
    }, 3000);
  }

  // ============================================================
  // VALIDATE CODE
  // ============================================================
  private validateCode(code: string) {
    console.log('[QR Scanner] Validating:', code);

    this.api.get<any>('documents/validate', { code }).subscribe({
      next: (res: any) => {
        console.log('[QR Scanner] Response:', res);

        const backendDoc = res?.data?.document;
        const isValid = res?.data?.valid === true;

        if (backendDoc) {
          this.parsedData.set({
            documentCode: backendDoc.documentCode || backendDoc.transactionNo || code,
            transactionCode: backendDoc.transactionCode,
            type: backendDoc.type,
            status: backendDoc.status,
            totalItems: backendDoc.totalItems,
            totalValue: backendDoc.totalValue,
            user: backendDoc.user,
            notes: backendDoc.notes,
          });
        }

        if (isValid) {
          this.validationStatus.set('valid');
          this.validationMessage.set('Dokumen valid');
          this.toast.success('Dokumen valid ✅');
          this.history.update((h) => [
            { code, valid: true, time: this.getTime() },
            ...h,
          ].slice(0, 10));
        } else {
          this.validationStatus.set('invalid');
          this.validationMessage.set(
            res?.data?.message || res?.message || 'Dokumen tidak terdaftar',
          );
          this.toast.warning('Dokumen tidak valid');
          this.history.update((h) => [
            { code, valid: false, time: this.getTime() },
            ...h,
          ].slice(0, 10));
        }
      },
      error: (err) => {
        console.error('[QR Scanner] Validation error:', err);
        this.validationStatus.set('invalid');
        this.validationMessage.set(
          err?.error?.message || 'Gagal validasi ke server',
        );
        this.toast.error('Gagal validasi dokumen');
        this.history.update((h) => [
          { code, valid: false, time: this.getTime() },
          ...h,
        ].slice(0, 10));
      },
    });
  }

  // ============================================================
  // PARSE QR DATA — support JSON & pipe-delimited
  // ============================================================
  private parseQRData(raw: string): ParsedQRData | null {
    try {
      const p = JSON.parse(raw);
      return {
        documentCode: p.c || p.documentCode || p.document_code,
        transactionCode: p.transactionCode || p.transaction_code,
        type: p.t || p.type,
        status: p.s || p.status,
        totalItems: p.i ?? p.totalItems ?? p.total_items,
        totalValue: p.v ?? p.totalValue ?? p.total_value,
        transactionDate: p.d ?? p.transactionDate ?? p.transaction_date,
        user: p.u ?? p.user,
        notes: p.n ?? p.notes,
        docNo: p.docNo,
      };
    } catch {
      // Bukan JSON
    }

    if (raw.includes('|')) {
      const parts = raw.split('|');
      return {
        documentCode: parts[0]?.trim(),
        status: parts[1]?.trim(),
        type: parts[2]?.trim(),
        totalItems: parts[3] ? Number(parts[3]) : undefined,
        totalValue: parts[4] ? Number(parts[4]) : undefined,
        user: parts[5]?.trim(),
        notes: parts[6]?.trim(),
      };
    }

    return null;
  }

  private getTime(): string {
    const d = new Date();
    return `${d.getHours().toString().padStart(2, '0')}:${d
      .getMinutes()
      .toString()
      .padStart(2, '0')}`;
  }

  clearResult() {
    this.scanResult.set('');
    this.parsedData.set(null);
    this.validationStatus.set('idle');
    this.validationMessage.set('');
    this.decodeStatus.set('');
    this.decodeOk.set(false);
  }

  clearHistory() {
    this.history.set([]);
  }

  // ============================================================
  // STATUS HELPERS
  // ============================================================
  currentStatus(): string {
    return String(this.parsedData()?.status ?? '').toLowerCase();
  }

  statusBadgeClass(status: string): string {
    const map: Record<string, string> = {
      draft: 'bg-slate-50 border-slate-300',
      pending: 'bg-amber-50 border-amber-200',
      approved: 'bg-green-50 border-green-200',
      rejected: 'bg-red-50 border-red-200',
      completed: 'bg-blue-50 border-blue-200',
      cancelled: 'bg-gray-100 border-gray-300',
    };
    return map[(status ?? '').toLowerCase()] ?? 'bg-green-50 border-green-200';
  }

  statusTextClass(status: string): string {
    const map: Record<string, string> = {
      draft: 'text-slate-700',
      pending: 'text-amber-700',
      approved: 'text-green-700',
      rejected: 'text-red-700',
      completed: 'text-blue-700',
      cancelled: 'text-gray-600',
    };
    return map[(status ?? '').toLowerCase()] ?? 'text-green-700';
  }

  statusSubTextClass(status: string): string {
    const map: Record<string, string> = {
      draft: 'text-slate-500',
      pending: 'text-amber-600',
      approved: 'text-green-600',
      rejected: 'text-red-600',
      completed: 'text-blue-600',
      cancelled: 'text-gray-500',
    };
    return map[(status ?? '').toLowerCase()] ?? 'text-green-600';
  }

  statusIcon(status: string): string {
    const map: Record<string, string> = {
      draft: '📝',
      pending: '⏳',
      approved: '✅',
      rejected: '❌',
      completed: '🏁',
      cancelled: '⛔',
    };
    return map[(status ?? '').toLowerCase()] ?? '✅';
  }

  statusLabel(status: string): string {
    const map: Record<string, string> = {
      draft: 'DOKUMEN VALID — DRAFT',
      pending: 'DOKUMEN VALID — PENDING',
      approved: 'DOKUMEN VALID — APPROVED',
      rejected: 'DOKUMEN VALID — REJECTED',
      completed: 'DOKUMEN VALID — COMPLETED',
      cancelled: 'DOKUMEN VALID — CANCELLED',
    };
    return map[(status ?? '').toLowerCase()] ?? 'DOKUMEN VALID';
  }
}