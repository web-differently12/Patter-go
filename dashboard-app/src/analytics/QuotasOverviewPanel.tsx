import React from 'react';
import type { TenantQuotas, CreditRate } from '../components/integrations/types';

export interface QuotasOverviewPanelProps {
  quotas: TenantQuotas;
}

export const CREDIT_RATES: CreditRate[] = [
  { service: 'Voix Standard (STT + TTS Standard)', unit: 'minute', creditsPerUnit: 15, note: 'Appel vocal entrant/sortant standard' },
  { service: 'Voix Ultra-HD (STT + TTS HD)', unit: 'minute', creditsPerUnit: 60, note: 'Haute fidelite vocale & Krisp NC' },
  { service: 'Meeting Assistant (Prise de notes & Synthese visio)', unit: 'minute', creditsPerUnit: 40, note: 'Bot Zoom / Meet / Teams / Webex' },
  { service: 'Live Meeting Agent (Agent vocal actif)', unit: 'minute', creditsPerUnit: 70, note: 'Agent vocal interactif en visio' },
  { service: 'Live Avatar Stream (Visio avec Avatar)', unit: 'minute', creditsPerUnit: 120, note: 'Rendu vidéo Simli / WaveSpeed' },
  { service: 'SMS Sortant', unit: 'SMS', creditsPerUnit: 35, note: 'Envoi SMS direct' },
  { service: 'Passerelle MCP Externe API', unit: 'appel outil', creditsPerUnit: 0, note: 'Gratuit, seule l\'action est debitee' },
];

