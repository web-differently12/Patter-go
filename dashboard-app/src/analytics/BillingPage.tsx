import React, { useState } from 'react';
import { BillingHeader } from './BillingHeader';
import { WalletOverviewCards } from './WalletOverviewCards';
import { LiveMeteringPanel } from './LiveMeteringPanel';
import { TopUpModal } from './TopUpModal';

export interface AgentCostBreakdown {
  agentName: string;
  campaign: string;
  minutesUsed: number;
  costEUR: number;
}

export const BillingPage: React.FC = () => {
  const [currency, setCurrency] = useState<'EUR' | 'USD'>('EUR');
  const [theme, setTheme] = useState<'dark' | 'light'>('dark');
  const [balanceEUR, setBalanceEUR] = useState<number>(1248.50);
  const [isTopUpOpen, setIsTopUpOpen] = useState(false);

  const agentCosts: AgentCostBreakdown[] = [
    { agentName: 'Inbound Support Fr', campaign: 'Support Client 24/7', minutesUsed: 4200, costEUR: 176.40 },
    { agentName: 'Outbound B2B Qualifier', campaign: 'Prospection Q1 SaaS', minutesUsed: 3100, costEUR: 130.20 },
    { agentName: 'WhatsApp Assistant', campaign: 'Relance Leads Inactifs', minutesUsed: 1800, costEUR: 45.00 },
  ];

  const handleConfirmTopUp = async (amount: number) => {
    setBalanceEUR((prev) => prev + amount);
  };

  const symbol = currency === 'EUR' ? '€' : '$';

  return (
    <div className={`space-y-6 p-6 min-h-screen font-sans transition-colors duration-200 ${
      theme === 'dark' ? 'bg-[#09090b] text-zinc-100' : 'bg-zinc-100 text-zinc-900'
    }`} data-theme={theme}>
      {/* 1. Header */}
      <BillingHeader
        currency={currency}
        setCurrency={setCurrency}
        theme={theme}
        setTheme={setTheme}
        onOpenTopUp={() => setIsTopUpOpen(true)}
      />

      {/* 2. Wallet Overview Cards */}
      <WalletOverviewCards
        balanceEUR={balanceEUR}
        currency={currency}
        satoshis={Math.round(balanceEUR * 1750)}
        estimatedDaysRemaining={67}
        agencyMarkupPercent={20}
      />

      {/* 3. Live Telemetry & Metering Panel */}
      <LiveMeteringPanel />

      {/* 4. Granular Agent Spend Table */}
      <div className="bg-zinc-900/50 backdrop-blur-md border border-zinc-800/80 rounded-xl p-6 hover:border-zinc-700/80 transition-all duration-150 space-y-4">
        <div className="flex justify-between items-center">
          <div>
            <h2 className="text-sm font-semibold tracking-tight text-zinc-100">Consommation Réelle par Agent & Campagne</h2>
            <p className="text-xs text-zinc-400 mt-0.5">Ventilation exacte des minutes télécoms et des tokens LLM consommés par vos sous-comptes.</p>
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
                <th className="p-3 text-right font-mono">Coût Total ({symbol})</th>
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
                      <button className="text-zinc-400 hover:text-zinc-100 transition-colors font-sans hover:underline">PDF →</button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>

      {/* 5. TopUp Modal */}
      <TopUpModal
        isOpen={isTopUpOpen}
        onClose={() => setIsTopUpOpen(false)}
        currency={currency}
        onConfirmTopUp={handleConfirmTopUp}
      />
    </div>
  );
};
