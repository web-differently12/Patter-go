import React, { useState } from 'react';

export interface CampaignStats {
  id: string;
  name: string;
  channel: 'Voice' | 'WhatsApp' | 'SMS';
  totalOutbound: number;
  liveAnswerRate: number; // percentage
  voicemailRate: number; // percentage
  unreachableRate: number; // percentage
  optOutCount: number; // Stop SMS / Block
  conversionRate: number; // percentage
}

export const CampaignsPage: React.FC = () => {
  const [selectedChannel, setChannel] = useState<'all' | 'Voice' | 'WhatsApp' | 'SMS'>('all');

  const campaigns: CampaignStats[] = [
    {
      id: 'cmp_01',
      name: 'Prospection Q1 - SaaS B2B',
      channel: 'Voice',
      totalOutbound: 4500,
      liveAnswerRate: 68.4,
      voicemailRate: 24.1,
      unreachableRate: 7.5,
      optOutCount: 14,
      conversionRate: 18.2,
    },
    {
      id: 'cmp_02',
      name: 'Relance Leads Inactifs',
      channel: 'WhatsApp',
      totalOutbound: 8200,
      liveAnswerRate: 91.2, // read rate
      voicemailRate: 0,
      unreachableRate: 2.1,
      optOutCount: 38,
      conversionRate: 24.6,
    },
  ];

  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 border-b border-slate-800 pb-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-white">Analytics des Campagnes Sortantes</h1>
          <p className="text-sm text-slate-400">Délivrabilité omnicanale (Téléphone, WhatsApp, SMS), taux de réponse et A/B Testing de scripts.</p>
        </div>
        <div className="flex bg-slate-900 border border-slate-800 rounded-lg p-1 text-xs">
          {(['all', 'Voice', 'WhatsApp', 'SMS'] as const).map((c) => (
            <button
              key={c}
              onClick={() => setChannel(c)}
              className={`px-3 py-1.5 rounded-md font-medium transition-colors ${selectedChannel === c ? 'bg-cyan-500/20 text-cyan-400' : 'text-slate-400 hover:text-white'}`}
            >
              {c === 'all' ? 'Tous' : c}
            </button>
          ))}
        </div>
      </div>

      {/* Deliverability & Response Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {campaigns
          .filter((c) => selectedChannel === 'all' || c.channel === selectedChannel)
          .map((cmp) => (
            <div key={cmp.id} className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
              <div className="flex justify-between items-start">
                <div>
                  <div className="flex items-center gap-2">
                    <span className={`text-[10px] font-bold px-2 py-0.5 rounded ${
                      cmp.channel === 'Voice' ? 'bg-indigo-500/20 text-indigo-400 border border-indigo-500/30' :
                      cmp.channel === 'WhatsApp' ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30' :
                      'bg-cyan-500/20 text-cyan-400 border border-cyan-500/30'
                    }`}>
                      {cmp.channel}
                    </span>
                    <h3 className="font-bold text-white text-base">{cmp.name}</h3>
                  </div>
                  <div className="text-xs text-slate-400 mt-1">{cmp.totalOutbound.toLocaleString()} contacts engagés</div>
                </div>
                <span className="text-xs font-semibold px-2.5 py-1 rounded bg-emerald-500/20 text-emerald-400 border border-emerald-500/30">
                  {cmp.conversionRate}% Conversion
                </span>
              </div>

              {/* Delivery Progress */}
              <div className="space-y-2 text-xs">
                <div className="flex justify-between text-slate-300 font-medium">
                  <span>Décroché Réel / Lecture</span>
                  <span className="text-emerald-400 font-bold">{cmp.liveAnswerRate}%</span>
                </div>
                <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden flex">
                  <div className="bg-emerald-400 h-full" style={{ width: `${cmp.liveAnswerRate}%` }}></div>
                  <div className="bg-amber-400 h-full" style={{ width: `${cmp.voicemailRate}%` }}></div>
                  <div className="bg-slate-600 h-full" style={{ width: `${cmp.unreachableRate}%` }}></div>
                </div>
                <div className="flex justify-between text-[11px] text-slate-400 pt-1">
                  <span>🟢 Décroché ({cmp.liveAnswerRate}%)</span>
                  {cmp.channel === 'Voice' && <span>🟡 Répondeur ({cmp.voicemailRate}%)</span>}
                  <span>⚫ Injoignable ({cmp.unreachableRate}%)</span>
                </div>
              </div>

              {/* Opt-out / STOP alerts */}
              <div className="flex justify-between items-center bg-slate-950 p-3 rounded-lg border border-slate-800 text-xs">
                <span className="text-slate-400">Désabonnements / STOP SMS / Blocages</span>
                <span className="text-rose-400 font-bold bg-rose-500/10 px-2 py-0.5 rounded border border-rose-500/20">
                  {cmp.optOutCount} opt-outs
                </span>
              </div>
            </div>
          ))}
      </div>

      {/* A/B Script Comparison */}
      <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-6 shadow-lg">
        <h2 className="text-lg font-bold text-white mb-2">Comparatif de Scripts / A-B Testing Prompts</h2>
        <p className="text-xs text-slate-400 mb-6">Comparaison des performances entre 2 déclinaisons de prompts / voix d'agents.</p>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-slate-950 p-4 rounded-xl border border-cyan-500/30 relative">
            <span className="absolute top-3 right-3 text-[10px] font-bold px-2 py-0.5 rounded bg-cyan-500/20 text-cyan-400 border border-cyan-500/40">Gagnant A</span>
            <h4 className="font-bold text-white text-sm mb-1">Variant A: Accroche Directe ROI</h4>
            <p className="text-xs text-slate-400 mb-4">"Bonjour, je vous appelle concernant la réduction de vos coûts d'acquisition de 40%..."</p>
            <div className="grid grid-cols-2 gap-2 text-xs">
              <div className="bg-slate-900 p-2 rounded text-slate-300">RDV Pris: <strong className="text-emerald-400">22.4%</strong></div>
              <div className="bg-slate-900 p-2 rounded text-slate-300">Interruptions: <strong className="text-slate-200">1.2 / appel</strong></div>
            </div>
          </div>

          <div className="bg-slate-950 p-4 rounded-xl border border-slate-800">
            <h4 className="font-bold text-white text-sm mb-1">Variant B: Accroche Consultative</h4>
            <p className="text-xs text-slate-400 mb-4">"Bonjour, nous étudions les défis actuels des équipes Sales sur le secteur SaaS..."</p>
            <div className="grid grid-cols-2 gap-2 text-xs">
              <div className="bg-slate-900 p-2 rounded text-slate-300">RDV Pris: <strong className="text-amber-400">14.1%</strong></div>
              <div className="bg-slate-900 p-2 rounded text-slate-300">Interruptions: <strong className="text-rose-400">2.8 / appel</strong></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
