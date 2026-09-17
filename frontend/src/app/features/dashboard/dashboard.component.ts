import { Component, OnInit, signal, computed, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { forkJoin } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import { AuthService } from '../../core/services/auth.service';
import { ToastService } from '../../core/services/toast.service';

interface StatCard {
  label: string;
  value: number;
  color: string;
  bg: string;
  icon: string;
  path: string;
  key: string;
}

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, RouterLink],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="space-y-6 fade-in">
      <!-- Header -->
      <div class="flex items-center justify-between flex-wrap gap-3">
        <div>
          <h2 class="text-2xl font-bold text-primary-900">Dashboard</h2>
          <p class="text-sm text-slate-500">Ringkasan sistem informasi ATK</p>
        </div>
        <div class="text-sm text-slate-500">
          {{ today | date: 'EEEE, dd MMMM yyyy' }}
        </div>
      </div>

      <!-- Stat Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        @for (s of filteredStats(); track s.label) {
          <a
            [routerLink]="s.path"
            class="bg-white rounded-2xl p-5 border border-primary-100 shadow-sm hover:shadow-md hover:-translate-y-0.5 transition-all cursor-pointer"
          >
            <div class="flex items-center justify-between">
              <div>
                <div class="text-xs font-medium text-slate-500 uppercase tracking-wide">
                  {{ s.label }}
                </div>
                <div class="text-3xl font-bold text-primary-900 mt-2">{{ s.value | number }}</div>
              </div>
              <div
                class="w-12 h-12 rounded-xl flex items-center justify-center text-white text-xl"
                [ngClass]="s.color"
              >
                <span [innerHTML]="s.icon"></span>
              </div>
            </div>
            <div
              class="mt-3 pt-3 border-t border-primary-50 flex items-center justify-between text-xs"
            >
              <span class="text-slate-400">Lihat detail</span>
              <span class="text-primary-600 font-medium">→</span>
            </div>
          </a>
        }
      </div>

      <!-- Charts Row -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Trend Chart -->
        <div class="lg:col-span-2 bg-white rounded-2xl p-6 border border-primary-100 shadow-sm">
          <div class="flex items-center justify-between mb-4">
            <h3 class="font-semibold text-primary-900">Trend Transaksi 7 Hari Terakhir</h3>
            <span class="text-xs text-slate-400">IN vs OUT</span>
          </div>
          @if (trendData().length > 0) {
            <div class="flex items-end gap-2 h-48">
              @for (d of trendData(); track d.date) {
                <div class="flex-1 flex flex-col items-center gap-1">
                  <div class="w-full flex flex-col justify-end h-40 gap-0.5">
                    <div
                      class="bg-primary-500 rounded-t transition-all"
                      [style.height.%]="(d.in / maxTrend()) * 100"
                      [title]="'IN: ' + d.in"
                    ></div>
                    <div
                      class="bg-blue-300 rounded-b transition-all"
                      [style.height.%]="(d.out / maxTrend()) * 100"
                      [title]="'OUT: ' + d.out"
                    ></div>
                  </div>
                  <div class="text-[10px] text-slate-500 whitespace-nowrap">{{ d.date }}</div>
                </div>
              }
            </div>
            <div class="flex items-center justify-center gap-4 mt-4 text-xs">
              <div class="flex items-center gap-1">
                <span class="w-3 h-3 bg-primary-500 rounded"></span> Masuk
              </div>
              <div class="flex items-center gap-1">
                <span class="w-3 h-3 bg-blue-300 rounded"></span> Keluar
              </div>
            </div>
          } @else {
            <div class="h-48 flex items-center justify-center text-slate-400 text-sm">
              Belum ada data
            </div>
          }
        </div>

        <!-- Category Chart -->
        <div class="bg-white rounded-2xl p-6 border border-primary-100 shadow-sm">
          <h3 class="font-semibold text-primary-900 mb-4">Stock per Kategori</h3>
          @if (categoryData().length > 0) {
            <div class="space-y-3">
              @for (c of categoryData(); track c.category) {
                <div>
                  <div class="flex items-center justify-between text-sm mb-1">
                    <span class="text-slate-700">{{ c.category }}</span>
                    <span class="font-semibold text-primary-700">{{ c.total | number }}</span>
                  </div>
                  <div class="h-2 bg-primary-50 rounded-full overflow-hidden">
                    <div
                      class="h-full bg-gradient-to-r from-primary-500 to-primary-400 transition-all duration-500"
                      [style.width.%]="(c.total / maxCategory()) * 100"
                    ></div>
                  </div>
                </div>
              }
            </div>
          } @else {
            <div class="h-48 flex items-center justify-center text-slate-400 text-sm">
              Belum ada data
            </div>
          }
        </div>
      </div>

      <!-- Lists Row -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Recent Transactions -->
        <div class="bg-white rounded-2xl p-6 border border-primary-100 shadow-sm">
          <div class="flex items-center justify-between mb-4">
            <h3 class="font-semibold text-primary-900">Transaksi Terbaru</h3>
            <a routerLink="/transactions" class="text-xs text-primary-600 hover:underline"
              >Lihat semua →</a
            >
          </div>
          <div class="space-y-2">
            @for (t of recentTransactions(); track t.id) {
              <div
                class="flex items-center justify-between p-3 rounded-lg bg-primary-50/60 hover:bg-primary-50 transition"
              >
                <div class="flex items-center gap-3">
                  <div
                    class="w-9 h-9 rounded-lg flex items-center justify-center text-xs font-bold"
                    [ngClass]="
                      t.type === 'IN' ? 'bg-green-100 text-green-700' : 'bg-blue-100 text-blue-700'
                    "
                  >
                    {{ t.type }}
                  </div>
                  <div>
                    <div class="text-sm font-medium text-primary-900">{{ t.transactionCode || t.transactionNo }}</div>
                    <div class="text-xs text-slate-500">
                      {{ t.items?.[0]?.master?.name || t.itemName || '-' }} • Qty {{ t.totalItems || t.quantity || 0 }}
                    </div>
                  </div>
                </div>
                <span
                  class="px-2 py-1 text-xs rounded-full font-medium"
                  [ngClass]="statusClass(t.status)"
                >
                  {{ t.status }}
                </span>
              </div>
            } @empty {
              <div class="text-center text-slate-400 py-8 text-sm">Belum ada transaksi</div>
            }
          </div>
        </div>

        <!-- Stock Lowest (Ascending paling rendah) -->
        <div class="bg-white rounded-2xl p-6 border border-primary-100 shadow-sm">
          <div class="flex items-center justify-between mb-4">
            <h3 class="font-semibold text-primary-900">Stock Lowest</h3>
            <div class="flex items-center gap-2">
              <span class="px-2 py-0.5 text-xs rounded-full bg-amber-100 text-amber-700 font-medium">
                {{ lowStock().length }} item
              </span>
              <a routerLink="/stock" class="text-xs text-primary-600 hover:underline">Lihat semua →</a>
            </div>
          </div>
          @if (lowStock().length === 0) {
            <div class="text-sm text-slate-400 mt-2">Tidak ada data stock</div>
          } @else {
            <div class="mt-3 space-y-2 max-h-44 overflow-y-auto pr-1">
              @for (item of lowStock(); track item.id) {
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
      </div>

      <!-- ============================================================ -->
      <!-- ✅ PENDING APPROVALS — HANYA untuk super_admin & manager    -->
      <!-- ============================================================ -->
      @if (canViewApprovals() && pendingApprovals().length > 0) {
        <div class="bg-amber-50 rounded-2xl p-6 border border-amber-200 shadow-sm">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <div
                class="w-8 h-8 rounded-lg bg-amber-500 text-white flex items-center justify-center"
              >
                ⏳
              </div>
              <h3 class="font-semibold text-amber-900">
                Menunggu Approval ({{ pendingApprovals().length }})
              </h3>
            </div>
            <a routerLink="/approvals" class="text-xs text-amber-700 hover:underline"
              >Proses sekarang →</a
            >
          </div>
          <div class="space-y-2">
            @for (a of pendingApprovals(); track a.id) {
              <div class="flex items-center justify-between p-3 rounded-lg bg-white">
                <div>
                  <div class="text-sm font-medium text-primary-900">{{ a.transactionCode || a.transactionNo }}</div>
                  <div class="text-xs text-slate-500">Diminta oleh {{ a.user?.full_name || a.user?.username || a.requestedBy || 'User' }}</div>
                </div>
              </div>
            }
          </div>
        </div>
      }
    </div>
  `,
})
export class DashboardComponent implements OnInit {
  today = new Date();

  allStats = signal<StatCard[]>([]);
  recentTransactions = signal<any[]>([]);
  lowStock = signal<any[]>([]);
  pendingApprovals = signal<any[]>([]);
  trendData = signal<{ date: string; in: number; out: number }[]>([]);
  categoryData = signal<{ category: string; total: number }[]>([]);
  filteredStats = signal<StatCard[]>([]);

  userRole = signal<string>('');

  // ✅ COMPUTED: Hanya super_admin & manager yang boleh lihat Pending Approval
  canViewApprovals = computed(() => {
    const role = this.userRole();
    return role === 'super_admin' || role === 'manager';
  });

  constructor(
    private api: ApiService,
    public auth: AuthService,
    private toast: ToastService,
  ) {}

  ngOnInit() {
    this.resolveUserRole();
    this.loadDashboardData();
  }

  // ============================================================
  // ✅ RESOLVE USER ROLE
  // ============================================================
  private resolveUserRole() {
    let role = this.extractRole(this.auth.currentUser());
    if (!role) {
      try {
        const stored = localStorage.getItem('atk_user');
        if (stored) role = this.extractRole(JSON.parse(stored));
      } catch (e) {}
    }
    role = String(role).toLowerCase().trim().replace(/\s+/g, '_');
    if (role === 'superadmin') role = 'super_admin';
    this.userRole.set(role);
    console.log('🔍 [Dashboard] Role user:', role);
  }

  private extractRole(user: any): string {
    if (!user) return '';
    const keys = ['roleName', 'role_name', 'userRole', 'user_role', 'role', 'Role', 'roles'];
    for (const key of keys) {
      const val = user[key];
      if (val === undefined || val === null) continue;
      if (typeof val === 'string' && val.trim() !== '') return val;
      if (typeof val === 'object' && !Array.isArray(val)) {
        const nested = val.name || val.title || val.roleName;
        if (nested && typeof nested === 'string') return nested;
      }
      if (Array.isArray(val) && val.length > 0) {
        const first = val[0];
        if (typeof first === 'string') return first;
        if (typeof first === 'object' && first.name) return first.name;
      }
    }
    return '';
  }

  // ============================================================
  // ✅ FILTER STAT CARDS
  // ============================================================
  private updateFilteredStats() {
    const role = this.userRole();
    const all = this.allStats();

    const allowedKeys: Record<string, string[]> = {
      super_admin: ['totalUser', 'itemStock', 'totalQty', 'monthTransactions', 'pendingApproval'],
      admin:       ['totalUser', 'itemStock', 'totalQty', 'monthTransactions'],
      manager:     ['itemStock', 'totalQty', 'monthTransactions', 'pendingApproval'],
      staff:       ['itemStock', 'totalQty', 'monthTransactions'],
      viewer:      ['monthTransactions'],
    };

    const allowed = allowedKeys[role] ?? allowedKeys['viewer'];
    const filtered = all.filter((s) => allowed.includes(s.key));
    this.filteredStats.set(filtered);
  }

  // ============================================================
  // ✅ LOAD DASHBOARD DATA
  // ============================================================
  loadDashboardData() {
    forkJoin({
      users: this.api.get<any>('users', { size: 1000 }),
      stocks: this.api.get<any>('stock', { size: 1000 }),
      transactions: this.api.get<any>('transactions', { size: 1000, sort: 'createdAt,desc' }),
    }).subscribe({
      next: (res) => {
        const extract = (r: any): any[] => {
          if (!r) return [];
          const root = r?.data ?? r;
          const inner = root?.data ?? root;
          if (Array.isArray(inner)) return inner;
          if (Array.isArray(root)) return root;
          if (Array.isArray(r)) return r;
          return [];
        };

        const users = extract(res.users);
        const stocks = extract(res.stocks);
        const transactions = extract(res.transactions);

        const totalUsers = users.length;
        const distinctItems = new Set(
          stocks.map((s) => String(s?.master?.code ?? s?.itemCode ?? ''))
        ).size;
        const totalQty = stocks.reduce((sum, s) => sum + Number(s?.quantity ?? 0), 0);

        const now = new Date();
        const monthTransactions = transactions.filter((t) => {
          const d = new Date(t.transactionDate || t.createdAt);
          return d.getMonth() === now.getMonth() && d.getFullYear() === now.getFullYear();
        }).length;

        const pendingList = transactions.filter((t) => {
          const st = String(t.status || '').toLowerCase();
          return st !== 'approved' && st !== 'rejected';
        });
        const pendingCount = pendingList.length;

        this.allStats.set([
          { key: 'totalUser', label: 'Total User', value: totalUsers, color: 'bg-primary-600', bg: 'bg-primary-50', icon: '👤', path: '/users' },
          { key: 'itemStock', label: 'Item Stock', value: distinctItems, color: 'bg-blue-500', bg: 'bg-blue-50', icon: '📦', path: '/stock' },
          { key: 'totalQty', label: 'Total Qty', value: totalQty, color: 'bg-indigo-500', bg: 'bg-indigo-50', icon: '🔢', path: '/stock' },
          { key: 'monthTransactions', label: 'Transaksi Bulan Ini', value: monthTransactions, color: 'bg-purple-500', bg: 'bg-purple-50', icon: '📊', path: '/transactions' },
          { key: 'pendingApproval', label: 'Pending Approval', value: pendingCount, color: 'bg-amber-500', bg: 'bg-amber-50', icon: '⏳', path: '/approvals' },
        ]);

        this.updateFilteredStats();
        this.recentTransactions.set(transactions.slice(0, 5));
        this.pendingApprovals.set(pendingList.slice(0, 5));

        const normalizedStocks = stocks.map((s) => ({
          id: s.id,
          quantity: Number(s.quantity ?? 0),
          itemName: s.master?.name || s.itemName || '-',
          itemCode: s.master?.code || s.itemCode || '-',
          category: s.master?.category || '-',
        }));

        const availableItems = normalizedStocks.filter((x) => x.quantity > 0);
        let lowestItems: any[] = [];
        if (availableItems.length > 0) {
          const lowestQty = Math.min(...availableItems.map((x) => x.quantity));
          lowestItems = availableItems
            .filter((x) => x.quantity === lowestQty)
            .map((x) => ({
              id: x.id,
              quantity: x.quantity,
              master: { code: x.itemCode, name: x.itemName, category: x.category }
            }));
        }
        this.lowStock.set(lowestItems.slice(0, 5));

        const catMap = new Map<string, number>();
        stocks.forEach((s) => {
          const cat = String(s?.master?.category ?? 'Lainnya');
          const qty = Number(s?.quantity ?? 0);
          catMap.set(cat, (catMap.get(cat) || 0) + qty);
        });
        const catArray = Array.from(catMap.entries())
          .map(([category, total]) => ({ category, total }))
          .sort((a, b) => b.total - a.total);
        this.categoryData.set(catArray);

        this.generateTrendData(transactions);
      },
      error: (err) => {
        console.error('[Dashboard] Error loading data:', err);
        this.toast.error('Gagal memuat data dashboard');
      },
    });
  }

  private generateTrendData(transactions: any[]) {
    const trend = [];
    const today = new Date();

    for (let i = 6; i >= 0; i--) {
      const d = new Date(today);
      d.setDate(today.getDate() - i);
      const dateStr = d.toISOString().slice(0, 10);

      const dayTx = transactions.filter((t) => {
        const txDate = String(t.transactionDate || t.createdAt || '').slice(0, 10);
        return txDate === dateStr;
      });

      const inCount = dayTx.filter((t) => String(t.type).toUpperCase() === 'IN').length;
      const outCount = dayTx.filter((t) => String(t.type).toUpperCase() === 'OUT').length;
      const displayDate = d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' });

      trend.push({ date: displayDate, in: inCount, out: outCount });
    }
    this.trendData.set(trend);
  }

  maxTrend(): number {
    const data = this.trendData();
    if (!data.length) return 1;
    return Math.max(...data.map((d) => Math.max(d.in, d.out)), 1);
  }

  maxCategory(): number {
    const data = this.categoryData();
    if (!data.length) return 1;
    return Math.max(...data.map((c) => c.total), 1);
  }

  statusClass(status: string): string {
    const map: Record<string, string> = {
      APPROVED: 'bg-green-100 text-green-700',
      PENDING: 'bg-amber-100 text-amber-700',
      REJECTED: 'bg-red-100 text-red-700',
      DRAFT: 'bg-slate-100 text-slate-700',
    };
    return map[String(status).toUpperCase()] ?? 'bg-slate-100 text-slate-700';
  }
}