export const QuotasOverviewPanel: React.FC<QuotasOverviewPanelProps> = ({ quotas }) => {
  const isUnlimitedWhatsApp = quotas.whatsappMessagesLimitMonth === -1;
  const isUnlimitedMCP = quotas.mcpClientsLimit === -1;

  const waMsgPct = isUnlimitedWhatsApp
    ? 0
    : Math.min(100, (quotas.whatsappMessagesUsedMonth / quotas.whatsappMessagesLimitMonth) * 100);

  const waSessPct = Math.min(100, (quotas.whatsappSessionsActive / quotas.whatsappSessionsLimit) * 100);

  const mcpPct = isUnlimitedMCP
    ? 0
    : Math.min(100, (quotas.mcpActiveClients / quotas.mcpClientsLimit) * 100);

  return (
    <div className="space-y-6 font-sans">
      {/* Quotas & Sessions Panel */}
      <div className="bg-zinc-900/50 backdrop-blur-md border border-zinc-800/80 rounded-xl p-5 hover:border-zinc-700/80 transition-all duration-150 space-y-4 text-xs">
        <div className="flex justify-between items-center border-b border-zinc-800/60 pb-3">
          <div>
            <h3 className="text-sm font-semibold tracking-tight text-zinc-100">Utilisation des Quotas et Sessions Multi-Tenant</h3>
            <p className="text-[11px] text-zinc-400 mt-0.5">Suivi en direct des limites de votre plan Stripe ({quotas.planTier.toUpperCase()}).</p>
          </div>
          <span className="text-[11px] font-mono text-zinc-400 bg-zinc-950 px-2.5 py-1 rounded border border-zinc-800">
            Reset cycle: {new Date(quotas.billingCycleResetAt).toLocaleDateString()}
          </span>
        </div>

        {/* 3 Meters Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-5">
          {/* Meter 1: WhatsApp Monthly Messages */}
          <div className="bg-zinc-950/80 border border-zinc-800/80 rounded-lg p-4 space-y-2">
            <div className="flex justify-between items-center text-zinc-300 font-medium">
              <span>Messages WhatsApp du Mois</span>
              <span className="text-[10px] bg-emerald-950 text-emerald-300 border border-emerald-800 px-2 py-0.5 rounded font-mono font-bold">
                Inclus dans le plan
              </span>
            </div>
            <div className="text-base font-semibold text-zinc-100 font-mono">
              {quotas.whatsappMessagesUsedMonth.toLocaleString()} / {isUnlimitedWhatsApp ? 'Illimite' : quotas.whatsappMessagesLimitMonth.toLocaleString()}
            </div>
            {/* Fine 4px bar */}
            <div className="w-full bg-zinc-800 h-1 rounded-full overflow-hidden">
              <div className="bg-zinc-200 h-full rounded-full" style={{ width: `${waMsgPct}%` }}></div>
            </div>
            <div className="text-[10px] text-zinc-500">
              {isUnlimitedWhatsApp ? 'Quota mensuel illimite' : `${(100 - waMsgPct).toFixed(1)}% restants. Puise dans le Wallet ensuite (5 Credits/msg).`}
            </div>
          </div>

          {/* Meter 2: WhatsApp Active Sessions */}
          <div className="bg-zinc-950/80 border border-zinc-800/80 rounded-lg p-4 space-y-2">
            <div className="flex justify-between items-center text-zinc-300 font-medium">
              <span>Instances WhatsApp Actives</span>
              <span className="text-[10px] bg-zinc-800 text-zinc-300 px-2 py-0.5 rounded font-mono">
                Evolution / WAHA
              </span>
            </div>
            <div className="text-base font-semibold text-zinc-100 font-mono">
              {quotas.whatsappSessionsActive} / {quotas.whatsappSessionsLimit} sessions
            </div>
            {/* Fine 4px bar */}
            <div className="w-full bg-zinc-800 h-1 rounded-full overflow-hidden">
              <div className="bg-zinc-200 h-full rounded-full" style={{ width: `${waSessPct}%` }}></div>
            </div>
            <div className="text-[10px] text-zinc-500">
              Guard session WhatsApp : Erreur 403 QuotaExceeded au-dela.
            </div>
          </div>

          {/* Meter 3: MCP External Clients */}
          <div className="bg-zinc-950/80 border border-zinc-800/80 rounded-lg p-4 space-y-2">
            <div className="flex justify-between items-center text-zinc-300 font-medium">
              <span>Clients MCP / API Actifs</span>
              <span className="text-[10px] bg-zinc-800 text-zinc-300 px-2 py-0.5 rounded font-mono">
                MCP Gateway
              </span>
            </div>
            <div className="text-base font-semibold text-zinc-100 font-mono">
              {quotas.mcpActiveClients} / {isUnlimitedMCP ? 'Illimite' : quotas.mcpClientsLimit}
            </div>
            {/* Fine 4px bar */}
            <div className="w-full bg-zinc-800 h-1 rounded-full overflow-hidden">
              <div className="bg-zinc-200 h-full rounded-full" style={{ width: `${mcpPct}%` }}></div>
            </div>
            <div className="text-[10px] text-zinc-500">
              Guard MCP Client : Connexions externes simultanees.
            </div>
          </div>
        </div>
      </div>

      {/* Unified Credit Rates Table */}
      <div className="bg-zinc-900/50 backdrop-blur-md border border-zinc-800/80 rounded-xl p-5 hover:border-zinc-700/80 transition-all duration-150 space-y-3 text-xs">
        <div className="flex justify-between items-center border-b border-zinc-800/60 pb-3">
          <div>
            <h3 className="text-sm font-semibold tracking-tight text-zinc-100">Bareme de Debit Unifie en Credits KallFlow</h3>
            <p className="text-[11px] text-zinc-400 mt-0.5">Equivalence stricte : 1 EUR = 1 000 Credits KallFlow.</p>
          </div>
          <span className="text-xs font-mono text-cyan-400 bg-zinc-950 px-2.5 py-1 rounded border border-zinc-800">
            Realtime Event Ledger Engine
          </span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs text-zinc-300">
            <thead className="bg-zinc-950/80 text-zinc-400 text-[11px] font-medium border-b border-zinc-800/80">
              <tr>
                <th className="p-2.5">Service KallFlow Engine</th>
                <th className="p-2.5">Unite</th>
                <th className="p-2.5 text-right font-mono">Tarif (Credits)</th>
                <th className="p-2.5 text-right font-mono">Equivalence EUR</th>
                <th className="p-2.5">Note</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800/60 font-mono text-[11px]">
              {CREDIT_RATES.map((rate, idx) => (
                <tr key={idx} className="hover:bg-zinc-800/40 transition-colors">
                  <td className="p-2.5 font-medium text-zinc-100 font-sans">{rate.service}</td>
                  <td className="p-2.5 text-zinc-400 font-sans">{rate.unit}</td>
                  <td className="p-2.5 text-right font-bold text-cyan-400">{rate.creditsPerUnit} Credits</td>
                  <td className="p-2.5 text-right text-zinc-300">{(rate.creditsPerUnit / 1000).toFixed(3)} EUR</td>
                  <td className="p-2.5 text-zinc-500 font-sans text-[10px]">{rate.note}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};
