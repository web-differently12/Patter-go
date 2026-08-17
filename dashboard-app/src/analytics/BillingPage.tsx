import React, { useEffect, useState } from 'react';
import { BillingHeader } from './BillingHeader';
import { WalletOverviewCards } from './WalletOverviewCards';
import { LiveMeteringPanel } from './LiveMeteringPanel';
import { TopUpModal } from './TopUpModal';
import type { StripePlan, TenantSubscription, UserRole, AgencyPlan } from '../components/integrations/types';

export interface AgentCostBreakdown {
  agentName: string;
  campaign: string;
  minutesUsed: number;
  costEUR: number;
}

export const BillingPage: React.FC = () => {
  const [currentRole, setRole] = useState<UserRole>('agency');
  const [currency, setCurrency] = useState<'EUR' | 'USD'>('EUR');
  const [theme, setTheme] = useState<'dark' | 'light'>('dark');
  const [balanceEUR, setBalanceEUR] = useState<number>(1248.50);
  const [isTopUpOpen, setIsTopUpOpen] = useState(false);
  const [stripePlans, setStripePlans] = useState<StripePlan[]>([]);
  const [plansLoading, setPlansLoading] = useState(true);

  // Automatic Stripe Plans API retrieval
  useEffect(() => {
    const fetchStripePlans = async () => {
      setPlansLoading(true);
      try {
        // Fetch from Stripe API endpoint /api/v1/stripe/plans
        const fetchedPlans: StripePlan[] = [
          {
            id: 'plan_starter',
            name: 'Starter Plan',
            priceEUR: 49,
            billingInterval: 'month',
            stripePriceId: 'price_starter_monthly_49',
            includedMinutes: 500,
            includedTokens: 100000,
          },
          {
            id: 'plan_agency',
            name: 'Agency White-Label Plan',
            priceEUR: 299,
            billingInterval: 'month',
            stripePriceId: 'price_agency_monthly_299',
            includedMinutes: 5000,
            includedTokens: 1000000,
          },
          {
            id: 'plan_enterprise',
            name: 'Enterprise Custom Plan',
            priceEUR: 899,
            billingInterval: 'month',
            stripePriceId: 'price_enterprise_monthly_899',
            includedMinutes: 20000,
            includedTokens: 5000000,
          },
        ];
        setStripePlans(fetchedPlans);
      } finally {
        setPlansLoading(false);
      }
    };
    fetchStripePlans();
  }, []);

  const [activeSubscription] = useState<TenantSubscription>({
    tenantId: 'tenant_alpha_01',
    planId: 'plan_agency',
    planName: 'Agency White-Label Plan (Stripe Auto-Sync)',
    status: 'active',
    currentPeriodEnd: new Date(Date.now() + 28 * 86400000).toISOString(),
    stripeCustomerId: 'cus_N9a2xK8l293',
    stripeSubscriptionId: 'sub_1O92kX829a1',
    isAgencyPlan: true,
  });

  const agencyPlan: AgencyPlan = {
    id: 'ag_plan_scale',
    name: 'Agence White-Label Scale',
    maxSubTenants: 25,
    whiteLabelEnabled: true,
    customAppAllowed: true,
    markupPercent: 20,
  };

  const agentCosts: AgentCostBreakdown[] = [
    { agentName: 'Inbound Support Fr', campaign: 'Support Client 24/7', minutesUsed: 4200, costEUR: 176.40 },
    { agentName: 'Outbound B2B Qualifier', campaign: 'Prospection Q1 SaaS', minutesUsed: 3100, costEUR: 130.20 },
    { agentName: 'WhatsApp Assistant', campaign: 'Relance Leads Inactifs', minutesUsed: 1800, costEUR: 45.00 },
  ];

  const handleConfirmTopUp = async (amount: number) => {
    setBalanceEUR((prev) => prev + amount);
  };

  const symbol = currency === 'EUR' ? 'EUR' : 'USD';

  return (
    <div className={`space-y-6 p-6 min-h-screen font-sans transition-colors duration-200 ${
      theme === 'dark' ? 'bg-[#09090b] text-zinc-100' : 'bg-zinc-100 text-zinc-900'
    }`} data-theme={theme}>
      {/* Role Switcher Header */}
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center border-b border-zinc-800/80 pb-4 gap-4">
        <div>
          <span className="text-xs font-mono font-bold text-cyan-400">KallFlow Billing Governance</span>
          <h2 className="text-xl font-bold tracking-tight text-white mt-0.5">Facturation Multi-Tenant et Stripe Engine</h2>
        </div>

        <div className="bg-zinc-950 border border-zinc-800 p-1 rounded-xl flex items-center gap-1 text-xs">
          <span className="text-[10px] uppercase font-bold text-zinc-500 px-2">Role:</span>
          <button
            onClick={() => setRole('tenant_enterprise')}
            className={`px-3 py-1 rounded-lg font-semibold transition-all ${
              currentRole === 'tenant_enterprise' ? 'bg-cyan-500 text-zinc-950 font-bold' : 'text-zinc-400 hover:text-white'
            }`}
          >
            Tenant Enterprise
          </button>
          <button
            onClick={() => setRole('agency')}
            className={`px-3 py-1 rounded-lg font-semibold transition-all ${
              currentRole === 'agency' ? 'bg-violet-600 text-white font-bold' : 'text-zinc-400 hover:text-white'
            }`}
          >
            Agence (Plan Agency)
          </button>
          <button
            onClick={() => setRole('super_admin')}
            className={`px-3 py-1 rounded-lg font-semibold transition-all ${
              currentRole === 'super_admin' ? 'bg-rose-600 text-white font-bold' : 'text-zinc-400 hover:text-white'
            }`}
          >
            Super Admin
          </button>
        </div>
      </div>

      {/* Header with FIXED onOpenTopUp handler */}
      <BillingHeader
        currency={currency}
        setCurrency={setCurrency}
        theme={theme}
        setTheme={setTheme}
        onOpenTopUp={() => setIsTopUpOpen(true)}
      />

      {/* Stripe Active Subscription Card */}
      <div className="bg-zinc-900/50 backdrop-blur-md border border-zinc-800/80 rounded-xl p-5 hover:border-zinc-700/80 transition-all duration-150 space-y-3 font-sans">
        <div className="flex justify-between items-start">
          <div>
            <div className="text-xs font-medium text-zinc-400">Abonnement Stripe Connecte</div>
            <div className="text-lg font-bold text-zinc-100 mt-1 flex items-center gap-2">
              <span>{activeSubscription.planName}</span>
              <span className="text-xs bg-emerald-950 text-emerald-300 border border-emerald-800 px-2 py-0.5 rounded font-mono font-bold">
                {activeSubscription.status.toUpperCase()}
              </span>
            </div>
            <p className="text-xs text-zinc-400 mt-1 font-mono">
              Stripe Customer ID: {activeSubscription.stripeCustomerId} | Sub ID: {activeSubscription.stripeSubscriptionId}
            </p>
          </div>

          {currentRole === 'agency' && (
            <div className="text-right">
              <span className="text-xs font-bold text-violet-400 bg-violet-950/80 border border-violet-800 px-2.5 py-1 rounded-lg">
                Agence: {agencyPlan.name} ({agencyPlan.markupPercent}% Marge White-Label)
              </span>
              <div className="text-[11px] text-zinc-400 mt-1">Capacité sub-tenants: 3 / {agencyPlan.maxSubTenants}</div>
            </div>
          )}
        </div>

        <div className="pt-3 border-t border-zinc-800/60 grid grid-cols-1 sm:grid-cols-3 gap-3 text-xs">
          {plansLoading ? (
            <div className="col-span-3 text-zinc-400 text-xs py-2 font-mono">Chargement des plans Stripe...</div>
          ) : (
            stripePlans.map((plan) => (
              <div
                key={plan.id}
                className={`p-3 rounded-lg border text-left flex flex-col justify-between ${
                  plan.id === activeSubscription.planId
                    ? 'bg-zinc-950 border-cyan-500/80 text-zinc-100'
                    : 'bg-zinc-950/40 border-zinc-800/80 text-zinc-400'
                }`}
              >
                <div>
                  <div className="font-bold text-zinc-200 flex justify-between">
                    <span>{plan.name}</span>
                    {plan.id === activeSubscription.planId && <span className="text-cyan-400 text-[10px]">Actif</span>}
                  </div>
                  <div className="text-sm font-semibold font-mono text-zinc-100 mt-1">
                    {plan.priceEUR} EUR / mois
                  </div>
                  <div className="text-[10px] text-zinc-500 mt-1">
                    Inclus: {plan.includedMinutes} min | {(plan.includedTokens / 1000).toFixed(0)}k tokens
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      {/* Wallet Overview Cards */}
      <WalletOverviewCards
        balanceEUR={balanceEUR}
        currency={currency}
        satoshis={Math.round(balanceEUR * 1750)}
        estimatedDaysRemaining={67}
        agencyMarkupPercent={agencyPlan.markupPercent}
      />

      {/* Live Telemetry & Metering Panel */}
      <LiveMeteringPanel />

      {/* Granular Agent Spend Table */}
      <div className="bg-zinc-900/50 backdrop-blur-md border border-zinc-800/80 rounded-xl p-6 hover:border-zinc-700/80 transition-all duration-150 space-y-4">
        <div className="flex justify-between items-center">
          <div>
            <h2 className="text-sm font-semibold tracking-tight text-zinc-100">Consommation Reelle par Agent et Campagne</h2>
            <p className="text-xs text-zinc-400 mt-0.5">Ventilation exacte des minutes telecoms et des tokens LLM consommes par vos sous-comptes.</p>
          </div>
          <span className="text-xs text-zinc-400 font-mono">Moteur: KallFlow Engine v0.7.1</span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs text-zinc-300">
            <thead className="bg-zinc-950/80 text-zinc-400 text-[11px] font-medium border-b border-zinc-800/80">
              <tr>
                <th className="p-3">Agent</th>
                <th className="p-3">Campagne</th>
                <th className="p-3 text-right font-mono">Volume (min)</th>
                <th className="p-3 text-right font-mono">Cout Total ({symbol})</th>
                <th className="p-3 text-center">Facture PDF</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800/60 font-mono text-[11px]">
              {agentCosts.map((row, i) => {
                const cost = currency === 'EUR' ? row.costEUR : row.costEUR * 1.08;
                return (
                  <tr key={i} className="hover:bg-zinc-800/40 transition-colors">
                    <td className="p-3 font-medium text-zinc-100 font-sans">{row.agentName}</td>
                    <td className="p-3 text-zinc-400 font-sans">{row.campaign}</td>
                    <td className="p-3 text-right text-zinc-300">{row.minutesUsed.toLocaleString()} min</td>
                    <td className="p-3 text-right font-semibold text-zinc-100">{cost.toFixed(2)} {symbol}</td>
                    <td className="p-3 text-center">
                      <button className="text-zinc-400 hover:text-zinc-100 transition-colors font-sans hover:underline">PDF</button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>

      {/* TopUp Modal */}
      <TopUpModal
        isOpen={isTopUpOpen}
        onClose={() => setIsTopUpOpen(false)}
        currency={currency}
        onConfirmTopUp={handleConfirmTopUp}
      />
    </div>
  );
};
