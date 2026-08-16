import React from 'react';

export interface AgentCostBreakdown {
  agentName: string;
  campaign: string;
  minutesUsed: number;
  costEUR: number;
}

export const BillingPage: React.FC = () => {
  const agentCosts: AgentCostBreakdown[] = [
    { agentName: 'Inbound Support Fr', campaign: 'Support Client 24/7', minutesUsed: 4200, costEUR: 176.40 },
    { agentName: 'Outbound B2B Qualifier', campaign: 'Prospection Q1 SaaS', minutesUsed: 3100, costEUR: 130.20 },
    { agentName: 'WhatsApp Assistant', campaign: 'Relance Leads Inactifs', minutesUsed: 1800, costEUR: 45.00 },
  ];

  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      {/* Header */}
      <div className="border-b border-slate-800 pb-4">
        <h1 className="text-2xl font-bold tracking-tight text-white">Facturation & Coûts Détaillés par Agent</h1>
        <p className="text-sm text-slate-400">Transparence totale sur la consommation de vos agents et rechargement de wallet unifié (Hyperswitch).</p>
      </div>

      {/* Wallet Balance & Payment Methods */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg flex flex-col justify-between">
          <div>
            <div className="text-xs uppercase tracking-wider text-slate-400 font-medium">Solde Wallet Unifié</div>
            <div className="text-4xl font-extrabold text-emerald-400 mt-2">1,248.50 €</div>
            <p className="text-xs text-slate-400 mt-2">Rechargement automatique quand le solde &lt; 100€</p>
          </div>
          <button className="mt-4 w-full bg-cyan-500 hover:bg-cyan-400 text-slate-950 font-bold py-2 px-4 rounded-lg text-xs transition-colors">
            + Recharger le Wallet (Hyperswitch)
          </button>
        </div>

        <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg col-span-2">
          <h3 className="font-bold text-white text-sm mb-3">Moyens de Paiement Intégrés (Hyperswitch)</h3>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 text-xs">
            <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 flex flex-col justify-between">
              <span className="font-semibold text-white">💳 Carte Bancaire</span>
              <span className="text-slate-400 text-[11px] mt-2">Visa / Mastercard (•••• 4242)</span>
            </div>
            <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 flex flex-col justify-between">
              <span className="font-semibold text-white">🏦 Prélèvement SEPA</span>
              <span className="text-slate-400 text-[11px] mt-2">IBAN FR76 •••• 8901</span>
            </div>
            <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 flex flex-col justify-between">
              <span className="font-semibold text-white">⚡ Bitcoin / Crypto</span>
              <span className="text-slate-400 text-[11px] mt-2">BTC / USDT via Hyperswitch</span>
            </div>
          </div>
        </div>
      </div>

      {/* Cost Breakdown per Agent */}
      <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-6 shadow-lg">
        <h2 className="text-lg font-bold text-white mb-4">Consommation Réelle par Agent & Campagne</h2>
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs text-slate-300">
            <thead className="bg-slate-950 text-slate-400 uppercase tracking-wider text-[10px]">
              <tr>
                <th className="p-3">Agent</th>
                <th className="p-3">Campagne</th>
                <th className="p-3 text-right">Minutes / Volume</th>
                <th className="p-3 text-right">Coût Total (€)</th>
                <th className="p-3 text-center">Facture PDF</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/60">
              {agentCosts.map((row, i) => (
                <tr key={i} className="hover:bg-slate-800/40">
                  <td className="p-3 font-semibold text-white">{row.agentName}</td>
                  <td className="p-3 text-slate-400">{row.campaign}</td>
                  <td className="p-3 text-right font-mono">{row.minutesUsed.toLocaleString()} min</td>
                  <td className="p-3 text-right font-bold text-cyan-400">{row.costEUR.toFixed(2)} €</td>
                  <td className="p-3 text-center">
                    <button className="text-cyan-400 hover:underline">Télécharger PDF</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};
