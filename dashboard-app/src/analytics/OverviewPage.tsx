import React, { useState } from 'react';

export interface OverviewKPIs {
  totalCalls: number;
  autonomousResolutionRate: number; // percentage
  bantQualificationRate: number; // percentage
  avgCostPerQualifiedLead: number; // in EUR
}

export interface FunnelStep {
  stage: string;
  durationRange: string;
  dropOffCount: number;
  dropOffRate: number;
  reason: string;
}

export const OverviewPage: React.FC = () => {
  const [dateRange, setRange] = useState<'7d' | '30d' | '90d'>('7d');
  const [channel, setChannel] = useState<'all' | 'voice' | 'whatsapp' | 'meeting'>('all');

  const kpis: OverviewKPIs = {
    totalCalls: 12480,
    autonomousResolutionRate: 84.2,
    bantQualificationRate: 31.5,
    avgCostPerQualifiedLead: 0.42,
  };

  const funnelSteps: FunnelStep[] = [
    { stage: 'Pitch Initial', durationRange: '0 - 10s', dropOffCount: 1420, dropOffRate: 11.3, reason: 'Pitch raté / Non intéressé' },
    { stage: 'Gestion Objections', durationRange: '10s - 30s', dropOffCount: 980, dropOffRate: 7.8, reason: 'Objection prix / timing non gérée' },
    { stage: 'Qualification BANT', durationRange: '30s - 2min', dropOffCount: 450, dropOffRate: 3.6, reason: 'Hors cible / Budget insuffisant' },
    { stage: 'Conclusion & Prise RDV', durationRange: '2min+', dropOffCount: 120, dropOffRate: 0.9, reason: 'Calendrier indisponible / Hésitation' },
  ];

  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      {/* Filters Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 border-b border-slate-800 pb-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-white">Vue d'Ensemble & Performance Métier</h1>
          <p className="text-sm text-slate-400">Suivi ROI, qualification BANT et efficacité de conversion des agents EchoFlow.</p>
        </div>
        <div className="flex items-center gap-3">
          <select
            value={channel}
            onChange={(e) => setChannel(e.target.value as any)}
            className="bg-slate-900 border border-slate-700 rounded-lg px-3 py-1.5 text-sm text-slate-200 focus:outline-none focus:border-cyan-500"
          >
            <option value="all">Tous les canaux</option>
            <option value="voice">Voix Téléphonie</option>
            <option value="whatsapp">WhatsApp / SMS</option>
            <option value="meeting">Réunions Visio</option>
          </select>
          <div className="flex bg-slate-900 border border-slate-800 rounded-lg p-1 text-xs">
            {(['7d', '30d', '90d'] as const).map((r) => (
              <button
                key={r}
                onClick={() => setRange(r)}
                className={`px-3 py-1 rounded-md transition-colors ${dateRange === r ? 'bg-cyan-500/20 text-cyan-400 font-semibold' : 'text-slate-400 hover:text-white'}`}
              >
                {r}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg">
          <div className="text-xs uppercase tracking-wider text-slate-400 font-medium">Total Conversions</div>
          <div className="text-3xl font-extrabold text-white mt-2">{kpis.totalCalls.toLocaleString()}</div>
          <div className="text-xs text-emerald-400 mt-2 flex items-center gap-1">
            <span>↑ +12.4%</span> <span className="text-slate-500">vs période précédente</span>
          </div>
        </div>

        <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg">
          <div className="text-xs uppercase tracking-wider text-slate-400 font-medium">Résolution Autonome IA</div>
          <div className="text-3xl font-extrabold text-emerald-400 mt-2">{kpis.autonomousResolutionRate}%</div>
          <div className="text-xs text-slate-400 mt-2">100% résolues sans transfert humain</div>
        </div>

        <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg">
          <div className="text-xs uppercase tracking-wider text-slate-400 font-medium">Qualification BANT / RDV</div>
          <div className="text-3xl font-extrabold text-cyan-400 mt-2">{kpis.bantQualificationRate}%</div>
          <div className="text-xs text-emerald-400 mt-2 flex items-center gap-1">
            <span>↑ +3.1%</span> <span className="text-slate-500">taux de validation</span>
          </div>
        </div>

        <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg">
          <div className="text-xs uppercase tracking-wider text-slate-400 font-medium">Coût / Lead Qualifié</div>
          <div className="text-3xl font-extrabold text-indigo-400 mt-2">{kpis.avgCostPerQualifiedLead.toFixed(2)} €</div>
          <div className="text-xs text-slate-400 mt-2">Économie de 85% vs call center traditionnel</div>
        </div>
      </div>

      {/* Drop-off Funnel */}
      <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-6 shadow-lg">
        <h2 className="text-lg font-bold text-white mb-2">Entonnoir de Décrochage (Drop-off Funnel)</h2>
        <p className="text-xs text-slate-400 mb-6">Analyse étape par étape des moments où les interlocuteurs mettent fin à l'échange.</p>

        <div className="space-y-4">
          {funnelSteps.map((step, idx) => (
            <div key={idx} className="bg-slate-950/60 border border-slate-800/80 rounded-lg p-4">
              <div className="flex justify-between items-center text-sm mb-2">
                <span className="font-semibold text-white flex items-center gap-2">
                  <span className="w-5 h-5 rounded-full bg-slate-800 text-xs flex items-center justify-center text-cyan-400">{idx + 1}</span>
                  {step.stage} <span className="text-xs text-slate-500">({step.durationRange})</span>
                </span>
                <span className="text-rose-400 font-medium text-xs bg-rose-500/10 px-2 py-0.5 rounded border border-rose-500/20">
                  {step.dropOffCount} décrochages ({step.dropOffRate}%)
                </span>
              </div>
              <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                <div
                  className="bg-gradient-to-r from-cyan-500 to-rose-500 h-full rounded-full"
                  style={{ width: `${Math.max(10, 100 - step.dropOffRate * 4)}%` }}
                />
              </div>
              <div className="text-xs text-slate-400 mt-2 flex justify-between">
                <span>Raison principale : <span className="text-slate-300">{step.reason}</span></span>
                <button className="text-cyan-400 hover:underline text-xs">Inspecter le sous-groupe →</button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
