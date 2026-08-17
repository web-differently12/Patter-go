import React from 'react';

export const CampaignManagerPage: React.FC = () => {
  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      <div className="border-b border-zinc-800 pb-4 flex justify-between items-center">
        <div>
          <span className="text-xs font-mono font-bold text-cyan-400">Omnichannel Outbound Engine</span>
          <h1 className="text-2xl font-bold tracking-tight text-white mt-0.5">Gestionnaire de Campagnes Sortantes</h1>
        </div>
        <button className="bg-cyan-500 hover:bg-cyan-400 text-zinc-950 font-bold px-4 py-2 rounded-xl text-xs transition-colors shadow-sm">
          + Lancer une Campagne (Voix / WhatsApp)
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 text-xs">
        <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 space-y-3">
          <div className="flex justify-between items-center">
            <h3 className="font-bold text-white text-sm">Campagne Voice Outbound #01</h3>
            <span className="px-2 py-0.5 rounded bg-emerald-950 text-emerald-300 border border-emerald-800 font-mono font-bold">EN COURS</span>
          </div>
          <p className="text-zinc-400">Prospection Q1 SaaS - Jitter anti-spam +/- 30%</p>
          <div className="text-zinc-300 font-mono">Progression: 4 210 / 5 000 contacts</div>
        </div>
      </div>
    </div>
  );
};
