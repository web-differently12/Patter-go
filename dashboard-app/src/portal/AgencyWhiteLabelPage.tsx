import React from 'react';

export const AgencyWhiteLabelPage: React.FC = () => {
  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      <div className="border-b border-zinc-800 pb-4 flex justify-between items-center">
        <div>
          <span className="text-xs font-mono font-bold text-cyan-400">Agency Reseller & Sub-Tenant Management</span>
          <h1 className="text-2xl font-bold tracking-tight text-white mt-0.5">Portail Agence Marque Blanche</h1>
        </div>
        <button className="bg-violet-600 hover:bg-violet-500 text-white font-bold px-4 py-2 rounded-xl text-xs transition-colors shadow-sm">
          + Ajouter un Sub-Tenant
        </button>
      </div>

      <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 shadow-xl space-y-4 text-xs">
        <div className="flex justify-between items-center border-b border-zinc-800/80 pb-3">
          <span className="font-semibold text-white">Vos Sub-Tenants (Capacite: 3 / 25)</span>
          <span className="text-xs font-mono text-violet-400">Marge Agence: +20%</span>
        </div>

        <div className="space-y-3 font-mono text-[11px]">
          <div className="p-3 bg-zinc-900/60 rounded-xl border border-zinc-800 flex justify-between items-center">
            <div>
              <div className="font-bold text-white">Client Alpha SARL</div>
              <div className="text-zinc-400 text-[10px]">ID: tenant_alpha_01 | Consommation: 176.40 €</div>
            </div>
            <button className="px-2.5 py-1 rounded bg-zinc-900 border border-zinc-800 text-zinc-300">Gérer Sub-Tenant →</button>
          </div>
        </div>
      </div>
    </div>
  );
};
