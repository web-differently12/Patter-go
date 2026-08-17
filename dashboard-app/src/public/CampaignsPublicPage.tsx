import React from 'react';

export const CampaignsPublicPage: React.FC = () => {
  return (
    <div className="min-h-screen bg-[#09090b] text-zinc-100 font-sans p-6 space-y-12 max-w-7xl mx-auto">
      {/* Header */}
      <div className="text-center space-y-3 pt-8">
        <span className="text-xs font-mono font-bold text-cyan-400 bg-zinc-950 px-3 py-1 rounded-full border border-zinc-800">
          Moteur de Campagnes Omnicanal
        </span>
        <h1 className="text-3xl sm:text-5xl font-extrabold tracking-tight text-zinc-100">
          WhatsApp, SMS, Voice & RCS Fallback
        </h1>
        <p className="text-xs sm:text-sm text-zinc-400 max-w-xl mx-auto">
          Dispatch haute performance avec validation E.164, anti-spam jitter dynamique (+/- 30%), et rotation round-robin multi-comptes.
        </p>
      </div>

      {/* Feature Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-6 shadow-xl space-y-3">
          <div className="text-cyan-400 font-mono font-bold text-sm">01 / Normalisation & Filtres</div>
          <h3 className="font-bold text-white text-base">Validation E.164 & Audiences</h3>
          <p className="text-xs text-zinc-400">
            Nettoyage automatique des numeros au format E.164 (+33 / +1) et filtrage dynamique par tags et categories d'audience.
          </p>
        </div>

        <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-6 shadow-xl space-y-3">
          <div className="text-cyan-400 font-mono font-bold text-sm">02 / Protections Anti-Bannissement</div>
          <h3 className="font-bold text-white text-base">Anti-Spam Jitter (+/- 30%)</h3>
          <p className="text-xs text-zinc-400">
            Calcul de delais aleatoires entre les envois pour simuler un comportement humain et eviter les filtres d'anti-spam.
          </p>
        </div>

        <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-6 shadow-xl space-y-3">
          <div className="text-cyan-400 font-mono font-bold text-sm">03 / Delivrabilite Maximale</div>
          <h3 className="font-bold text-white text-base">Round-Robin & Fallback RCS</h3>
          <p className="text-xs text-zinc-400">
            Rotation intelligente des sessions WhatsApp actives et bascule automatique sur SMS/RCS si le canal principal est indisponible.
          </p>
        </div>
      </div>
    </div>
  );
};
